package slogqlite

import (
	"context"
	"database/sql"
	"io"
	"log/slog"
	"strings"

	_ "embed"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema.sql
var schema string

type SqliteHandler struct {
	db      *sql.DB
	handler slog.Handler
}

func (s *SqliteHandler) Enabled(c context.Context, l slog.Level) bool {
	return s.handler.Enabled(c, l)
}
func (s *SqliteHandler) Handle(c context.Context, r slog.Record) error {
	tx, err := s.db.BeginTx(c, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.Exec("insert into logs (level, msg, time) values ($1, $2, $3)", r.Level.String(), r.Message, r.Time)
	if err != nil {
		return err
	}
	r.Attrs(func(a slog.Attr) bool {
		_, err = tx.Exec("insert into log_attrs (log_id, key, value) values ((select id from logs where time = $1) ,$2, $3)", r.Time, a.Key, a.Value.String())
		return err == nil
	})
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return s.handler.Handle(c, r)
}
func (s *SqliteHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &SqliteHandler{s.db, s.handler.WithAttrs(attrs)}
}
func (s *SqliteHandler) WithGroup(name string) slog.Handler {
	return &SqliteHandler{s.db, s.handler.WithGroup(name)}
}

func New(w io.Writer, dbConn string, opts *slog.HandlerOptions) (*SqliteHandler, error) {
	handler := &SqliteHandler{handler: slog.NewTextHandler(w, opts)}
	var err error
	handler.db, err = sql.Open("sqlite3", dbConn)
	if err != nil {
		return handler, err
	}
	for _, table := range strings.Split(schema, ";\n") {
		_, err := handler.db.Exec(table)
		if err != nil {
			return handler, err
		}
	}
	return handler, nil
}
