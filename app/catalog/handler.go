package catalog

import (
	"context"
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

// HandleGetByCode returns a specific product or none
func (h *CatalogHandler) HandleGetByCode(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	res, err := h.repository.GetByCode(ctx, code)
	if err != nil {
		h.logger.Error("failed to get product by code", slog.String("code", code))
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())

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

	api.OKResponse(w, product)
}

// HandleGet returns products to the client based off of a filter to page results nicely
func (h *CatalogHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	filter := &ProductFilter{}

	if err := filter.parse(r); err != nil {
		slog.Error("failed to parse query parameters", slog.Any("error", err))
		api.ErrorResponse(w, http.StatusBadRequest, "invalid query parameters")

		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	res, err := h.repository.GetPaged(ctx, filter.Page, filter.Limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	api.OKResponse(w, response)
}
