package roles

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
	db, err := gorm.Open(sqlite.Open("/tmp/gorm_roles.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	repo := NewRepository(db)

	err = db.AutoMigrate(&Role{})
	if err != nil {
		panic(err)
	}
	return db, repo, logger, func() {
		logger.Sync()
		_ = os.Remove("/tmp/gorm_roles.db")
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
		roleReq RoleRequest
		checkFunc func(req RoleRequest, res *Role) error
		wantErr   bool
	}{
		{
			name: "create new role object",
			roleReq: RoleRequest{
				Title:   "Google",
			},
			checkFunc: func(req RoleRequest, res *Role) error {
				if !utiles.IsValidUUID(res.UUID) {
					return fmt.Errorf("response uuid is not valid")
				}

				if req.Title != res.Title {
					return fmt.Errorf("response is not same as request")
				}

				return nil
			},
			wantErr: false,
		},
		{
			name: "create role object with duplicate address",
			roleReq: RoleRequest{
				Title:   "Google",
			},
			checkFunc: nil,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.Create(tt.roleReq)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.checkFunc == nil {
				return
			}
			if err := tt.checkFunc(tt.roleReq, got); err != nil {
				t.Errorf("Repository.Create() error = %v", err)
				return
			}
		})
	}

}

func TestRepository_Update(t *testing.T) {

	_, repo, _, cleaner := dbTestSetup()
	defer cleaner()

	d1 := RoleRequest{
		Title:   "admin",
	}
	d2 := RoleRequest{
		Title:   "owner",
	}
	role1, _ := repo.Create(d1)
	role2, _ := repo.Create(d2)

	tests := []struct {
		name      string
		uuid      string
		roleReq RoleRequest
		checkFunc func(req RoleRequest, res *Role) error
		wantErr   bool
	}{
		{
			name: "update role title",
			uuid: role1.UUID,
			roleReq: RoleRequest{
				Title:   "administrator",
			},
			checkFunc: func(req RoleRequest, res *Role) error {
				if req.Title != res.Title {
					return fmt.Errorf("response is not same as request, req : %v, res: %v", req, res)
				}
				return nil
			},
			wantErr: false,
		},
		{
			name: "update role address to a already exist duplicate",
			uuid: role2.UUID,
			roleReq: RoleRequest{
				Title:   "administrator",
			},
			checkFunc: nil,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.Update(tt.uuid, tt.roleReq)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.checkFunc == nil {
				return
			}
			if err := tt.checkFunc(tt.roleReq, got); err != nil {
				t.Errorf("Repository.Update() error = %v", err)
				return
			}
		})
	}
}

func TestRepository_Retrieve(t *testing.T) {

	_, repo, _, cleaner := dbTestSetup()
	defer cleaner()

	d1 := RoleRequest{
		Title:   "admin",
	}
	d2 := RoleRequest{
		Title:   "owner",
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

	page1Items, total, err := repo.Retrieve(offset, pageLimit)
	if err != nil {
		t.Errorf("Repository.Retrieve() error = %v", err)
		return
	}

	if len(page1Items) != pageLimit {
		t.Error("Repository.Retrieve() return wrong number of result in on page")
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

	d1 := RoleRequest{
		Title:   "admin",
	}
	res, _ := repo.Create(d1)

	item, err := repo.GetByUUID(res.UUID)
	if err != nil {
		t.Errorf("Repository.GetByUUID() error = %v", err)
		return
	}

	if item.Title != d1.Title {
		t.Errorf("Repository.GetByUUID() return wrong item = %v", item)
		return
	}
}

func TestRepository_Get(t *testing.T) {

	_, repo, _, cleaner := dbTestSetup()
	defer cleaner()

	d1 := RoleRequest{
		Title:   "admin",
	}
	res, _ := repo.Create(d1)

	item, err := repo.Get("uuid = ?", res.UUID)
	if err != nil {
		t.Errorf("Repository.Get() error = %v", err)
		return
	}

	if item.Title != d1.Title {
		t.Errorf("Repository.Get() return wrong item = %v", item)
		return
	}
}

func TestRepository_Delete(t *testing.T) {

	_, repo, _, cleaner := dbTestSetup()
	defer cleaner()

	d1 := RoleRequest{
		Title:   "owner",
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
