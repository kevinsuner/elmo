package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
    viper.SetDefault("Language", "en")
    viper.SetDefault("ContentDir", "content")
    viper.SetDefault("PublicDir", "public")
    viper.SetDefault("ThemesDir", "themes")
    viper.SetDefault("Theme", "sesame")

    viper.SetDefault(
        "ThemeDir",
        fmt.Sprintf("%s/%s", viper.GetString("ThemesDir"), viper.GetString("Theme")),
    )

    viper.SetDefault(
        "PartialsDir",
        fmt.Sprintf("%s/partials", viper.GetString("ThemeDir")),
    )

    viper.SetDefault(
        "PublicPostsDir",
        fmt.Sprintf("%s/posts", viper.GetString("PublicDir")),
    )

    viper.SetDefault(
        "ContentPostsDir",
        fmt.Sprintf("%s/posts", viper.GetString("ContentDir")),
    )

    viper.SetDefault("LogLevel", log.InfoLevel)
    viper.Set("Logger", 
        log.NewWithOptions(
            os.Stderr,
            log.Options{
                ReportCaller: true,
                ReportTimestamp: true,
                TimeFormat: time.Kitchen,
                Level: viper.Get("LogLevel").(log.Level),
            },
        ),
    )
}

var rootCmd = &cobra.Command{
    Use: "elmo",
    Short: "A minimalist static web site generator",
    Long: `A minimalist open-source static web page generator, that lives
in your terminal, and is built for individuals that prioritize
content over features.
    `,
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
