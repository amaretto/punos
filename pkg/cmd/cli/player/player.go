package player

import (
	"github.com/spf13/cobra"
)

// NewCommand create command
func NewCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "player",
		Short: "start player",
		Long:  `start player`,
		Run: func(cmd *cobra.Command, args []string) {
			app := New()
			app.Start()
		},
	}
	return c
}
