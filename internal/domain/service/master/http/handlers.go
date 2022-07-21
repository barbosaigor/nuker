package http

import (
	"encoding/json"
	"net/http"

	"github.com/barbosaigor/nuker/internal/domain/model"
	"github.com/barbosaigor/nuker/pkg/metrics"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func (m *master) newWorkerWithID(c *fiber.Ctx) error {
	wb, err := m.parseWorkerBody(c)
	if err != nil {
		return err
	}

	wb.ID = utils.ImmutableString(c.Params("id"))

	m.orchSvc.AddWorker(wb.ID, wb.Weight)

	resp, err := json.Marshal(wb)
	if err != nil {
		return err
	}

	c.Send(resp)
	return nil
}

func (m *master) newWorker(c *fiber.Ctx) error {
	wb, err := m.parseWorkerBody(c)
	if err != nil {
		return err
	}

	wb.ID = uuid.NewString()

	m.orchSvc.AddWorker(wb.ID, wb.Weight)

	resp, err := json.Marshal(wb)
	if err != nil {
		return err
	}

	c.Set("Content-Type", "application/json")
	c.Send(resp)

	return nil
}

func (m *master) deleteWorker(c *fiber.Ctx) error {
	workerID := utils.ImmutableString(c.Params("id"))

	m.orchSvc.DelWorker(workerID)

	return nil
}

func (m *master) getWorkload(c *fiber.Ctx) error {
	workerID := utils.ImmutableString(c.Params("id"))

	if !m.orchSvc.HasWorker(workerID) {
		log.Debugf("worker %q not found", workerID)
		return c.SendStatus(http.StatusNotFound)
	}

	m.orchSvc.FlushWorker(workerID)

	wls := m.orchSvc.TakeWorkloads(workerID)
	if !m.done && len(wls) == 0 {
		return c.SendStatus(http.StatusNoContent)
	}

	var op model.WorkerOp
	if m.done && len(wls) == 0 {
		op = model.Detach
	} else {
		op = model.Assignment
	}

	lc := model.LaborContract{
		Operation: op,
		Workloads: wls,
	}

	resp, err := json.Marshal(lc)
	if err != nil {
		return err
	}

	c.Set("Content-Type", "application/json")
	c.Send(resp)

	return nil
}

func (m *master) addMetrics(metChan chan<- *metrics.NetworkMetrics) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		workerID := utils.ImmutableString(c.Params("id"))

		if !m.orchSvc.HasWorker(workerID) {
			return c.SendStatus(http.StatusNotFound)
		}

		var met metrics.NetworkMetrics

		err := c.BodyParser(&met)
		if err != nil {
			return err
		}

		metChan <- &met

		return nil
	}
}
