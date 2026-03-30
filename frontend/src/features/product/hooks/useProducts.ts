import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import api from "@/lib/axios";
import type { Product, CreateProductPayload } from "../types";

// Query keys — disimpan di satu tempat agar konsisten
export const productKeys = {
  all: ["products"] as const,
  detail: (id: number) => ["products", id] as const,
};

// GET semua products
export const useProducts = () =>
  useQuery({
    queryKey: productKeys.all,
    queryFn: () => api.get<Product[]>("/products").then((r) => r.data),
  });

// GET satu product by id
export const useProduct = (id: number) =>
  useQuery({
    queryKey: productKeys.detail(id),
    queryFn: () => api.get<Product>(`/products/${id}`).then((r) => r.data),
    enabled: !!id,
  });

// POST — tambah product baru
export const useCreateProduct = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: CreateProductPayload) =>
      api.post<Product>("/products", payload).then((r) => r.data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: productKeys.all });
    },
  });
};

// DELETE — hapus product
export const useDeleteProduct = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: number) => api.delete(`/products/${id}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: productKeys.all });
    },
  });
};
