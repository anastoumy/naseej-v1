package utils

// PaginationInput holds pagination parameters.
type PaginationInput struct {
	Limit  int
	Offset int
}

// DefaultPagination returns default pagination settings.
func DefaultPagination() PaginationInput {
	return PaginationInput{Limit: 20, Offset: 0}
}

// Validate ensures pagination values are within acceptable ranges.
func (p *PaginationInput) Validate() {
	if p.Limit <= 0 {
		p.Limit = 20
	}
	if p.Limit > 100 {
		p.Limit = 100
	}
	if p.Offset < 0 {
		p.Offset = 0
	}
}
