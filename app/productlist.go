package app

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	transipdomain "github.com/transip/gotransip/v6/domain"
	transiphaip "github.com/transip/gotransip/v6/haip"
	transipvps "github.com/transip/gotransip/v6/vps"
	"go.uber.org/zap"
	"time"
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

	pl.app.Logger.Debug("Initialising UI", zap.String("element", "productList"))

	rootDir := "Products"
	root := tview.NewTreeNode(rootDir).
		SetColor(tcell.ColorRed)

	pl.treeView = tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)

	pl.treeView.SetBorder(true)

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

		pl.app.Logger.Debug("vps product selected", zap.Bool("expanded", vps.IsExpanded()))

		if !vps.IsExpanded() {
			return
		}

		pl.app.Logger.Debug("fetching vpses")
		now := time.Now()
		vpses, err := pl.vpsRepo.GetAll()
		if err != nil {
			pl.app.Logger.Error("error fetching vpses", zap.Error(err), zap.Duration("duration", time.Since(now)))
			panic(err)
		}

		pl.app.Logger.Debug("done fetching vpses", zap.Duration("duration", time.Since(now)))

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

		pl.app.Logger.Debug("bigstorage product selected", zap.Bool("expanded", vps.IsExpanded()))

		if !bigstorage.IsExpanded() {
			return
		}

		pl.app.Logger.Debug("fetching bigstorages")
		now := time.Now()
		bigStorages, err := pl.bigstorageRepo.GetAll()
		if err != nil {
			pl.app.Logger.Error("error fetching bigstorages", zap.Error(err), zap.Duration("duration", time.Since(now)))
			panic(err)
		}

		pl.app.Logger.Debug("done fetching bigstorages", zap.Duration("duration", time.Since(now)))

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

		pl.app.Logger.Debug("haip product selected", zap.Bool("expanded", vps.IsExpanded()))

		if !haip.IsExpanded() {
			return
		}

		pl.app.Logger.Debug("fetching haips")
		now := time.Now()
		haips, err := pl.haipRepo.GetAll()
		if err != nil {
			pl.app.Logger.Error("error fetching haips", zap.Error(err), zap.Duration("duration", time.Since(now)))
			panic(err)
		}

		pl.app.Logger.Debug("done fetching haips", zap.Duration("duration", time.Since(now)))

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

		pl.app.Logger.Debug("domains product selected", zap.Bool("expanded", vps.IsExpanded()))

		if !domain.IsExpanded() {
			return
		}

		pl.app.Logger.Debug("fetching domains")
		now := time.Now()
		domains, err := pl.domainRepo.GetAll()
		if err != nil {
			pl.app.Logger.Error("error fetching domains", zap.Error(err), zap.Duration("duration", time.Since(now)))
			panic(err)
		}

		pl.app.Logger.Info("done fetching domains", zap.Duration("duration", time.Since(now)))

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

	pl.treeView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		pl.app.Logger.Debug("key pressed", zap.String("widget", "productlist"), zap.Int32("rune", event.Rune()), zap.Int16("key", int16(event.Key())))

		switch event.Key() {
		case tcell.KeyTAB:
			pl.app.Logger.Debug("detected TAB press", zap.String("widget", "productlist"))

			pl.app.Logger.Debug("checking for active ProductInfo window", zap.Bool("present", pl.app.productInfo.currentView != nil))
			if pl.app.productInfo.currentView != nil {
				pl.app.Logger.Debug("changing focus to ProductInfo")
				pl.app.productInfo.currentView.Focus()
			}
		}
		return event
	})
}
