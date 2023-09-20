package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/mattn/go-isatty"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/youen/xixo/pkg/xixo"
)

//nolint:gochecknoglobals
var (
	name      string // provisioned by ldflags
	version   string // provisioned by ldflags
	commit    string // provisioned by ldflags
	buildDate string // provisioned by ldflags
	builtBy   string // provisioned by ldflags

	verbosity string
	jsonlog   bool
	debug     bool
	colormode string

	subscribers map[string]string
)

func main() {
	cobra.OnInitialize(initLog)

	rootCmd := &cobra.Command{ //nolint:exhaustruct
		Use:   fmt.Sprintf("%v real-data-file.jsonl", name),
		Short: "Masked Input Metrics Output",
		Long:  `XIXO is a purpose-built tool designed for edit XML file with stream process like PIMO`,
		Example: `
		to apply "jq -c '.bar |= ascii_upcase'" shell on all text of 'bar' elements of 'foo' entity:
			$ echo '<foo><bar>a</bar></foo>' | xixo --subscribers foo="jq -c '.bar |= ascii_upcase'"
			<foo><bar>A<bar></foo>
		`,
		Version: fmt.Sprintf(`%v (commit=%v date=%v by=%v)
Copyright (C) 2021 CGI France
License GPLv3: GNU GPL version 3 <https://gnu.org/licenses/gpl.html>.
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.`, version, commit, buildDate, builtBy),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			log.Info().
				Str("verbosity", verbosity).
				Bool("log-json", jsonlog).
				Bool("debug", debug).
				Str("color", colormode).
				Msg("start XIXO")
		},
		Args: cobra.ExactArgs(0),

		Run: func(cmd *cobra.Command, args []string) {
			if err := run(cmd); err != nil {
				log.Fatal().Err(err).Msg("end XIXO")
			}
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			log.Info().Int("return", 0).Msg("end XIXO")
		},
	}

	rootCmd.PersistentFlags().StringVarP(&verbosity, "verbosity", "v", "warn",
		"set level of log verbosity : none (0), error (1), warn (2), info (3), debug (4), trace (5)")
	rootCmd.PersistentFlags().BoolVar(&jsonlog, "log-json", false, "output logs in JSON format")
	rootCmd.PersistentFlags().StringVar(&colormode, "color", "auto", "use colors in log outputs : yes, no or auto")
	rootCmd.PersistentFlags().StringToStringVar(
		&subscribers, "subscribers", map[string]string{},
		"subscribers shell for matching elements",
	)

	if err := rootCmd.Execute(); err != nil {
		log.Err(err).Msg("error when executing command")
		os.Exit(1)
	}
}

func run(_ *cobra.Command) error {
	driver := xixo.NewDriver(os.Stdin, os.Stdout, subscribers)

	err := driver.Stream()
	if err != nil {
		log.Err(err).Msg("Error during processing")

		return err
	}

	return nil
}

func initLog() {
	color := false

	switch strings.ToLower(colormode) {
	case "auto":
		if isatty.IsTerminal(os.Stdout.Fd()) && runtime.GOOS != "windows" {
			color = true
		}
	case "yes", "true", "1", "on", "enable":
		color = true
	}

	if jsonlog {
		log.Logger = zerolog.New(os.Stderr)
	} else {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, NoColor: !color}) //nolint:exhaustruct
	}

	if debug {
		log.Logger = log.Logger.With().Caller().Logger()
	}

	setVerbosity()
}

func setVerbosity() {
	switch verbosity {
	case "trace", "5":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case "debug", "4":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info", "3":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn", "2":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error", "1":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.Disabled)
	}
}
