// Package ymirblog implement data to object struct and object-relational mapping.
// using ent framework makes it easy to build and maintain with large data-models.
// Schema As Code - model any database schema as Go objects.
// Multi Storage Driver - supports MySQL, MariaDB, TiDB, PostgresSQL, CockroachDB, SQLite and Gremlin.
// # This manifest was generated by ymir. DO NOT EDIT.
package ymirblog

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"ariga.io/atlas/sql/migrate"
	entDialect "entgo.io/ent/dialect"
	sqlDialect "entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/schema"
	"github.com/rs/zerolog/log"

	"gitlab.playcourt.id/dedenurr12/ymirblog/pkg/persist/ymirblog/ent"

	"gitlab.playcourt.id/dedenurr12/ymirblog/pkg/infrastructure"
)

var (
	// ErrNotFound is the error returned by driver if a resource cannot be found.
	ErrNotFound = errors.New("not found")

	// ErrAlreadyExists is the error returned by driver if a resource ID is taken during a creation.
	ErrAlreadyExists = errors.New("ID already exists")
)

// Database is data of instances.
type Database struct {
	*ent.Client
	dialect   string
	driver    *sqlDialect.Driver
	txOptions *sql.TxOptions
}

// WithDriver sets sql driver and dialect.
func WithDriver(s *sqlDialect.Driver, dialect string) func(*Database) {
	return func(d *Database) {
		d.driver = s
		d.dialect = dialect
	}
}

// WithTxIsolationLevel sets correct isolation level for database transactions.
func WithTxIsolationLevel(level sql.IsolationLevel) func(*Database) {
	return func(h *Database) {
		h.txOptions = &sqlDialect.TxOptions{Isolation: level}
	}
}

// Driver - wrapper sql driver to ent client.
func Driver(cfg ...func(database *Database)) *Database {
	db := &Database{}
	for _, f := range cfg {
		f(db)
	}
	opts := make([]ent.Option, 0)
	if strings.EqualFold(infrastructure.Envs.App.Environment, infrastructure.Development) {
		driverWithDebugContext := entDialect.DebugWithContext(db.driver,
			func(ctx context.Context, i ...any) {
				log.Debug().Str("query", fmt.Sprintf("%v", i)).Msg("driverWithDebugContext")
			})
		opts = append(opts, ent.Driver(driverWithDebugContext))
	} else {
		opts = append(opts, ent.Driver(db.driver))
	}
	db.Client = ent.NewClient(opts...)

	temp := filepath.Join(os.TempDir(), "migrations")
	if err := embedWriteTemp(temp, db.dialect); err != nil {
		log.Error().Err(err).Msg("migrate schema is failed")
		//nolint: revive
		os.Exit(1)
	}

	dir, err := migrate.NewLocalDir(temp)
	if err != nil {
		log.Error().Err(err).Msg("migrate schema is failed")
		//nolint: revive
		os.Exit(1)
	}

	options := []schema.MigrateOption{
		schema.WithDir(dir),
		schema.WithDropColumn(true),
		schema.WithDropIndex(true),
	}
	ctx := context.Background()
	if err = db.Client.Schema.Create(ctx, options...); err != nil {
		log.Error().Err(err).Msg("migrate schema is failed")
		//nolint: revive
		os.Exit(1)
	}
	if err = os.RemoveAll(temp); err != nil {
		log.Error().Err(err).Msg("failed to remove temporary directory for database migrations")
	}

	return db
}

// WithTransaction is a wrapper to begin transaction with defined options.
func (d *Database) WithTransaction(ctx context.Context, fn func(ctx context.Context, tx *ent.Tx) error) error {
	tx, err := d.Client.BeginTx(ctx, d.txOptions)
	if err != nil {
		return err
	}

	if err = fn(ctx, tx); err != nil {
		if errRoll := tx.Rollback(); errRoll != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, errRoll)
		}
		return d.ConvertDBError("execute:", err)
	}
	if err = tx.Commit(); err != nil {
		if errRoll := tx.Rollback(); errRoll != nil {
			return d.ConvertDBError("rollback:", errRoll)
		}
		return d.ConvertDBError("commit is failed", err)
	}
	return nil
}

// BeginTx is a wrapper to begin transaction with defined options.
func (d *Database) BeginTx(ctx context.Context) (*ent.Tx, error) {
	return d.Client.BeginTx(ctx, d.txOptions)
}

// ConvertDBError set wrapper db error.
func (*Database) ConvertDBError(t string, err error) error {
	if ent.IsNotFound(err) {
		return ErrNotFound
	}

	if ent.IsConstraintError(err) {
		return ErrAlreadyExists
	}

	return fmt.Errorf(t, err)
}