import { client } from "./client";
import type {
  ApiResponse,
  Product,
  CreateProductRequest,
  UpdateProductRequest,
} from "../types/api";

export const productsApi = {
  list: async () => {
    const res = await client.get<ApiResponse<Product[]>>("/api/v1/products");
    return res.data;
  },

  getById: async (id: string) => {
    const res = await client.get<ApiResponse<Product>>(
      `/api/v1/products/${id}`,
    );
    return res.data;
  },

  create: async (data: CreateProductRequest) => {
    const res = await client.post<ApiResponse<Product>>(
      "/api/v1/products",
      data,
    );
    return res.data;
  },

  update: async (id: string, data: UpdateProductRequest) => {
    const res = await client.put<ApiResponse<Product>>(
      `/api/v1/products/${id}`,
      data,
    );
    return res.data;
  },

  delete: async (id: string) => {
    await client.delete(`/api/v1/products/${id}`);
  },
};
