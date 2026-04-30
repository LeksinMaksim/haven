import type { Todo, ApiResponse } from "./types";

const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8080/api";

async function fetchApi<T>(
  endpoint: string,
  options?: RequestInit,
): Promise<T> {
  const response = await fetch(`${API_URL}${endpoint}`, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...options?.headers,
    },
  });

  const json: ApiResponse<T> = await response.json();

  if (!response.ok || json.error) {
    throw new Error(json.error || `HTTP error! status: ${response.status}`);
  }

  return json.data as T;
}

export const todoApi = {
  getTodos: () => fetchApi<Todo[]>("/todos"),
  createTodo: (title: string, priority: string = "medium", due_date?: string) =>
    fetchApi<Todo>("/todos", {
      method: "POST",
      body: JSON.stringify({
        title,
        priority,
        // "2025-05-01" → "2025-05-01T00:00:00Z"
        due_date: due_date ? `${due_date}T00:00:00Z` : undefined,
      }),
    }),
  toggleDone: (id: number, done: boolean) =>
    fetchApi<Todo>(`/todos/${id}`, {
      method: "PATCH",
      body: JSON.stringify({ done }),
    }),
  deleteTodo: async (id: number) => {
    const response = await fetch(`${API_URL}/todos/${id}`, {
      method: "DELETE",
    });
    if (!response.ok) throw new Error("Failed to delete");
  },
};
