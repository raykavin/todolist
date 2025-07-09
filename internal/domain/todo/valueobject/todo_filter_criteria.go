package valueobject

import "time"

// TodoFilterCriteria representa critérios de filtragem para queries de Todo.
// Faz parte do domínio, pois define como a entidade pode ser consultada.
type TodoFilterCriteria struct {
	UserID      int64
	Status      []string
	Priority    []string
	Tags        []string
	IsOverdue   *bool
	DueDateFrom *time.Time
	DueDateTo   *time.Time
	SearchTerm  string
}

// HasStatusFilter indica se há filtro por status
func (f *TodoFilterCriteria) HasStatusFilter() bool {
	return len(f.Status) > 0
}

// HasPriorityFilter indica se há filtro por prioridade
func (f *TodoFilterCriteria) HasPriorityFilter() bool {
	return len(f.Priority) > 0
}

// HasTagsFilter indica se há filtro por tags
func (f *TodoFilterCriteria) HasTagsFilter() bool {
	return len(f.Tags) > 0
}

// HasDueDateRange indica se há filtro por intervalo de data de vencimento
func (f *TodoFilterCriteria) HasDueDateRange() bool {
	return f.DueDateFrom != nil || f.DueDateTo != nil
}

// HasSearchTerm indica se há termo de pesquisa
func (f *TodoFilterCriteria) HasSearchTerm() bool {
	return f.SearchTerm != ""
}
