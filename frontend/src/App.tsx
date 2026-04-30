import { useEffect, useState } from "react";
import type { Todo } from "./types";
import { todoApi } from "./api";

// ─── helpers ────────────────────────────────────────────────────────────────

const PRIORITY_COLOR: Record<string, string> = {
  high: "text-red-600",
  medium: "text-amber-600",
  low: "text-blue-600",
};

const HARD_SHADOW = { boxShadow: "4px 4px 0 #9ca3af" } as const;

function formatUptime(seconds: number) {
  const h = Math.floor(seconds / 3600)
    .toString()
    .padStart(2, "0");
  const m = Math.floor((seconds % 3600) / 60)
    .toString()
    .padStart(2, "0");
  const s = (seconds % 60).toString().padStart(2, "0");
  return `${h}:${m}:${s}`;
}

function formatDueDate(iso: string) {
  return new Date(iso).toLocaleDateString("ru-RU", {
    day: "2-digit",
    month: "2-digit",
    year: "2-digit",
  });
}

function dueDateStatus(iso: string): "ok" | "soon" | "overdue" {
  const diff = new Date(iso).getTime() - Date.now();
  if (diff < 0) return "overdue";
  if (diff < 1000 * 60 * 60 * 24 * 2) return "soon";
  return "ok";
}

const DUE_DATE_CLS = {
  ok: "text-green-700",
  soon: "text-amber-600",
  overdue: "text-red-600",
};

// ─── Window component ────────────────────────────────────────────────────────

function Window({
  title,
  right,
  children,
  className = "",
}: {
  title: string;
  right?: string;
  children: React.ReactNode;
  className?: string;
}) {
  return (
    <div
      className={`border border-gray-800 bg-white ${className}`}
      style={HARD_SHADOW}
    >
      <div className="bg-gray-900 px-3 py-1 flex items-center justify-between">
        <span className="text-gray-100 tracking-wider text-xs">{title}</span>
        {right && <span className="text-gray-500 text-xs">{right}</span>}
      </div>
      {children}
    </div>
  );
}

// ─── App ─────────────────────────────────────────────────────────────────────

