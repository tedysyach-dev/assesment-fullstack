package marketplace

import (
	"context"
	"net/http"

	"github.com/sirupsen/logrus"
)

func (c *Client) ShipOrder(ctx context.Context, req ShipOrderRequest) (*BaseResponse[ShipOrderResponse], error) {
	path := "/logistic/ship"

	log := c.logger.WithFields(logrus.Fields{
		"shop_id": c.shopID,
		"path":    path,
	})

	log.Info("[marketplace] create ship order")

	if err := c.EnsureToken(ctx); err != nil {
		log.WithError(err).Error("[marketplace] failed to ensure token")
		return nil, err
	}

	var result BaseResponse[ShipOrderResponse]
	body := map[string]interface{}{
		"order_sn":   req.OrderSn,
		"channel_id": req.ChannelId,
	}

	err := c.DoRequest(ctx, http.MethodPost, path, body, &result)
	if err != nil {
		log.WithError(err).Error("[marketplace] ship order failed")
		return nil, err
	}

	log.Info("[marketplace] order shipped successfully")

	return &result, nil
}

func (c *Client) LogisticChannel(ctx context.Context) (*BaseResponse[[]LogisticChannelResponse], error) {
	path := "/logistic/channels"

	log := c.logger.WithFields(logrus.Fields{
		"shop_id": c.shopID,
		"path":    path,
	})

	log.Info("[marketplace] get logistic channel")

	if err := c.EnsureToken(ctx); err != nil {
		log.WithError(err).Error("[marketplace] failed to ensure token")
		return nil, err
	}

	var result BaseResponse[[]LogisticChannelResponse]
	err := c.DoRequest(ctx, http.MethodGet, path, nil, &result)
	if err != nil {
		log.WithError(err).Error("[marketplace] fetch logistic channel failed")
		return nil, err
	}

	log.Info("[marketplace] fetch logistic channel successfully")

	return &result, nil
}
