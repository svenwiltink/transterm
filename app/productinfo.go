package app

import (
	"fmt"
	"github.com/rivo/tview"
	"strconv"
)

type ProductInfo struct {
	app *Application

	grid *tview.Flex
}

func (pi *ProductInfo) init() {
	pi.grid = tview.NewFlex()
	pi.grid.SetDirection(tview.FlexRow)
}

func (pi *ProductInfo) ShowVpsFunc(vpsName string) func() {

	return func() {
		pi.grid.Clear()

		vps, err := pi.app.vpsRepo.GetByName(vpsName)
		if err != nil {
			panic(err)
		}

		info := tview.NewTable()
		info.SetTitle("Overview").SetBorder(true)

		info.SetSelectable(false, false)

		info.SetCellSimple(0, 0, "Name").
			SetCellSimple(0, 1, vps.Name)

		info.SetCellSimple(1, 0, "Description").
			SetCellSimple(1, 1, vps.Description)

		info.SetCellSimple(2, 0, "Product").
			SetCellSimple(2, 1, vps.ProductName)

		info.SetCellSimple(3, 0, "Availability zone").
			SetCellSimple(3, 1, vps.AvailabilityZone)

		info.SetCellSimple(4, 0, "CPUs").
			SetCellSimple(4, 1, strconv.Itoa(vps.CPUs))

		info.SetCellSimple(5, 0, "Disk size").
			SetCellSimple(5, 1, fmt.Sprintf("%dG", vps.DiskSize / 1024 / 1024))

		info.SetCellSimple(6, 0, "Memory").
			SetCellSimple(6, 1, fmt.Sprintf("%dG", vps.MemorySize / 1024 / 1024))

		pi.grid.AddItem(info, 0, 1, false)

		network := tview.NewTable()
		network.SetTitle("Network").SetBorder(true)

		network.SetCellSimple(0, 0, "IP").
			SetCellSimple(0, 1, "Subnet").
			SetCellSimple(0, 2, "Gateway").
			SetCellSimple(0, 3, "Reverse DNS")

		ips, err := pi.app.vpsRepo.GetIPAddresses(vpsName)
		if err != nil {
			panic(err)
		}

		for i, ip := range ips {
			network.SetCellSimple(i+1, 0, ip.Address.String())
			text, _ := ip.SubnetMask.MarshalText()
			network.SetCellSimple(i+1, 1, string(text))
			network.SetCellSimple(i+1, 2, ip.Gateway.String())
			network.SetCellSimple(i+1, 3, ip.ReverseDNS)
		}

		pi.grid.AddItem(network, 0, 2, false)
	}
}
