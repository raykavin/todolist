package valueobject

// TodoStatistics representa estatísticas dos todos,
// encapsulando métricas importantes para o domínio.
type TodoStatistics struct {
	Total          int64
	ByStatus       map[string]int64
	ByPriority     map[string]int64
	Overdue        int64
	DueToday       int64
	DueThisWeek    int64
	CompletedToday int64
	CompletionRate float64
}

// NewTodoStatistics cria uma nova instância de estatísticas
func NewTodoStatistics(
	total int64,
	byStatus map[string]int64,
	byPriority map[string]int64,
	overdue int64,
	dueToday int64,
	dueThisWeek int64,
	completedToday int64,
	completionRate float64,
) TodoStatistics {
	return TodoStatistics{
		Total:          total,
		ByStatus:       byStatus,
		ByPriority:     byPriority,
		Overdue:        overdue,
		DueToday:       dueToday,
		DueThisWeek:    dueThisWeek,
		CompletedToday: completedToday,
		CompletionRate: completionRate,
	}
}
