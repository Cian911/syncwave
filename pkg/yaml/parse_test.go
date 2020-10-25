package yaml

import (
	"fmt"
	"testing"
)

func TestParseMain(t *testing.T) {
	t.Run("It parses a yaml file.", func(t *testing.T) {
		file := "config_test.yaml"
		nodes, err := ParseMain(file)

		if err != nil {
			fmt.Errorf("Error: %v", err)
		}

		fmt.Printf("Nodes: %#v\n", nodes.MasterN.H)
	})
}
