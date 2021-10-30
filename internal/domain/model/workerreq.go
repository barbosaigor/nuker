package model

type LaborContract struct {
	Operation WorkerOp   `json:"operation"`
	Workloads []Workload `json:"workloads"`
}

type WorkerOp string

const (
	Detach     = WorkerOp("detach")
	Assignment = WorkerOp("assignment")
)

type WorkerBody struct {
	ID     string `json:"id"`
	Weight int    `json:"weight"`
}
