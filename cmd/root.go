package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
    "github.com/spf13/viper"
)

func init() {
    viper.SetDefault("ContentDir", "content")
    viper.SetDefault("LayoutDir", "layout")
    viper.SetDefault("PublicDir", "public")
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
