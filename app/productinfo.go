package app

import (
	"fmt"
	"github.com/rivo/tview"
	"go.uber.org/zap"
	"strconv"
	"time"
)


type Focuser interface {
	Focus()
}

type ProductInfo struct {
	app *Application

	grid *tview.Flex

	currentView Focuser
}

func (pi *ProductInfo) init() {
	pi.grid = tview.NewFlex()
	pi.grid.SetDirection(tview.FlexRow)
}

func (pi *ProductInfo) ShowVpsFunc(vpsName string) func() {

	return func() {
		info := &VpsInfo{app: pi.app}
		pi.currentView = info

		info.ShowVps(pi.grid, vpsName)
	}
}

type VpsInfo struct {
	app *Application
	overview *tview.Table
	network *tview.Table
}

func (v *VpsInfo) ShowVps(grid *tview.Flex, vpsName string) {
	grid.Clear()

	v.app.Logger.Debug("fetching vps", zap.String("name", vpsName))
	now := time.Now()
	vps, err := v.app.vpsRepo.GetByName(vpsName)
	if err != nil {
		v.app.Logger.Error("error fetching vps", zap.Error(err), zap.String("name", vpsName), zap.Duration("duration", time.Since(now)))
		panic(err)
	}

	v.app.Logger.Debug("done fetching vps", zap.String("name", vpsName), zap.Duration("duration", time.Since(now)))

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
		SetCellSimple(5, 1, fmt.Sprintf("%dG", vps.DiskSize / 1024 / 1024))

	v.overview.SetCellSimple(6, 0, "Memory").
		SetCellSimple(6, 1, fmt.Sprintf("%dG", vps.MemorySize / 1024 / 1024))

	grid.AddItem(v.overview, 0, 1, false)

	v.network = tview.NewTable().SetSelectable(false, false)
	v.network.SetTitle("Network").SetBorder(true)

	v.network.SetCellSimple(0, 0, "IP").
		SetCellSimple(0, 1, "Subnet").
		SetCellSimple(0, 2, "Gateway").
		SetCellSimple(0, 3, "Reverse DNS")

	v.app.Logger.Debug("fetching ip data", zap.String("name", vpsName))
	now = time.Now()
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

	grid.AddItem(v.network, 0, 2, false)
}

func (v *VpsInfo) Focus() {
	v.app.tviewApp.SetFocus(v.network)
}
