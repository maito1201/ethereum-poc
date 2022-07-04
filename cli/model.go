package main

import (
	"math/big"
	"strconv"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/maito1201/ethrium-poc/cli/coin"
	"github.com/maito1201/ethrium-poc/cli/nft"
)

type model struct {
	list            list.Model
	items           []item
	choice          string
	quitting        bool
	progress        int
	account         string
	execType        string
	transactionType string
	coinClient      *coin.Client
	nftClient       *nft.Client
}

func newModel() model {
	items := []list.Item{
		item(user1PrvKey),
		item(user2PrvKey),
	}
	l := newList("Select Account", items)
	return model{list: *l}
}

func newList(title string, items []list.Item) *list.Model {
	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = title
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	return &l
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
					switch m.execType {
					case "Coin":
						m.executeCoinTransaction()
						return m, tea.Quit
					case "NFT":
						m.executeNftTransaction()
						if m.transactionType == "Mint" {
							return m, tea.Quit
						}
					}
				case 3:
					m.executeTransferNft()
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
	l := newList("Select Token Type", items)
	m.list = *l
}

func (m *model) execute1() {
	m.progress = m.progress + 1
	m.execType = m.choice
	switch m.execType {
	case "Coin":
		m.coinClient = initCoinClient(m.account)
	case "NFT":
		m.nftClient = initNftClient(m.account)
	default:
		panic("not implemented error")
	}
	items := []list.Item{
		item("Mint"),
		item("TransferToFriend"),
		item("TransferToMe"),
	}
	l := newList("Select Function", items)
	m.list = *l
}

func (m *model) executeCoinTransaction() {
	m.progress = m.progress + 1
	switch m.choice {
	case "Mint":
		mintCoin(m.coinClient)
	case "TransferToFrined":
		transferCoinToFrined(m.coinClient)
	case "TransferToMe":
		transferCoinToMe(m.coinClient)
	default:
		panic("not implemented error")
	}
}

func (m *model) executeNftTransaction() {
	m.progress = m.progress + 1
	m.transactionType = m.choice
	switch m.choice {
	case "Mint":
		mintNft(m.nftClient)
	case "TransferToFriend", "TransferToMe":
		items := []list.Item{
			item("0"),
			item("1"),
			item("2"),
			item("3"),
			item("4"),
			item("5"),
		}
		l := newList("Select TokenID", items)
		m.list = *l
	default:
		panic("not implemented error")
	}
}

func (m *model) executeTransferNft() {
	m.progress = m.progress + 1
	id, err := strconv.Atoi(m.choice)
	if err != nil {
		panic(err)
	}
	tokenID := big.NewInt(int64(id))
	switch m.transactionType {
	case "TransferToFriend":
		transferNftToFrined(m.nftClient, tokenID)
	case "TransferToMe":
		transferNftToMe(m.nftClient, tokenID)
	default:
		panic("not implemented error")
	}
}
