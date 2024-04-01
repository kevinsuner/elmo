package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// TODO: Setup goreleaser :)

const (
    ContentDir  string = "content"
    PublicDir   string = "public"
    ThemesDir   string = "themes"
    
    DefaultTheme string = "https://github.com/kevinsuner/elmo-thumbalina.git"
)

var (
    contentPostsDir string = fmt.Sprintf("%s/posts", ContentDir)
    publicPostsDir  string = fmt.Sprintf("%s/posts", PublicDir)
    themeDir        string
    partialsDir     string

    logger *log.Logger
)

func init() {
    viper.SetDefault("language", "en")
    viper.SetDefault("theme", "elmo-thumbalina")
    viper.SetDefault("debug", false)

    themeDir        = fmt.Sprintf("%s/%s", ThemesDir, viper.GetString("theme"))
    partialsDir     = fmt.Sprintf("%s/partials", themeDir)

    logLevel := log.InfoLevel
    reportCaller := false
    if viper.GetBool("debug") {
        reportCaller = true
        logLevel = log.DebugLevel
    }

    logger = log.NewWithOptions(
        os.Stderr,
        log.Options{
            ReportCaller: reportCaller,
            ReportTimestamp: true,
            TimeFormat: time.Kitchen,
            Level: logLevel,
        },
    )
}

var rootCmd = &cobra.Command{
    Use: "elmo",
    Short: "A minimalist static web site generator",
    Long: `A minimalist open-source static web page generator
that lives in your terminal, and is built for individuals
that prioritize content over features.`,
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
