package punos

import (
	"fmt"
	"io"
	"os"

	"github.com/amaretto/punos/pkg/cmd/cli/controller"
	"github.com/amaretto/punos/pkg/cmd/cli/player"
	"github.com/amaretto/punos/pkg/cmd/cli/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// NewCommand create command
func NewCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "punos [subcommand]",
		Short: "dj platform",
		Long:  `dj platform`,
		//		Args: func(cmd *cobra.Command, args []string) error {
		//			if len(args) < 1 {
		//				return errors.New("requires source mp3 path")
		//			}
		//			src = args[0]
		//			return nil
		//		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("hoge")
			//if err := saveWaveImage(); err != nil {
			//	fmt.Println(err)
			//	os.Exit(1)
			//}
		},
	}

	// for logging by logrus
	var v string
	c.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		file, err := os.OpenFile("punos.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		if err := setUpLogs(file, v); err != nil {
			return err
		}
		return nil
	}
	c.PersistentFlags().StringVarP(&v, "verbosity", "v", logrus.WarnLevel.String(), "Log level (debug, info, warn, error, fatal, panic")

	c.AddCommand(
		player.NewCommand(),
	)

	c.AddCommand(
		server.NewCommand(),
	)

	c.AddCommand(
		controller.NewCommand(),
	)

	return c
}

func setUpLogs(out io.Writer, level string) error {
	logrus.SetOutput(out)
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}
	logrus.SetLevel(lvl)
	return nil
}

//func saveWaveImage() error {
//	w := waveform.NewWaveformer()
//	w.MusicPath = src
//	if err := w.SaveWaveImage(dst); err != nil {
//		return err
//	}
//	return nil
//}
