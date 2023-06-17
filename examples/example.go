package main

import (
	"database/sql"
	"fmt"
	"github.com/derivatan/si"
	"strings"
	// _ "github.com/lib/pq"
)

func main() {
	// Db connection
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=public sslmode=disable")
	if err != nil {
		panic(err)
	}

	// Configuration
	si.InitSecretIngredient(si.NewSQLDB(db))
	si.AddModel[contact]("contacts")
	si.AddModel[artist]("artists")
	si.AddModel[album]("albums")

	// Get all albums
	fmt.Println("\nExample 1 - Get all albums")
	albums, err := si.Query[album]().OrderBy("year", true).Get()
	if err != nil {
		panic(err)
	}
	for _, album := range albums {
		fmt.Println(album.Year, album.Name)
	}

	// Get all artists along with their albums. Only two queries will be executed, no matter how many results the first query have.
	fmt.Println("\nExample 2 - Get all artists with albums")
	artists, err := si.Query[artist]().With(func(m artist, r []artist) error {
		return m.Albums().Execute(r)
	}).Get()
	if err != nil {
		panic(err)
	}
	for _, artist := range artists {
		var albumNames []string
		for _, album := range artist.Albums().MustGet() {
			albumNames = append(albumNames, album.Name)
		}
		fmt.Println(artist.Name, ": ", strings.Join(albumNames, ", "))
	}

	// Find will return one entity instead of a list.
	fmt.Println("\nExample 3 - Find one artist")
	pinkFloyd, err := si.Query[artist]().Where("name", "=", "Pink Floyd").Find()
	if err != nil {
		panic(err)
	}
	fmt.Println(pinkFloyd.ID)

	// Relations can also be fetched on demand.
	fmt.Println("\nExample 4 - Find a relation from artist.")
	pinkFloydContact, err := pinkFloyd.Contact().Find()
	if err != nil {
		panic(err)
	}
	fmt.Println(pinkFloydContact.Email)

	// Save pink floyd contacts with a new phone number.
	fmt.Println("\nExample 5 - Save")
	pinkFloydContact.Phone = "321"
	err = si.Save(pinkFloydContact)
	if err != nil {
		panic(err)
	}

}
