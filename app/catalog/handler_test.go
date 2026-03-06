package catalog

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/app/database"
	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestGetAll(t *testing.T) {
	handler := NewCatalogHandler(&CatalogHandlerOpts{Repository: &mockrepo{}, Logger: slog.New(slog.NewTextHandler(io.Discard, nil))})

	t.Run("test get products", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/catalog", nil)
		response := []Product{}

		handler.HandleGet(recorder, request)

		decoder := json.NewDecoder(recorder.Body)

		err := decoder.Decode(&response)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, 10, len(response))
		assert.Equal(t, "PRD009", response[10].Code)
	})

	t.Run("test get products with pagination", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/catalog?offset=0&limit=3", nil)
		response := []Product{}

		handler.HandleGet(recorder, request)

		decoder := json.NewDecoder(recorder.Body)

		err := decoder.Decode(&response)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, 3, len(response))
		assert.Equal(t, "PRD000", response[0].Code)
	})

	t.Run("test get products with wrong offset parameter", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/catalog?offset=-10&limit=1", nil)

		handler.HandleGet(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Result().StatusCode)
	})

	t.Run("test get products with wrong limit parameter", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/catalog?offset=0&limit=-1", nil)

		handler.HandleGet(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Result().StatusCode)
	})
}

// I wish I could use mockery, that would make this test much cleaner
type mockrepo struct{}

func (m *mockrepo) GetAll(ctx context.Context) ([]models.Product, error) {
	data := make([]models.Product, 30)

	for i := 0; i < 30; i++ {
		data[i] = models.Product{ID: uint(i), Code: fmt.Sprintf("PRD%03d", i), Price: decimal.NewFromInt(int64(i)), Variants: []models.Variant{}}
	}

	return data, nil
}

func (m *mockrepo) GetPaged(ctx context.Context, page, size int) ([]models.Product, error) {
	data := make([]models.Product, 30)

	for i := 0; i < 30; i++ {
		data[i] = models.Product{ID: uint(i), Code: fmt.Sprintf("PRD%03d", i), Price: decimal.NewFromInt(int64(i)), Variants: []models.Variant{}}
	}

	offset := page * size
	if offset >= len(data) {
		return data[len(data)-1-size : len(data)-1], nil
	}

	return data[offset : offset+size], nil
}

func (m *mockrepo) GetWithFilters(ctx context.Context, filters []database.Filter) ([]models.Product, error) {
	data := make([]models.Product, 30)

	for i := 0; i < 30; i++ {
		data[i] = models.Product{ID: uint(i), Code: fmt.Sprintf("PRD%03d", i), Price: decimal.NewFromInt(int64(i)), Variants: []models.Variant{}}
	}

	offset := page * size
	if offset >= len(data) {
		return data[len(data)-1-size : len(data)-1], nil
	}

	return data[offset : offset+size], nil
}

func (m *mockrepo) GetByCode(ctx context.Context, code string) (*models.Product, error) {
	return nil, nil
}

func (m *mockrepo) Save(ctx context.Context, data *models.Product) error {
	return nil
}
