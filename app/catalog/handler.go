package catalog

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/shopspring/decimal"
)

type CatalogHandler struct {
	repository models.Repository[models.Product]
	logger     *slog.Logger
}

type CatalogHandlerOpts struct {
	Repository models.Repository[models.Product]
	Logger     *slog.Logger
}

func NewCatalogHandler(opts *CatalogHandlerOpts) *CatalogHandler {
	return &CatalogHandler{
		repository: opts.Repository,
		logger:     opts.Logger,
	}
}

// GetByCode godoc
// @Summary Returns a specific product from the catalog or none if not found
// @Tags catalog
// @Produce json
// @Param code path string true "product code"
// @Success 200 {object} Product
// @Failure 404 {object} api.Response
// @Router /catalog/{code} [get]
func (h *CatalogHandler) HandleGetByCode(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	res, err := h.repository.GetByCode(ctx, code)
	if err != nil {
		h.logger.Error("failed to get product by code", slog.String("code", code))
		api.RespondError(w, api.Response{Status: http.StatusInternalServerError, Error: err.Error()})

		return
	}

	if res == nil {
		api.RespondError(w, api.Response{Status: http.StatusNotFound, Message: fmt.Sprintf("product with code [%s] not found", code)})

		return
	}

	variants := []Variant{}
	for _, variant := range res.Variants {
		price := variant.Price.InexactFloat64()

		if variant.Price.Equal(decimal.NewFromInt(0)) {
			price = res.Price.InexactFloat64()
		}

		variants = append(variants, Variant{
			Name:  variant.Name,
			Price: price,
		})
	}

	product := Product{
		Code:     res.Code,
		Price:    res.Price.InexactFloat64(),
		Category: res.Category.Code,
		Variants: variants,
	}

	api.RespondOK(w, product)
}

// HandleGet godoc
// @Summary Returns products in a paged format
// @Tags catalog
// @Param page query int false "page number (index at 0)"
// @Param limit query int false "number of items per page"
// @Produce json
// @Success 200 {object} []Product
// @Failure 404 {object} api.Response
// @Router /catalog [get]
func (h *CatalogHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	filter := &ProductFilter{}

	if err := filter.parse(r); err != nil {
		slog.Error("failed to parse query parameters", slog.Any("error", err))
		api.RespondError(w, api.Response{Status: http.StatusBadRequest, Message: "invalid query parameter(s)"})

		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	res, err := h.repository.GetPaged(ctx, filter.Page, filter.Limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(res) == 0 {
		api.RespondError(w, api.Response{Status: http.StatusNotFound, Message: fmt.Sprintf("products were not found")})

		return
	}

	products := make([]Product, len(res))

	for i, p := range res {
		products[i] = Product{
			Code:     p.Code,
			Price:    p.Price.InexactFloat64(),
			Category: p.Category.Code,
			Variants: []Variant{},
		}
	}

	response := Response{
		Products: products,
		Filter:   &ResponseFilter{Offset: filter.Page, Limit: filter.Limit},
	}

	api.RespondOK(w, response)
}
