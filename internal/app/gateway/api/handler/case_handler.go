package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/domain/entity"
	case_usecase "github.com/leonardo-gmuller/lexflow-ai/internal/app/domain/usecase/case"
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/gateway/api/handler/schema"
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/gateway/api/rest"
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/gateway/api/rest/response"
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/gateway/postgres"
)

const (
	casePattern = "/cases"
)

func (h *Handler) caseSetupRoutes(router chi.Router) {
	router.Route(casePattern, func(r chi.Router) {
		r.Post("/", h.createCase())
		r.Get("/", h.listAllCases())
	})
}

func (h *Handler) createCase() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		var resp *response.Response

		var in schema.CreateCaseRequest

		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			resp = response.BadRequest(err, fmt.Errorf("invalid payload: %w", err).Error())
			rest.SendJSON(w, resp.Status, resp.Payload, resp.Headers) //nolint:errcheck

			return
		}

		if err := in.Validate(); err != nil {
			resp = response.BadRequest(err, fmt.Sprintf("validation error: %v", err))
			rest.SendJSON(w, resp.Status, resp.Payload, resp.Headers) //nolint:errcheck

			return
		}

		var newCase *entity.Case

		err := h.app.DB.WithTx(r.Context(), func(ctx context.Context, db postgres.DBTX) error {
			uc := h.app.NewCaseUseCase(db)

			result, err := uc.CreateCase(ctx, in.Name, in.Description)
			if err != nil {
				return err
			}

			newCase = result

			return nil
		})

		if err != nil {
			resp = response.InternalServerError(err)
			rest.SendJSON(w, resp.Status, resp.Payload, resp.Headers) //nolint:errcheck

			return
		}

		resp = response.Created(schema.CreateCaseResponse{
			ID:          newCase.ID.String(),
			Name:        newCase.Name,
			Description: newCase.Description,
		})
		rest.SendJSON(w, resp.Status, resp.Payload, resp.Headers)
	}
}

func (h *Handler) listAllCases() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queryStrings := r.URL.Query()

		input := schema.ListCasesRequest{}

		h.getPaginationParams(queryStrings, &input.Pagination)
		input.Search = h.readString(queryStrings, "search", "")

		if err := input.Pagination.Validate(schema.ValidSortCasesFields); err != nil {
			resp := response.BadRequest(err, err.Error())
			rest.SendJSON(w, resp.Status, resp.Payload, resp.Headers) //nolint:errcheck
			return
		}

		caseUseCase := h.app.NewCaseUseCase(h.app.DB.Pool)

		result, err := caseUseCase.ListAllCases(r.Context(), case_usecase.ListAllCasesInput{
			Pagination: input.Pagination,
			Search:     input.Search,
		})

		if err != nil {
			resp := response.InternalServerError(err)
			rest.SendJSON(w, resp.Status, resp.Payload, resp.Headers) //nolint:errcheck
			return
		}

		metadata := schema.Meta{
			CurrentPage:  input.Pagination.Page,
			ItemsPerPage: input.Pagination.ItemsPerPage,
			TotalItems:   int(result.TotalItems),
		}

		output := schema.ListCasesOutput{
			Items:    schema.ConvertListCaseToResponse(result.Cases),
			Metadata: metadata,
		}

		resp := response.OK(output)
		rest.SendJSON(w, resp.Status, resp.Payload, resp.Headers) //nolint:errcheck

	}
}
