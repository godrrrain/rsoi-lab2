package storage

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Library struct {
	ID          int    `json:"id"`
	Library_uid string `json:"library_uid"`
	Name        string `json:"name"`
	City        string `json:"city"`
	Address     string `json:"address"`
}

type Book struct {
	ID              int    `json:"id"`
	Book_uid        string `json:"book_uid"`
	Name            string `json:"name"`
	Author          string `json:"author"`
	Genre           string `json:"genre"`
	Condition       string `json:"condition"`
	Available_count int    `json:"available_count"`
}

type Storage interface {
	// Insert(e *Person)
	// Get(id int) (Person, error)
	// Update(e *Person) error
	// Delete(id int) error
	// GetAll() []Person
	GetLibrariesByCity(ctx context.Context, city string) ([]Library, error)
	GetBooksByLibraryUid(ctx context.Context, libraryUid string) ([]Book, error)
}

type postgres struct {
	db *pgxpool.Pool
}

func NewPgStorage(ctx context.Context, connString string) (*postgres, error) {
	var pgInstance *postgres
	var pgOnce sync.Once
	pgOnce.Do(func() {
		db, err := pgxpool.New(ctx, connString)
		if err != nil {
			fmt.Printf("Unable to create connection pool: %v\n", err)
			return
		}

		pgInstance = &postgres{db}
	})

	return pgInstance, nil
}

func (pg *postgres) Ping(ctx context.Context) error {
	return pg.db.Ping(ctx)
}

func (pg *postgres) Close() {
	pg.db.Close()
}

func (pg *postgres) GetLibrariesByCity(ctx context.Context, city string) ([]Library, error) {
	query := fmt.Sprintf(`SELECT id, library_uid, name, city, address FROM library WHERE city = '%s'`, city)

	rows, err := pg.db.Query(ctx, query)

	var libraries []Library

	if err != nil {
		return libraries, fmt.Errorf("unable to query: %w", err)
	}
	defer rows.Close()

	libraries, err = pgx.CollectRows(rows, pgx.RowToStructByName[Library])
	if err != nil {
		fmt.Printf("CollectRows error: %v", err)
		return libraries, err
	}

	return libraries, nil
}

func (pg *postgres) GetBooksByLibraryUid(ctx context.Context, libraryUid string) ([]Book, error) {
	query := fmt.Sprintf(`SELECT books.*, library_books.available_count from library_books, books, library where library.library_uid = '%s'`, libraryUid)

	rows, err := pg.db.Query(ctx, query)

	var books []Book

	if err != nil {
		return books, fmt.Errorf("unable to query: %w", err)
	}
	defer rows.Close()

	books, err = pgx.CollectRows(rows, pgx.RowToStructByName[Book])
	if err != nil {
		fmt.Printf("CollectRows error: %v", err)
		return books, err
	}

	return books, nil
}
