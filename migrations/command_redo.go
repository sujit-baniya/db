package migrations

import (
	"flag"
	"fmt"
	"github.com/sujit-baniya/db"
	"strings"
)

type RedoCommand struct {
}

func (c *RedoCommand) Help() string {
	helpText := `
Usage: verify-rest redo [options] ...

  Reapply the last migration.

Options:

  -config=dbconfig.yml   Configuration file to use.
  -env="development"     Environment.
  -dryrun                Don't apply migrations, just print them.

`
	return strings.TrimSpace(helpText)
}

func (c *RedoCommand) Synopsis() string {
	return "Reapply the last migration"
}

func (c *RedoCommand) Run(args []string) int {
	var dryrun bool

	cmdFlags := flag.NewFlagSet("redo", flag.ContinueOnError)
	cmdFlags.Usage = func() { ui.Output(c.Help()) }
	cmdFlags.BoolVar(&dryrun, "dryrun", false, "Don't apply migrations, just print them.")

	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}
	curBD, _ := db.DB.DB()

	source := FileMigrationSource{
		Dir: DefaultConfig.Dir,
	}

	migrations, _, err := PlanMigration(curBD, db.DefaultDialect, source, Down, 1)
	if err != nil {
		ui.Error(fmt.Sprintf("Migration (redo) failed: %v", err))
		return 1
	} else if len(migrations) == 0 {
		ui.Output("Nothing to do!")
		return 0
	}

	if dryrun {
		PrintMigration(migrations[0], Down)
		PrintMigration(migrations[0], Up)
	} else {
		_, err := ExecMax(curBD, db.DefaultDialect, source, Down, 1)
		if err != nil {
			ui.Error(fmt.Sprintf("Migration (down) failed: %s", err))
			return 1
		}

		_, err = ExecMax(curBD, db.DefaultDialect, source, Up, 1)
		if err != nil {
			ui.Error(fmt.Sprintf("Migration (up) failed: %s", err))
			return 1
		}

		ui.Output(fmt.Sprintf("Reapplied migration %s.", migrations[0].Id))
	}

	return 0
}
