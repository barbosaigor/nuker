package runner

import (
	"context"
	"fmt"
	"os"

	"github.com/barbosaigor/nuker/internal/domain/service/pipeline"
	"github.com/barbosaigor/nuker/pkg/config"
	pkgrunner "github.com/barbosaigor/nuker/pkg/runner"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var VERSION string = "0.1.0"

type runner struct {
	pipeline pipeline.Pipeline
}

func New(pipeline pipeline.Pipeline) pkgrunner.Runner {
	return &runner{
		pipeline: pipeline,
	}
}

func (r *runner) Run(ctx context.Context) error {
	return r.execCli(ctx)
}

func (r *runner) execCli(ctx context.Context) error {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Nuker version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("nuker version " + VERSION)
		},
	}

	var dryRunFlagExecCmd bool
	var minFlagExecCmd int
	var maxFlagExecCmd int
	var durationFlagExecCmd int
	var methodFlagExecCmd string
	execCmd := &cobra.Command{
		Use:     "exec [ENDPOINT] [FLAGS]",
		Short:   "Exec executes a simple http/https inline plan config",
		Example: "nuker exec http://my-api.com/product/v2/123 --min 15 --max 20 --duration 10",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.Config{
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
										Min:  minFlagExecCmd,
										Max:  maxFlagExecCmd,
										Network: config.Network{
											Host:   args[0],
											Method: methodFlagExecCmd,
										},
										Duration: durationFlagExecCmd,
									},
								},
							},
						},
					},
				},
			}

			if dryRunFlagExecCmd {
				log.Info().Msgf("plan:\n%+v", cfg)
				return nil
			}

			return r.pipeline.Run(ctx, cfg)
		},
	}

	execCmd.Flags().BoolVar(&dryRunFlagExecCmd, "dry-run", false, "dry run verifies your plan config, but do not run")
	execCmd.Flags().IntVar(&minFlagExecCmd, "min", 0, "min defines minimum request count")
	execCmd.Flags().IntVar(&maxFlagExecCmd, "max", 0, "max defines maximum request count (default equals to min)")
	execCmd.Flags().IntVar(&durationFlagExecCmd, "duration", 0, "duration defines how long you test should run")
	execCmd.Flags().StringVar(&methodFlagExecCmd, "method", "GET", "method defines http rest method (default GET)")
	execCmd.MarkFlagRequired("min")
	execCmd.MarkFlagRequired("duration")

	runCmd := &cobra.Command{
		Use:     "run [PLAN FILE]",
		Short:   "Run executes a plan config file",
		Example: "nuker run my-plan.yaml",
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return r.runPipelineFromFile(ctx, args[0])
		},
	}

	cli := &cobra.Command{
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

	cli.AddCommand(versionCmd, execCmd, runCmd)

	return cli.Execute()
}

func (r *runner) runPipelineFromFile(ctx context.Context, cfgFName string) error {
	data, err := os.ReadFile(cfgFName)
	if err != nil {
		return err
	}

	cfg, err := config.YamlUnmarshal(data)
	if err != nil {
		return err
	}

	return r.pipeline.Run(ctx, *cfg)
}
