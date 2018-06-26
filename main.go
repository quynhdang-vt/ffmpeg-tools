package main

import (
	"errors"
	"github.com/spf13/cobra"
)

/*
Given inputfile, --> segments into non-silence
OUTPUT:  result file with
      start, end, filename

 */
func main() {
	var timingFile, outputFile, outputAudioDirectory string
	var cmdSplit = &cobra.Command{
		Use:   "split inputfile ",
		Short: "split inputfile ",
		Long: "split:  split an audio file into segments of audios that are non-silence. Must give timing information about start, end. OutputFile when succeeded will contain the list of segment filenames",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("requires at least one file")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {

			for _, v := range args {
				worker :=  Extractor  {
					inputFile:v,
					outputAudioDir:outputFile,
					timingFile:timingFile,
					outputFile:outputFile,
				}
				worker.Extract()
			}
		},
	}
	cmdSplit.Flags().StringVarP(&timingFile, "timingFile", "t", "", "timing file containing start,end")
	cmdSplit.Flags().StringVarP(&outputFile, "output", "o", "", "output file")
	cmdSplit.Flags().StringVarP(&outputAudioDirectory, "directory", "d", "", "directory to output audio segments")



	var rootCmd = &cobra.Command{Use: "ffmpeg-tools"}

	rootCmd.AddCommand(cmdSplit)
	rootCmd.Execute()
}
