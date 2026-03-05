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

func NewCategoryHandler(opts *CategoryHandlerOpts) *CategoryHandler {
	return &CategoryHandler{
		repository: opts.Repository,
		logger:     opts.Logger,
	}
}

// GetCategory godoc
// @Summary Get all available categories
// @Tags categories
// @Produce json
// @Success 200 {object} Category
// @Failure 404 {object} api.Response
// @Router /categories [get]
func (c *CategoryHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	res, err := c.repository.GetAll(ctx)
	if err != nil {
		c.logger.Error("failed to get all categories", slog.Any("error", err))
		api.RespondError(w, api.Response{Status: http.StatusInternalServerError, Error: err.Error()})

		return
	}

	if len(res) == 0 {
		api.RespondError(w, api.Response{Status: http.StatusNotFound, Message: "no available category"})

		return
	}

	categories := []Category{}

	for _, category := range res {
		categories = append(categories, Category{
			Code: category.Code,
			Name: category.Name,
		})
	}

	api.RespondOK(w, categories)
}

// CreateCategory godoc
// @Summary Create a new category
// @Tags categories
// @Accept json
// @Success 201
// @Failure 400 {object} api.Response
// @Failure 503 {object} api.Response
// @Router /categories [post]
func (c *CategoryHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	data, err := io.ReadAll(r.Body)
	if err != nil {
		c.logger.Error("failed to read request body", slog.Any("error", err))
		api.RespondError(w, api.Response{Status: http.StatusInternalServerError, Error: err.Error()})

		return
	}
	r.Body.Close()

	category := Category{}

	// the validation should've been happening at the decoder level, or using struct tags to clear the
	// possible clutter of validation logic
	err = json.Unmarshal(data, &category)
	if err != nil || (category.Code == "" && category.Name == "") {
		c.logger.Error("failed to parse request object", slog.Any("error", err), slog.String("data", string(data)))
		api.RespondError(w, api.Response{Status: http.StatusBadRequest, Message: "invalid input model"})

		return
	}

	err = c.repository.Save(ctx, &models.Category{Code: category.Code, Name: category.Name})
	if err != nil {
		c.logger.Error("failed to persist data", slog.Any("error", err))
		api.RespondError(w, api.Response{Status: http.StatusServiceUnavailable, Message: "failed to persist data, try again"})

		return
	}

	api.RespondCustom(w, http.StatusCreated, nil)
}
