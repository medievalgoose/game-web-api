package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getGenres(c *gin.Context) {
	db := openSqlConnection()
	defer db.Close()
	allGenresData := listAllGenres(db)
	c.IndentedJSON(http.StatusOK, allGenresData)

}

func listAllGenres(db *sql.DB) []genre {
	listAllGenreQuery := "SELECT * FROM genres;"

	rows, err := db.Query(listAllGenreQuery)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var allGenres []genre

	for rows.Next() {
		var currentGenre genre

		err := rows.Scan(&currentGenre.ID, &currentGenre.Name)
		if err != nil {
			log.Fatal(err)
		}

		allGenres = append(allGenres, currentGenre)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	db.Close()
	return allGenres
}

func postGenre(c *gin.Context) {
	var newGenre genre

	if err := c.BindJSON(&newGenre); err != nil {
		log.Fatal(err)
	}

	insertGenreQuery := "INSERT INTO genres (name) VALUES ($1);"

	db := openSqlConnection()

	// res := 0
	result, err := db.Exec(insertGenreQuery, newGenre.Name)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
	// fmt.Printf("Newly created genre ID: %v", res)

	defer db.Close()
}

func updateGenre(c *gin.Context) {
	var updatedGenre genre

	if err := c.BindJSON(&updatedGenre); err != nil {
		log.Fatal(err)
	}

	checkGenreQuery := "SELECT id FROM genres WHERE id = $1;"
	Id := 0

	db := openSqlConnection()
	err := db.QueryRow(checkGenreQuery, updatedGenre.ID).Scan(&Id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Genre not found"})
		return
	}

	updateGenreQuery := "UPDATE genres SET name = $1 WHERE id = $2;"
	_, err = db.Exec(updateGenreQuery, updatedGenre.Name, updatedGenre.ID)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
}

func deleteGenre(c *gin.Context) {
	genreId := c.Param("id")

	db := openSqlConnection()
	defer db.Close()

	checkGenreValidityQuery := "SELECT id FROM genres WHERE id = $1;"
	res := db.QueryRow(checkGenreValidityQuery, genreId)
	if err := res.Err(); err != nil {
		log.Fatal(err)
	}

	deleteGenreQuery := "DELETE FROM genres WHERE id = $1 RETURNING *;"
	var deletedGenre genre
	err := db.QueryRow(deleteGenreQuery, genreId).Scan(&deletedGenre.ID, &deletedGenre.Name)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, deletedGenre)
}

func getListOfGamesByGenre(c *gin.Context) {
	requestedGenre := c.Param("genreName")

	var relevantGameList []game

	db := openSqlConnection()
	defer db.Close()

	selectRelevantGamesQuery := "SELECT g.id, g.name, g.price, g.genre_id, gr.id, gr.name FROM games g JOIN genres gr ON g.genre_id = gr.id WHERE LOWER(gr.name) = LOWER($1);"
	rows, err := db.Query(selectRelevantGamesQuery, requestedGenre)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var relevantGame game

		err := rows.Scan(&relevantGame.ID, &relevantGame.Name, &relevantGame.Price, &relevantGame.GenreID, &relevantGame.Genre.ID, &relevantGame.Genre.Name)
		if err != nil {
			log.Fatal(err)
		}

		relevantGameList = append(relevantGameList, relevantGame)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, relevantGameList)
}
