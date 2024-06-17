package postgres

import (
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgconn"
)

func log(tag pgconn.CommandTag, msg string) {
	slog.Info(fmt.Sprintf("postgres query: %s", msg),
		"status", tag.String(),
		"rows-affected", tag.RowsAffected(),
	)
}
