package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/barbosaigor/nuker/internal/domain/model"
	"github.com/barbosaigor/nuker/internal/domain/service/requester"
	"github.com/barbosaigor/nuker/pkg/metrics"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

type Worker interface {
	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error
	Do(ctx context.Context) error
}

type worker struct {
	id        string
	weight    int
	requester requester.Requester

	masterURI string
	client    *resty.Client
}

func New(ID, masterURI string, weight int, requester requester.Requester) Worker {
	return &worker{
		id:        ID,
		weight:    weight,
		requester: requester,
		masterURI: masterURI,
		client:    resty.New(),
	}
}

func (w *worker) Connect(ctx context.Context) error {
	log.Trace("connecting to master node")

	workerBody := model.WorkerBody{Weight: w.weight}

	req, err := w.client.R().
		SetBody(workerBody).
		Post(w.workerEndpoint())
	if err != nil {
		return fmt.Errorf("error to connect to master node: %w", err)
	}

	if req.IsError() {
		log.Errorf("error to connect to master node, non ok status code: %v", req.StatusCode())
		return errors.New("non ok status code from master")
	}

	var wr model.WorkerBody

	err = json.Unmarshal(req.Body(), &wr)
	if err != nil {
		log.Errorf("error to unmarshal master response body: %v", err)
		return err
	}

	if w.id == "" && wr.ID != "" {
		w.id = wr.ID
	}

	log.
		WithField("id", w.id).
		Debug("worker connected to master")

	return nil
}

func (w *worker) Disconnect(ctx context.Context) error {
	_, err := w.client.R().Delete(w.workerEndpoint())

	return err
}

func (w worker) workerEndpoint() string {
	var wID string
	if w.id != "" {
		wID = "/" + w.id
	}

	return w.masterURI + "/worker" + wID
}

// Do makes requests to master server, to get pipeline events
func (w worker) Do(ctx context.Context) error {
	metChan := make(chan *metrics.NetworkMetrics)
	defer close(metChan)
	go w.sendMetrics(metChan)

	for {
		req, err := w.client.R().
			Get(w.workerEndpoint())

		if err != nil {
			return fmt.Errorf("error to connect to master node: %w", err)
		}

		if req.StatusCode() == http.StatusNotFound {
			err := w.Connect(ctx)
			if err != nil {
				log.
					WithField("worker", w.id).
					Errorf("error to connect to master node: %v", err)
			}
			continue
		}

		if req.IsError() {
			log.
				WithField("worker", w.id).
				Tracef("non ok status code from master: %v", req.StatusCode())
			<-time.After(time.Second)
			continue
		}

		if req.StatusCode() == http.StatusNoContent {
			log.
				WithField("worker", w.id).
				Trace("no workload remaining")
			<-time.After(time.Second)
			continue
		}

		log.
			WithField("worker", w.id).
			Tracef("master's body: %s", req.Body())

		var wr model.LaborContract
		err = json.Unmarshal(req.Body(), &wr)
		if err != nil {
			log.
				WithField("worker", w.id).
				Errorf("error to decode master's response body: %v", err)
			<-time.After(time.Second)
			continue
		}

		if wr.Operation == model.Detach {
			log.
				WithField("worker", w.id).
				Infof("detaching from master")
			return nil
		}

		select {
		case <-time.After(time.Second):
		case <-w.assignWl(ctx, wr.Workload, metChan):
		}
	}
}

func (w worker) sendMetrics(metChan <-chan *metrics.NetworkMetrics) error {
	for m := range metChan {
		log.Tracef("sending metrics to master: %v", m)

		req, err := w.client.R().
			SetBody(m).
			Post(w.workerEndpoint() + "/metrics")

		if err != nil {
			log.Errorf("error to connect master node: %v", err)
			continue
		}

		if req.IsError() {
			log.Errorf("non success status code: %v, when sent metrics. Body: %s", req.StatusCode(), req.Body())
		}
	}

	return nil
}

func (w worker) assignWl(ctx context.Context, wl model.Workload, metChan chan<- *metrics.NetworkMetrics) <-chan struct{} {
	log.
		WithField("worker", w.id).
		Tracef("request count: %d", wl.RequestsCount)
	done := make(chan struct{})

	go func() {
		_ = w.requester.Assign(ctx, wl, metChan)
		done <- struct{}{}
	}()

	return done
}
