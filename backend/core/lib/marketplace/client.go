package marketplace

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	URL        string `mapstructure:"url"`
	PartnerID  string `mapstructure:"partnerId"`
	PartnerKey string `mapstructure:"partnerKey"`
}

type Client struct {
	baseURL    string
	partnerID  string
	partnerKey string
	shopID     string
	httpClient *http.Client
	apiToken   string
	logger     *logrus.Logger
	store      SessionStore
	sessionKey string
}

type MarketplaceClientOption func(*Client)

func NewClient(viper *viper.Viper, opts ...MarketplaceClientOption) (*Client, error) {
	cfg := Config{
		URL:        viper.GetString("mock.url"),
		PartnerID:  viper.GetString("mock.partnerId"),
		PartnerKey: viper.GetString("mock.partnerKey"),
	}

	c := &Client{
		baseURL:    cfg.URL,
		partnerID:  cfg.PartnerID,
		partnerKey: cfg.PartnerKey,
		logger:     logrus.New(),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}

func WithSessionStore(store SessionStore, sessionKey string) MarketplaceClientOption {
	return func(c *Client) {
		c.store = store
		c.sessionKey = sessionKey
	}
}

func WithBaseURL(baseURL string) MarketplaceClientOption {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

func WithHTTPClient(httpClient *http.Client) MarketplaceClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

func WithTimeout(timeout time.Duration) MarketplaceClientOption {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

func WithToken(token string) MarketplaceClientOption {
	return func(c *Client) {
		c.apiToken = token
	}
}

func WithShopID(shopID string) MarketplaceClientOption {
	return func(c *Client) {
		c.shopID = shopID
	}
}

func WithLogger(logger *logrus.Logger) MarketplaceClientOption {
	return func(c *Client) {
		c.logger = logger
	}
}

func isForbiddenError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "api error status 403")
}

// retryDelay returns exponential backoff duration for a given attempt (0-indexed).
// attempt 0 → 500ms, attempt 1 → 1s, attempt 2 → 2s
func retryDelay(attempt int) time.Duration {
	return time.Duration(500<<attempt) * time.Millisecond
}

func (c *Client) EnsureToken(ctx context.Context) error {
	log := c.logger.WithFields(logrus.Fields{
		"session_key": c.sessionKey,
		"shop_id":     c.shopID,
	})

	log.Info("[marketplace] ensuring token")

	session, err := c.store.Get(ctx, c.sessionKey)
	if err != nil {
		log.WithError(err).Error("[marketplace] failed to get session")
		return err
	}

	if session == nil {
		log.Warn("[marketplace] session not found, not authenticated")
		at, err := c.Authorize(ctx, "shopee-123", "pm", "http://localhost:3000/callback")
		if err != nil {
			log.Warn("[bootstrap] authorize failed: ", err)
			return err
		}

		_, err = c.GetToken(ctx, at.Data.Code)
		if err != nil {
			log.Warn("[bootstrap] get token failed: ", err)
			return err
		}

		session, err = c.store.Get(ctx, c.sessionKey)
		if err != nil {
			log.WithError(err).Error("[marketplace] failed to get session")
			return err
		}

	} else if time.Now().After(session.ExpiredAt) {
		log.Warn("[marketplace] token expired, refreshing")

		resp, err := c.RefreshToken(ctx, session.AccessToken, session.RefreshToken)
		if err != nil {
			if isForbiddenError(err) {
				log.Warn("[marketplace] refresh token got 403, re-authenticating from scratch")

				at, authErr := c.Authorize(ctx, "shopee-123", "pm", "http://localhost:3000/callback")
				if authErr != nil {
					log.WithError(authErr).Error("[marketplace] re-authorize failed")
					return authErr
				}

				_, authErr = c.GetToken(ctx, at.Data.Code)
				if authErr != nil {
					log.WithError(authErr).Error("[marketplace] get token after re-auth failed")
					return authErr
				}

				session, authErr = c.store.Get(ctx, c.sessionKey)
				if authErr != nil {
					log.WithError(authErr).Error("[marketplace] failed to get session after re-auth")
					return authErr
				}

				log.Info("[marketplace] re-auth successful, session restored")
			} else {
				log.WithError(err).Error("[marketplace] failed to refresh token")
				return err
			}
		} else {
			session.AccessToken = resp.Data.AccessToken
			session.RefreshToken = resp.Data.RefreshToken
			session.ExpiredAt = time.Now().Add(
				time.Duration(resp.Data.ExpiresIn) * time.Second,
			)

			if err = c.store.Set(ctx, c.sessionKey, session); err != nil {
				log.WithError(err).Error("[marketplace] failed to update session after refresh")
				return err
			}

			log.WithField("expires_at", session.ExpiredAt).
				Info("[marketplace] token refreshed and session updated")
		}
	} else {
		log.WithField("expires_at", session.ExpiredAt).
			Info("[marketplace] token still valid")
	}

	c.apiToken = session.AccessToken

	return nil
}

func (c *Client) DoRequest(
	ctx context.Context,
	method string,
	path string,
	body interface{},
	result interface{},
) error {
	const maxAttempts = 3

	url := c.baseURL + path

	log := c.logger.WithFields(logrus.Fields{
		"method": method,
		"url":    url,
	})

	var lastErr error

	for attempt := 0; attempt < maxAttempts; attempt++ {
		if attempt > 0 {
			delay := retryDelay(attempt - 1)
			log.WithFields(logrus.Fields{
				"attempt": attempt + 1,
				"delay":   delay,
			}).Warn("[marketplace] retrying request")

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
		}

		log.WithField("attempt", attempt+1).Info("[marketplace] sending request")

		var reqBody io.Reader
		if body != nil {
			jsonData, err := json.Marshal(body)
			if err != nil {
				log.WithError(err).Error("[marketplace] failed to marshal request body")
				return fmt.Errorf("marshal body error: %w", err)
			}
			reqBody = bytes.NewBuffer(jsonData)
		}

		req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
		if err != nil {
			log.WithError(err).Error("[marketplace] failed to create request")
			return fmt.Errorf("create request error: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		if c.apiToken != "" {
			req.Header.Set("Authorization", "Bearer "+c.apiToken)
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			log.WithError(err).Warn("[marketplace] http call failed, will retry")
			lastErr = fmt.Errorf("http call error: %w", err)
			continue
		}

		respBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			log.WithError(err).Error("[marketplace] failed to read response body")
			return fmt.Errorf("read body error: %w", err)
		}

		log.WithField("status_code", resp.StatusCode).Info("[marketplace] response received")

		switch resp.StatusCode {

		case http.StatusUnauthorized: // 401
			log.Warn("[marketplace] got 401, refreshing token and retrying")
			if tokenErr := c.EnsureToken(ctx); tokenErr != nil {
				log.WithError(tokenErr).Error("[marketplace] token refresh after 401 failed")
				return fmt.Errorf("token refresh failed after 401: %w", tokenErr)
			}
			lastErr = fmt.Errorf("api error status 401")
			continue

		case http.StatusTooManyRequests: // 429
			log.WithField("attempt", attempt+1).
				Warn("[marketplace] got 429 rate limit, backing off")
			lastErr = fmt.Errorf("api error status 429: rate limited")
			continue

		case http.StatusInternalServerError: // 500
			log.WithFields(logrus.Fields{
				"attempt": attempt + 1,
				"body":    string(respBody),
			}).Warn("[marketplace] got 500, will retry")
			lastErr = fmt.Errorf("api error status 500: %s", string(respBody))
			continue

		default:
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				if result != nil {
					if err := json.Unmarshal(respBody, result); err == nil {
						return fmt.Errorf("api error status %d", resp.StatusCode)
					}
				}
				return fmt.Errorf("api error status %d: %s", resp.StatusCode, string(respBody))
			}
		}

		// Success — unmarshal and return
		if result != nil {
			if err = json.Unmarshal(respBody, result); err != nil {
				log.WithError(err).Error("[marketplace] failed to unmarshal response")
				return fmt.Errorf("unmarshal error: %w", err)
			}
		}

		return nil
	}

	log.WithField("attempts", maxAttempts).
		Error("[marketplace] all retry attempts exhausted")

	return fmt.Errorf("max retries exceeded after %d attempts: %w", maxAttempts, lastErr)
}
