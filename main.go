package main

import (
	"elmo/cmd"
	"fmt"
	"os"
)

func main() {
    commander := cmd.Init()
    if len(os.Args) > 1 {
        if command, ok := commander[os.Args[1]]; ok {
            var (
                help        = command.Flags.Bool("help", false, "send help")
                project     = command.Flags.String("project", "", "project name")
                theme       = command.Flags.String("theme", "", "theme name")
                themeUrl    = command.Flags.String("theme-url", "", "theme url")
            )

            command.Flags.Parse(os.Args[2:])

            if *help {
                fmt.Println(command.Use)
                return
            }

            command.Run(map[string]string{
                "project":      *project,
                "theme":        *theme,
                "theme-url":    *themeUrl,
            })

            return
        } else {
            fmt.Println(commander["main"].Use)
        }
    } else {
        fmt.Println(commander["main"].Use)
    }
    
    // cmd.Execute()
}
