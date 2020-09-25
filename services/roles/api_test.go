package roles

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
	db, err := gorm.Open(sqlite.Open("/tmp/api_roles.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	repo := NewRepository(db)

	NewService(r, repo, logger)

	err = db.AutoMigrate(&Role{})
	if err != nil {
		panic(err)
	}
	return r, repo, func() {
		logger.Sync()
		_ = os.Remove("/tmp/api_roles.db")
	}
}

func marshal(data RoleRequest) *bytes.Buffer {
	reqBody, err := json.Marshal(data)

	if err != nil {
		print(err)
	}
	return bytes.NewBuffer(reqBody)
}

func unmarshalSingle(res string) (*Role, error) {

	type data struct {
		Data Role `json:"data"`
	}

	var jsonData data
	err := json.Unmarshal([]byte(res), &jsonData)
	if err != nil {
		return nil, err
	}
	return &jsonData.Data, nil
}

func unmarshalList(res string) ([]Role, error) {
	type data struct {
		Data []Role `json:"data"`
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
		request        RoleRequest
		checkFunc      func(res string) error
		wantStatusCode int
	}{
		{
			name: "create role object with correct data",
			request: RoleRequest{
				Title:   "admin",
			},
			checkFunc: func(res string) error {
				role, err := unmarshalSingle(res)
				if err != nil {
					return err
				}

				if !utiles.IsValidUUID(role.UUID) || role.Title != "admin"{
					return fmt.Errorf("invalid role data")
				}
				return nil
			},
			wantStatusCode: 200,
		},
		{
			name: "create role object with duplicate data",
			request: RoleRequest{
				Title:   "admin",
			},
			checkFunc:      nil,
			wantStatusCode: 400,
		},
		{
			name: "create role object with incorrect data",
			request: RoleRequest{
			},
			checkFunc:      nil,
			wantStatusCode: 400,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, err := http.NewRequest("POST", "/roles", marshal(tt.request))
			if err != nil {
				t.Errorf("role CreateHandler, request failed, error = %v", err)
				return
			}
			router.ServeHTTP(w, req)
			if w.Code != tt.wantStatusCode {
				t.Errorf("role CreateHandler failed, status code = %d", w.Code)
				return
			}

			if tt.checkFunc == nil {
				return
			}

			if err := tt.checkFunc(w.Body.String()); err != nil {
				t.Errorf("role CreateHandler returned invalid data = %v", w.Body.String())
				return
			}
		})
	}
}

func TestService_UpdateHandler(t *testing.T) {
	router, repo, cleaner := apiTestSetup()
	defer cleaner()

	role1, err := repo.Create(RoleRequest{
		Title:   "admin",
	})
	if err != nil {
		panic(err)
	}

	_, err = repo.Create(RoleRequest{
		Title:   "owner",
	})
	if err != nil {
		panic(err)
	}

	tests := []struct {
		name           string
		uuid           string
		request        RoleRequest
		checkFunc      func(res string) error
		wantStatusCode int
	}{
		{
			name: "update role object with correct data",
			uuid: role1.UUID,
			request: RoleRequest{
				Title:   "manager",
			},
			checkFunc: func(res string) error {
				role, err := unmarshalSingle(res)
				if err != nil {
					return err
				}

				if role.Title != "manager" {
					return fmt.Errorf("invalid role data")
				}
				return nil
			},
			wantStatusCode: 200,
		},
		{
			name: "update role object with duplicate data",
			uuid: role1.UUID,
			request: RoleRequest{
				Title:   "owner",
			},
			checkFunc:      nil,
			wantStatusCode: 400,
		},
		{
			name: "update role object with incorrect data",
			uuid: role1.UUID,
			request: RoleRequest{
			},
			checkFunc:      nil,
			wantStatusCode: 400,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			url := fmt.Sprintf("/roles/%s", tt.uuid)
			req, err := http.NewRequest("PUT", url, marshal(tt.request))
			if err != nil {
				t.Errorf("role UpdateHandler, request failed, error = %v", err)
				return
			}
			router.ServeHTTP(w, req)
			if w.Code != tt.wantStatusCode {
				t.Errorf("role UpdateHandler failed, status code = %d", w.Code)
				return
			}

			if tt.checkFunc == nil {
				return
			}

			if err := tt.checkFunc(w.Body.String()); err != nil {
				t.Errorf("role UpdateHandler returned invalid data = %v", w.Body.String())
				return
			}
		})
	}
}

