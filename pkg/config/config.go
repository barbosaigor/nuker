package config

type Network struct {
	Protocol string
	Host     string
	Path     string
	Method   string
	Headers  map[string]string
	Body     string
}

type Container struct {
	Name     string
	Duration int
	HoldFor  int
	Min      int
	Max      int
	Network  Network
	Payload  string
	Filename string
}

type Step struct {
	Name       string
	Containers []Container
}

type Stage struct {
	Name  string
	Steps []Step
}

type Global struct {
	Host string
}

type Config struct {
	Name   string
	Global Global
	Stages []Stage
}
