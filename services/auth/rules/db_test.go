package rules

import (
	"fmt"
	"os"
	"projectmanager/pkg/utiles"
	"reflect"
	"testing"

	"gorm.io/driver/sqlite"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func dbTestSetup() (*gorm.DB, *Repository, *zap.Logger, func()) {
	db, err := gorm.Open(sqlite.Open("/tmp/gorm_rules.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	repo := NewRepository(db)

	err = db.AutoMigrate(&Rule{})
	if err != nil {
		panic(err)
	}
	return db, repo, logger, func() {
		logger.Sync()
		_ = os.Remove("/tmp/gorm_rules.db")
	}
}

func TestNewRepository(t *testing.T) {

	db, _, _, cleaner := dbTestSetup()
	defer cleaner()

	type args struct {
		db *gorm.DB
	}
	tests := []struct {
		name string
		args args
		want *Repository
	}{
		{
			"create new repo with nil",
			args{db: nil},
			&Repository{DB: nil},
		},
		{
			"create new repo with db",
			args{db: db},
			&Repository{DB: db},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRepository(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepository_Create(t *testing.T) {

	_, repo, _, cleaner := dbTestSetup()
	defer cleaner()

	tests := []struct {
		name      string
		ruleReq RuleRequest
		checkFunc func(req RuleRequest, res *Rule) error
		wantErr   bool
	}{
		{
			name: "create new rule object",
			ruleReq: RuleRequest{
				Subject:  "bob",
				Domain:   "task.com",
				Resource: "sprints",
				Action:   "get",
				Object:   "*",
			},
			checkFunc: func(req RuleRequest, res *Rule) error {
				if !utiles.IsValidUUID(res.UUID) {
					return fmt.Errorf("response uuid is not valid")
				}

				return nil
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.Create(tt.ruleReq)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.checkFunc == nil {
				return
			}
			if err := tt.checkFunc(tt.ruleReq, got); err != nil {
				t.Errorf("Repository.Create() error = %v", err)
				return
			}
		})
	}

}

func TestRepository_Update(t *testing.T) {

	_, repo, _, cleaner := dbTestSetup()
	defer cleaner()

	d1 := RuleRequest{
		Subject:  "bob",
		Domain:   "task.com",
		Resource: "sprints",
		Action:   "get",
		Object:   "*",
	}
	rule1, _ := repo.Create(d1)

	tests := []struct {
		name      string
		uuid      string
		ruleReq RuleRequest
		checkFunc func(req RuleRequest, res *Rule) error
		wantErr   bool
	}{
		{
			name: "update rule title",
			uuid: rule1.UUID,
			ruleReq: RuleRequest{
				Subject:  "bob",
				Domain:   "api.task.com",
				Resource: "sprints",
				Action:   "get",
				Object:   "*",
			},
			checkFunc: func(req RuleRequest, res *Rule) error {
				if req.Domain != res.Domain {
					return fmt.Errorf("response is not same as request, req : %v, res: %v", req, res)
				}
				return nil
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.Update(tt.uuid, tt.ruleReq)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.checkFunc == nil {
				return
			}
			if err := tt.checkFunc(tt.ruleReq, got); err != nil {
				t.Errorf("Repository.Update() error = %v", err)
				return
			}
		})
	}
}

func TestRepository_Retrieve(t *testing.T) {

	_, repo, _, cleaner := dbTestSetup()
	defer cleaner()

	d1 := RuleRequest{
		Subject:  "bob",
		Domain:   "task.com",
		Resource: "sprints",
		Action:   "get",
		Object:   "*",
	}
	d2 := RuleRequest{
		Subject:  "rob",
		Domain:   "task.com",
		Resource: "sprints",
		Action:   "get",
		Object:   "*",
	}
	repo.Create(d1)
	repo.Create(d2)

	_, total, err := repo.Retrieve(0, 10)
	if err != nil {
		t.Errorf("Repository.Retrieve() error = %v", err)
		return
	}

	if total != 2 {
		t.Errorf("Repository.Retrieve() return wrong number of result = %d", total)
		return
	}
}

func TestRepository_RetrievePagination(t *testing.T) {

	_, repo, _, cleaner := dbTestSetup()
	defer cleaner()

	const (
		totalItems = 4
		pageLimit  = 2
		offset     = 0
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

	page1Items, total, err := repo.Retrieve(offset, pageLimit)
	if err != nil {
		t.Errorf("Repository.Retrieve() error = %v", err)
		return
	}

	if len(page1Items) != pageLimit {
		t.Errorf("Repository.Retrieve() return wrong number of result in one page, %d", len(page1Items))
		return
	}

	if total != totalItems {
		t.Errorf("Repository.Retrieve() return wrong number of result = %d", total)
		return
	}

	page2Items, total, err := repo.Retrieve(offset, pageLimit)
	if err != nil {
		t.Errorf("Repository.Retrieve() error = %v", err)
		return
	}

	if len(page2Items) != pageLimit {
		t.Error("Repository.Retrieve() return wrong number of result in on page")
		return
	}
}

func TestRepository_GetByUUID(t *testing.T) {

	_, repo, _, cleaner := dbTestSetup()
	defer cleaner()

	d1 := RuleRequest{
		Subject:  "admin",
		Domain:   "task.com",
		Resource: "sprints",
		Action:   "delete",
		Object:   "*",
	}
	res, _ := repo.Create(d1)

	item, err := repo.GetByUUID(res.UUID)
	if err != nil {
		t.Errorf("Repository.GetByUUID() error = %v", err)
		return
	}

	if item.Subject != d1.Subject && item.Domain != d1.Domain && item.Resource != d1.Resource && item.Action != d1.Action || item.Object != d1.Object {
		t.Errorf("Repository.GetByUUID() return wrong item = %v", item)
		return
	}
}

func TestRepository_Get(t *testing.T) {

	_, repo, _, cleaner := dbTestSetup()
	defer cleaner()

	d1 := RuleRequest{
		Subject:  "admin",
		Domain:   "task.com",
		Resource: "sprints",
		Action:   "delete",
		Object:   "*",
	}
	res, _ := repo.Create(d1)

	item, err := repo.Get("uuid = ?", res.UUID)
	if err != nil {
		t.Errorf("Repository.Get() error = %v", err)
		return
	}

	if item.Subject != d1.Subject && item.Domain != d1.Domain && item.Resource != d1.Resource && item.Action != d1.Action || item.Object != d1.Object {
		t.Errorf("Repository.Get() return wrong item = %v", item)
		return
	}
}

func TestRepository_Delete(t *testing.T) {

	_, repo, _, cleaner := dbTestSetup()
	defer cleaner()

	d1 := RuleRequest{
		Subject:  "admin",
		Domain:   "task.com",
		Resource: "sprints",
		Action:   "delete",
		Object:   "*",
	}
	res, _ := repo.Create(d1)

	err := repo.Delete(res.UUID)
	if err != nil {
		t.Errorf("Repository.Delete() error = %v", err)
		return
	}

	item, err := repo.GetByUUID(res.UUID)
	if err == nil {
		t.Errorf("Repository.GetByUUID() returned deleted item %v", item)
		return
	}
}
