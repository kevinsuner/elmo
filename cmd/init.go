package cmd

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"os/exec"
)

func init() {
    commander["init"] = Command{
        Flags: flag.NewFlagSet("init", flag.ExitOnError),
        Use: `Creates a new directory using the given project name,
with a set of sub-folders such as content, posts and themes,
which are required for the program to work.
 
If a theme is provided, it will be cloned inside the themes
folder using the <git> command.`,
        Run: initialize,
    }
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

func initialize(params map[string]string) {
    projectName := params["project"]
    _, err := os.Stat(projectName)
    if err != nil {
        if !errors.Is(err, fs.ErrNotExist) {
            logger.Error("os.Stat", "error", err.Error())
            os.Exit(1)
        }
    }

    err = os.MkdirAll(fmt.Sprintf("%s/%s", projectName, contentPostsDir), os.ModePerm)
    if err != nil {
        logger.Error("os.MkdirAll", "error", err.Error())
        os.Exit(1)
    }

    err = os.Mkdir(fmt.Sprintf("%s/%s", projectName, ThemesDir), os.ModePerm)
    if err != nil {
        logger.Error("os.Mkdir", "error", err.Error())
        os.Exit(1)
    }

    err = cloneTheme(params["theme-url"], projectName)
    if err != nil {
        logger.Error("cloneTheme", "error", err.Error())
        os.Exit(1)
    }
}
