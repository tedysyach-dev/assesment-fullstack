import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import api from "@/lib/axios";
import type { Order, OrderDetail, ShipOrderPayload } from "../types";
import type { ApiResponse } from "@/types/response";

export const orderKeys = {
  all: ["orders"] as const,
  detail: (orderSn: string) => ["orders", orderSn] as const,
};

//Query
export const useOrders = () =>
  useQuery({
    queryKey: orderKeys.all,
    queryFn: () =>
      api.get<{ resource: Order[] }>("/order").then((r) => r.data.resource),
  });

export const useOrderDetail = (orderSn: string) =>
  useQuery({
    queryKey: orderKeys.detail(orderSn),
    queryFn: () =>
      api
        .get<{ resource: OrderDetail }>(`/order/${orderSn}`)
        .then((r) => r.data.resource),
    enabled: !!orderSn,
  });

//Mutations
export const usePickOrder = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (orderSn: string) =>
      api.post<ApiResponse<null>>(`/order/${orderSn}/pick`),

    onSuccess: (_, orderSn) => {
      // refresh list orders
      queryClient.invalidateQueries({
        queryKey: orderKeys.all,
      });

      // refresh detail order
      queryClient.invalidateQueries({
        queryKey: orderKeys.detail(orderSn),
      });
    },
  });
};

export const usePackOrder = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (orderSn: string) =>
      api.post<ApiResponse<null>>(`/order/${orderSn}/pack`),

    onSuccess: (_, orderSn) => {
      // refresh list
      queryClient.invalidateQueries({
        queryKey: orderKeys.all,
      });

      // refresh detail
      queryClient.invalidateQueries({
        queryKey: orderKeys.detail(orderSn),
      });
    },
  });
};

export const useShipOrder = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ orderSn, channelId }: ShipOrderPayload) =>
      api.post<ApiResponse<null>>(`/order/${orderSn}/ship`, {
        channelId,
      }),

    onSuccess: (_, { orderSn }) => {
      // refresh list
      queryClient.invalidateQueries({
        queryKey: orderKeys.all,
      });

      // refresh detail
      queryClient.invalidateQueries({
        queryKey: orderKeys.detail(orderSn),
      });
    },
  });
};
