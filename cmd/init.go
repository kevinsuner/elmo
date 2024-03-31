package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
    rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
    Use: "init <project-name>",
    Short: "placeholder",
    Long: `placeholder`,
    Run: initialize,
}

func initialize(cmd *cobra.Command, args []string) {
    if len(args[0]) == 0 {
        log.Fatal("len(args[0])", "error", "empty argument")
    }

    _, err := os.Stat(args[0])
    switch {
    case os.IsNotExist(err):
        break
    default:
        log.Fatal("os.Stat", "error", err.Error())
    }

    if err = os.MkdirAll(
        fmt.Sprintf(
            "%s/%s",
            args[0],
            viper.GetString("ContentPostsDir"),
        ),
        os.ModePerm,
    ); err != nil {
        log.Fatal("os.MkdirAll", "error", err.Error())
    }
}
