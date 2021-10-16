package model

type LaborContract struct {
	Operation WorkerOp `json:"operation"`
	Workload  Workload `json:"workload"`
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
