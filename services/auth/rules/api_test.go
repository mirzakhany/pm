package rules

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
	db, err := gorm.Open(sqlite.Open("/tmp/api_rules.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	repo := NewRepository(db)

	NewService(r, repo, logger)

	err = db.AutoMigrate(&Rule{})
	if err != nil {
		panic(err)
	}
	return r, repo, func() {
		logger.Sync()
		_ = os.Remove("/tmp/api_rules.db")
	}
}

func marshal(data RuleRequest) *bytes.Buffer {
	reqBody, err := json.Marshal(data)

	if err != nil {
		print(err)
	}
	return bytes.NewBuffer(reqBody)
}

func unmarshalSingle(res string) (*Rule, error) {

	type data struct {
		Data Rule `json:"data"`
	}

	var jsonData data
	err := json.Unmarshal([]byte(res), &jsonData)
	if err != nil {
		return nil, err
	}
	return &jsonData.Data, nil
}

func unmarshalList(res string) ([]Rule, error) {
	type data struct {
		Data []Rule `json:"data"`
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
		request        RuleRequest
		checkFunc      func(res string) error
		wantStatusCode int
	}{
		{
			name: "create rule object with correct data",
			request: RuleRequest{
				Subject:  "admin",
				Domain:   "task.com",
				Resource: "sprints",
				Action:   "delete",
				Object:   "*",
			},
			checkFunc: func(res string) error {
				rule, err := unmarshalSingle(res)
				if err != nil {
					return err
				}

				if !utiles.IsValidUUID(rule.UUID) || rule.Subject != "admin" {
					return fmt.Errorf("invalid rule data")
				}
				return nil
			},
			wantStatusCode: 200,
		},
		{
			name:           "create rule object with incorrect data",
			request:        RuleRequest{},
			checkFunc:      nil,
			wantStatusCode: 400,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, err := http.NewRequest("POST", "/rules", marshal(tt.request))
			if err != nil {
				t.Errorf("rule CreateHandler, request failed, error = %v", err)
				return
			}
			router.ServeHTTP(w, req)
			if w.Code != tt.wantStatusCode {
				t.Errorf("rule CreateHandler failed, status code = %d", w.Code)
				return
			}

			if tt.checkFunc == nil {
				return
			}

			if err := tt.checkFunc(w.Body.String()); err != nil {
				t.Errorf("rule CreateHandler returned invalid data = %v", w.Body.String())
				return
			}
		})
	}
}

func TestService_UpdateHandler(t *testing.T) {
	router, repo, cleaner := apiTestSetup()
	defer cleaner()

	rule1, err := repo.Create(RuleRequest{
		Subject:  "admin",
		Domain:   "task.com",
		Resource: "sprints",
		Action:   "delete",
		Object:   "*",
	})
	if err != nil {
		panic(err)
	}

	tests := []struct {
		name           string
		uuid           string
		request        RuleRequest
		checkFunc      func(res string) error
		wantStatusCode int
	}{
		{
			name: "update rule object with correct data",
			uuid: rule1.UUID,
			request: RuleRequest{
				Subject:  "manager",
				Domain:   "task.com",
				Resource: "sprints",
				Action:   "delete",
				Object:   "*",
			},
			checkFunc: func(res string) error {
				rule, err := unmarshalSingle(res)
				if err != nil {
					return err
				}

				if rule.Subject != "manager" {
					return fmt.Errorf("invalid rule data")
				}
				return nil
			},
			wantStatusCode: 200,
		},
		{
			name:           "update rule object with incorrect data",
			uuid:           rule1.UUID,
			request:        RuleRequest{},
			checkFunc:      nil,
			wantStatusCode: 400,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			url := fmt.Sprintf("/rules/%s", tt.uuid)
			req, err := http.NewRequest("PUT", url, marshal(tt.request))
			if err != nil {
				t.Errorf("rule UpdateHandler, request failed, error = %v", err)
				return
			}
			router.ServeHTTP(w, req)
			if w.Code != tt.wantStatusCode {
				t.Errorf("rule UpdateHandler failed, status code = %d", w.Code)
				return
			}

			if tt.checkFunc == nil {
				return
			}

			if err := tt.checkFunc(w.Body.String()); err != nil {
				t.Errorf("rule UpdateHandler returned invalid data = %v", w.Body.String())
				return
			}
		})
	}
}

