package ui

import (
	"time"
	"vkpdeveloper/statuspages-cli/utils"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type pageItem struct {
	title     string
	pageId    string
	subdomain string
}

type refreshComponents int

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func (i pageItem) Title() string       { return i.title }
func (i pageItem) Description() string { return i.subdomain }
func (i pageItem) FilterValue() string { return i.title }

type ViewModel struct {
	needAuthentication bool
	textInput          textinput.Model
	pagesList          list.Model
	config             *utils.AppConfig
	client             *utils.StatusPageClient
	selectedPage       *pageItem
	componentsTable    table.Model
	updateChan         chan int
}

var docStyle = lipgloss.NewStyle().Margin(1, 2)

func initNewViewModel(config *utils.AppConfig, client *utils.StatusPageClient) *ViewModel {

	var model = &ViewModel{
		client:             client,
		config:             config,
		needAuthentication: false,
		pagesList:          list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0),
		componentsTable:    table.New(),
		updateChan:         make(chan int),
	}

	if config.APIKey == "" {

		apiKeyTextModel := textinput.New()
		apiKeyTextModel.Placeholder = "?What's your API Key"
		apiKeyTextModel.Focus()
		apiKeyTextModel.CharLimit = 40

		model.textInput = apiKeyTextModel
		model.needAuthentication = true

		return model
	}

	return model
}

func (m ViewModel) Init() tea.Cmd {
	return nil
}

func (m ViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	m.textInput, _ = m.textInput.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.needAuthentication {
				apiKey := m.textInput.Value()
				m.config.WriteConfig(apiKey)
				m.config.ReadConfig()
				m.needAuthentication = false
				return m, tea.Quit
			} else {
				if m.selectedPage == nil {
					selectedPage := m.pagesList.SelectedItem().(pageItem)
					m.selectedPage = &selectedPage

					go func() {
						for {
							time.Sleep(10 * time.Second)
							cli.Send(refreshComponents(1))
						}
					}()
				}
			}
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.pagesList.SetSize(msg.Width-h, msg.Height-v)
	case refreshComponents:
		components, err := m.client.GetPageComponents(m.selectedPage.pageId)

		if err != nil {
			panic(err)
		}

		m.componentsTable.SetColumns([]table.Column{
			{
				Title: "Name",
				Width: 50,
			},
			{
				Title: "Status",
				Width: 20,
			},
		})

		rows := []table.Row{}

		for _, component := range *components {
			rows = append(rows, table.Row{component.Name, component.Status})
		}

		m.componentsTable.SetRows([]table.Row{})
		m.componentsTable.SetRows(rows)

		return m, cmd
	case error:
		return m, nil
	}

	if !m.needAuthentication && m.selectedPage == nil {
		pages, err := m.client.GetPages()

		if err != nil {
			panic(err)
		}

		var items = []list.Item{}

		for _, page := range *pages {
			items = append(items, pageItem{
				title:     page.Name,
				pageId:    page.Id,
				subdomain: page.Subdomain,
			})
		}

		m.pagesList.SetItems(items)
		m.pagesList.Title = "Select status page"
	}

	if !m.needAuthentication && m.selectedPage != nil {
		components, err := m.client.GetPageComponents(m.selectedPage.pageId)

		if err != nil {
			panic(err)
		}

		m.componentsTable.SetColumns([]table.Column{
			{
				Title: "Name",
				Width: 50,
			},
			{
				Title: "Status",
				Width: 20,
			},
		})

		rows := []table.Row{}

		for _, component := range *components {
			rows = append(rows, table.Row{component.Name, component.Status})
		}

		m.componentsTable.SetRows(rows)
	}

	m.pagesList, _ = m.pagesList.Update(msg)
	m.componentsTable, cmd = m.componentsTable.Update(msg)

	return m, cmd
}

func (m ViewModel) View() string {
	if m.needAuthentication {
		return m.textInput.View()
	} else {
		if m.selectedPage == nil {
			return docStyle.Render(m.pagesList.View())
		} else {
			return baseStyle.Render(m.componentsTable.View())
		}
	}
}

var cli *tea.Program

func Run(config *utils.AppConfig, apiClient *utils.StatusPageClient) {
	cli = tea.NewProgram(initNewViewModel(config, apiClient), tea.WithAltScreen())
	if _, err := cli.Run(); err != nil {
		panic(err)
	}
}
