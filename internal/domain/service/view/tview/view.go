package tview

import (
	"fmt"
	"time"

	"github.com/barbosaigor/nuker/internal/domain/service/view"
	"github.com/barbosaigor/nuker/pkg/metrics"
	"github.com/rivo/tview"
	"github.com/sirupsen/logrus"
)

type tView struct {
	mr       *metrics.MetricRate
	app      *tview.Application
	textView *tview.TextView
	txt      string
	quit     chan bool
}

func New() (view.View, error) {
	app := tview.NewApplication()
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	textView.SetBorder(true)

	go func() {
		err := app.SetRoot(textView, true).EnableMouse(false).Run()
		if err != nil {
			panic(err)
		}
	}()

	fmt.Fprint(textView, "setting up pipeline...\n")

	vw := &tView{
		mr:       &metrics.MetricRate{},
		app:      app,
		textView: textView,
		quit:     make(chan bool),
	}

	go func() {
		for {
			select {
			case <-vw.quit:
				return
			case <-time.After(100 * time.Millisecond):
				if vw.txt == "" {
					continue
				}
				textView.Clear()
				fmt.Fprint(textView, vw.txt)
			}
		}
	}()

	return vw, nil
}

func (vw *tView) SetMetric(mr *metrics.MetricRate) {
	// logrus.Error("Setting metric: %v", mr)

	var txt string
	// vw.textView.Clear()
	if mr == nil {
		txt = fmt.Sprintf("%v\n", mr)
	} else {
		txt = fmt.Sprintf(
			// "min: %d, max: %d\n"+
			"iterations........: [::b]total=[green]%d[white]\tsuccess=[green]%d[white]\tfailed=[green]%d[white][::-]\t\n"+
				"success ratio.....: [::b][green]%f[white][::-]\n"+
				"request time......: [::b]"+
				"max=[green]%v[white]\t"+
				"min=[green]%v[white]\t"+
				"avg=[green]%v[white][::-]\n",
			// -1, -1,
			mr.Total, mr.Success, mr.Failed,
			mr.AvgSuccess,
			mr.MaxTime,
			mr.AvgTime,
			mr.MinTime)
	}

	// fmt.Fprint(vw.textView, txt)

	vw.txt = txt
}

func (vw tView) ShutDown() {
	logrus.Error("Shutting down...")
	vw.app.Stop()
}
