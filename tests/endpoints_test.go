package endpoints_test

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-kv/internal/config"
	"go-kv/internal/handlers"
	"go-kv/internal/services"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

// MockKVService implements the KVServiceInterface for testing
type MockKVService struct {
	data map[string]interface{}
	mu   sync.Mutex
}

func NewMockKVService() *MockKVService {
	return &MockKVService{
		data: make(map[string]interface{}),
	}
}

func (m *MockKVService) Get(key string) (interface{}, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	value, exists := m.data[key]
	if !exists {
		return nil, nil // simulate key not found
	}
	return value, nil
}

func (m *MockKVService) Put(key string, value interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
	return nil
}

func (m *MockKVService) Delete(key string) error {
    //Lock the field
	m.mu.Lock()
	defer m.mu.Unlock()
	_, exists := m.data[key]
	if !exists {
		return fmt.Errorf("key not found")
	}
	delete(m.data, key)
	return nil
}

func (m *MockKVService) ListKeys() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	keys := make([]string, 0, len(m.data))
	for key := range m.data {
		keys = append(keys, key)
	}
	return keys
}

func setupRouter(kvService services.KVServiceInterface) *gin.Engine {
	cfg := config.LoadConfig()
	router := gin.Default()
	router.Use(gin.Recovery())
	api := router.Group(cfg.APIPath) //Obviously dont need router groups but just for future dev purposes
	{
		api.GET("/:key", handlers.HandleGet(kvService))
		api.PUT("/:key", handlers.HandlePut(kvService))
		api.DELETE("/:key", handlers.HandleDelete(kvService))
		api.GET("/", handlers.HandleListKeys(kvService))
	}
	return router
}

func TestHandleGet(t *testing.T) {
	kvService := NewMockKVService()
	router := setupRouter(kvService)

	// Test key not found
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/nonexistent", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}

	// Test successful retrieval
	kvService.Put("key1", "value1")
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/key1", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestHandlePut(t *testing.T) {
	kvService := NewMockKVService()
	router := setupRouter(kvService)

	// Test successful put
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/key1", bytes.NewBuffer([]byte("testValue")))
	router.ServeHTTP(w, req)

	if w.Code != http.StatusAccepted {
		t.Errorf("Expected status 202, got %d", w.Code)
	}

	// Verify value is set
	kvService.Put("key1", "value1")
	value, _ := kvService.Get("key1")
	if value != "value1" {
		t.Errorf("Expected value 'value1', got '%v'", value)
	}
}

func TestHandleDelete(t *testing.T) {
	kvService := NewMockKVService()
	router := setupRouter(kvService)

	// Test delete key
	kvService.Put("key1", "value1")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/key1", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Test delete nonexistent key
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/nonexistent", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestHandleListKeys(t *testing.T) {
	kvService := NewMockKVService()
	router := setupRouter(kvService)

	// Test listing keys
	kvService.Put("key1", "value1")
	kvService.Put("key2", "value2")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestParallelRequests(t *testing.T) {
	kvService := NewMockKVService()
	kvService.Put("key1", "value1")

	t.Run("parallel puts", func(t *testing.T) {
		t.Parallel()

		wg := sync.WaitGroup{}
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(n int) {
				defer wg.Done()
				kvService.Put("key1", fmt.Sprintf("value%d", n)) // Overwrite the same key
			}(i)
		}
		wg.Wait()

		// Ensure 'key1' can be retrieved
		value, err := kvService.Get("key1")
		if err != nil {
			t.Fatalf("Failed to get value for 'key1': %v", err)
		}

		// Ensure that the key is still accessible
		if value == "" {
			t.Errorf("Expected value for 'key1' to not be empty, got '%s'", value)
		}

		// Verify that the key is in the list
		keys := kvService.ListKeys()
		if len(keys) == 0 || keys[0] != "key1" {
			t.Errorf("Expected 'key1' to be present in the store, but it was not found.")
		}
	})

	t.Run("parallel gets", func(t *testing.T) {
		t.Parallel()

		wg := sync.WaitGroup{}
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, err := kvService.Get("key1") // Check if 'key1' is accessible
				if err != nil {
					t.Errorf("Failed to get value for 'key1': %v", err)
				}
			}()
		}
		wg.Wait()
	})
}
