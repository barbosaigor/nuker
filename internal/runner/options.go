package runner

import "github.com/barbosaigor/nuker/internal/cli"

type Operation int

const (
	exec = Operation(iota)
	run
	employee
	unknow
)

type Options struct {
	Op Operation
}

func LoadCfg() (Options, error) {
	return Options{
		Op: opFromCli(),
	}, nil
}

func opFromCli() Operation {
	if cli.IsExec {
		return exec
	}

	if cli.IsRun {
		return run
	}

	if cli.IsWorker {
		return employee
	}

	return unknow
}
