package cmd

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(image_b64)
}

var image_b64 = &cobra.Command{
	Use:   "image-encode [image-path]",
	Short: "Encode image to base64.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		encoded, err := convert64(args[0])
		if encoded == "" {
			cmd.PrintErrln("error:", err)
			os.Exit(1)
		}

		_, err = cmd.OutOrStdout().Write([]byte(encoded))
		if err != nil {
			cmd.PrintErrln("error:", err)
		}
	},
}

func convert64(image_file string) (image64 string, err error) {
	image, err := os.ReadFile(image_file)
	if err != nil {
		//slog.Error("something went wrong", "error:", err)
		return "", fmt.Errorf("%v", err)

	}
	image64 = "data:image/png;base64," + base64.StdEncoding.EncodeToString(image)
	return
}
