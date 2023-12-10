package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type game struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	GenreID int    `json:"genre_id"`
	Price   int    `json:"price"`
	Genre   genre
}

type genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type platform struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
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

	// Game Routes
	router.GET("/games", getGames)
	router.POST("/games", postGames)
	// colon indicates a path parameter.
	router.GET("/games/:id", getGamesById)
	router.PUT("/games/", updateGame)
	router.DELETE("/games/:id/delete", deleteGame)

	// Genre Routes
	router.GET("/genres/", getGenres)
	router.GET("/genres/:genreName", getListOfGamesByGenre)
	router.POST("/genres/", postGenre)
	router.PUT("/genres/", updateGenre)
	router.DELETE("/genres/:id/delete", deleteGenre)

	// Platform Routes
	router.GET("/platforms/", getPlatforms)
	router.GET("/platforms/:id/games", getGamesByPlatform)
	router.POST("/platforms/", postPlatform)
	router.PUT("/platforms/", updatePlatform)
	router.DELETE("/platforms/:id/delete", deletePlatform)

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
