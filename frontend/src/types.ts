export type Priority = "low" | "medium" | "high";

export interface Todo {
  id: number;
  title: string;
  description?: string;
  done: boolean;
  priority: string;
  due_date?: string;
  created_at: string;
  updated_at: string;
}

export interface ApiResponse<T> {
  data?: T;
  error?: string;
}
