package cmd

import (
	"fmt"
	"log/slog"
	"os"

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

    viper.SetDefault("LogLevel", slog.LevelInfo)
    viper.Set("Logger", 
        slog.New(
            slog.NewTextHandler(
                os.Stderr,
                &slog.HandlerOptions{
                    Level: viper.Get("LogLevel").(slog.Level),
                },
            ),
        ),
    )
}

var rootCmd = &cobra.Command{
    Use: "elmo",
    Short: "placeholder",
    Long: `placeholder`,
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("placeholder") 
    },
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
