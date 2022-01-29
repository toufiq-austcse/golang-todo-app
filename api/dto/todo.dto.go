package dto

type CreateTodoDto struct {
	Task        string `json:"task" binding:"required"`
	IsCompleted bool   `json:"is_completed" binding:"required"`
}

type UpdateTodoDto struct {
	Task        string `json:"task,omitempty" `
	IsCompleted bool   `json:"is_completed,omitempty"`
}
