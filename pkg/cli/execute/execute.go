package execute

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewCommand() (c *cobra.Command) {
	c = &cobra.Command{
		Use:   "execute",
		Short: "Execute command on remote hosts",
		Long:  "TODO: Write a longer message for this.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("I am being executed.")
		},
	}

	return
}
