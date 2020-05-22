package punos

import (
	"fmt"

	//"github.com/amaretto/punos/pkg/punos"
	"github.com/amaretto/punos/pkg/cmd/cli/turntable"
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

	c.AddCommand(
		turntable.NewCommand(),
	)

	return c
}

//func saveWaveImage() error {
//	w := waveform.NewWaveformer()
//	w.MusicPath = src
//	if err := w.SaveWaveImage(dst); err != nil {
//		return err
//	}
//	return nil
//}
