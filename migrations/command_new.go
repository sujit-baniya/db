package migrations

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"
	"time"
)

var templateContent = `
-- +migrate Up

-- +migrate Down
`
var tpl *template.Template

func init() {
	tpl = template.Must(template.New("new_migration").Parse(templateContent))
}

type NewCommand struct {
}

func (c *NewCommand) Help() string {
	helpText := `
Usage: verify-rest new [options] name

  Create a new a database migration.

Options:

  -config=dbconfig.yml   Configuration file to use.
  -env="development"     Environment.
  name                   The name of the migration
`
	return strings.TrimSpace(helpText)
}

func (c *NewCommand) Synopsis() string {
	return "Create a new migration"
}

func (c *NewCommand) Run(args []string) int {
	cmdFlags := flag.NewFlagSet("new", flag.ContinueOnError)
	cmdFlags.Usage = func() { ui.Output(c.Help()) }

	if len(args) < 1 {
		err := errors.New("A name for the migration is needed")
		ui.Error(err.Error())
		return 1
	}

	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	if err := CreateMigration(cmdFlags.Arg(0)); err != nil {
		ui.Error(err.Error())
		return 1
	}
	return 0
}

func CreateMigration(name string) error {
	name = strings.ToLower(name)
	if _, err := os.Stat(DefaultConfig.Dir); os.IsNotExist(err) {
		return err
	}
	query := SmartMigration(name)
	if query != "" {
		tpl = template.Must(template.New("new_migration").Parse(query))
	}
	fileName := fmt.Sprintf("%s-%s.sql", time.Now().Format("20060102150405"), strings.TrimSpace(name))
	pathName := path.Join(DefaultConfig.Dir, fileName)
	f, err := os.Create(pathName)

	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	if err := tpl.Execute(f, nil); err != nil {
		return err
	}

	ui.Output(fmt.Sprintf("Created migration %s", pathName))
	return nil
}

func SmartMigration(migrationName string) string {
	nameParts := strings.Split(migrationName, `_`)
	upQuery := ""
	downQuery := ""
	if nameParts[len(nameParts)-1] == "table" {
		switch nameParts[0] {
		case "create":
			tableName := strings.Join(nameParts[1:(len(nameParts)-1)], `_`)
			createSequence := "CREATE SEQUENCE IF NOT EXISTS " + tableName + "_id_seq;"
			upQuery = createSequence + "CREATE TABLE IF NOT EXISTS " + tableName + `
(
id int8 NOT NULL DEFAULT nextval('` + tableName + `_id_seq'::regclass) PRIMARY KEY, 
is_active bool default false,
created_at timestamptz,
updated_at timestamptz,
deleted_at timestamptz
)` + ";"
			dropSequenceQuery := "DROP SEQUENCE IF EXISTS " + tableName + "_seq;"
			downQuery = dropSequenceQuery + "DROP TABLE IF EXISTS " + tableName + ";"
		case "drop":
			tableName := strings.Join(nameParts[1:(len(nameParts)-1)], `_`)
			dropSequenceQuery := "DROP SEQUENCE IF EXISTS " + tableName + "_seq;"
			createSequence := "CREATE SEQUENCE IF NOT EXISTS " + tableName + "_id_seq;"
			upQuery = dropSequenceQuery + "DROP TABLE IF EXISTS " + tableName + ";"
			downQuery = createSequence + "CREATE TABLE IF NOT EXISTS " + tableName + `
(
id int8 NOT NULL DEFAULT nextval('` + tableName + `_id_seq'::regclass) PRIMARY KEY, 
is_active bool default false,
created_at timestamptz,
updated_at timestamptz,
deleted_at timestamptz
)` + ";"
		case "add":
			for i, part := range nameParts {
				if part == "in" {
					field := strings.Join(nameParts[1:i], `_`)
					tableName := strings.Join(nameParts[(i+1):(len(nameParts)-1)], `_`)
					upQuery = "ALTER TABLE " + tableName + " ADD COLUMN " + field + " VARCHAR(200)" + ";"
					downQuery = "ALTER TABLE " + tableName + " DROP COLUMN " + field + ";"
					break
				}
			}
		case "remove":
			for i, part := range nameParts {
				if part == "from" {
					field := strings.Join(nameParts[1:i], `_`)
					tableName := strings.Join(nameParts[(i+1):(len(nameParts)-1)], `_`)
					upQuery = "ALTER TABLE " + tableName + " DROP COLUMN " + field + ";"
					downQuery = "ALTER TABLE " + tableName + " ADD COLUMN " + field + " VARCHAR(200)" + ";"
					break
				}
			}
		case "rename":
			for i, part := range nameParts {
				if part == "in" {
					oldTableName := strings.Join(nameParts[1:i], `_`)
					newTableName := strings.Join(nameParts[(i+1):(len(nameParts)-1)], `_`)
					upQuery = "ALTER TABLE " + oldTableName + " RENAME TO " + newTableName + ";"
					downQuery = "ALTER TABLE " + newTableName + " RENAME TO " + oldTableName + ";"
					break
				}
			}
		case "alter", "change":
			for i, part := range nameParts {
				if part == "in" {
					field := strings.Join(nameParts[1:i], `_`)
					tableName := strings.Join(nameParts[(i+1):(len(nameParts)-1)], `_`)
					upQuery = "ALTER TABLE " + tableName + " ALTER COLUMN " + field + " VARCHAR(200)" + ";"
					downQuery = "ALTER TABLE " + tableName + " ALTER COLUMN " + field + " VARCHAR(200)" + ";"
					break
				}
			}
		}
	}
	query := fmt.Sprintf(`
-- +migrate Up
%s

-- +migrate Down
%s
`, upQuery, downQuery)
	return query
}
