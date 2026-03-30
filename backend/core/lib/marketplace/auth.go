package marketplace

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/sirupsen/logrus"
)

func (c *Client) generateTimestamp() int64 {
	return time.Now().Unix()
}

func (c *Client) generateSign(path string, timestamp int64, extra string) string {
	base := fmt.Sprintf("%s%s%d%s", c.partnerID, path, timestamp, extra)

	h := hmac.New(sha256.New, []byte(c.partnerKey))
	h.Write([]byte(base))

	return hex.EncodeToString(h.Sum(nil))
}

func (c *Client) Authorize(
	ctx context.Context,
	shopID string,
	state string,
	redirectURL string,
) (*BaseResponse[Authorize], error) {

	path := "/oauth/authorize"

	log := c.logger.WithFields(logrus.Fields{
		"shop_id": shopID,
		"path":    path,
	})

	log.Info("[marketplace] initiating authorization")

	timestamp := c.generateTimestamp()
	sign := c.generateSign(path, timestamp, shopID) // fix: pakai shopID param

	query := url.Values{}
	query.Set("shop_id", shopID)
	query.Set("state", state)
	query.Set("partner_id", c.partnerID)
	query.Set("timestamp", fmt.Sprintf("%d", timestamp))
	query.Set("sign", sign)
	query.Set("redirect", redirectURL)

	var result BaseResponse[Authorize]

	err := c.DoRequest(ctx, http.MethodGet, path+"?"+query.Encode(), nil, &result)
	if err != nil {
		log.WithError(err).Error("[marketplace] authorization request failed")
		return nil, err
	}

	log.WithField("state", state).Info("[marketplace] authorization request sent")

	return &result, nil
}

func (c *Client) GetToken(
	ctx context.Context,
	code string,
) (*BaseResponse[Token], error) {

	path := "/oauth/token"

	log := c.logger.WithFields(logrus.Fields{
		"shop_id":     c.shopID,
		"session_key": c.sessionKey,
		"path":        path,
	})

	log.Info("[marketplace] exchanging code for token")

	timestamp := c.generateTimestamp()
	sign := c.generateSign(path, timestamp, code)

	query := url.Values{}
	query.Set("partner_id", c.partnerID)
	query.Set("timestamp", fmt.Sprintf("%d", timestamp))
	query.Set("sign", sign)

	body := tokenRequest{
		GrantType: "authorization_code",
		Code:      code,
	}

	var result BaseResponse[Token]

	err := c.DoRequest(ctx, http.MethodPost, path+"?"+query.Encode(), body, &result)
	if err != nil {
		log.WithError(err).Error("[marketplace] get token request failed")
		return nil, err
	}

	// auto-save session setelah dapat token
	session := &Session{
		AccessToken:  result.Data.AccessToken,
		RefreshToken: result.Data.RefreshToken,
		ExpiredAt:    time.Now().Add(time.Duration(result.Data.ExpiresIn) * time.Second),
	}

	if err := c.store.Set(ctx, c.sessionKey, session); err != nil {
		log.WithError(err).Error("[marketplace] failed to save session after get token")
		return nil, err
	}

	c.apiToken = result.Data.AccessToken

	log.WithField("expires_at", session.ExpiredAt).
		Info("[marketplace] token acquired and session saved")

	return &result, nil
}

func (c *Client) RefreshToken(
	ctx context.Context,
	accessToken string,
	refreshToken string,
) (*BaseResponse[Token], error) {

	path := "/oauth/token"

	log := c.logger.WithFields(logrus.Fields{
		"shop_id": c.shopID,
		"path":    path,
	})

	log.Info("[marketplace] refreshing token")

	timestamp := c.generateTimestamp()
	sign := c.generateSign(path, timestamp, accessToken)

	query := url.Values{}
	query.Set("partner_id", c.partnerID)
	query.Set("timestamp", fmt.Sprintf("%d", timestamp))
	query.Set("sign", sign)

	body := tokenRequest{
		GrantType:    "refresh_token",
		RefreshToken: refreshToken,
	}

	var result BaseResponse[Token]

	err := c.DoRequest(ctx, http.MethodPost, path+"?"+query.Encode(), body, &result)
	if err != nil {
		log.WithError(err).Error("[marketplace] refresh token request failed")
		return nil, err
	}

	c.apiToken = result.Data.AccessToken

	log.Info("[marketplace] token refreshed successfully")

	return &result, nil
}
