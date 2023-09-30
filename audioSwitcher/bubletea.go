package audioSwitcher

import (
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/samber/lo"
)

const listHeight = 14

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

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
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type item string

func (i item) FilterValue() string { return "" }

type ChoiceDevice struct {
	Index int
	Name  string
}

type model struct {
	list     list.Model
	choice   ChoiceDevice
	quitting bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				selectDevice, _, _ := lo.FindIndexOf(deviceList, func(device map[string]interface{}) bool {
					return device["Name"].(string) == string(i)
				})
				if selectDevice == nil {
					return m, nil
				}
				m.choice = ChoiceDevice{
					Index: selectDevice["Index"].(int),
					Name:  selectDevice["Name"].(string),
				}
			}
			return m, tea.Quit
		default:
			atoi, err := strconv.Atoi(keypress)
			if err == nil {
				device, find := lo.Find(deviceList, func(device map[string]interface{}) bool {
					return device["Index"].(int) == atoi
				})
				if find {
					m.choice = ChoiceDevice{
						Index: device["Index"].(int),
						Name:  device["Name"].(string),
					}
				}
				return m, nil
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.choice.Name != "" {
		err := SetPlaybackDevices(m.choice.Index)
		if err != nil {
			return quitTextStyle.Render(fmt.Sprintf("%s 으로 변경 실패", m.choice.Name))
		}
		SelectedDeviceData = append(SelectedDeviceData, SelectedDevice{
			Name:  m.choice.Name,
			Count: 1,
		})

		err = saveData()
		if err != nil {
			log.Fatalln("설정 저장 실패")
		}
		return quitTextStyle.Render(fmt.Sprintf("%s 으로 변경이 완료 되었습니다", m.choice.Name))
	}
	return "\n" + m.list.View()
}