func TestService_GetHandler(t *testing.T) {
	router, repo, cleaner := apiTestSetup()
	defer cleaner()

	role1, err := repo.Create(RoleRequest{
		Title:   "admin",
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
			name: "get role object with correct uuid",
			uuid: role1.UUID,
			checkFunc: func(res string) error {
				role, err := unmarshalSingle(res)
				if err != nil {
					return err
				}

				if role.Title != role1.Title {
					return fmt.Errorf("invalid role data")
				}
				return nil
			},
			wantStatusCode: 200,
		},
		{
			name:           "get role object with incorrect data",
			uuid:           uuid.New().String(),
			checkFunc:      nil,
			wantStatusCode: 404,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			url := fmt.Sprintf("/roles/%s", tt.uuid)
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				t.Errorf("role GetHandler, request failed, error = %v", err)
				return
			}
			router.ServeHTTP(w, req)
			if w.Code != tt.wantStatusCode {
				t.Errorf("role GetHandler failed, status code = %d", w.Code)
				return
			}

			if tt.checkFunc == nil {
				return
			}

			if err := tt.checkFunc(w.Body.String()); err != nil {
				t.Errorf("role GetHandler returned invalid data = %v", w.Body.String())
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
	)

	items := []RoleRequest{
		{
			Title:   "admin",
		},
		{
			Title:   "owner",
		},
		{
			Title:   "writer",
		},
		{
			Title:   "viewer",
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
			name: "get all roles",
			url:  "/roles",
			checkFunc: func(res string) error {
				roles, err := unmarshalList(res)
				if err != nil {
					return err
				}

				if len(roles) != totalItems {
					return fmt.Errorf("returned wrong number of roles, want %d get %d", totalItems, len(roles))
				}

				return nil
			},
			wantStatusCode: 200,
		},
		{
			name: "get roles with limit",
			url:  "/roles?limit=2",
			checkFunc: func(res string) error {
				roles, err := unmarshalList(res)
				if err != nil {
					return err
				}

				if len(roles) != pageLimit {
					return fmt.Errorf("returned wrong number of roles, want %d get %d", pageLimit, len(roles))
				}
				return nil
			},
			wantStatusCode: 200,
		},
		{
			name: "get roles with limit and offset",
			url:  "/roles?limit=2&offset=2",
			checkFunc: func(res string) error {
				roles, err := unmarshalList(res)
				if err != nil {
					return err
				}

				if len(roles) != pageLimit {
					return fmt.Errorf("returned wrong number of roles, want %d get %d", pageLimit, len(roles))
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
				t.Errorf("role RetrieveHandler, request failed, error = %v", err)
				return
			}
			router.ServeHTTP(w, req)
			if w.Code != 200 {
				t.Errorf("role RetrieveHandler failed, status code = %d", w.Code)
				return
			}
			if tt.checkFunc == nil {
				return
			}
			if err := tt.checkFunc(w.Body.String()); err != nil {
				t.Errorf("role RetrieveHandler returned invalid data = %v", w.Body.String())
				return
			}
		})
	}
}

func TestService_DeleteHandler(t *testing.T) {
	router, repo, cleaner := apiTestSetup()
	defer cleaner()

	role1, err := repo.Create(RoleRequest{
		Title:   "owner",
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
			name:           "delete role object with correct uuid",
			uuid:           role1.UUID,
			wantStatusCode: 204,
		},
		{
			name:           "delete role object without uuid",
			wantStatusCode: 404,
		},
		{
			name:           "delete role object with incorrect uuid",
			uuid:           uuid.New().String(),
			wantStatusCode: 404,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			url := fmt.Sprintf("/roles/%s", tt.uuid)
			if tt.uuid == "" {
				url = "/roles"
			}
			req, err := http.NewRequest("DELETE", url, nil)
			if err != nil {
				t.Errorf("role DeleteHandler, request failed, error = %v", err)
				return
			}
			router.ServeHTTP(w, req)
			if w.Code != tt.wantStatusCode {
				t.Errorf("role DeleteHandler failed, status code = %d", w.Code)
				return
			}
		})
	}
}
