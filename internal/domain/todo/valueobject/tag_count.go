package valueobject

// TagCount representa a contagem de uso de uma tag.
// Faz parte do domínio pois expressa uma métrica relevante no negócio.
type TagCount struct {
	Tag   string
	Count int64
}

// NewTagCount cria um novo TagCount
func NewTagCount(tag string, count int64) TagCount {
	return TagCount{
		Tag:   tag,
		Count: count,
	}
}
