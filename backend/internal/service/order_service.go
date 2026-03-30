package service

import (
	"context"
	"strings"
	"time"
	"wms/core/base"
	"wms/core/errors"
	"wms/core/lib/marketplace"
	"wms/internal/converter"
	"wms/internal/entity"
	"wms/internal/model"
	"wms/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
)

type OrderService struct {
	DB                  *bun.DB
	Log                 *logrus.Logger
	Validate            *validator.Validate
	MarketplaceClient   *marketplace.Client
	OrderRepository     *repository.OrderRepository
	OrderItemRepository *repository.OrderItemRepository
}

func NewOrderService(
	db *bun.DB,
	logger *logrus.Logger,
	validate *validator.Validate,
	marketplaceClient *marketplace.Client,
	orderRepository *repository.OrderRepository,
	orderItemRepository *repository.OrderItemRepository,
) *OrderService {
	return &OrderService{
		DB:                  db,
		Log:                 logger,
		Validate:            validate,
		MarketplaceClient:   marketplaceClient,
		OrderRepository:     orderRepository,
		OrderItemRepository: orderItemRepository,
	}
}

// SyncOrders: fetch dari marketplace & simpan ke DB
func (s *OrderService) SyncOrders(ctx context.Context) error {
	res, err := s.MarketplaceClient.OrderList(ctx)
	if err != nil {
		s.Log.WithError(err).Error("[order] failed to fetch from marketplace")
		return err
	}

	var orders []entity.Order
	var orderItems []entity.OrderItem

	for _, v := range *res.Data {
		newOrder := entity.Order{
			ID:                    uuid.New().String(),
			OrderSN:               v.OrderSn,
			ShopID:                v.ShopID,
			MarketplaceStatus:     v.Status,
			ShippingStatus:        v.ShippingStatus,
			TotalAmount:           v.TotalAmount,
			RawMarketplacePayload: v.ToJson(),
			CreatedAt:             v.CreatedAt,
			UpdatedAt:             time.Now(),
		}

		s.mapWmsStatus(&newOrder, v)

		for _, i := range v.Items {
			orderItems = append(orderItems, entity.OrderItem{
				ID:        uuid.New().String(),
				OrderSn:   v.OrderSn,
				SKU:       i.Sku,
				Price:     i.Price,
				Quantity:  i.Quantity,
				CreatedAt: v.CreatedAt,
			})
		}

		orders = append(orders, newOrder)
	}

	return base.ExecuteInTransaction(ctx, s.DB, func(ctx context.Context) error {
		if err := s.OrderRepository.UpsertBulk(ctx, &orders,
			[]string{"order_sn"},
			[]string{"marketplace_status", "shipping_status"},
		); err != nil {
			s.Log.WithError(err).Error("[order] upsert orders failed")
			return err
		}

		s.Log.Infof("[order] upserted %d orders", len(orders))

		// Delete semua items lama berdasarkan order_id, lalu insert ulang
		orderIDs := make([]string, len(orders))
		for i, o := range orders {
			orderIDs[i] = o.OrderSN
		}

		if err := s.OrderItemRepository.DeleteByOrderSNs(ctx, orderIDs); err != nil {
			return err
		}

		if len(orderItems) > 0 {
			return s.OrderItemRepository.CreateBulk(ctx, &orderItems)
		}
		return nil
	})
}

func (s *OrderService) WebhookOrderStatus(ctx context.Context, req *model.WebhookOrderStatusBody) error {
	order := entity.Order{}
	if err := s.OrderRepository.FindOne(ctx, &order, base.WithWhere("order_sn = ?", req.Data.OrderSn)); err != nil {
		return errors.NewNotFoundError("order")
	}

	order.MarketplaceStatus = req.Data.Status
	order.UpdatedAt = time.Now()

	if err := s.OrderRepository.Update(ctx, &order); err != nil {
		return errors.NewInternalError(err)
	}

	return nil
}

func (s *OrderService) WebhookShippingStatus(ctx context.Context, req *model.WebhookShippingStatusBody) error {
	order := entity.Order{}
	if err := s.OrderRepository.FindOne(ctx, &order, base.WithWhere("order_sn = ?", req.Data.OrderSn)); err != nil {
		return errors.NewNotFoundError("order")
	}

	if order.ShippingStatus == req.Data.ShippingState {
		return nil
	}

	order.ShippingStatus = req.Data.ShippingState
	order.UpdatedAt = time.Now()

	if err := s.OrderRepository.Update(ctx, &order); err != nil {
		return errors.NewInternalError(err)
	}

	return nil
}

func (s *OrderService) OrderList(ctx context.Context, wmsStatus string) []model.OrderResponse {

	opts := []base.QueryOption{
		base.WithOrder("updated_at DESC"),
	}

	if wmsStatus != "" {
		wmsStatusFilter := strings.Split(wmsStatus, ",")

		opts = append(opts, base.WithWhere("wms_status IN (?)", wmsStatusFilter))
	}
	order := []entity.Order{}
	s.OrderRepository.FindAll(ctx, &order, opts...)

	return converter.OrderListToResponse(&order)
}

