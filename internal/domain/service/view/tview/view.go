package tview

import (
	"fmt"
	"sync"
	"time"

	"github.com/barbosaigor/nuker/internal/domain/service/view"
	"github.com/barbosaigor/nuker/pkg/metrics"
	"github.com/rivo/tview"
	"github.com/sirupsen/logrus"
)

type tView struct {
	mr   *metrics.MetricRate
	app  *tview.Application
	txt  string
	quit chan bool
}

func New() view.View {
	return &tView{
		mr:   &metrics.MetricRate{},
		app:  tview.NewApplication(),
		quit: make(chan bool),
	}
}

func (vw *tView) SetMetric(mr *metrics.MetricRate) {
	logrus.Tracef("setting metric: %v", mr)

	var txt string
	if mr == nil {
		txt = fmt.Sprintf("%v\n", mr)
	} else {
		txt = fmt.Sprintf(
			"requests........: [::b]total=[green]%d[white]\tsuccess=[green]%d[white]\tfailed=[green]%d[white][::-]\t\n"+
				"success ratio.....: [::b][green]%f[white][::-]\n"+
				"request time......: [::b]%%"+
				"max=[green]%v[white]\t"+
				"min=[green]%v[white]\t"+
				"avg=[green]%v[white][::-]\n",
			// -1, -1,
			mr.Total, mr.Success, mr.Failed,
			mr.AvgSuccess*100,
			mr.MaxTime,
			mr.AvgTime,
			mr.MinTime)
	}

	vw.txt = txt
}

func (vw *tView) Start() {
	logrus.Trace("starting tView...")

	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetChangedFunc(func() {
			vw.app.Draw()
		})
	textView.SetBorder(true)

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		err := vw.app.SetRoot(textView, true).EnableMouse(false).Run()
		if err != nil {
			panic(err)
		}
	}()

	fmt.Fprint(textView, "setting up pipeline...\n")

	go func() {
		defer wg.Done()
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

	wg.Wait()
}

func (vw tView) ShutDown() {
	logrus.Trace("shutting down...")
	vw.app.Stop()
}
