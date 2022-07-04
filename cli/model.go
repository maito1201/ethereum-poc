package main

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/maito1201/ethrium-poc/cli/coin"
)

type model struct {
	list       list.Model
	items      []item
	choice     string
	quitting   bool
	progress   int
	account    string
	execType   string
	coinClient *coin.Client
}

func newModel() model {
	items := []list.Item{
		item(user1PrvKey),
		item(user2PrvKey),
	}

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Select your Account"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	return model{list: l}
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
				m.choice = string(i)
				switch m.progress {
				case 0:
					m.execute0()
				case 1:
					m.execute1()
				case 2:
					m.executeCoinTransaction()
					return m, tea.Quit
				}
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.quitting {
		return quitTextStyle.Render("Bye.")
	}
	return "\n" + m.list.View()
}

func (m *model) execute0() {
	m.progress = m.progress + 1
	m.account = m.choice
	items := []list.Item{
		item("NFT"),
		item("Coin"),
	}
	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	l.Title = "Select Token Type"
	m.list = l
}

func (m *model) execute1() {
	m.progress = m.progress + 1
	m.execType = m.choice
	switch m.execType {
	case "Coin":
		m.coinClient = initCoinClient(m.account)
		items := []list.Item{
			item("Mint"),
			item("Transfer"),
		}
		l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
		l.SetShowStatusBar(false)
		l.SetFilteringEnabled(false)
		l.Styles.Title = titleStyle
		l.Styles.PaginationStyle = paginationStyle
		l.Styles.HelpStyle = helpStyle
		l.Title = "Select Function"
		m.list = l
	case "NFT":
		panic("not implemented")
	}
}

func (m *model) executeCoinTransaction() {
	m.progress = m.progress + 1
	switch m.choice {
	case "Mint":
		mintCoin(m.coinClient)
	case "Transfer":
		transferCoin(m.coinClient)
	}
}
