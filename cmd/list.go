package cmd

import (
	"fmt"
	"os"

	"github.com/anjolaoluwaakindipe/testcli/internal/pkg/views"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "get the lists of files and folder in a specific directory",
	Run: func(cmd *cobra.Command, args []string) {

		var viewArg views.DriveDirectoryArgs
		if len(args) > 0 && args[0] != "" {
			viewArg = views.DriveDirectoryArgs{Directory: args[0]}
		}
		view := views.InitDriveDirectoryView(viewArg)
		if _, err := tea.NewProgram(view).Run(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}
	},
}
