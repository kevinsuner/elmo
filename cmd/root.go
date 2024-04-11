package cmd

import (
	"flag"
	"log/slog"
	"os"
)

// TODO: Setup goreleaser :)

type Command struct {
    Flags   *flag.FlagSet
    Use     string
    Run     func(params map[string]string)
}

type Commander map[string]Command

var (
    commander   = make(Commander, 0)
    logger      = slog.New(slog.NewTextHandler(os.Stderr, nil))
)

func Init() Commander {
    commander["root"] = Command{
        Flags: nil,
        Use: `A minimalist open-source static web page generator
that lives in your terminal, and is built for individuals
that prioritize content over features.`,
        Run: func(params map[string]string) {},
    }

    return commander
}
