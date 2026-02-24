import axios from "axios";
import type { ApiError } from "../types/api";

export const client = axios.create({
  baseURL: import.meta.env.VITE_API_URL ?? "http://localhost:8080",
  withCredentials: true, // critical — sends httpOnly cookie with every request
  headers: {
    "Content-Type": "application/json",
  },
});

// Response interceptor — handle errors globally
client.interceptors.response.use(
  (response) => response,
  (error) => {
    const apiError: ApiError = error.response?.data ?? {
      error: "Something went wrong",
    };
    return Promise.reject(apiError);
  },
);
