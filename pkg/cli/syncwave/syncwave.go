package syncwave

import (
	"github.com/cian911/raspberry-pi-provisioner/pkg/cli/execute"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func New() (c *cobra.Command) {
	c = &cobra.Command{
		Use:   "syncwave",
		Short: "syncwave helps to automate the management of your infrastructre.",
		Long:  "syncwave helps to automate the management of your local infrastructure. Built and maintained by Cian Gallagher - @Cian911",
	}

	// Define flags
	c.PersistentFlags().StringP("config", "c", "", "Pass configuration path/file.")
	c.PersistentFlags().StringP("scenario", "s", "", "Pass scenario path/file.")

	// Mark flags as required
	c.MarkFlagRequired("scenario")
	c.MarkFlagRequired("config")

	viper.BindPFlag("config", c.PersistentFlags().Lookup("config"))
	viper.BindPFlag("scenario", c.PersistentFlags().Lookup("scenario"))

	c.AddCommand(execute.NewCommand())

	return
}
