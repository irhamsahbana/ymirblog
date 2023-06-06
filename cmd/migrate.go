// Package cmd is the command surface of ymirblog cli tool provided by kubuskotak.
// # This manifest was generated by ymir. DO NOT EDIT.
package cmd

import (
	"entgo.io/ent/dialect"
	"github.com/spf13/cobra"
	ymirBlogSchema "gitlab.playcourt.id/dedenurr12/ymirblog/pkg/persist/ymirblog/diff"
)

type migrateOptions struct {
	Name    string
	Dialect string
	DSN     string
}

// Dialect is the func to convert dsn dialect.
func Dialect(d string) string {
	switch d {
	case dialect.Postgres:
		return dialect.Postgres
	case dialect.MySQL:
		return dialect.MySQL
	default:
		return dialect.SQLite
	}
}

func newMigrateCmd() *cobra.Command {
	m := &migrateOptions{}
	cmd := &cobra.Command{
		Use:   `migrate`,
		Short: "Print version info",
		RunE: func(cmd *cobra.Command, args []string) error {
			return m.Run(cmd, args)
		},
	}
	cmd.Flags().StringVarP(&m.Name, "filename", "n", "init", "migrate -n init")
	cmd.Flags().StringVarP(&m.Dialect, "dialect", "q", "sqlite3", "migrate -q mysql")
	cmd.Flags().StringVarP(&m.DSN, "data source name", "s", "sqlite://.data/blogs.migration.db?_fk=1", "migrate -s 'sqlite://.data/blogs.migration.db?_fk=1'")

	return cmd
}

// Run is ent func to migrate.
func (m *migrateOptions) Run(cmd *cobra.Command, _ []string) error {
	switch m.Dialect {
	case dialect.SQLite, dialect.MySQL:
		if err := ymirBlogSchema.SchemaMigrate(m.Name, Dialect(m.Dialect), m.DSN); err != nil {
			return err
		}
	default:
		return cmd.Usage()
	}
	return nil
}
