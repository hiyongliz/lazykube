package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const listHeight = 14

var (
    titleStyle        = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
    itemStyle         = lipgloss.NewStyle().PaddingLeft(2)
    selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
    paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
    helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
    quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
    borderStyle       = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#404040"))
	selectedBorderStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("62"))
)

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
    i, ok := listItem.(item)
    if !ok {
        return
    }

    str := fmt.Sprintf("%d. %s", index+1, i)

    fn := itemStyle.Render
    if index == m.Index() {
        fn = func(s ...string) string {
            return selectedItemStyle.Render(strings.Join(s, " "))
        }
    }

    fmt.Fprint(w, fn(str))
}

type model struct {
    list       list.Model
    list2      list.Model
    choice     string
    quitting   bool
    activeList int // 0 for list, 1 for list2
}

func (m model) Init() tea.Cmd {
    return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.list.SetWidth(msg.Width / 2)
        m.list2.SetWidth(msg.Width / 2)
		m.list.SetHeight(msg.Height - 2)
		m.list2.SetHeight(msg.Height - 2)
        return m, nil

    case tea.KeyMsg:
        switch keypress := msg.String(); keypress {
        case "q", "ctrl+c":
            m.quitting = true
            return m, tea.Quit

        case "enter":
            if m.activeList == 0 {
                i, ok := m.list.SelectedItem().(item)
                if ok {
                    m.choice = string(i)
                }
            } else {
                i, ok := m.list2.SelectedItem().(item)
                if ok {
                    m.choice = string(i)
                }
            }
            return m, tea.Quit

        case "tab":
            m.activeList = (m.activeList + 1) % 2
        }
    }

    var cmd tea.Cmd
    if m.activeList == 0 {
        m.list, cmd = m.list.Update(msg)
    } else {
        m.list2, cmd = m.list2.Update(msg)
    }
    return m, cmd
}

func (m model) View() string {
    if m.choice != "" {
        return quitTextStyle.Render(fmt.Sprintf("%s? Sounds good to me.", m.choice))
    }
    if m.quitting {
        return quitTextStyle.Render("Not hungry? That’s cool.")
    }

    var listView, list2View string
    if m.activeList == 0 {
        listView = selectedBorderStyle.Render(m.list.View())
        list2View = borderStyle.Render(m.list2.View())
    } else {
        list2View = selectedBorderStyle.Render(m.list2.View())
        listView = borderStyle.Render(m.list.View())
    }

    return "\n" + lipgloss.JoinHorizontal(lipgloss.Top, listView, list2View)
}

func main() {
    items := []list.Item{
        item("Ramen"),
        item("Tomato Soup"),
        item("Hamburgers"),
        item("Cheeseburgers"),
        item("Currywurst"),
        item("Okonomiyaki"),
        item("Pasta"),
        item("Fillet Mignon"),
        item("Caviar"),
        item("Just Wine"),
    }

    items2 := []list.Item{
        item("Water"),
        item("Soda"),
        item("Tea"),
        item("Coffee"),
        item("Juice"),
        item("Milk"),
        item("Beer"),
        item("Wine"),
        item("Whiskey"),
        item("Cocktail"),
    }

    const defaultWidth = 20

    l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
    l.Title = "What do you want for dinner?"
    l.SetShowStatusBar(false)
    l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
    l.Styles.Title = titleStyle
    l.Styles.PaginationStyle = paginationStyle
    l.Styles.HelpStyle = helpStyle

    l2 := list.New(items2, itemDelegate{}, defaultWidth, listHeight)
    l2.Title = "What do you want to drink?"
    l2.SetShowStatusBar(false)
    l2.SetFilteringEnabled(false)
	l2.SetShowHelp(false)
    l2.Styles.Title = titleStyle
    l2.Styles.PaginationStyle = paginationStyle
    l2.Styles.HelpStyle = helpStyle

    m := model{list: l, list2: l2}

    if _, err := tea.NewProgram(m).Run(); err != nil {
        fmt.Println("Error running program:", err)
        os.Exit(1)
    }
}
