package schema

import (
	"fmt"

	"github.com/leonardo-gmuller/lexflow-ai/internal/app/domain/dto"
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/domain/entity"
)

// -------------------------------CREATE-------------------------------
type CreateCaseRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

func (r *CreateCaseRequest) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("name is required")
	}
	return nil
}

type CreateCaseResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// -------------------------------LIST-------------------------------
type ListCasesRequest struct {
	Search     string                     `json:"search"` // termo livre
	Order      string                     `json:"order"`  // "id" | "email" | "created_at"
	Desc       bool                       `json:"desc"`   // true => DESC
	Pagination ListCasesRequestPagination `json:"pagination"`
}

type ListCasesRequestPagination = dto.Pagination

// Output da listagem paginada
type ListCasesOutput = PaginatedResponse[ListCasesResponse]

type ListCasesResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
}

func ConvertListCaseToResponse(cases []entity.Case) []ListCasesResponse {
	items := make([]ListCasesResponse, 0, len(cases))
	for _, c := range cases {
		items = append(items, ListCasesResponse{
			ID:          c.ID.String(),
			Name:        c.Name,
			Description: c.Description,
			CreatedAt:   c.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}
	return items
}

var ValidSortCasesFields = map[string]bool{
	"created_at": true,
	"name":       true,
}
