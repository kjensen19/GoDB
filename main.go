package main

import (
	"context"
	"database/sql"
	"fmt"

	// "log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

func main() {
	//Environment Variables are set in cmd line with $export KEY:VALUE, then
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	var greeting string
	err = dbpool.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(2)
	}

	fmt.Printf(greeting)
	albums, err := albumsByArtist("John Coltrane", dbpool)
	if err != nil {
		fmt.Fprintf(os.Stderr, "find album failed: %v\n", err)
		os.Exit(3)
	}
	fmt.Printf("Albums found: %v\n", albums)

	alb, err := albumByID(2, dbpool)
	if err != nil {
		fmt.Fprintf(os.Stderr, "album by id failed: %v\n", err)
		os.Exit(3)
	}
	fmt.Printf("Album found: %v\n", alb)
}

// albumsByArtist queries for albums that have the specified artist name.
func albumsByArtist(name string, dbpool *pgxpool.Pool) ([]Album, error) {
	// An albums slice to hold data from returned rows.
	var albums []Album

	rows, err := dbpool.Query(context.Background(), "select * from album where artist=$1", name)
	if err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
		}
		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	return albums, nil
}

func albumByID(id int64, dbpool *pgxpool.Pool) (Album, error) {
	var alb Album

	row := dbpool.QueryRow(context.Background(), "SELECT * FROM album WHERE id=$1", id)
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		if err == sql.ErrNoRows {
			return alb, fmt.Errorf("AlbumsByID %d: no such album", id)
		}
		return alb, fmt.Errorf("AlbumsByID %d: %v", id, err)
	}
	return alb, nil
}
