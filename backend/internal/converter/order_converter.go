package converter

import (
	"wms/internal/entity"
	"wms/internal/model"
)

func OrderListToResponse(order *[]entity.Order) []model.OrderResponse {
	var response []model.OrderResponse

	for _, v := range *order {
		response = append(response, model.OrderResponse{
			ID:                v.ID,
			OrderSN:           v.OrderSN,
			ShopID:            v.ShopID,
			MarketplaceStatus: v.MarketplaceStatus,
			ShippingStatus:    v.ShippingStatus,
			WmsStatus:         v.WmsStatus,
			TrackingNumber:    v.TrackingNumber,
			TotalAmount:       v.TotalAmount,
			CreatedAt:         v.CreatedAt,
			UpdatedAt:         v.UpdatedAt,
		})
	}

	return response
}

func OrderDetailToResponse(order *entity.Order, orderItem *[]entity.OrderItem) model.OrderResponse {
	var response model.OrderResponse

	orderItemMap := make(map[string][]model.OrderItemRes)
	for _, v := range *orderItem {
		orderItemMap[v.OrderSn] = append(orderItemMap[v.OrderSn], model.OrderItemRes{
			ID:        v.ID,
			OrderID:   v.OrderSn,
			SKU:       v.SKU,
			Quantity:  v.Quantity,
			Price:     v.Price,
			CreatedAt: v.CreatedAt,
		})
	}

	response = model.OrderResponse{
		ID:                order.ID,
		OrderSN:           order.OrderSN,
		ShopID:            order.ShopID,
		MarketplaceStatus: order.MarketplaceStatus,
		ShippingStatus:    order.ShippingStatus,
		WmsStatus:         order.WmsStatus,
		TrackingNumber:    order.TrackingNumber,
		OrderItem:         orderItemMap[order.OrderSN],
		TotalAmount:       order.TotalAmount,
		CreatedAt:         order.CreatedAt,
		UpdatedAt:         order.UpdatedAt,
	}

	return response
}
