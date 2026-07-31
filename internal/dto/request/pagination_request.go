package request

type PaginationRequest struct {
	Page    int `form:"page,default=1"`
	PerPage int `form:"per_page,default=20"`
}

func (p *PaginationRequest) Offset() int {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PerPage < 1 || p.PerPage > 100 {
		p.PerPage = 20
	}
	return (p.Page - 1) * p.PerPage
}
