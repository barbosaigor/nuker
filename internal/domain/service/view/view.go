package view

import (
	"fmt"
	"time"

	"github.com/barbosaigor/nuker/pkg/metrics"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sirupsen/logrus"
)

type View interface {
	SetMetric(mr *metrics.MetricRate)
	ShutDown()
}

type view struct {
	mr *metrics.MetricRate
	p  *tea.Program
}

func New() (View, error) {
	vw := &view{
		mr: &metrics.MetricRate{},
	}

	p := tea.NewProgram(vw)
	vw.p = p
	go func() {
		if err := p.Start(); err != nil {
			logrus.Fatal(err)
		}
	}()

	return vw, nil
}

func (vw *view) SetMetric(mr *metrics.MetricRate) {
	// logrus.Error("Setting metric: %v", mr)
	*vw.mr = *mr
}

type tickMsg struct{}

func tick() tea.Msg {
	time.Sleep(500 * time.Millisecond)
	return tickMsg{}
}

func (view) Init() tea.Cmd {
	return tick
}

func (vw *view) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		return vw, tea.Quit
	case tickMsg:
		return vw, tick
	default:
		return vw, nil
	}
}

func (vw view) View() string {
	if vw.mr == nil {
		return fmt.Sprintf("%v\n", vw.mr)
	}
	return fmt.Sprintf("success: %d/%d\n", vw.mr.Success, vw.mr.Total)
}

func (vw view) ShutDown() {
	vw.p.Quit()
}
