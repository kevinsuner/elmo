package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"os/exec"
	"time"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Name, Theme string

func init() {
    rootCmd.AddCommand(initCmd)
    initCmd.Flags().StringVar(&Name, "name", "", "project name")
    initCmd.Flags().StringVar(&Theme, "theme", "", "e.g. https://github.com/kevinsuner/elmo-erlosung.git")
    initCmd.MarkFlagRequired("name")
}

var initCmd = &cobra.Command{
    Use: "init",
    Short: "Initializes a new project with the given name",
    Long: `Creates a new directory with the given project name,
with a set of sub-folders such as content, posts and themes,
that are required for the program to work.
    
If a theme is provided, it will be cloned inside the themes
folder using the <git> command.`,
    Run: initialize,
}

func initialize(cmd *cobra.Command, args []string) {
    ts := time.Now()
    log := viper.Get("Logger").(*log.Logger)

    _, err := os.Stat(Name)
    if err != nil {
        if !errors.Is(err, fs.ErrNotExist) {
            log.Fatal("os.Stat", "error", err.Error())
        }
    }

    if err := os.MkdirAll(
        fmt.Sprintf(
            "%s/%s",
            Name,
            viper.GetString("ContentPostsDir"),
        ),
        os.ModePerm,
    ); err != nil {
        log.Fatal("os.MkdirAll", "error", err.Error())
    }

    if err := os.Mkdir(
        fmt.Sprintf(
            "%s/%s",
            Name,
            viper.GetString("ThemesDir"),
        ),
        os.ModePerm,
    ); err != nil {
        log.Fatal("os.Mkdir", "error", err.Error())
    }

    if cmd.Flag("theme").Changed {
        _, err = url.ParseRequestURI(cmd.Flag("theme").Value.String())
        if err != nil {
            log.Fatal("url.ParseRequestURI", "error", err.Error())
        }

        command := exec.Command("git", "clone", cmd.Flag("theme").Value.String())
        command.Dir = fmt.Sprintf("%s/%s", Name, viper.GetString("ThemesDir"))
        if err := command.Run(); err != nil {
            log.Fatal("exec.Command", "error", err.Error())
        }
    }

    log.SetReportCaller(false)
    log.Info("Done!", "took", fmt.Sprintf("%dms", time.Since(ts).Milliseconds()))
    log.Info("Created project", "name", Name)
    if cmd.Flag("theme").Changed { log.Info("Downloaded theme", "url", Theme) }
}
