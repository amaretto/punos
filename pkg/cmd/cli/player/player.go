package player

import (
	"github.com/spf13/cobra"
)

type Options struct {
	confPath string
}

var (
	o = &Options{}
)

// NewCommand create command
func NewCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "player",
		Short: "start player",
		Long:  `start player`,
		Run: func(cmd *cobra.Command, args []string) {
			app := New(o.confPath)
			app.Start()
		},
	}

	c.Flags().StringVarP(&o.confPath, "conf", "c", "~/.punos/conf", "config path")

	return c
}
