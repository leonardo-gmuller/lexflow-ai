package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/domain/entity"
	case_usecase "github.com/leonardo-gmuller/lexflow-ai/internal/app/domain/usecase/case"
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/domain/usecase/document"
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
		r.Post("/{id}/documents", h.uploadDocumentToCase())
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

func (h *Handler) uploadDocumentToCase() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		caseID := chi.URLParam(r, "id")

		defer r.Body.Close()

		var resp *response.Response

		if err := r.ParseMultipartForm(20 << 20); err != nil {
			resp = response.BadRequest(err, "invalid multipart form")
			rest.SendJSON(w, resp.Status, resp.Payload, resp.Headers)
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			resp = response.BadRequest(nil, "file is required")
			rest.SendJSON(w, resp.Status, resp.Payload, resp.Headers)
			return
		}
		defer file.Close()

		content, readErr := io.ReadAll(file)
		if readErr != nil {
			resp = response.InternalServerError(fmt.Errorf("failed to read file: %w", readErr))
			rest.SendJSON(w, resp.Status, resp.Payload, resp.Headers)
			return
		}

		mimeType := header.Header.Get("Content-Type")
		if mimeType == "" || mimeType != "application/pdf" {
			resp = response.BadRequest(nil, "invalid file type: only PDF files are allowed")
			rest.SendJSON(w, resp.Status, resp.Payload, resp.Headers)
			return
		}

		var doc *entity.Document

		err = h.app.DB.WithTx(r.Context(), func(ctx context.Context, db postgres.DBTX) error {
			uc := h.app.NewDocumentUsecase(db)

			result, err := uc.UploadDocument(ctx, document.UploadDocumentInput{
				CaseID:   caseID,
				FileName: header.Filename,
				Content:  content,
			})
			if err != nil {
				return err
			}

			doc = result
			return nil
		})

		if err != nil {
			resp = response.InternalServerError(err)
			rest.SendJSON(w, resp.Status, resp.Payload, resp.Headers) //nolint:errcheck
			return
		}

		resp = response.Created(schema.UploadDocumentToCaseResponse{
			ID:   doc.ID.String(),
			Name: doc.FileName,
		})

		rest.SendJSON(w, resp.Status, resp.Payload, resp.Headers) //nolint:errcheck
	}
}
