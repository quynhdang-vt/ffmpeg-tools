package main

import (
	"errors"
	"fmt"
	"github.com/quynhdang-vt/vt-engine-tools/model"
	"github.com/spf13/cobra"
)

func main() {
	var segmentDuration string
	var idTag string
	var cmdMerge = &cobra.Command{
		Use:   "split-on-silence input_audio_file ...",
		Short: "split-on-silence audio_file.",
		Long: "split-on-silence:  split an audio file into segments of audios that are non-silence. The segments must have at least " +
			"duration >= threshold (default 30 seconds ",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("requires at least one file")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			transcriptionOutput := model.NewEngineOutput()
			for i, v := range args {
				vlf, err := model.ParseVLF(v)
				if err == nil {
					transcriptionOutput.AddVLF(fmt.Sprintf("%s_%d", idTag, i), vlf)
				}
			}
			transcriptionOutput.Sort()
			// output to outputFile
			transcriptionOutput.Write(outputFile)
		},
	}
	var rootCmd = &cobra.Command{Use: "ffmpeg-tools"}
	cmdMerge.Flags().StringVarP(&segmentDuration, "duration", "d", "30", "minimum segment duration")

	rootCmd.AddCommand(cmdMerge)
	rootCmd.Execute()
}
