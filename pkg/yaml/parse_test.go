package yaml

import (
	"fmt"
	"testing"
)

func TestParseConfigFile(t *testing.T) {
	t.Run("It parses a config yaml file.", func(t *testing.T) {
		file := "config_test.yaml"
		nodes, err := ParseConfigFile(file)

		if err != nil {
			fmt.Errorf("Error: %v", err)
		}

		fmt.Printf("Nodes: %#v\n", nodes.MasterN.H)
	})
}

func TestParseScenarioFile(t *testing.T) {
	t.Run("It parsed a scenario config yaml file", func(t *testing.T) {
		file := "scenario_test.yaml"
		scenario, err := ParseScenarioFile(file)

		if err != nil {
			fmt.Errorf("Error: %v\n", err)
		}

		fmt.Printf("Scenario: %v\n", scenario.S.ScenarioTasks)
	})
}
