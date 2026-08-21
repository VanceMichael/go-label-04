package domain

type PageRequest struct {
	Limit  int
	Cursor string
	Sort   string
}
type Page[T any] struct {
	Items      []T
	NextCursor string
	Total      int
}

func (p PageRequest) Normalized() PageRequest {
	if p.Limit < 1 {
		p.Limit = 50
	}
	if p.Limit > 200 {
		p.Limit = 200
	}
	if p.Sort == "" {
		p.Sort = "created_at"
	}
	return p
}
