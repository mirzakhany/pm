package domains

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
	db, err := gorm.Open(sqlite.Open("/tmp/gorm_domains.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	repo := NewRepository(db)

	err = db.AutoMigrate(&Domain{})
	if err != nil {
		panic(err)
	}
	return db, repo, logger, func() {
		logger.Sync()
		_ = os.Remove("/tmp/gorm_domains.db")
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
		domainReq DomainRequest
		checkFunc func(req DomainRequest, res *Domain) error
		wantErr   bool
	}{
		{
			name: "create new domain object",
			domainReq: DomainRequest{
				Title:   "Google",
				Address: "https://google.com",
			},
			checkFunc: func(req DomainRequest, res *Domain) error {
				if !utiles.IsValidUUID(res.UUID) {
					return fmt.Errorf("response uuid is not valid")
				}

				if req.Title != res.Title || req.Address != req.Address {
					return fmt.Errorf("response is not same as request")
				}

				return nil
			},
			wantErr: false,
		},
		{
			name: "create domain object with duplicate address",
			domainReq: DomainRequest{
				Title:   "Google",
				Address: "https://google.com",
			},
			checkFunc: nil,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.Create(tt.domainReq)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.checkFunc == nil {
				return
			}
			if err := tt.checkFunc(tt.domainReq, got); err != nil {
				t.Errorf("Repository.Create() error = %v", err)
				return
			}
		})
	}

}

func TestRepository_Update(t *testing.T) {

	_, repo, _, cleaner := dbTestSetup()
	defer cleaner()

	d1 := DomainRequest{
		Title:   "aws",
		Address: "https://aws.com",
	}
	d2 := DomainRequest{
		Title:   "Fb",
		Address: "https://fb.com",
	}
	domain1, _ := repo.Create(d1)
	domain2, _ := repo.Create(d2)

	tests := []struct {
		name      string
		uuid      string
		domainReq DomainRequest
		checkFunc func(req DomainRequest, res *Domain) error
		wantErr   bool
	}{
		{
			name: "update aws.com domain title",
			uuid: domain1.UUID,
			domainReq: DomainRequest{
				Title:   "aws website",
				Address: "https://amazon.com",
			},
			checkFunc: func(req DomainRequest, res *Domain) error {
				if req.Title != res.Title || req.Address != req.Address {
					return fmt.Errorf("response is not same as request, req : %v, res: %v", req, res)
				}
				return nil
			},
			wantErr: false,
		},
		{
			name: "update fb address to a already exist duplicate",
			uuid: domain2.UUID,
			domainReq: DomainRequest{
				Title:   "Fb",
				Address: "https://amazon.com",
			},
			checkFunc: nil,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.Update(tt.uuid, tt.domainReq)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.checkFunc == nil {
				return
			}
			if err := tt.checkFunc(tt.domainReq, got); err != nil {
				t.Errorf("Repository.Update() error = %v", err)
				return
			}
		})
	}
}

func TestRepository_Retrieve(t *testing.T) {

	_, repo, _, cleaner := dbTestSetup()
	defer cleaner()

	d1 := DomainRequest{
		Title:   "aws",
		Address: "https://aws.com",
	}
	d2 := DomainRequest{
		Title:   "Fb",
		Address: "https://fb.com",
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

	d1 := DomainRequest{
		Title:   "aws",
		Address: "https://aws.com",
	}
	res, _ := repo.Create(d1)

	item, err := repo.GetByUUID(res.UUID)
	if err != nil {
		t.Errorf("Repository.GetByUUID() error = %v", err)
		return
	}

	if item.Address != d1.Address || item.Title != d1.Title {
		t.Errorf("Repository.GetByUUID() return wrong item = %v", item)
		return
	}
}

func TestRepository_Get(t *testing.T) {

	_, repo, _, cleaner := dbTestSetup()
	defer cleaner()

	d1 := DomainRequest{
		Title:   "aws",
		Address: "https://aws.com",
	}
	res, _ := repo.Create(d1)

	item, err := repo.Get("uuid = ?", res.UUID)
	if err != nil {
		t.Errorf("Repository.Get() error = %v", err)
		return
	}

	if item.Address != d1.Address || item.Title != d1.Title {
		t.Errorf("Repository.Get() return wrong item = %v", item)
		return
	}
}

func TestRepository_Delete(t *testing.T) {

	_, repo, _, cleaner := dbTestSetup()
	defer cleaner()

	d1 := DomainRequest{
		Title:   "aws",
		Address: "https://aws.com",
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
