package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getGames(c *gin.Context) {

	db := openSqlConnection()
	sqlGamesData = listAllGamesData(db)
	defer db.Close()

	name := c.Query("name")

	// if name != "" {
	// 	for _, game := range sqlGamesData {
	// 		if strings.Contains(strings.ToLower(game.Name), strings.ToLower(name)) {
	// 			c.IndentedJSON(http.StatusOK, game)
	// 			return
	// 		}
	// 	}
	// }

	if name != "" {
		var requestedGame game
		getOneGameQuery := "SELECT * FROM games WHERE LOWER(name) LIKE '%' || LOWER($1) || '%';"
		row := db.QueryRow(getOneGameQuery, name)
		row.Scan(&requestedGame.ID, &requestedGame.Name, &requestedGame.Price, &requestedGame.GenreID)

		if err := row.Err(); err != nil {
			log.Fatal(err)
		}

		if requestedGame.Name != "" {
			// c.IndentedJSON(http.StatusOK, requestedGame)
			c.JSON(http.StatusOK, requestedGame)
		} else {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Object not found"})
		}
		return
	}

	c.IndentedJSON(http.StatusOK, sqlGamesData)
}

func postGames(c *gin.Context) {
	var newGame game

	if err := c.BindJSON(&newGame); err != nil {
		fmt.Println("Hello world, this the error message")
		log.Fatal(err)
	}

	/*
		games = append(games, newGame)
		c.IndentedJSON(http.StatusCreated, newGame)
	*/

	// SQL INSERT STATEMENT
	insertStatement := `
		INSERT INTO games (name, price, genre_id)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	db := openSqlConnection()
	defer db.Close()

	id := 0

	err := db.QueryRow(insertStatement, newGame.Name, newGame.Price, newGame.GenreID).Scan(&id)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Newly created id: %v", id)

}

func updateGame(c *gin.Context) {
	var updatedGame game

	if err := c.BindJSON(&updatedGame); err != nil {
		log.Fatal(err)
	}

	checkGameExistQuery := "SELECT id FROM games WHERE id = $1;"

	db := openSqlConnection()
	defer db.Close()

	checkId := 0
	err := db.QueryRow(checkGameExistQuery, updatedGame.ID).Scan(&checkId)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Game not found"})
		// log.Fatal(err)
		return
	}

	updateGameQuery := "UPDATE games SET name = $1, price = $2, genre_id = $3 WHERE id = $4;"
	_, err = db.Exec(updateGameQuery, updatedGame.Name, updatedGame.Price, updatedGame.GenreID, updatedGame.ID)

	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Game data updated"})
}

func getGamesById(c *gin.Context) {

	db := openSqlConnection()
	sqlGamesData = listAllGamesData(db)
	defer db.Close()

	idInt, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		log.Fatal(err)
	}

	for _, game := range sqlGamesData {
		if game.ID == idInt {
			c.IndentedJSON(http.StatusOK, game)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Object not found"})
}

func listAllGamesData(db *sql.DB) []game {
	// selectionQuery := "SELECT * FROM games;"
	selectionQueryV2 := "SELECT g.id, g.name, price, g.genre_id, n.id, n.name AS \"genre\" FROM games g JOIN genres n ON g.genre_id = n.id; "

	rows, err := db.Query(selectionQueryV2)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var dbGamesData []game

	for rows.Next() {
		var newGame game

		if err := rows.Scan(&newGame.ID, &newGame.Name, &newGame.Price, &newGame.GenreID, &newGame.Genre.ID, &newGame.Genre.Name); err != nil {
			log.Fatal(err)
		}

		dbGamesData = append(dbGamesData, newGame)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	// db.Close()

	return dbGamesData
}

func deleteGame(c *gin.Context) {
	gameId := c.Param("id")

	db := openSqlConnection()
	defer db.Close()

	checkGameValidityQuery := "SELECT id FROM games WHERE id = $1;"
	res := db.QueryRow(checkGameValidityQuery, gameId)
	if err := res.Err(); err != nil {
		log.Fatal(err)
	}

	deleteGameQuery := "DELETE FROM games WHERE id = $1 RETURNING *;"

	var deletedGame game

	err := db.QueryRow(deleteGameQuery, gameId).Scan(&deletedGame.ID, &deletedGame.Name, &deletedGame.Price, &deletedGame.GenreID)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, deletedGame)
}
