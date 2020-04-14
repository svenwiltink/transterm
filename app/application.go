package app

import (
	"github.com/rivo/tview"
	"github.com/transip/gotransip/v6"
	"github.com/transip/gotransip/v6/authenticator"
	transipdomain "github.com/transip/gotransip/v6/domain"
	transiphaip "github.com/transip/gotransip/v6/haip"
	transipvps "github.com/transip/gotransip/v6/vps"
)

type Config struct {
	AccountName    string
	PrivateKeyPath string
}

type Application struct {
	vpsRepo        transipvps.Repository
	bigstorageRepo transipvps.BigStorageRepository
	haipRepo       transiphaip.Repository
	domainRepo     transipdomain.Repository

	productList *ProductList

	grid *tview.Grid
}

func NewApplication(config Config) *Application {
	tokenCache, err := authenticator.NewFileTokenCache(".token-cache")
	if err != nil {
		panic(err)
	}

	client, err := gotransip.NewClient(gotransip.ClientConfiguration{
		AccountName:    config.AccountName,
		PrivateKeyPath: config.PrivateKeyPath,
		TestMode:       true,
		TokenCache:     tokenCache,
	})

	return &Application{
		vpsRepo:        transipvps.Repository{Client: client},
		bigstorageRepo: transipvps.BigStorageRepository{Client: client},
		haipRepo:       transiphaip.Repository{Client: client},
		domainRepo:     transipdomain.Repository{Client: client},
	}
}

func (a *Application) init() {
	a.productList = &ProductList{
		vpsRepo:        a.vpsRepo,
		bigstorageRepo: a.bigstorageRepo,
		domainRepo:     a.domainRepo,
		haipRepo:       a.haipRepo,
	}

	a.productList.init()

	a.grid = tview.NewGrid().
		SetColumns(-1, -3).
		SetBorders(true)

	text := tview.NewTextView().SetText("LEKKER BELLEN MET HEM")

	a.grid.AddItem(a.productList.treeView, 0, 0, 1, 1, 0, 0, true)
	a.grid.AddItem(text, 0, 1, 1, 1, 0, 0, false)
}

func (a *Application) Run() {
	a.init()
	if err := tview.NewApplication().SetRoot(a.grid, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
