package views

import (
	"log"
	"strings"

	"github.com/anjolaoluwaakindipe/testcli/internal/pkg/components"
	"github.com/anjolaoluwaakindipe/testcli/internal/pkg/services"
	"github.com/anjolaoluwaakindipe/testcli/internal/pkg/utils"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type DriveDirectoryArgs struct {
	Directory string
}

type DriveDirectoryView struct {
	driveDirectoryArgs DriveDirectoryArgs
	keys               *listKeyBindings
	list               list.Model
	currentDirectory   []string
	folderQueue        []string
}

func (ddv *DriveDirectoryView) Init() tea.Cmd {
	// Google drive folder structure
	if ddv.driveDirectoryArgs.Directory == "" {
		ddv.currentDirectory = append(ddv.currentDirectory, "/")
	}

	folderItems := ddv.GetFolderItemList("")

	// List configuration
	var (
		defaultWidth = 20
		listHeight   = 25
		keyBindings  = newListKeyBindings()
	)
	ddv.list = list.New(folderItems, components.ItemDelegate{}, defaultWidth, listHeight)
	ddv.list.Title = "Current Directory = " + strings.Join(ddv.currentDirectory, "")
	ddv.list.SetShowStatusBar(false)
	ddv.list.Styles.Title = components.TitleStyle
	ddv.list.Styles.PaginationStyle = components.PaginationStyle
	ddv.list.Styles.HelpStyle = components.HelpStyle
	ddv.list.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			keyBindings.downloadFileOrFolder,
			keyBindings.navigateToFolder,
			keyBindings.goBackToPreviousFolder,
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
		switch {
		case key.Matches(msg, ddv.keys.navigateToFolder):
			val, ok := ddv.list.SelectedItem().(components.Item)
			if !ok || !val.IsFolderOrShortcut() {
				break
			}
			ddv.currentDirectory = append(ddv.currentDirectory, val.Title())
			ddv.UpdateListTitle()
			ddv.folderQueue = append(ddv.folderQueue, val.Id())
			return ddv, ddv.SendMessage(
				utils.FolderNavigationMsg{FolderId: val.Id(), Name: val.Title()},
			)
		case key.Matches(msg, ddv.keys.goBackToPreviousFolder):
			if len(ddv.folderQueue) == 0 {
				break
			}
			ddv.folderQueue = ddv.folderQueue[:len(ddv.folderQueue)-1]
			var folderItems []list.Item
			if len(ddv.folderQueue) > 0 {
				folderItems = ddv.GetFolderItemList(ddv.folderQueue[len(ddv.folderQueue)-1])
			} else {
				folderItems = ddv.GetFolderItemList("")

			}
			if len(ddv.currentDirectory) > 1 {
				ddv.currentDirectory = ddv.currentDirectory[:len(ddv.currentDirectory)-1]
				ddv.UpdateListTitle()
			}
			cmd := ddv.list.SetItems(folderItems)
			return ddv, cmd
		}
	case utils.FolderNavigationMsg:
		folderItems := ddv.GetFolderItemList(msg.FolderId)
		cmd := ddv.list.SetItems(folderItems)
		return ddv, cmd
	}

	var cmd tea.Cmd
	ddv.list, cmd = ddv.list.Update(msg)
	cmds = append(cmds, cmd)

	return ddv, tea.Batch(cmds...)
}

func (ddv *DriveDirectoryView) SendMessage(msg tea.Msg) tea.Cmd {
	return func() tea.Msg {
		return msg
	}
}

func (ddv *DriveDirectoryView) GetFolderItemList(folderId string) []list.Item {
	folderItems := make([]list.Item, 0)
	files, err := services.InitDriveService().GetDirectoryList(folderId)
	if err != nil {
		log.Fatal("Error in getting folder information: ", err)
	}
	for _, val := range files {
		folderItems = append(folderItems, components.NewListItem(val.Name, val.Size, val.MimeType, val.Id, val.DocumentType))
	}
	
	return folderItems
}

func (ddv *DriveDirectoryView) View() string {
	// s += "\n\n"
	// s += "Press Enter to navigate \n"
	// s += "Press D to download\n"
	// s += "Press up or down arrow to scroll\n"

	return ddv.list.View()
}

func (ddv *DriveDirectoryView) UpdateListTitle() {
	ddv.list.Title = "Current Directory: " + strings.Join(ddv.currentDirectory, "")
}

func InitDriveDirectoryView(args DriveDirectoryArgs) *DriveDirectoryView {
	return &DriveDirectoryView{driveDirectoryArgs: args, keys: newListKeyBindings()}
}

// Key bindings
type listKeyBindings struct {
	downloadFileOrFolder   key.Binding
	navigateToFolder       key.Binding
	goBackToPreviousFolder key.Binding
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
		goBackToPreviousFolder: key.NewBinding(
			key.WithKeys("b"),
			key.WithHelp("b", "back to previous folder"),
		),
	}
}
