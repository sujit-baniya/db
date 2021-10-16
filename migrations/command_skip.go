package migrations

import (
	"flag"
	"fmt"
	"github.com/sujit-baniya/db"
	"strings"
)

type SkipCommand struct {
}

func (c *SkipCommand) Help() string {
	helpText := `
Usage: verify-rest skip [options] ...

  Set the database level to the most recent version available, without actually running the migrations.

Options:

  -config=dbconfig.yml   Configuration file to use.
  -env="development"     Environment.
  -limit=0               Limit the number of migrations (0 = unlimited).

`
	return strings.TrimSpace(helpText)
}

func (c *SkipCommand) Synopsis() string {
	return "Sets the database level to the most recent version available, without running the migrations"
}

func (c *SkipCommand) Run(args []string) int {
	var limit int
	var dryrun bool

	cmdFlags := flag.NewFlagSet("up", flag.ContinueOnError)
	cmdFlags.Usage = func() { ui.Output(c.Help()) }
	cmdFlags.IntVar(&limit, "limit", 0, "Max number of migrations to skip.")

	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	err := SkipMigrations(Up, dryrun, limit)
	if err != nil {
		ui.Error(err.Error())
		return 1
	}

	return 0
}

func SkipMigrations(dir MigrationDirection, dryrun bool, limit int) error {
	curBD, _ := db.DB.DB()

	source := FileMigrationSource{
		Dir: DefaultConfig.Dir,
	}

	n, err := SkipMax(curBD, db.DefaultDialect, source, dir, limit)
	if err != nil {
		return fmt.Errorf("Migration failed: %s", err)
	}

	switch n {
	case 0:
		ui.Output("All migrations have already been applied")
	case 1:
		ui.Output("Skipped 1 migration")
	default:
		ui.Output(fmt.Sprintf("Skipped %d migrations", n))
	}

	return nil
}
