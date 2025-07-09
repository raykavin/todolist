package dto

// TodoStatistic represents todo statistics
type TodoStatistics struct {
	Total          int64            `json:"total"`
	ByStatus       map[string]int64 `json:"by_status"`
	ByPriority     map[string]int64 `json:"by_priority"`
	Overdue        int64            `json:"overdue"`
	DueToday       int64            `json:"due_today"`
	DueThisWeek    int64            `json:"due_this_week"`
	CompletedToday int64            `json:"completed_today"`
	CompletionRate float64          `json:"completion_rate"`
}
