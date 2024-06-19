package postgres

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lindell/go-stderrs/stderrs"
	"github.com/lindell/sr-restored/domain"
	"github.com/pkg/errors"
)

type Config struct {
	PostgresURL string
}

type DB struct {
	db *pgxpool.Pool
	sq squirrel.StatementBuilderType
}

func NewDB(config Config) (*DB, error) {
	conn, err := pgxpool.New(context.Background(), config.PostgresURL)

	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	if err != nil {
		return nil, err
	}

	return &DB{
		db: conn,
		sq: sq,
	}, nil
}

func (db *DB) GetProgram(ctx context.Context, programID int) (domain.Program, error) {
	program, err := db.getProgram(ctx, programID)
	if err != nil {
		return domain.Program{}, errors.WithMessage(err, "could not fetch program from db")
	}

	program.Episodes, err = db.getEpisodes(ctx, programID)
	if err != nil {
		return domain.Program{}, errors.WithMessage(err, "could not fetch episodes from db")
	}

	return program, nil
}

func (db *DB) getEpisodes(ctx context.Context, programID int) ([]domain.Episode, error) {
	rows, err := db.db.Query(ctx, `
		SELECT
			id,
			program_id,
			title,
			description,
			url,
			publish_date,
			image_url,
			content_type,
			file_url,
			file_duration,
			file_bytes
		FROM episodes WHERE program_id = $1`, programID)
	if err != nil {
		return nil, errors.WithMessage(err, "could not select episodes")
	}

	episodes := []domain.Episode{}
	for rows.Next() {
		episode := &domain.Episode{}
		err := rows.Scan(
			&episode.ID,
			&episode.ProgramID,
			&episode.Title,
			&episode.Description,
			&episode.URL,
			&episode.PublishDate,
			&episode.ImageURL,
			&episode.ContentType,
			&episode.FileURL,
			&episode.FileDurationSeconds,
			&episode.FileBytes,
		)
		if err != nil {
			return nil, err
		}

		episodes = append(episodes, *episode)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return episodes, nil
}

func (db *DB) InsertEpisodes(ctx context.Context, episodes []domain.Episode) error {
	if len(episodes) == 0 {
		return nil
	}

	query := db.sq.Insert("episodes").Columns(
		"id",
		"program_id",
		"title",
		"description",
		"url",
		"publish_date",
		"image_url",
		"content_type",
		"file_url",
		"file_duration",
		"file_bytes",
	).Suffix("ON CONFLICT (id) DO NOTHING")

	for _, episode := range episodes {
		query = query.Values(
			episode.ID,
			episode.ProgramID,
			episode.Title,
			episode.Description,
			episode.URL,
			episode.PublishDate,
			episode.ImageURL,
			episode.ContentType,
			episode.FileURL,
			episode.FileDurationSeconds,
			episode.FileBytes,
		)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return errors.WithMessage(err, "could not build query")
	}

	tag, err := db.db.Exec(ctx, sql, args...)
	log(tag, "insert episodes")
	return err
}

func (db *DB) getProgram(ctx context.Context, programID int) (domain.Program, error) {
	rows, err := db.db.Query(ctx, `
		SELECT
			id,
			name,
			description,
			email,
			copyright,
			url,
			image_url
		FROM programs WHERE id = $1`, programID)
	if err != nil {
		return domain.Program{}, errors.WithMessage(err, "could not select episodes")
	}

	if !rows.Next() {
		return domain.Program{}, stderrs.NewNotFound("could not find program")
	}

	program := domain.Program{}
	err = rows.Scan(
		&program.ID,
		&program.Name,
		&program.Description,
		&program.Email,
		&program.Copyright,
		&program.URL,
		&program.ImageURL,
	)
	if err != nil {
		return domain.Program{}, err
	}

	return program, nil
}

func (db *DB) InsertProgram(ctx context.Context, program domain.Program) error {
	query := db.sq.Insert("programs").Columns(
		"id",
		"name",
		"description",
		"email",
		"copyright",
		"url",
		"image_url",
	).Values(
		program.ID,
		program.Name,
		program.Description,
		program.Email,
		program.Copyright,
		program.URL,
		program.ImageURL,
	).Suffix("ON CONFLICT (id) DO NOTHING")

	sql, args, err := query.ToSql()
	if err != nil {
		return errors.WithMessage(err, "could not build query")
	}

	tag, err := db.db.Exec(ctx, sql, args...)
	log(tag, "insert program")
	return err
}
