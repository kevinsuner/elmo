package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"os/exec"
	"time"

	"github.com/spf13/cobra"
)

var ProjectName, ThemeURL string

func init() {
    rootCmd.AddCommand(initCmd)
    initCmd.Flags().StringVar(&ProjectName, "name", "", "project name")
    initCmd.Flags().StringVar(&ThemeURL, "theme-url", "", "e.g. https://github.com/kevinsuner/elmo-thumbalina.git")
    initCmd.MarkFlagRequired("name")
}

var initCmd = &cobra.Command{
    Use: "init",
    Short: "Initialize a new project using the given name",
    Long: `Creates a new directory using the given project name,
with a set of sub-folders such as content, posts and themes,
which are required for the program to work.
 
If a theme is provided, it will be cloned inside the themes
folder using the <git> command, otherwise it will fallback
to cloning the default theme (https://github.com/kevinsuner/thumbalina).

The user is responsible to let the program know which theme
should it use, via its configuration file <elmo.toml>.`,
    Run: initialize,
}

func cloneTheme(themeUrl, projectName string) error {
    _, err := url.ParseRequestURI(themeUrl)
    if err != nil {
        return err     
    }

    command := exec.Command("git", "clone", themeUrl)
    command.Dir = fmt.Sprintf("%s/%s", projectName, ThemesDir)
    if err := command.Run(); err != nil {
        return err
    }

    return nil
}

func initialize(cmd *cobra.Command, args []string) {
    ts := time.Now()
    _, err := os.Stat(ProjectName)
    if err != nil {
        if !errors.Is(err, fs.ErrNotExist) {
            logger.Fatal("os.Stat", "error", err.Error())
        }
    }

    err = os.MkdirAll(fmt.Sprintf("%s/%s", ProjectName, contentPostsDir), os.ModePerm)
    if err != nil {
        logger.Fatal("os.MkdirAll", "error", err.Error())
    }

    err = os.Mkdir(fmt.Sprintf("%s/%s", ProjectName, ThemesDir), os.ModePerm)
    if err != nil {
        logger.Fatal("os.Mkdir", "error", err.Error())
    }

    themeUrl := DefaultTheme
    if cmd.Flag("theme-url").Changed {
        themeUrl = cmd.Flag("theme-url").Value.String()
        if err := cloneTheme(themeUrl, ProjectName); err != nil {
            logger.Fatal("cloneTheme", "error", err.Error())
        }
    } else {
        if err := cloneTheme(themeUrl, ProjectName); err != nil {
            logger.Fatal("cloneTheme", "error", err.Error())
        }
    }

    logger.Info("Done!", "took", fmt.Sprintf("%dms", time.Since(ts).Milliseconds()))
    logger.Info("Created project", "name", ProjectName)
    logger.Info("Downloaded theme", "url", themeUrl)
}
