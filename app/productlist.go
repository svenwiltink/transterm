package app

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	transipdomain "github.com/transip/gotransip/v6/domain"
	transiphaip "github.com/transip/gotransip/v6/haip"
	transipvps "github.com/transip/gotransip/v6/vps"
)

type ProductList struct {
	app            *Application
	treeView       *tview.TreeView
	vpsRepo        transipvps.Repository
	bigstorageRepo transipvps.BigStorageRepository
	haipRepo       transiphaip.Repository
	domainRepo     transipdomain.Repository
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

	bigstorage := tview.NewTreeNode("BigStorage").
		SetReference("BigStorage").
		SetSelectable(true).
		SetColor(tcell.ColorGreen)

	haip := tview.NewTreeNode("Haip").
		SetReference("Haip").
		SetSelectable(true).
		SetColor(tcell.ColorGreen)

	domain := tview.NewTreeNode("Domain").
		SetReference("Domain").
		SetSelectable(true).
		SetColor(tcell.ColorGreen)

	root.AddChild(vps)
	root.AddChild(bigstorage)
	root.AddChild(haip)
	root.AddChild(domain)

	vps.SetSelectedFunc(func() {
		vps.SetExpanded(!vps.IsExpanded())

		if !vps.IsExpanded() {
			return
		}

		vpses, err := pl.vpsRepo.GetAll()
		if err != nil {
			panic(err)
		}

		vps.ClearChildren()

		for _, v := range vpses {
			node := tview.NewTreeNode(fmt.Sprintf("%s (%s)", v.Name, v.Description)).
				SetReference(v.Name).
				SetSelectable(true).
				SetColor(tcell.ColorGreen).
				SetSelectedFunc(pl.app.productInfo.ShowVpsFunc(v.Name))
			vps.AddChild(node)
		}
	})

	vps.SetExpanded(false)

	bigstorage.SetSelectedFunc(func() {
		bigstorage.SetExpanded(!bigstorage.IsExpanded())

		if !bigstorage.IsExpanded() {
			return
		}

		bigStorages, err := pl.bigstorageRepo.GetAll()
		if err != nil {
			panic(err)
		}

		bigstorage.ClearChildren()

		for _, b := range bigStorages {
			node := tview.NewTreeNode(fmt.Sprintf("%s (%s)", b.Name, b.Description)).
				SetReference(b.Name).
				SetSelectable(true).
				SetColor(tcell.ColorGreen)
			bigstorage.AddChild(node)
		}
	})

	bigstorage.SetExpanded(false)

	haip.SetSelectedFunc(func() {
		haip.SetExpanded(!haip.IsExpanded())

		if !haip.IsExpanded() {
			return
		}

		haips, err := pl.haipRepo.GetAll()
		if err != nil {
			panic(err)
		}

		haip.ClearChildren()

		for _, h := range haips {
			node := tview.NewTreeNode(fmt.Sprintf("%s (%s)", h.Name, h.Description)).
				SetReference(h.Name).
				SetSelectable(true).
				SetColor(tcell.ColorGreen)
			haip.AddChild(node)
		}
	})

	haip.SetExpanded(false)

	domain.SetSelectedFunc(func() {
		domain.SetExpanded(!domain.IsExpanded())

		if !domain.IsExpanded() {
			return
		}

		domains, err := pl.domainRepo.GetAll()
		if err != nil {
			panic(err)
		}

		domain.ClearChildren()

		for _, d := range domains {
			node := tview.NewTreeNode(d.Name).
				SetReference(d.Name).
				SetSelectable(true).
				SetColor(tcell.ColorGreen)
			domain.AddChild(node)
		}
	})

	domain.SetExpanded(false)
}
