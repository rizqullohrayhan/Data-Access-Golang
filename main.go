package main

import (
    "database/sql"
    "fmt"
    "log"

    "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type Album struct {
	ID		int64
	Title	string
	Artist	string
	Price	float32
}

func main()  {
	// Capture connection properties.
	cfg := mysql.Config{
		User:					"root",
		Passwd:					"",
		Net:					"tcp",
		Addr: 					"127.0.0.1:3306",
		DBName: 				"recordings",
		AllowNativePasswords: 	true,
	}
	// Get a dataabse handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	albums, err := albumsByArtist("John Coltrane")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Result from func albumsByArtist:")
	fmt.Printf("Albums found: %v\n", albums)

	// Hard-code ID 2 ere to test the query.
	alb, err := albumByID(2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Result from func albumsByID:")
	fmt.Printf("Album found: %v\n", alb)

	_, err = addAlbum(Album{
		Title: "The Modern Sund of Betty Carter",
		Artist: "Betty Carter",
		Price: 49.99,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Album succesfully added!")

	updated, err := editAlbum(2, Album{
		Title: "Small Steps",
		Artist: "John Coltrane",
		Price: 63.99,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d Album updated\n", updated)

	destroyed, err := destroyAlbum(1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d Album deleted\n", destroyed)
}

// albumsByArtist queries for albums that have
// the specified artist name.
func albumsByArtist(name string) ([]Album, error) {
	// An albums slice to hold data from returne rows.
	var albums []Album

	rows, err := db.Query("SELECT * FROM album WHERE artist = ?", name)
	if err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("albumByArtist %q: %v", name, err)
		}
		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	return albums, nil
}

// albumByID queries for the album with the specified ID.
func albumByID(id int64) (Album, error) {
	// An album to hold data from the returned row.
	var alb Album

	row := db.QueryRow("SELECT * FROM album WHERE id = ?", id)
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		if err == sql.ErrNoRows{
			return alb, fmt.Errorf("albumsById %d: no such album", id)
		}
		return alb, fmt.Errorf("albumsById %d: %v", id, err)
	}

	return alb, nil
}

// addAlbum adds the specified album to the database,
// returning the album ID of the new entry
func addAlbum(alb Album) (int64, error) {
	result, err := db.Exec("INSERT INTO album (title, artist, price) VALUES (?, ?, ?)", alb.Title, alb.Artist, alb.Price)
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	return id, nil
}

// editAlbum modifying the specified album from database
// based on ID album,
// returning the album ID of the effected row
func editAlbum(id int64, alb Album) (int64, error) {
	result, err := db.Exec("UPDATE album SET title=?, artist=?, price=? WHERE id=?", alb.Title, alb.Artist, alb.Price, id)
	if err != nil {
		return 0, fmt.Errorf("editAlbum: %v", err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("editAlbum: %v", err)
	}
	if count == 0 {
		return 0, fmt.Errorf("editAlbum %d: no such album", id)
	}
	return count, nil
}

// editAlbum modifying the specified album from database
// based on ID album,
// returning the album ID of the effected row
func destroyAlbum(id int64) (int64, error) {
	result, err := db.Exec("DELETE FROM album WHERE id=?", id)
	if err != nil {
		return 0, fmt.Errorf("editAlbum: %v", err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("editAlbum: %v", err)
	}
	if count == 0 {
		return 0, fmt.Errorf("editAlbum %d: no such album", id)
	}
	return count, nil
}