func (s *OrderService) OrderDetail(ctx context.Context, orderSn string) (*model.OrderResponse, error) {
	order := entity.Order{}
	s.OrderRepository.FindOne(ctx, &order, base.WithWhere("order_sn = ?", orderSn))

	if order.ID == "" {
		return nil, errors.NewNotFoundError("order")
	}

	orderItem := []entity.OrderItem{}
	s.OrderItemRepository.FindAll(ctx, &orderItem, base.WithWhere("order_sn = ?", orderSn))
	res := converter.OrderDetailToResponse(&order, &orderItem)
	return &res, nil
}

func (s *OrderService) PickOrder(ctx context.Context, orderSn string) error {
	order := entity.Order{}
	if err := s.OrderRepository.FindOne(ctx, &order, base.WithWhere("order_sn = ?", orderSn)); err != nil {
		return errors.NewNotFoundError("order")
	}

	if order.WmsStatus != entity.WmsStatusReadyToPick {
		return errors.NewBadRequestError("action is forbiden", map[string]any{
			"wms_status": order.WmsStatus,
		})
	}

	order.WmsStatus = entity.WmsStatusPicking
	order.UpdatedAt = time.Now()

	if err := s.OrderRepository.Update(ctx, &order); err != nil {
		return errors.NewInternalError(err)
	}

	return nil
}

func (s *OrderService) PackOrder(ctx context.Context, orderSn string) error {
	order := entity.Order{}
	if err := s.OrderRepository.FindOne(ctx, &order, base.WithWhere("order_sn = ?", orderSn)); err != nil {
		return errors.NewNotFoundError("order")
	}

	if order.WmsStatus != entity.WmsStatusPicking {
		return errors.NewBadRequestError("action is forbiden", map[string]any{
			"wms_status": order.WmsStatus,
		})
	}

	order.WmsStatus = entity.WmsStatusPacked
	order.UpdatedAt = time.Now()

	if err := s.OrderRepository.Update(ctx, &order); err != nil {
		return errors.NewInternalError(err)
	}

	return nil
}

func (s *OrderService) ShipOrder(ctx context.Context, orderSn string, req *model.ShipOrderRequest) (*model.ShipOrderResponse, error) {

	if err := s.Validate.Struct(req); err != nil {
		return nil, errors.NewValidationError(err)
	}

	order := entity.Order{}
	if err := s.OrderRepository.FindOne(ctx, &order, base.WithWhere("order_sn = ?", orderSn)); err != nil {
		return nil, errors.NewNotFoundError("order")
	}

	if order.WmsStatus != entity.WmsStatusPacked {
		return nil, errors.NewBadRequestError("action is forbiden", map[string]any{
			"wms_status": order.WmsStatus,
		})
	}

	logisticChannels, err := s.MarketplaceClient.LogisticChannel(ctx)
	if err != nil {
		return nil, err
	}

	channelExists := false
	for _, v := range *logisticChannels.Data {
		if v.ID == req.ChannelId {
			channelExists = true
			break
		}
	}

	if !channelExists {
		return nil, errors.NewBadRequestError("channel id not found", map[string]any{
			"order_sn":   orderSn,
			"channel_id": req.ChannelId,
		})
	}

	shipRequest := marketplace.ShipOrderRequest{
		OrderSn:   order.OrderSN,
		ChannelId: req.ChannelId,
	}

	res, err := s.MarketplaceClient.ShipOrder(ctx, shipRequest)
	if err != nil {
		return nil, err
	}

	order.WmsStatus = entity.WmsStatusShiped
	order.ShippingStatus = res.Data.ShippingStatus
	order.TrackingNumber = &res.Data.TrackingNo
	order.UpdatedAt = time.Now()

	if err := s.OrderRepository.Update(ctx, &order); err != nil {
		return nil, errors.NewInternalError(err)
	}

	result := model.ShipOrderResponse{
		OrderSn:        order.OrderSN,
		WmsStatus:      order.WmsStatus,
		ShippingStatus: res.Data.ShippingStatus,
		TrackingNumber: res.Data.TrackingNo,
	}

	return &result, nil
}

func (s *OrderService) mapWmsStatus(order *entity.Order, v marketplace.Order) {
	tracking := v.TrackingNumber
	switch v.ShippingStatus {
	case "awaiting_pickup", "label_created":
		order.WmsStatus = entity.WmsStatusReadyToPick
		order.TrackingNumber = nil
	case "shipped":
		order.WmsStatus = entity.WmsStatusShiped
		order.TrackingNumber = &tracking
	case "delivered":
		order.WmsStatus = entity.WmsStatusShiped
		order.TrackingNumber = &tracking
	case "cancelled":
		order.WmsStatus = entity.WmsStatusShiped
		order.TrackingNumber = &tracking
	default:
		order.WmsStatus = entity.WmsStatusReadyToPick
		order.TrackingNumber = nil
	}
}
