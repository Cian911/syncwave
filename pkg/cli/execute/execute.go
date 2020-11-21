package execute

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cian911/raspberry-pi-provisioner/pkg/printer"
	"github.com/cian911/raspberry-pi-provisioner/pkg/ssh"
	"github.com/cian911/raspberry-pi-provisioner/pkg/yaml"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ExecutionResult struct {
	Host   string
	Task   string
	Output string
	Status error
}

func NewCommand() (c *cobra.Command) {
	c = &cobra.Command{
		Use:   "execute",
		Short: "Execute command on remote hosts",
		Long:  "TODO: Write a longer message for this.",
		Run: func(cmd *cobra.Command, args []string) {
			configFile := viper.GetString("config-file")

			if configFile == "" {
				// Need to check as default if one exists in the current directory
				log.Fatal("You must pass a config file.")
			}

			parsedConfigFile, err := yaml.ParseConfigFile(configFile)

			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}

			scenarioFile := viper.GetString("scenario")

			if scenarioFile == "" {
				log.Fatal("You must pass a valid scanrio to be executed.")
			}

			parsedScenarioFile, err := yaml.ParseScenarioFile(scenarioFile)

			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}

			// results := make(chan error, 10)
			results := make(chan ExecutionResult)

			timeout := time.After(15 * time.Second)

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
				},
			)

			for _, task := range tasks {
				for _, host := range workerNodes {
					go func(host, task string) {
						output, status := ssh.Execute(host, task)
						res := ExecutionResult{
							host,
							task,
							output,
							status,
						}

						results <- res
					}(host, task)
				}
			}

			totalResultsCount := totalPrintableResults(len(workerNodes)+len(masterNodes), len(tasks))

			for i := 0; i < totalResultsCount; i++ {
				select {
				case res := <-results:
					data := []string{res.Host, res.Task}
					table.Rich(data, []tablewriter.Colors{tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiGreenColor}, tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiGreenColor}})
					table.Append([]string{
						res.Host,
						res.Task,
					})
				case <-timeout:
					fmt.Println("Timeout.")
					table.Render()
					return
				}
			}

			table.Render()
		},
	}

	c.Flags().StringP("config-file", "", "", "Pass a configuration file to syncwave to be parsed.")
	c.Flags().StringP("scenario", "", "", "Pass the scenario as an argument you which to run.")

	viper.BindPFlag("config-file", c.Flags().Lookup("config-file"))
	viper.BindPFlag("scenario", c.Flags().Lookup("scenario"))

	return
}

func totalPrintableResults(nodesCount, tasksCount int) int {
	return nodesCount * tasksCount
}

func hasCommandFailed(status error) bool {
	if status != nil {
		return true
	}

	return false
}
