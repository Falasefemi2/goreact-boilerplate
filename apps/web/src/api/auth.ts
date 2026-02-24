import { client } from "./client";
import type {
  ApiResponse,
  LoginRequest,
  RegisterRequest,
  User,
} from "../types/api";

export const authApi = {
  register: async (data: RegisterRequest) => {
    const res = await client.post<ApiResponse<{ message: string }>>(
      "/api/v1/auth/register",
      data,
    );
    return res.data;
  },

  login: async (data: LoginRequest) => {
    const res = await client.post<ApiResponse<{ message: string }>>(
      "/api/v1/auth/login",
      data,
    );
    return res.data;
  },

  logout: async () => {
    const res = await client.post<ApiResponse<{ message: string }>>(
      "/api/v1/auth/logout",
    );
    return res.data;
  },

  me: async () => {
    const res = await client.get<ApiResponse<User>>("/api/v1/auth/me");
    return res.data;
  },
};
