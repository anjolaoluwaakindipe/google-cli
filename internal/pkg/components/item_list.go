package components

import (
	"fmt"
	"io"
	"strings"

	"github.com/anjolaoluwaakindipe/testcli/internal/pkg/utils"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	TitleStyle        = lipgloss.NewStyle().MarginLeft(2)
	ItemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	SelectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	FolderStyle       = lipgloss.NewStyle().PaddingLeft(4).Foreground(lipgloss.Color("#0000ff"))
	PaginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	HelpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	QuitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type Item struct {
	title        string
	size         int
	mimeType     string
	documentType utils.DocumentType
	id           string
}

func (i Item) Title() string                    { return i.title }
func (i Item) Size() int                        { return i.size }
func (i Item) MimeType() string                 { return i.mimeType }
func (i Item) DocumentType() utils.DocumentType { return i.documentType }
func (i Item) Id() string                       { return i.id }
func (i Item) IsFolderOrShortcut() bool {
	return i.DocumentType() == utils.Folder || i.DocumentType() == utils.Shortcut
}

func (i Item) FilterValue() string { return i.title }

func NewListItem(title string, size int, mimeType string, id string, documentType utils.DocumentType) Item {
	return Item{title: title, size: size, mimeType: mimeType, id: id, documentType: documentType}
}

type ItemDelegate struct{}

func (d ItemDelegate) Height() int                             { return 1 }
func (d ItemDelegate) Spacing() int                            { return 0 }
func (d ItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d ItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(Item)
	if !ok {
		return
	}
	sizeStr := ""
	if i.size > 0 {
		sizeStr = fmt.Sprintf("(%v)", utils.FormatBytes(i.size))
	}
	str := fmt.Sprintf("%d. %s", index+1, i.title) + sizeStr

	fn := ItemStyle.Render
	if i.documentType == utils.Folder || i.documentType == utils.Shortcut {
		fn = FolderStyle.Render
	}
	if index == m.Index() {
		fn = func(s ...string) string {
			return SelectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}
