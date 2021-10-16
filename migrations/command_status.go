package migrations

import (
	"flag"
	"fmt"
	"github.com/sujit-baniya/db"
	"os"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
)

type StatusCommand struct {
}

func (c *StatusCommand) Help() string {
	helpText := `
Usage: verify-rest status [options] ...

  Show migration status.

Options:

  -config=dbconfig.yml   Configuration file to use.
  -env="development"     Environment.

`
	return strings.TrimSpace(helpText)
}

func (c *StatusCommand) Synopsis() string {
	return "Show migration status"
}

func (c *StatusCommand) Run(args []string) int {
	cmdFlags := flag.NewFlagSet("status", flag.ContinueOnError)
	cmdFlags.Usage = func() { ui.Output(c.Help()) }

	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}
	curBD, _ := db.DB.DB()
	source := FileMigrationSource{
		Dir: DefaultConfig.Dir,
	}
	migrations, err := source.FindMigrations()
	if err != nil {
		ui.Error(err.Error())
		return 1
	}

	records, err := GetMigrationRecords(curBD, db.DefaultDialect)
	if err != nil {
		ui.Error(err.Error())
		return 1
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Migration", "Applied"})
	table.SetColWidth(60)

	rows := make(map[string]*statusRow)

	for _, m := range migrations {
		rows[m.Id] = &statusRow{
			Id:       m.Id,
			Migrated: false,
		}
	}

	for _, r := range records {
		if rows[r.Id] == nil {
			ui.Warn(fmt.Sprintf("Could not find migration file: %v", r.Id))
			continue
		}

		rows[r.Id].Migrated = true
		rows[r.Id].AppliedAt = r.AppliedAt
	}

	for _, m := range migrations {
		if rows[m.Id] != nil && rows[m.Id].Migrated {
			table.Append([]string{
				m.Id,
				rows[m.Id].AppliedAt.String(),
			})
		} else {
			table.Append([]string{
				m.Id,
				"no",
			})
		}
	}

	table.Render()

	return 0
}

type statusRow struct {
	Id        string
	Migrated  bool
	AppliedAt time.Time
}
