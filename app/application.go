package app

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/transip/gotransip/v6"
	"github.com/transip/gotransip/v6/authenticator"
	transipdomain "github.com/transip/gotransip/v6/domain"
	transiphaip "github.com/transip/gotransip/v6/haip"
	transipvps "github.com/transip/gotransip/v6/vps"
	"go.uber.org/zap"
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

	Logger *zap.Logger

	productList *ProductList
	productInfo *ProductInfo

	grid     *tview.Grid
	tviewApp *tview.Application
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

	var err error
	a.Logger, err = NewLogger()

	if err != nil {
		panic(err)
	}

	a.productList = &ProductList{
		app:            a,
		vpsRepo:        a.vpsRepo,
		bigstorageRepo: a.bigstorageRepo,
		domainRepo:     a.domainRepo,
		haipRepo:       a.haipRepo,
	}

	a.productList.init()

	a.productInfo = &ProductInfo{
		app: a,
	}

	a.productInfo.init()

	a.grid = tview.NewGrid().
		SetColumns(-1, -3).
		SetBorders(true)

	a.grid.AddItem(a.productList.treeView, 0, 0, 1, 1, 0, 0, true)
	a.grid.AddItem(a.productInfo.grid, 0, 1, 1, 1, 0, 0, true)
}

func (a *Application) Run() {
	a.init()
	a.tviewApp = tview.NewApplication()

	a.tviewApp.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		a.Logger.Debug("received input")
		return event
	})
	if err := a.tviewApp.SetRoot(a.grid, true).EnableMouse(false).Run(); err != nil {
		panic(err)
	}
}

func NewLogger() (*zap.Logger, error){
	cfg := zap.NewDevelopmentConfig()
	cfg.OutputPaths = []string{
		"transterm.log",
	}

	cfg.ErrorOutputPaths = []string{
		"transterm.log",
	}

	return cfg.Build()
}
