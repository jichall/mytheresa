package category

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/models"
)

type CategoryHandler struct {
	repository models.Repository[models.Category]
	logger     *slog.Logger
}

type CategoryHandlerOpts struct {
	Repository models.Repository[models.Category]
	Logger     *slog.Logger
}

func NewCatalogHandler(opts *CategoryHandlerOpts) *CategoryHandler {
	return &CategoryHandler{
		repository: opts.Repository,
		logger:     opts.Logger,
	}
}

// HandleGet returns all categories in the database
func (c *CategoryHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	res, err := c.repository.GetAll(ctx)
	if err != nil {
		c.logger.Error("failed to get all categories", slog.Any("error", err))
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())

		return
	}

	categories := []Category{}

	for _, category := range res {
		categories = append(categories, Category{
			Code: category.Code,
			Name: category.Name,
		})
	}

	api.OKResponse(w, categories)
}

// HandleCreate creates a new category and persists it into the database
func (c *CategoryHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	data, err := io.ReadAll(r.Body)
	if err != nil {
		c.logger.Error("failed to read request body", slog.Any("error", err))
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())

		return
	}
	r.Body.Close()

	category := Category{}

	err = json.Unmarshal(data, &category)
	if err != nil {
		c.logger.Error("failed to parse request object", slog.Any("error", err), slog.String("data", string(data)))
		api.ErrorResponse(w, http.StatusBadRequest, "invalid input model")

		return
	}

	err = c.repository.Save(ctx, &models.Category{Code: category.Code, Name: category.Name})
	if err != nil {
		c.logger.Error("failed to persist data", slog.Any("error", err))
		api.ErrorResponse(w, http.StatusServiceUnavailable, "failed to persist data, try again")

		return
	}

	api.Response(w, http.StatusCreated, nil)
}
