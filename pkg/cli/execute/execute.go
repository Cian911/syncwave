package execute

import (
	"log"
	"os"
	"sync"

	"github.com/cian911/raspberry-pi-provisioner/pkg/printer"
	"github.com/cian911/raspberry-pi-provisioner/pkg/ssh"
	"github.com/cian911/raspberry-pi-provisioner/pkg/yaml"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ExecutionResult struct {
	Host   string
	Task   string
	Stdout string
	Stderr string
	Status error
}

func NewCommand() (c *cobra.Command) {
	c = &cobra.Command{
		Use:   "execute",
		Short: "Execute command on remote hosts",
		Long:  "Pass a hosts configuration file and scenario file in order to execute tasks on remote hosts.",
		Run: func(cmd *cobra.Command, args []string) {
			configFile := viper.GetString("config")
			parsedConfigFile, err := yaml.ParseConfigFile(configFile)

			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}

			scenarioFile := viper.GetString("scenario")
			parsedScenarioFile, err := yaml.ParseScenarioFile(scenarioFile)

			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}

			results := make(chan ExecutionResult)

			masterNodes := make(map[string]string)
			workerNodes := make(map[string]string)

			for _, v := range parsedConfigFile.MasterN.H {
				masterNodes[v.Hostname] = v.IPAddress
			}

			for _, v := range parsedConfigFile.WorkerN.W {
				workerNodes[v.Hostname] = v.IPAddress
			}

			tasks := make(map[string]string)

			for _, v := range parsedScenarioFile.S.ScenarioTasks {
				tasks[v.TaskName] = v.TaskCMD
			}

			table := printer.NewTable(os.Stdout,
				[]string{
					"Host",
					"Task",
					"Output",
				},
			)

			var wg_cmd sync.WaitGroup
			var wg_processing sync.WaitGroup

			for _, task := range tasks {
				for _, host := range workerNodes {
					wg_cmd.Add(1)
					go func(host, task string) {
						defer wg_cmd.Done()
						stdout, stderr, status := ssh.Execute(host, task)
						res := ExecutionResult{
							host,
							task,
							stdout,
							stderr,
							status,
						}

						results <- res
					}(host, task)
				}
			}

			wg_processing.Add(1)
			go func() {
				defer wg_processing.Done()
				for res := range results {
					data := []string{res.Host, res.Task, res.Stdout, res.Stderr}
					if res.Status == nil {
						table = printer.SuccessStyle(table, data)
					} else {
						table = printer.ErrorStyle(table, data)
					}
				}
			}()
			wg_cmd.Wait()
			close(results)
			wg_processing.Wait()
			table.Render()
		},
	}

	return
}
