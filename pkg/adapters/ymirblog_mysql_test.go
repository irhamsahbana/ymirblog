// Package adapters are the glue between components and external sources.
// # This manifest was generated by ymir. DO NOT EDIT.
package adapters

import (
	"testing"

	sqlEnt "entgo.io/ent/dialect/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"gitlab.playcourt.id/dedenurr12/ymirblog/pkg/infrastructure"
)

func TestWithYmirBlogMySQL(t *testing.T) {
	is := assert.New(t)

	db, mock, err := sqlmock.New()
	if err != nil {
		is.Failf("failed to open stub db", "%v", err)
	}

	is.NotNil(db, "mock db is null")
	is.NotNil(mock, "sqlmock is null")

	YmirBlogMySQLOpen = func(dialect, source string) (*sqlEnt.Driver, error) {
		return sqlEnt.NewDriver(dialect, sqlEnt.Conn{ExecQuerier: db}), nil
	}

	infrastructure.Configuration(
		infrastructure.WithPath("../.."),
		infrastructure.WithFilename("config.yaml"),
	).Initialize()

	adapter := &Adapter{}
	adapter.Sync(
		WithYmirBlogMySQL(&YmirBlogMySQL{
			NetworkDB: NetworkDB{
				Database:    infrastructure.Envs.YmirBlogMySQL.Database,
				User:        infrastructure.Envs.YmirBlogMySQL.User,
				Password:    infrastructure.Envs.YmirBlogMySQL.Password,
				Host:        infrastructure.Envs.YmirBlogMySQL.Host,
				Port:        infrastructure.Envs.YmirBlogMySQL.Port,
				MaxIdleCons: infrastructure.Envs.DB.MaxIdleCons,
			},
		}),
	)

	mock.ExpectClose()

	// Asserts
	is.Nil(adapter.UnSync())
	is.Nil(mock.ExpectationsWereMet())
}

