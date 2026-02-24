// mirrors Go's response wrapper
export interface ApiResponse<T> {
  data: T;
}

export interface ApiError {
  error: string;
  fields?: ValidationError[];
}

export interface ValidationError {
  field: string;
  message: string;
}

// mirrors Go's User struct
export interface User {
  id: string;
  email: string;
  role: string;
  created_at: string;
  updated_at: string;
}

// mirrors Go's Product struct
export interface Product {
  id: string;
  user_id: string;
  name: string;
  description: string | null;
  price: string;
  stock: number;
  created_at: string;
  updated_at: string;
}

// request types
export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
}

export interface CreateProductRequest {
  name: string;
  description?: string;
  price: string;
  stock: number;
}

export interface UpdateProductRequest {
  name?: string;
  description?: string;
  price?: string;
  stock?: number;
}
