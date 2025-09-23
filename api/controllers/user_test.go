package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/brando1998/docubot-api/mocks"
	"github.com/brando1998/docubot-api/models"
)

func setupMockRouter() *gin.Engine {
	// Creamos variables para los emails con sus punteros
	johnEmail := "john@example.com"
	janeEmail := "jane@example.com"

	mock := &mocks.MockClientRepo{
		CreateClientFunc: func(user *models.Client) error {
			user.ID = 1
			return nil
		},
		GetClientByIDFunc: func(id uint) (*models.Client, error) {
			return &models.Client{
				ID:    1,
				Name:  "John Doe",
				Email: &johnEmail, // Ahora es un puntero
			}, nil
		},
		GetClientByPhoneFunc: func(phone string) (*models.Client, error) {
			return &models.Client{
				ID:    2,
				Name:  "Jane Doe",
				Email: &janeEmail, // Ahora es un puntero
			}, nil
		},
		GetOrCreateClientFunc: func(phone, name, email string) (*models.Client, error) {
			var emailPtr *string
			if email != "" {
				emailPtr = &email
			}
			return &models.Client{
				ID:    3,
				Name:  name,
				Email: emailPtr, // Puede ser nil si email está vacío
			}, nil
		},
	}
	SetClientRepo(mock)

	r := gin.Default()
	r.POST("/users", CreateClient)
	r.GET("/users/profile", GetClientByID)
	r.POST("/users/get-or-create", GetOrCreateClient)
	r.GET("/users/:phone", GetClientByPhone)
	return r
}

// El resto de las funciones de test pueden permanecer igual ya que
// trabajan con el JSON y el marshaling/unmarshaling manejará la conversión
// entre string y *string automáticamente

func TestCreateClient(t *testing.T) {
	r := setupMockRouter()
	user := map[string]string{
		"name":  "Brando",
		"email": "brando@example.com",
	}
	body, _ := json.Marshal(user)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Usuario creado")
}

func TestGetProfile(t *testing.T) {
	r := setupMockRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/profile", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "John Doe")
}

func TestGetOrCreateClient(t *testing.T) {
	r := setupMockRouter()
	user := map[string]string{
		"phone": "1234567890",
		"name":  "Alice",
		"email": "alice@example.com",
	}
	body, _ := json.Marshal(user)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users/get-or-create", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Alice")
}

func TestGetClientByPhone(t *testing.T) {
	r := setupMockRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/1234567890", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Jane Doe")
}
