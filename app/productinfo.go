package app

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	transipvps "github.com/transip/gotransip/v6/vps"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type Focuser interface {
	Focus()
}

type ProductInfo struct {
	app *Application

	grid *tview.Grid

	currentView Focuser
}

func (pi *ProductInfo) init() {
	pi.grid = tview.NewGrid()
}

func (pi *ProductInfo) ShowVpsFunc(vpsName string) func() {

	return func() {
		info := &VpsInfo{app: pi.app}
		pi.currentView = info

		info.ShowVps(pi.grid, vpsName)
	}
}

type VpsInfo struct {
	app       *Application
	overview  *tview.Table
	network   *tview.Table
	backups   *tview.Table
	snapshots *tview.Table
}

func (v *VpsInfo) ShowVps(grid *tview.Grid, vpsName string) {
	grid.Clear()

	grid.SetColumns(50, -1)
	grid.SetMinSize(20, 0)

	v.app.Logger.Debug("fetching vps", zap.String("name", vpsName))
	now := time.Now()
	vps, err := v.app.vpsRepo.GetByName(vpsName)
	if err != nil {
		v.app.Logger.Error("error fetching vps", zap.Error(err), zap.String("name", vpsName), zap.Duration("duration", time.Since(now)))
		panic(err)
	}

	v.app.Logger.Debug("done fetching vps", zap.String("name", vpsName), zap.Duration("duration", time.Since(now)))

	v.createNetwork(vpsName)
	v.createSnapshots(vpsName)
	v.createBackups(vpsName)
	v.createOverview(vps)

	// vertical layout first
	grid.AddItem(v.overview, 0, 0, 1, 2, 0, 0, true)
	grid.AddItem(v.backups, 1, 0, 1, 2, 0, 0, true)
	grid.AddItem(v.snapshots, 2, 0, 1, 2, 0, 0, true)
	grid.AddItem(v.network, 3, 0, 1, 2, 0, 0, true)

	// horizontal after 100 px
	grid.AddItem(v.overview, 0, 0, 1, 1, 0, 100, true)
	grid.AddItem(v.backups, 1, 0, 1, 1, 0, 100, true)
	grid.AddItem(v.snapshots, 2, 0, 1, 1, 0, 100, true)
	grid.AddItem(v.network, 0, 1, 3, 1, 0, 100, true)
}

func (v *VpsInfo) createNetwork(vpsName string) {
	v.network = tview.NewTable().SetSelectable(false, false)
	v.network.SetTitle("Network").SetBorder(true)

	v.network.SetCellSimple(0, 0, "IP").
		SetCellSimple(0, 1, "Subnet").
		SetCellSimple(0, 2, "Gateway").
		SetCellSimple(0, 3, "Reverse DNS")

	v.app.Logger.Debug("fetching ip data", zap.String("name", vpsName))
	now := time.Now()
	ips, err := v.app.vpsRepo.GetIPAddresses(vpsName)
	if err != nil {
		v.app.Logger.Error("error fetching ip data", zap.Error(err), zap.String("name", vpsName), zap.Duration("duration", time.Since(now)))
		panic(err)
	}

	v.app.Logger.Debug("done fetching ip data", zap.String("name", vpsName), zap.Duration("duration", time.Since(now)))

	for i, ip := range ips {
		v.network.SetCellSimple(i+1, 0, ip.Address.String())
		text, _ := ip.SubnetMask.MarshalText()
		v.network.SetCellSimple(i+1, 1, string(text))
		v.network.SetCellSimple(i+1, 2, ip.Gateway.String())
		v.network.SetCellSimple(i+1, 3, ip.ReverseDNS)
	}

	v.network.SetInputCapture(v.nextInputHandler("ip data", v.app.productList.treeView))
}

