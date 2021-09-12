package cli

import (
	"fmt"

	"github.com/barbosaigor/nuker/pkg/config"
	"github.com/spf13/cobra"
)

var VERSION string = "0.0.0"

// Exec Command flags
var DryRunFlagExecCmd bool
var MinFlagExecCmd int
var MaxFlagExecCmd int
var DurationFlagExecCmd int
var MethodFlagExecCmd string

// Global flags
var Verbose bool
var Quiet bool
var NoLogFile bool
var LogFile string

// Exec flags
var Master bool
var Worker bool

var IsExec bool
var IsRun bool

var Args []string

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Nuker version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("nuker version " + VERSION)
	},
}

var ExecCmd = &cobra.Command{
	Use:     "exec [ENDPOINT] [FLAGS]",
	Short:   "Exec executes a simple http/https inline plan config",
	Example: "nuker exec http://my-api.com/product/v2/123 --min 15 --max 20 --duration 10",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		IsExec = true
		Args = args
	},
}

var RunCmd = &cobra.Command{
	Use:     "run [PLAN FILE]",
	Short:   "Run executes a plan config file",
	Example: "nuker run my-plan.yaml",
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		IsRun = true
		Args = args
	},
}

var Cli = &cobra.Command{
	Use:   "nuker",
	Short: "Nuker is a CLI tool for load testing",
	Long: "Nuker is a CLI tool for load testing, with a " +
		"powerful configuration file (but easy) for planning your tests.",
	Example: "nuker run my-plan.yaml",
	Version: VERSION,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

func ExecCli() error {

	ExecCmd.Flags().BoolVar(&DryRunFlagExecCmd, "dry-run", false, "dry run verifies your plan config, but do not run")
	ExecCmd.Flags().IntVar(&MinFlagExecCmd, "min", 0, "min defines minimum request count")
	ExecCmd.Flags().IntVar(&MaxFlagExecCmd, "max", 0, "max defines maximum request count (default equals to min)")
	ExecCmd.Flags().IntVar(&DurationFlagExecCmd, "duration", 0, "duration defines how long you test should run")
	ExecCmd.Flags().StringVar(&MethodFlagExecCmd, "method", "GET", "method defines http rest method (default GET)")
	ExecCmd.MarkFlagRequired("min")
	ExecCmd.MarkFlagRequired("duration")

	Cli.PersistentFlags().BoolVar(&Verbose, "verbose", false, "verbose shows detailed logs")
	Cli.PersistentFlags().BoolVar(&NoLogFile, "disable-log-file", false, "disable-log-file doesn't create log file")
	Cli.PersistentFlags().StringVar(&LogFile, "log-file", "", "log-file defines log file name")

	Cli.PersistentFlags().BoolVar(&Master, "master", false, "master makes nuker a master application, awaiting for workers come out")
	Cli.PersistentFlags().BoolVar(&Worker, "worker", false, "worker makes nuker a worker, and need to connect to master")

	Cli.AddCommand(VersionCmd, ExecCmd, RunCmd)

	return Cli.Execute()
}

func BuildExecCmdCfg() *config.Config {
	if len(Args) == 0 {
		return nil
	}

	return &config.Config{
		Name: "exec",
		Stages: []config.Stage{
			{
				Name: "exec stage",
				Steps: []config.Step{
					{
						Name: "exec step",
						Containers: []config.Container{
							{
								Name: "exec container",
								Min:  MinFlagExecCmd,
								Max:  MaxFlagExecCmd,
								Network: config.Network{
									Host:   Args[0],
									Method: MethodFlagExecCmd,
								},
								Duration: DurationFlagExecCmd,
							},
						},
					},
				},
			},
		},
	}
}
