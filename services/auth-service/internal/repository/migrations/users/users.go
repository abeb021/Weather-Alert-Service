package users

import (
    "embed"

    "github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/database/postgres"
    "github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed *.sql
var files embed.FS

func Run(dbURL string) error {
    source, err := iofs.New(files, ".")
    if err != nil {
        return err
    }

    m, err := migrate.NewWithSourceInstance("iofs", source, dbURL)
    if err != nil {
        return err
    }
    defer m.Close()

    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        return err
    }
    return nil
}