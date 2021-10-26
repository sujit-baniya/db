package migrations

import (
	"fmt"
	"os"

	"github.com/mitchellh/cli"
)

var ui cli.Ui

func SetMigrationUi() {
	i := &cli.BasicUi{Writer: os.Stdout}
	ui = &cli.ColoredUi{
		Ui:          i,
		OutputColor: cli.UiColorGreen,
		ErrorColor:  cli.UiColorRed,
		InfoColor:   cli.UiColorBlue,
		WarnColor:   cli.UiColorYellow,
	}
}

func MigrationMain() int {
	SetMigrationUi()
	migrateIndex := 2
	for i, val := range os.Args {
		if val == "--migrate" {
			migrateIndex = i + 1
		}
	}

	cmd := &cli.CLI{
		Args: os.Args[migrateIndex:],
		Commands: map[string]cli.CommandFactory{
			"up": func() (cli.Command, error) {
				return &UpCommand{}, nil
			},
			"down": func() (cli.Command, error) {
				return &DownCommand{}, nil
			},
			"redo": func() (cli.Command, error) {
				return &RedoCommand{}, nil
			},
			"status": func() (cli.Command, error) {
				return &StatusCommand{}, nil
			},
			"new": func() (cli.Command, error) {
				return &NewCommand{}, nil
			},
			"skip": func() (cli.Command, error) {
				return &SkipCommand{}, nil
			},
		},
		HelpFunc: cli.BasicHelpFunc("verify-rest"),
		Version:  "1.0.0",
	}

	exitCode, err := cmd.Run()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err.Error())
		return 1
	}

	return exitCode
}
