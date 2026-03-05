package category

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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
		Logger:     nil,
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
}

type mockrepo struct{}

func (m *mockrepo) GetAll(ctx context.Context) ([]models.Category, error) {
	return []models.Category{{ID: 0, Code: "acsssshually", Name: "acsssshually"}}, nil
}

func (m *mockrepo) GetPaged(ctx context.Context, page, limit int) ([]models.Category, error) {
	return []models.Category{}, nil
}

func (m *mockrepo) GetByCode(ctx context.Context, code string) (*models.Category, error) {
	return nil, nil
}

func (m *mockrepo) Save(ctx context.Context, data *models.Category) error {
	return nil
}