func TestService_GetHandler(t *testing.T) {
	router, repo, cleaner := apiTestSetup()
	defer cleaner()

	rule1, err := repo.Create(RuleRequest{
		Subject:  "manager",
		Domain:   "task.com",
		Resource: "sprints",
		Action:   "delete",
		Object:   "*",
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
			name: "get rule object with correct uuid",
			uuid: rule1.UUID,
			checkFunc: func(res string) error {
				rule, err := unmarshalSingle(res)
				if err != nil {
					return err
				}

				if rule.Subject != rule1.Subject {
					return fmt.Errorf("invalid rule data")
				}
				return nil
			},
			wantStatusCode: 200,
		},
		{
			name:           "get rule object with incorrect data",
			uuid:           uuid.New().String(),
			checkFunc:      nil,
			wantStatusCode: 404,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			url := fmt.Sprintf("/rules/%s", tt.uuid)
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				t.Errorf("rule GetHandler, request failed, error = %v", err)
				return
			}
			router.ServeHTTP(w, req)
			if w.Code != tt.wantStatusCode {
				t.Errorf("rule GetHandler failed, status code = %d", w.Code)
				return
			}

			if tt.checkFunc == nil {
				return
			}

			if err := tt.checkFunc(w.Body.String()); err != nil {
				t.Errorf("rule GetHandler returned invalid data = %v", w.Body.String())
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

	items := []RuleRequest{
		{
			Subject:  "bob",
			Domain:   "task.com",
			Resource: "sprints",
			Action:   "get",
			Object:   "*",
		},
		{
			Subject:  "rob",
			Domain:   "gitlab.com",
			Resource: "issues",
			Action:   "get",
			Object:   "*",
		},
		{
			Subject:  "bob",
			Domain:   "github.com",
			Resource: "repository",
			Action:   "get",
			Object:   "*",
		},
		{
			Subject:  "admin",
			Domain:   "task.com",
			Resource: "sprints",
			Action:   "delete",
			Object:   "*",
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
			name: "get all rules",
			url:  "/rules",
			checkFunc: func(res string) error {
				rules, err := unmarshalList(res)
				if err != nil {
					return err
				}

				if len(rules) != totalItems {
					return fmt.Errorf("returned wrong number of rules, want %d get %d", totalItems, len(rules))
				}

				return nil
			},
			wantStatusCode: 200,
		},
		{
			name: "get rules with limit",
			url:  "/rules?limit=2",
			checkFunc: func(res string) error {
				rules, err := unmarshalList(res)
				if err != nil {
					return err
				}

				if len(rules) != pageLimit {
					return fmt.Errorf("returned wrong number of rules, want %d get %d", pageLimit, len(rules))
				}
				return nil
			},
			wantStatusCode: 200,
		},
		{
			name: "get rules with limit and offset",
			url:  "/rules?limit=2&offset=2",
			checkFunc: func(res string) error {
				rules, err := unmarshalList(res)
				if err != nil {
					return err
				}

				if len(rules) != pageLimit {
					return fmt.Errorf("returned wrong number of rules, want %d get %d", pageLimit, len(rules))
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
				t.Errorf("rule RetrieveHandler, request failed, error = %v", err)
				return
			}
			router.ServeHTTP(w, req)
			if w.Code != 200 {
				t.Errorf("rule RetrieveHandler failed, status code = %d", w.Code)
				return
			}
			if tt.checkFunc == nil {
				return
			}
			if err := tt.checkFunc(w.Body.String()); err != nil {
				t.Errorf("rule RetrieveHandler returned invalid data = %v", w.Body.String())
				return
			}
		})
	}
}

func TestService_DeleteHandler(t *testing.T) {
	router, repo, cleaner := apiTestSetup()
	defer cleaner()

	rule1, err := repo.Create(RuleRequest{
		Subject:  "admin",
		Domain:   "task.com",
		Resource: "sprints",
		Action:   "delete",
		Object:   "*",
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
			name:           "delete rule object with correct uuid",
			uuid:           rule1.UUID,
			wantStatusCode: 204,
		},
		{
			name:           "delete rule object without uuid",
			wantStatusCode: 404,
		},
		{
			name:           "delete rule object with incorrect uuid",
			uuid:           uuid.New().String(),
			wantStatusCode: 404,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			url := fmt.Sprintf("/rules/%s", tt.uuid)
			if tt.uuid == "" {
				url = "/rules"
			}
			req, err := http.NewRequest("DELETE", url, nil)
			if err != nil {
				t.Errorf("rule DeleteHandler, request failed, error = %v", err)
				return
			}
			router.ServeHTTP(w, req)
			if w.Code != tt.wantStatusCode {
				t.Errorf("rule DeleteHandler failed, status code = %d", w.Code)
				return
			}
		})
	}
}