export default function App() {
  const [todos, setTodos] = useState<Todo[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [newTitle, setNewTitle] = useState("");
  const [priority, setPriority] = useState<"low" | "medium" | "high">("medium");
  const [dueDate, setDueDate] = useState("");
  const [selectedId, setSelectedId] = useState<number | null>(null);
  const [uptime, setUptime] = useState(0);
  const [now, setNow] = useState(new Date());

  useEffect(() => {
    todoApi
      .getTodos()
      .then((data) => setTodos(data || []))
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
  }, []);

  useEffect(() => {
    const id = setInterval(() => {
      setUptime((s) => s + 1);
      setNow(new Date());
    }, 1000);
    return () => clearInterval(id);
  }, []);

  const handleCreate = async () => {
    if (!newTitle.trim()) return;
    try {
      const todo = await todoApi.createTodo(
        newTitle,
        priority,
        dueDate || undefined,
      );
      setTodos((prev) => [...prev, todo]);
      setNewTitle("");
      setDueDate("");
      setError(null);
    } catch (err: any) {
      setError(err.message);
    }
  };

  const handleToggle = async (todo: Todo) => {
    try {
      const updated = await todoApi.toggleDone(todo.id, !todo.done);
      setTodos((prev) => prev.map((t) => (t.id === todo.id ? updated : t)));
    } catch (err: any) {
      setError(err.message);
    }
  };

  const handleDelete = async (id: number) => {
    try {
      await todoApi.deleteTodo(id);
      setTodos((prev) => prev.filter((t) => t.id !== id));
      setSelectedId(null);
    } catch (err: any) {
      setError(err.message);
    }
  };

  const done = todos.filter((t) => t.done).length;
  const total = todos.length;
  const progress = total ? Math.round((done / total) * 100) : 0;

  const dateStr = now.toLocaleDateString("ru-RU", {
    weekday: "short",
    day: "2-digit",
    month: "2-digit",
    year: "numeric",
  });
  const timeStr = now.toLocaleTimeString("ru-RU");

  return (
    <div className="min-h-screen bg-stone-100 text-gray-900 font-mono p-4 flex flex-col gap-3">
      {/* Top bar */}
      <div className="border border-gray-800 bg-white" style={HARD_SHADOW}>
        <div className="px-4 py-2 flex items-center justify-between">
          <span className="font-bold tracking-widest text-gray-900">
            ~/Haven
          </span>
          <span className="text-gray-400 text-xs">v1.0.0</span>
        </div>
        <div className="border-t border-gray-200 px-4 py-1 flex gap-6 text-xs text-gray-500">
          <span>
            UPTIME:{" "}
            <span className="text-gray-800 tabular-nums">
              {formatUptime(uptime)}
            </span>
          </span>
          <span>
            TASKS:{" "}
            <span className="text-gray-800">
              {done}/{total}
            </span>
          </span>
          <span className="text-green-700">STATUS: ONLINE</span>
          <span className="ml-auto tabular-nums text-gray-700">
            {dateStr} {timeStr}
          </span>
        </div>
      </div>

      {/* Main */}
      <div className="flex gap-3 flex-1">
        {/* Task list */}
        <Window
          title="[ TASKS ]"
          right="CLICK TO SELECT"
          className="flex-1 flex flex-col"
        >
          <div className="flex-1 overflow-y-auto py-1">
            {loading ? (
              <p className="px-3 py-2 text-gray-400 text-sm animate-pulse">
                Loading core modules...
              </p>
            ) : todos.length === 0 ? (
              <p className="px-3 py-2 text-gray-400 text-sm">
                No tasks found. Create one.
              </p>
            ) : (
              <ul>
                {todos.map((todo) => {
                  const ds = todo.due_date
                    ? dueDateStatus(todo.due_date)
                    : null;
                  const isSelected = selectedId === todo.id;
                  return (
                    <li
                      key={todo.id}
                      onClick={() => setSelectedId(isSelected ? null : todo.id)}
                      className={`flex items-center gap-3 px-3 py-1.5 text-sm cursor-pointer transition-colors border-b border-gray-100 ${
                        isSelected
                          ? "bg-gray-900 text-gray-100"
                          : "hover:bg-gray-100"
                      }`}
                    >
                      <span
                        onClick={(e) => {
                          e.stopPropagation();
                          handleToggle(todo);
                        }}
                        className={`select-none ${
                          isSelected
                            ? "text-gray-300 hover:text-white"
                            : "text-gray-500 hover:text-gray-900"
                        }`}
                      >
                        [{todo.done ? "X" : "\u00a0"}]
                      </span>

                      <span
                        className={
                          todo.done ? "line-through text-gray-400" : ""
                        }
                      >
                        {todo.title}
                      </span>

                      {todo.due_date && ds && (
                        <span
                          className={`text-xs ${isSelected ? "text-gray-400" : DUE_DATE_CLS[ds]}`}
                        >
                          DL:{formatDueDate(todo.due_date)}
                          {ds === "overdue"
                            ? " (!)"
                            : ds === "soon"
                              ? " (~)"
                              : ""}
                        </span>
                      )}

                      <span
                        className={`ml-auto text-xs ${
                          isSelected
                            ? "text-gray-400"
                            : (PRIORITY_COLOR[todo.priority] ?? "text-gray-600")
                        }`}
                      >
                        {todo.priority.toUpperCase()}
                      </span>

                      {isSelected && (
                        <button
                          onClick={(e) => {
                            e.stopPropagation();
                            handleDelete(todo.id);
                          }}
                          className="text-red-400 hover:text-red-200 text-xs select-none"
                        >
                          [DEL]
                        </button>
                      )}
                    </li>
                  );
                })}
              </ul>
            )}
          </div>
        </Window>

        {/* Sidebar */}
        <div className="w-56 flex flex-col gap-3">
          {/* New task */}
          <Window title="[ NEW TASK ]">
            <div className="p-3 flex flex-col gap-2">
              <label className="text-xs text-gray-500">TITLE:</label>
              <input
                value={newTitle}
                onChange={(e) => setNewTitle(e.target.value)}
                onKeyDown={(e) => e.key === "Enter" && handleCreate()}
                placeholder="task description..."
                className="border border-gray-400 focus:border-gray-900 text-gray-900 placeholder:text-gray-300 px-2 py-1 text-xs w-full outline-none transition-colors bg-white"
              />

              <label className="text-xs text-gray-500 mt-1">PRIORITY:</label>
              <div className="flex gap-1">
                {(["low", "medium", "high"] as const).map((p) => (
                  <button
                    key={p}
                    onClick={() => setPriority(p)}
                    className={`flex-1 py-1 text-xs border transition-colors ${
                      priority === p
                        ? "border-gray-900 bg-gray-900 text-white"
                        : "border-gray-300 text-gray-500 hover:border-gray-600 hover:text-gray-700"
                    }`}
                  >
                    {p === "medium" ? "MED" : p.toUpperCase()}
                  </button>
                ))}
              </div>

              <label className="text-xs text-gray-500 mt-1">DUE DATE:</label>
              <input
                type="date"
                value={dueDate}
                onChange={(e) => setDueDate(e.target.value)}
                className="border border-gray-400 focus:border-gray-900 text-gray-900 px-2 py-1 text-xs w-full outline-none transition-colors bg-white"
              />

              <button
                onClick={handleCreate}
                disabled={!newTitle.trim()}
                className="mt-1 border border-gray-800 text-gray-800 hover:bg-gray-900 hover:text-white disabled:opacity-30 py-1.5 text-xs tracking-widest transition-colors"
              >
                &gt; ADD TASK
              </button>
            </div>
          </Window>

          {/* Stats */}
          <Window title="[ STATS ]">
            <div className="p-3 flex flex-col gap-1.5 text-xs">
              {[
                { label: "TOTAL", value: total, cls: "text-gray-900" },
                { label: "DONE", value: done, cls: "text-green-700" },
                {
                  label: "PENDING",
                  value: total - done,
                  cls: "text-amber-600",
                },
              ].map(({ label, value, cls }) => (
                <div key={label} className="flex justify-between">
                  <span className="text-gray-400">{label}</span>
                  <span className={cls}>{value}</span>
                </div>
              ))}
              <div className="border-t border-gray-200 mt-2 pt-2">
                <div className="text-gray-400 mb-1.5">PROGRESS {progress}%</div>
                <div className="border border-gray-400 h-2 bg-white">
                  <div
                    className="h-full bg-gray-800 transition-all duration-500"
                    style={{ width: `${progress}%` }}
                  />
                </div>
              </div>
            </div>
          </Window>
        </div>
      </div>

      {/* Status bar */}
      <div
        className="border border-gray-800 bg-white px-3 py-1.5 text-xs"
        style={HARD_SHADOW}
      >
        {error ? (
          <span className="text-red-600">[ERR] {error}</span>
        ) : (
          <span className="text-gray-400">
            Ready · [X] toggle · [DEL] remove · DL:(!) overdue · DL:(~) expires
            soon
          </span>
        )}
      </div>
    </div>
  );
}
