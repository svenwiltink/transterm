package main

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/transip/gotransip/v6"
	"github.com/transip/gotransip/v6/authenticator"
	transiphaip "github.com/transip/gotransip/v6/haip"
	transipvps "github.com/transip/gotransip/v6/vps"
)

var app *Application

type Application struct {
	vpsRepo  transipvps.Repository
	haipRepo transiphaip.Repository

	productList *ProductList
}

type Config struct {
	AccountName    string
	PrivateKeyPath string
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
		vpsRepo:  transipvps.Repository{Client: client},
		haipRepo: transiphaip.Repository{Client: client},
	}
}

func (a *Application) init() {
	a.productList = new(ProductList)
	a.productList.init()
}

func (a *Application) Run() {
	a.init()
	if err := tview.NewApplication().SetRoot(a.productList.treeView, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

type ProductList struct {
	treeView *tview.TreeView
}

func (pl *ProductList) init() {

	rootDir := "Products"
	root := tview.NewTreeNode(rootDir).
		SetColor(tcell.ColorRed)

	pl.treeView = tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)

	vps := tview.NewTreeNode("Vps").
		SetReference("Vps").
		SetSelectable(true).
		SetColor(tcell.ColorGreen)

	haip := tview.NewTreeNode("Haip").
		SetReference("Haip").
		SetSelectable(true).
		SetColor(tcell.ColorGreen)

	root.AddChild(vps)
	root.AddChild(haip)

	vps.SetSelectedFunc(func() {
		vps.SetExpanded(!vps.IsExpanded())

		if !vps.IsExpanded() {
			return
		}

		vpses, err := app.vpsRepo.GetAll()
		if err != nil {
			panic(err)
		}

		vps.ClearChildren()

		for _, v := range vpses {
			node := tview.NewTreeNode(v.Description).
				SetReference(v.Name).
				SetSelectable(true).
				SetColor(tcell.ColorGreen)
			vps.AddChild(node)
		}
	})

	vps.SetExpanded(false)

	haip.SetSelectedFunc(func() {
		haip.SetExpanded(!haip.IsExpanded())

		if !haip.IsExpanded() {
			return
		}

		haips, err := app.haipRepo.GetAll()
		if err != nil {
			panic(err)
		}

		haip.ClearChildren()

		for _, h := range haips {
			node := tview.NewTreeNode(h.Description).
				SetReference(h.Name).
				SetSelectable(true).
				SetColor(tcell.ColorGreen)
			haip.AddChild(node)
		}
	})

	haip.SetExpanded(false)
}

// Show a navigable tree view of the current directory.
func main() {

	app = NewApplication(Config{
		AccountName:    "swiltink",
		PrivateKeyPath: "transip.key",
	})

	app.Run()
}
