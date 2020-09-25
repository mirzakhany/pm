package domains

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"projectmanager/pkg/utiles"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func apiTestSetup() (*gin.Engine, *Repository, func()) {
	db, err := gorm.Open(sqlite.Open("/tmp/api_domains.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	repo := NewRepository(db)

	NewService(r, repo, logger)

	err = db.AutoMigrate(&Domain{})
	if err != nil {
		panic(err)
	}
	return r, repo, func() {
		logger.Sync()
		_ = os.Remove("/tmp/api_domains.db")
	}
}

func marshal(data DomainRequest) *bytes.Buffer {
	reqBody, err := json.Marshal(data)

	if err != nil {
		print(err)
	}
	return bytes.NewBuffer(reqBody)
}

func unmarshalSingle(res string) (*Domain, error) {

	type data struct {
		Data Domain `json:"data"`
	}

	var jsonData data
	err := json.Unmarshal([]byte(res), &jsonData)
	if err != nil {
		return nil, err
	}
	return &jsonData.Data, nil
}

func unmarshalList(res string) ([]Domain, error) {
	type data struct {
		Data []Domain `json:"data"`
	}

	var jsonData data
	err := json.Unmarshal([]byte(res), &jsonData)
	if err != nil {
		return nil, err
	}
	return jsonData.Data, nil
}

func TestService_CreateHandler(t *testing.T) {
	router, _, cleaner := apiTestSetup()
	defer cleaner()
	tests := []struct {
		name           string
		request        DomainRequest
		checkFunc      func(res string) error
		wantStatusCode int
	}{
		{
			name: "create domain object with correct data",
			request: DomainRequest{
				Title:   "github.com",
				Address: "https://github.com",
			},
			checkFunc: func(res string) error {
				domain, err := unmarshalSingle(res)
				if err != nil {
					return err
				}

				if !utiles.IsValidUUID(domain.UUID) || domain.Title != "github.com" || domain.Address != "https://github.com" {
					return fmt.Errorf("invalid domain data")
				}
				return nil
			},
			wantStatusCode: 200,
		},
		{
			name: "create domain object with duplicate data",
			request: DomainRequest{
				Title:   "pages.github.com",
				Address: "https://github.com",
			},
			checkFunc:      nil,
			wantStatusCode: 400,
		},
		{
			name: "create domain object with incorrect data",
			request: DomainRequest{
				Address: "https://google.com",
			},
			checkFunc:      nil,
			wantStatusCode: 400,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, err := http.NewRequest("POST", "/domains", marshal(tt.request))
			if err != nil {
				t.Errorf("domain CreateHandler, request failed, error = %v", err)
				return
			}
			router.ServeHTTP(w, req)
			if w.Code != tt.wantStatusCode {
				t.Errorf("domain CreateHandler failed, status code = %d", w.Code)
				return
			}

			if tt.checkFunc == nil {
				return
			}

			if err := tt.checkFunc(w.Body.String()); err != nil {
				t.Errorf("domain CreateHandler returned invalid data = %v", w.Body.String())
				return
			}
		})
	}
}

func TestService_UpdateHandler(t *testing.T) {
	router, repo, cleaner := apiTestSetup()
	defer cleaner()

	domain1, err := repo.Create(DomainRequest{
		Title:   "google.com",
		Address: "https://google.com",
	})
	if err != nil {
		panic(err)
	}

	_, err = repo.Create(DomainRequest{
		Title:   "aws.com",
		Address: "https://aws.com",
	})
	if err != nil {
		panic(err)
	}

	tests := []struct {
		name           string
		uuid           string
		request        DomainRequest
		checkFunc      func(res string) error
		wantStatusCode int
	}{
		{
			name: "update domain object with correct data",
			uuid: domain1.UUID,
			request: DomainRequest{
				Title:   "page.github.com",
				Address: "https://github.com",
			},
			checkFunc: func(res string) error {
				domain, err := unmarshalSingle(res)
				if err != nil {
					return err
				}

				if domain.Title != "page.github.com" || domain.Address != "https://github.com" {
					return fmt.Errorf("invalid domain data")
				}
				return nil
			},
			wantStatusCode: 200,
		},
		{
			name: "update domain object with duplicate data",
			uuid: domain1.UUID,
			request: DomainRequest{
				Title:   "github.com",
				Address: "https://aws.com",
			},
			checkFunc:      nil,
			wantStatusCode: 400,
		},
		{
			name: "update domain object with incorrect data",
			uuid: domain1.UUID,
			request: DomainRequest{
				Address: "https://google.com",
			},
			checkFunc:      nil,
			wantStatusCode: 400,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			url := fmt.Sprintf("/domains/%s", tt.uuid)
			req, err := http.NewRequest("PUT", url, marshal(tt.request))
			if err != nil {
				t.Errorf("domain UpdateHandler, request failed, error = %v", err)
				return
			}
			router.ServeHTTP(w, req)
			if w.Code != tt.wantStatusCode {
				t.Errorf("domain UpdateHandler failed, status code = %d", w.Code)
				return
			}

			if tt.checkFunc == nil {
				return
			}

			if err := tt.checkFunc(w.Body.String()); err != nil {
				t.Errorf("domain UpdateHandler returned invalid data = %v", w.Body.String())
				return
			}
		})
	}
}

