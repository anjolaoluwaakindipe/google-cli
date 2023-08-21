package views

import (
	"log"

	"github.com/anjolaoluwaakindipe/testcli/internal/pkg/components"
	"github.com/anjolaoluwaakindipe/testcli/internal/pkg/services"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type DriveDirectoryArgs struct {
	Directory string
}

type DriveDirectoryView struct {
	driveDirectoryArgs DriveDirectoryArgs
	folderItemList     []list.Item
	list               list.Model
	currentDirectory   string
}

func (ddv *DriveDirectoryView) Init() tea.Cmd {
	// Google drive folder structure
	if ddv.driveDirectoryArgs.Directory == "" {
		ddv.currentDirectory = "/"
	}
	ddv.folderItemList = make([]list.Item, 0)
	files, err := services.InitDriveService().GetDirectoryList("")
	if err != nil {
		log.Fatal("Error in getting folder information: ", err)
	}
	for _, val := range files {
		ddv.folderItemList = append(ddv.folderItemList, components.NewListItem(val.Name, val.Size, val.MimeType, val.Id, val.DocumentType))
	}

	// List configuration
	var (
		defaultWidth = 20
		listHeight   = 25
		keyBindings  = newListKeyBindings()
	)
	ddv.list = list.New(ddv.folderItemList, components.ItemDelegate{}, defaultWidth, listHeight)
	ddv.list.Title = "Current Directory = " + ddv.currentDirectory
	ddv.list.SetShowStatusBar(false)
	ddv.list.Styles.Title = components.TitleStyle
	ddv.list.Styles.PaginationStyle = components.PaginationStyle
	ddv.list.Styles.HelpStyle = components.HelpStyle
	ddv.list.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			keyBindings.downloadFileOrFolder,
			keyBindings.navigateToFolder,
		}
	}

	return nil
}

func (ddv *DriveDirectoryView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		ddv.list.SetWidth(msg.Width)
		return ddv, nil
	case tea.KeyMsg:
		if ddv.list.FilterState() == list.Filtering {
			break
		}
		switch msg.Type {
		case tea.KeyEnter:

		}
	}

	var cmd tea.Cmd
	ddv.list, cmd = ddv.list.Update(msg)
	cmds = append(cmds, cmd)

	return ddv, tea.Batch(cmds...)
}

func (ddv *DriveDirectoryView) View() string {
	// s += "\n\n"
	// s += "Press Enter to navigate \n"
	// s += "Press D to download\n"
	// s += "Press up or down arrow to scroll\n"

	return ddv.list.View()
}

func InitDriveDirectoryView(args DriveDirectoryArgs) *DriveDirectoryView {
	return &DriveDirectoryView{driveDirectoryArgs: args}
}

// Key bindings
type listKeyBindings struct {
	downloadFileOrFolder key.Binding
	navigateToFolder     key.Binding
}

func newListKeyBindings() *listKeyBindings {
	return &listKeyBindings{
		downloadFileOrFolder: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "download file/folder"),
		),
		navigateToFolder: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "navigate to folder"),
		),
	}
}
