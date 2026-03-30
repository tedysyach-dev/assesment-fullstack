package marketplace

import (
	"context"
	"net/http"
	"net/url"

	"github.com/sirupsen/logrus"
)

func (c *Client) OrderList(
	ctx context.Context,
) (*BaseResponse[[]Order], error) {

	path := "/order/list"

	log := c.logger.WithFields(logrus.Fields{
		"shop_id": c.shopID,
		"path":    path,
	})

	log.Info("[marketplace] fetching order list")

	if err := c.EnsureToken(ctx); err != nil {
		log.WithError(err).Error("[marketplace] failed to ensure token")
		return nil, err
	}

	var result BaseResponse[[]Order]

	err := c.DoRequest(ctx, http.MethodGet, path, nil, &result)
	if err != nil {
		log.WithError(err).Error("[marketplace] fetch order list failed")
		return nil, err
	}

	log.Info("[marketplace] order list fetched successfully")

	return &result, nil
}

func (c *Client) OrderDetail(ctx context.Context, orderSn string) (*BaseResponse[Order], error) {
	path := "/order/detail"

	log := c.logger.WithFields(logrus.Fields{
		"shop_id": c.shopID,
		"path":    path,
	})

	log.Info("[marketplace] fetching order detail")

	if err := c.EnsureToken(ctx); err != nil {
		log.WithError(err).Error("[marketplace] failed to ensure token")
		return nil, err
	}

	query := url.Values{}
	query.Set("order_sn", orderSn)

	var result BaseResponse[Order]

	err := c.DoRequest(ctx, http.MethodGet, path+"?"+query.Encode(), nil, &result)
	if err != nil {
		log.WithError(err).Error("[marketplace] fetch order detail failed")
		return nil, err
	}

	log.Info("[marketplace] order detail fetched successfully")

	return &result, nil
}
