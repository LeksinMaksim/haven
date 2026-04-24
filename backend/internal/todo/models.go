package todo

import "time"

type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHight  Priority = "high"
)

type Todo struct {
	ID          int64      `json:"id" db:"id"`
	Title       string     `json:"title" db:"title"`
	Description string     `json:"description,omitempty" db:"description"`
	Done        bool       `json:"done" db:"done"`
	Priority    Priority   `json:"priority" db:"priority"`
	DueDate     *time.Time `json:"due_date,omitempty" db:"due_date"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

type CreateRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Priority    Priority   `json:"priority"`
	DueDate     *time.Time `json:"due_date"`
}

func (r *CreateRequest) Validate() error {
	if r.Title == "" {
		return ErrEmptyTitle
	}
	if r.Priority == "" {
		r.Priority = PriorityMedium
	}

	return nil
}

type UpdateRequest struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	Done        *bool      `json:"done"`
	Priority    *Priority  `json:"priority"`
	DueDate     *time.Time `json:"due_date"`
}
