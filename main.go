package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/duolok/blue-jay/engine"
	"github.com/lucasb-eyer/go-colorful"
)

const (
	progressBarWidth  = 71
	progressFullChar  = "█"
	progressEmptyChar = "░"
	dotChar           = " • "
)

var (
	keywordStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	subtleStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	ticksStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("79"))
	checkboxStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	selectedLineStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("218")).Bold(true) 
	progressEmpty    = subtleStyle.Render(progressEmptyChar)
	dotStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("236")).Render(dotChar)
	mainStyle        = lipgloss.NewStyle().MarginLeft(2)

	ramp = makeRampStyles("#B14FFF", "#00FFA3", progressBarWidth)
)

type model struct {
	Choice    int
	Chosen    bool
	Ticks     int
	Frames    int
	Progress  float64
	Loaded    bool
	Quitting  bool
	TextInput textinput.Model
	Err       error
	Games     [][]string
	Cursor    int
	Searching bool
	ViewState string
	Paginator paginator.Model
}

type (
	tickMsg  struct{}
	frameMsg struct{}
)

func main() {
	initialModel := initModel()
	p := tea.NewProgram(initialModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
	}
}

func initModel() model {
	ti := textinput.New()
	ti.Placeholder = "Search for games"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	p := paginator.New()
	p.Type = paginator.Dots
	p.PerPage = 10
	p.ActiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "235", Dark: "252"}).Render("•")
	p.InactiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "250", Dark: "238"}).Render("•")

	return model{
		Choice:    0,
		Chosen:    false,
		Ticks:     10,
		Frames:    0,
		Progress:  0,
		Loaded:    false,
		Quitting:  false,
		TextInput: ti,
		Err:       nil,
		Games:     [][]string{},
		Cursor:    0,
		Searching: false,
		ViewState: "choices",
		Paginator: p,
	}
}

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

func frame() tea.Cmd {
	return tea.Tick(time.Second/60, func(time.Time) tea.Msg {
		return frameMsg{}
	})
}

func (m model) Init() tea.Cmd {
	return tick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "ctrl+c" {
			m.Quitting = true
			return m, tea.Quit
		}
	}

	switch m.ViewState {
	case "choices":
		return updateChoices(msg, m)
	case "search":
		return updateSearch(msg, m)
	case "results":
		return updateResults(msg, m)
	}
	return m, nil
}

func (m model) View() string {
	if m.Quitting {
		return "\n  See you later!\n\n"
	}

	switch m.ViewState {
	case "choices":
		return mainStyle.Render("\n" + choicesView(m) + "\n\n")
	case "search":
		return mainStyle.Render("\n" + searchView(m) + "\n\n")
	case "results":
		return mainStyle.Render("\n" + gameChoiceView(m) + "\n\n")
	}

	return ""
}

func updateChoices(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			m.Choice++
			if m.Choice > 3 {
				m.Choice = 3
			}
		case "k", "up":
			m.Choice--
			if m.Choice < 0 {
				m.Choice = 0
			}
		case "enter":
			if m.Choice == 0 {
				m.ViewState = "search"
				return m, nil
			}
			if m.Choice == 3 {
				m.Quitting = true
				return m, tea.Quit
			}
			m.Chosen = true
			return m, frame()
		}
	case tea.WindowSizeMsg:
		m.TextInput.Width = msg.Width
	}

	return m, nil
}

func updateSearch(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// Trigger game search
			searchTerm := m.TextInput.Value()
			searchGames(searchTerm)
			games, err := loadGames("games.csv")
			if err != nil {
				fmt.Println("Some error has occurred")
			}
			m.Games = games
			m.Paginator.SetTotalPages((len(m.Games) + m.Paginator.PerPage - 1) / m.Paginator.PerPage)
			m.ViewState = "results"
			return m, nil
		case "esc":
			m.ViewState = "choices"
			return m, nil
		}
	}
	m.TextInput, cmd = m.TextInput.Update(msg)
	return m, cmd
}

func updateResults(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "j", "down":
            m.Cursor++
            if m.Cursor >= len(m.Games) {
                m.Cursor = 0
            }
        case "k", "up":
            m.Cursor--
            if m.Cursor < 0 {
                m.Cursor = len(m.Games) - 1
            }
        case "enter":
            if len(m.Games) > 0 {
                openInBrowser(m.Games[m.Cursor][2])
                // Set the view back to choices instead of quitting
                m.ViewState = "choices"
                return m, nil
            }
        case "esc":
            m.ViewState = "search"
            return m, nil
        case "left", "h":
            m.Paginator.PrevPage()
        case "right", "l":
            m.Paginator.NextPage()
        }
    }

    return m, nil
}


func openInBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		fmt.Printf("unsupported platform")
		return
	}

	if err != nil {
		fmt.Printf("failed to open browser: %v\n", err)
	}
}

func choicesView(m model) string {
	c := m.Choice

	tpl := `
		  ____  _                      _             
		 | __ )| |_   _  ___          | | __ _ _   _ 
		 |  _ \| | | | |/ _ \_____ _  | |/ _  | | | |
		 | |_) | | |_| |  __/_____| |_| | (_| | |_| |
		 |____/|_|\__,_|\___|      \___/ \__,_|\__, |
											   |___/ 
		`

	tpl += "\n\n%s"
	tpl += "\n\n"
	tpl += subtleStyle.Render("j/k, up/down: select") + dotStyle +
		subtleStyle.Render("enter: choose") + dotStyle +
		subtleStyle.Render("q, esc: quit")

	choices := fmt.Sprintf(
		"%s\n%s\n%s\n%s",
		checkbox("Search For Games", c == 0),
		checkbox("Last Search", c == 1),
		checkbox("Settings", c == 2),
		checkbox("Quit", c == 3),
	)
	return fmt.Sprintf(tpl, choices)
}

func searchView(m model) string {
	return fmt.Sprintf(
		"Enter game name to search and press Enter:\n\n%s\n\n(esc to return)",
		m.TextInput.View(),
	)
}

func gameChoiceView(m model) string {
	var s strings.Builder

	s.WriteString("Select a game from the results:\n\n")

	start, end := m.Paginator.GetSliceBounds(len(m.Games))
	start = start + 1
	for i := start; i < end; i++ {
		if m.Games[i][1] != "" {
			if m.Cursor == i {
				s.WriteString(selectedLineStyle.Render("[x] " + m.Games[i][0] + " ━━━━━━━━━━ " + m.Games[i][1]))
			} else {
				s.WriteString(subtleStyle.Render("[ ] " + m.Games[i][0] + " ━━━━━━━━━━ " + m.Games[i][1]))
			}

			s.WriteString("\n")
		}
	}

	s.WriteString("\n" + m.Paginator.View() + "\n")
	s.WriteString(subtleStyle.Render("j/k, up/down: select") + dotStyle +
		subtleStyle.Render("h/l, left/right: page") + dotStyle +
		subtleStyle.Render("enter: choose"))

	return s.String()
}

func checkbox(label string, checked bool) string {
	if checked {
		return checkboxStyle.Render("[x] " + label)
	}
	return fmt.Sprintf("[ ] %s", label)
}

func progressbar(percent float64) string {
	w := float64(progressBarWidth)

	fullSize := int(math.Round(w * percent))
	var fullCells string
	for i := 0; i < fullSize; i++ {
		fullCells += ramp[i].Render(progressFullChar)
	}

	emptySize := int(w) - fullSize
	emptyCells := strings.Repeat(progressEmpty, emptySize)

	return fmt.Sprintf("%s%s %3.0f", fullCells, emptyCells, math.Round(percent*100))
}

func makeRampStyles(colorA, colorB string, steps float64) (s []lipgloss.Style) {
	cA, _ := colorful.Hex(colorA)
	cB, _ := colorful.Hex(colorB)

	for i := 0.0; i < steps; i++ {
		c := cA.BlendLuv(cB, i/steps)
		s = append(s, lipgloss.NewStyle().Foreground(lipgloss.Color(colorToHex(c))))
	}
	return
}

func colorToHex(c colorful.Color) string {
	return fmt.Sprintf("#%s%s%s", colorFloatToHex(c.R), colorFloatToHex(c.G), colorFloatToHex(c.B))
}

func colorFloatToHex(f float64) (s string) {
	s = strconv.FormatInt(int64(f*255), 16)
	if len(s) == 1 {
		s = "0" + s
	}
	return
}

func searchGames(query string) {
	scraperNames := engine.LoadScrapers()
	var wg sync.WaitGroup

	for _, scraperName := range scraperNames {
		wg.Add(1)
		go func(scraperName string) {
			defer wg.Done()
			engine.Search([]string{scraperName}, query, &wg)
		}(scraperName)
	}
	wg.Wait()
}

func loadGames(filename string) ([][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var games [][]string
	for _, record := range records {
		games = append(games, []string{record[0], record[1], record[2]})
	}
	return games, nil
}