func (v *VpsInfo) createOverview(vps transipvps.Vps) {
	v.overview = tview.NewTable().SetSelectable(false, false)
	v.overview.SetTitle("Overview").SetBorder(true)

	v.overview.SetCellSimple(0, 0, "Name").
		SetCellSimple(0, 1, vps.Name)

	v.overview.SetCellSimple(1, 0, "Description").
		SetCellSimple(1, 1, vps.Description)

	v.overview.SetCellSimple(2, 0, "Product").
		SetCellSimple(2, 1, vps.ProductName)

	v.overview.SetCellSimple(3, 0, "Availability zone").
		SetCellSimple(3, 1, vps.AvailabilityZone)

	v.overview.SetCellSimple(4, 0, "CPUs").
		SetCellSimple(4, 1, strconv.Itoa(vps.CPUs))

	v.overview.SetCellSimple(5, 0, "Disk size").
		SetCellSimple(5, 1, fmt.Sprintf("%dG", vps.DiskSize/1024/1024))

	v.overview.SetCellSimple(6, 0, "Memory").
		SetCellSimple(6, 1, fmt.Sprintf("%dG", vps.MemorySize/1024/1024))

	v.overview.SetInputCapture(v.nextInputHandler("overview", v.backups))
}

func (v *VpsInfo) createBackups(vpsName string) {
	v.backups = tview.NewTable().SetSelectable(false, false)
	v.backups.SetTitle("Backups").SetBorder(true)

	v.backups.SetCellSimple(0, 0, "Date").
		SetCellSimple(0, 1, "Zone").
		SetCellSimple(0, 2, "Size").
		SetCellSimple(0, 3, "Status")

	v.app.Logger.Debug("fetching backup data", zap.String("name", vpsName))
	start := time.Now()
	backups, err := v.app.vpsRepo.GetBackups(vpsName)
	if err != nil {
		v.app.Logger.Error("error fetching backup data", zap.String("name", vpsName), zap.Error(err), zap.Duration("duration", time.Since(start)))
		panic(err)
	}

	v.app.Logger.Debug("done fetching backup data", zap.String("name", vpsName), zap.Duration("duration", time.Since(start)))

	for i, b := range backups {
		v.backups.
			SetCellSimple(i+1, 0, b.DateTimeCreate.In(time.Local).Format("Jan 02 15:04:05")).
			SetCellSimple(i+1, 1, b.AvailabilityZone).
			SetCellSimple(i+1, 2, fmt.Sprintf("%dG", b.DiskSize/1024/1024)).
			SetCellSimple(i+1, 3, string(b.Status))
	}

	v.backups.SetInputCapture(v.nextInputHandler("vps backups", v.snapshots))
}

func (v *VpsInfo) createSnapshots(vpsName string) {
	v.snapshots = tview.NewTable().SetSelectable(false, false)
	v.snapshots.SetTitle("Snapshots").SetBorder(true)

	v.snapshots.SetCellSimple(0, 0, "Date").
		SetCellSimple(0, 1, "Description").
		SetCellSimple(0, 2, "Size").
		SetCellSimple(0, 3, "Status")

	v.app.Logger.Debug("fetching snapshot data", zap.String("name", vpsName))
	start := time.Now()
	snapshots, err := v.app.vpsRepo.GetSnapshots(vpsName)
	if err != nil {
		v.app.Logger.Error("error fetching snapshot data", zap.String("name", vpsName), zap.Error(err), zap.Duration("duration", time.Since(start)))
		panic(err)
	}

	v.app.Logger.Debug("done fetching snapshot data", zap.String("name", vpsName), zap.Duration("duration", time.Since(start)))

	for i, s := range snapshots {
		v.snapshots.
			SetCellSimple(i+1, 0, s.DateTimeCreate.In(time.Local).Format("Jan 02 15:04:05")).
			SetCellSimple(i+1, 1, s.Description).
			SetCellSimple(i+1, 2, fmt.Sprintf("%dG", s.DiskSize/1024/1024)).
			SetCellSimple(i+1, 3, string(s.Status))
	}

	v.snapshots.SetInputCapture(v.nextInputHandler("vps snapshots", v.network))
}

func (v *VpsInfo) nextInputHandler(widgetName string, next tview.Primitive) func(event *tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		v.app.Logger.Debug("key pressed", zap.String("widget", widgetName), zap.Int32("rune", event.Rune()), zap.Int16("key", int16(event.Key())))

		switch event.Key() {
		case tcell.KeyTAB:
			v.app.Logger.Debug("detected TAB press", zap.String("widget", widgetName))
			v.app.Logger.Debug("changing focus")
			v.app.tviewApp.SetFocus(next)
		}
		return event
	}
}

func (v *VpsInfo) Focus() {
	v.app.tviewApp.SetFocus(v.overview)
}
