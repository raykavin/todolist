package errors

import "fmt"

type AppError struct {
	Code    string
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Erros comuns do sistema de tarefas

var (
	ErrTodoNotFound  = &AppError{Code: "TODO_NOT_FOUND", Message: "Todo not found"}
	ErrUnauthorized  = &AppError{Code: "UNAUTHORIZED", Message: "Unauthorized access"}
	ErrInvalidInput  = &AppError{Code: "INVALID_INPUT", Message: "Invalid input data"}
	ErrQuotaExceeded = &AppError{Code: "QUOTA_EXCEEDED", Message: "Todo quota exceeded"}
)
