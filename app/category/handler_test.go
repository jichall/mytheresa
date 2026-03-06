package category

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/app/database"
	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/stretchr/testify/assert"
)

func TestHandleGet(t *testing.T) {

	handler := NewCategoryHandler(&CategoryHandlerOpts{
		Repository: &mockrepo{},
		Logger:     nil,
	})

	t.Run("test get successfully", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		expected := []models.Category([]models.Category{models.Category{ID: 0, Code: "acsssshually", Name: "acsssshually"}})

		decoder := json.NewDecoder(recorder.Body)

		handler.HandleGet(recorder, &http.Request{})

		models := []models.Category{}
		decoder.Decode(&models)

		assert.Equal(t, expected, models)
	})
}

func TestHandlePost(t *testing.T) {
	handler := NewCategoryHandler(&CategoryHandlerOpts{
		Repository: &mockrepo{},
		Logger:     slog.New(slog.NewTextHandler(io.Discard, nil)),
	})

	t.Run("test post successfully", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		expected := ""

		payload := Category{
			Code: "acshualy",
			Name: "Acshualy",
		}

		data, err := json.Marshal(payload)
		assert.Nil(t, err)

		request := httptest.NewRequest(http.MethodPost, "/category", bytes.NewReader(data))

		handler.HandleCreate(recorder, request)

		assert.Equal(t, expected, recorder.Body.String())
	})

	t.Run("test post with a bad request", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		expected := "{\"status\":400,\"message\":\"invalid input model\"}\n"

		data := []byte(`{"bad formatted js":on}`)

		request := httptest.NewRequest(http.MethodPost, "/category", bytes.NewReader(data))

		handler.HandleCreate(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Result().StatusCode)
		assert.Equal(t, expected, recorder.Body.String())
	})

	t.Run("test post with a valid json but invalid model", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		expected := "{\"status\":400,\"message\":\"invalid input model\"}\n"

		payload := struct{ Bogus []string }{Bogus: []string{"this is a BAD BAD thing", "so BAD"}}
		data, err := json.Marshal(payload)
		assert.Nil(t, err)

		request := httptest.NewRequest(http.MethodPost, "/category", bytes.NewReader(data))

		handler.HandleCreate(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Result().StatusCode)
		assert.Equal(t, expected, recorder.Body.String())
	})
}

// I wish I could use mockery, that would make this test much cleaner
type mockrepo struct{}

func (m *mockrepo) GetAll(ctx context.Context) ([]models.Category, error) {
	return []models.Category{{ID: 0, Code: "acsssshually", Name: "acsssshually"}}, nil
}

func (m *mockrepo) GetPaged(ctx context.Context, page, limit int) ([]models.Category, error) {
	return []models.Category{}, nil
}

func (m *mockrepo) GetWithFilters(ctx context.Context, filters []database.Filter) ([]models.Category, error) {
	return []models.Category{}, nil
}

func (m *mockrepo) GetByCode(ctx context.Context, code string) (*models.Category, error) {
	return nil, nil
}

func (m *mockrepo) Save(ctx context.Context, data *models.Category) error {
	return nil
}
