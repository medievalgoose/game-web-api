package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type game struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Genre string `json:"genre"`
	Price int    `json:"price"`
}

var sqlGamesData []game

// Const containing the database data.
const (
	db_host  = "localhost"
	port     = 5432
	user     = "postgres"
	password = "root"
	db_name  = "GameWebApi"
)

func main() {
	router := gin.Default()
	router.GET("/games", getGames)
	router.POST("/games", postGames)

	// colon indicates a path parameter.
	router.GET("/games/:id", getGamesById)
	router.Run("localhost:8080")
}

func openSqlConnection() *sql.DB {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", db_host, port, user, password, db_name)

	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()

	if err != nil {
		log.Fatal(err)
	}

	return db
}

func getGames(c *gin.Context) {

	db := openSqlConnection()
	sqlGamesData = listAllGamesData(db)
	defer db.Close()

	name := c.Query("name")

	if name != "" {
		for _, game := range sqlGamesData {
			if strings.Contains(strings.ToLower(game.Name), strings.ToLower(name)) {
				c.IndentedJSON(http.StatusOK, game)
				return
			}
		}
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
		INSERT INTO games (name, genre, price)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	db := openSqlConnection()
	defer db.Close()

	id := 0

	err := db.QueryRow(insertStatement, newGame.Name, newGame.Genre, newGame.Price).Scan(&id)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Newly created id: %v", id)

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
	selectionQuery := "SELECT * FROM games;"

	rows, err := db.Query(selectionQuery)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var dbGamesData []game

	for rows.Next() {
		var newGame game

		if err := rows.Scan(&newGame.ID, &newGame.Name, &newGame.Genre, &newGame.Price); err != nil {
			log.Fatal(err)
		}

		dbGamesData = append(dbGamesData, newGame)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	db.Close()

	return dbGamesData
}
