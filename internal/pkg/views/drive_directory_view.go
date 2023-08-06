package views

import (
	"github.com/anjolaoluwaakindipe/testcli/internal/pkg/services"
	"github.com/pterm/pterm"
)

type DriveDirectoryArgs struct {
	Directory string
}

type DriveDirectoryView struct {
	dirveDirectoryArgs DriveDirectoryArgs
}

func (ddv *DriveDirectoryView) View() {
	directory, _ := services.InitDriveService().GetDirectoryList("")
	result, _ := pterm.DefaultInteractiveSelect.
		WithOptions(directory).
		Show()
	pterm.Info.Printfln("You answered: %s, %s", result, ddv.dirveDirectoryArgs.Directory)
}

func InitDriveDirectoryView(args DriveDirectoryArgs) *DriveDirectoryView {
	return &DriveDirectoryView{args}

}