func TestService_GetHandler(t *testing.T) {
	router, repo, cleaner := apiTestSetup()
	defer cleaner()

	domain1, err := repo.Create(DomainRequest{
		Title:   "google.com",
		Address: "https://google.com",
	})
	if err != nil {
		panic(err)
	}

	tests := []struct {
		name           string
		uuid           string
		checkFunc      func(res string) error
		wantStatusCode int
	}{
		{
			name: "get domain object with correct uuid",
			uuid: domain1.UUID,
			checkFunc: func(res string) error {
				domain, err := unmarshalSingle(res)
				if err != nil {
					return err
				}

				if domain.Title != domain1.Title || domain.Address != domain1.Address {
					return fmt.Errorf("invalid domain data")
				}
				return nil
			},
			wantStatusCode: 200,
		},
		{
			name:           "create domain object with incorrect data",
			uuid:           uuid.New().String(),
			checkFunc:      nil,
			wantStatusCode: 404,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			url := fmt.Sprintf("/domains/%s", tt.uuid)
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				t.Errorf("domain GetHandler, request failed, error = %v", err)
				return
			}
			router.ServeHTTP(w, req)
			if w.Code != tt.wantStatusCode {
				t.Errorf("domain GetHandler failed, status code = %d", w.Code)
				return
			}

			if tt.checkFunc == nil {
				return
			}

			if err := tt.checkFunc(w.Body.String()); err != nil {
				t.Errorf("domain GetHandler returned invalid data = %v", w.Body.String())
				return
			}
		})
	}
}

func TestService_RetrieveHandler(t *testing.T) {

	router, repo, cleaner := apiTestSetup()
	defer cleaner()
	const (
		totalItems = 4
		pageLimit  = 2
		offset     = 0
	)

	items := []DomainRequest{
		{
			Title:   "aws",
			Address: "https://aws.com",
		},
		{
			Title:   "google",
			Address: "https://google.com",
		},
		{
			Title:   "svt",
			Address: "https://svt.com",
		},
		{
			Title:   "golang",
			Address: "https://golang.org",
		},
	}

	for _, newItem := range items {
		repo.Create(newItem)
	}

	tests := []struct {
		name           string
		url            string
		checkFunc      func(res string) error
		wantStatusCode int
	}{
		{
			name: "get all domains",
			url:  "/domains",
			checkFunc: func(res string) error {
				domains, err := unmarshalList(res)
				if err != nil {
					return err
				}

				if len(domains) != totalItems {
					return fmt.Errorf("returned wrong number of domains, want %d get %d", totalItems, len(domains))
				}

				return nil
			},
			wantStatusCode: 200,
		},
		{
			name: "get domains with limit",
			url:  "/domains?limit=2",
			checkFunc: func(res string) error {
				domains, err := unmarshalList(res)
				if err != nil {
					return err
				}

				if len(domains) != pageLimit {
					return fmt.Errorf("returned wrong number of domains, want %d get %d", pageLimit, len(domains))
				}
				return nil
			},
			wantStatusCode: 200,
		},
		{
			name: "get domains with limit and offset",
			url:  "/domains?limit=2&offset=2",
			checkFunc: func(res string) error {
				domains, err := unmarshalList(res)
				if err != nil {
					return err
				}

				if len(domains) != pageLimit {
					return fmt.Errorf("returned wrong number of domains, want %d get %d", pageLimit, len(domains))
				}

				return nil
			},
			wantStatusCode: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, err := http.NewRequest("GET", tt.url, nil)
			if err != nil {
				t.Errorf("domain RetrieveHandler, request failed, error = %v", err)
				return
			}
			router.ServeHTTP(w, req)
			if w.Code != 200 {
				t.Errorf("domain RetrieveHandler failed, status code = %d", w.Code)
				return
			}
			if tt.checkFunc == nil {
				return
			}
			if err := tt.checkFunc(w.Body.String()); err != nil {
				t.Errorf("domain RetrieveHandler returned invalid data = %v", w.Body.String())
				return
			}
		})
	}
}

func TestService_DeleteHandler(t *testing.T) {
	router, repo, cleaner := apiTestSetup()
	defer cleaner()

	domain1, err := repo.Create(DomainRequest{
		Title:   "google.com",
		Address: "https://google.com",
	})
	if err != nil {
		panic(err)
	}

	tests := []struct {
		name           string
		uuid           string
		wantStatusCode int
	}{
		{
			name:           "delete domain object with correct uuid",
			uuid:           domain1.UUID,
			wantStatusCode: 204,
		},
		{
			name:           "delete domain object without uuid",
			wantStatusCode: 404,
		},
		{
			name:           "delete domain object with incorrect uuid",
			uuid:           uuid.New().String(),
			wantStatusCode: 404,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			url := fmt.Sprintf("/domains/%s", tt.uuid)
			if tt.uuid == "" {
				url = "/domains"
			}
			req, err := http.NewRequest("DELETE", url, nil)
			if err != nil {
				t.Errorf("domain DeleteHandler, request failed, error = %v", err)
				return
			}
			router.ServeHTTP(w, req)
			if w.Code != tt.wantStatusCode {
				t.Errorf("domain DeleteHandler failed, status code = %d", w.Code)
				return
			}
		})
	}
}
