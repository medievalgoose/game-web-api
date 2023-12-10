package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

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

	// Genre Routes
	router.GET("/genres/", getGenres)
	router.GET("/genres/:genreName", getListOfGamesByGenre)
	router.POST("/genres/", postGenre)
	router.PUT("/genres/", updateGenre)

	// Platform Routes
	router.GET("/platforms/", getPlatforms)
	router.DELETE("/platforms/:id/delete", deletePlatform)
	router.GET("/platforms/:id/games", getGamesByPlatform)
	router.POST("/platforms/", postPlatform)
	router.PUT("/platforms/", updatePlatform)

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

func getPlatforms(c *gin.Context) {
	var allPlatforms []platform

	db := openSqlConnection()
	defer db.Close()

	getAllPlatformsQuery := "SELECT * FROM platforms;"
	rows, err := db.Query(getAllPlatformsQuery)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var newPlatform platform

		err := rows.Scan(&newPlatform.ID, &newPlatform.Name)
		if err != nil {
			log.Fatal(err)
		}

		allPlatforms = append(allPlatforms, newPlatform)
	}

	c.JSON(http.StatusOK, allPlatforms)
}

func deletePlatform(c *gin.Context) {
	platformId := c.Param("id")

	db := openSqlConnection()
	defer db.Close()

	checkPlatformValidityQuery := "SELECT id FROM platforms WHERE id = $1;"
	row := db.QueryRow(checkPlatformValidityQuery, platformId)

	if err := row.Err(); err != nil {
		log.Fatal(err)
	}

	deletePlatformQuery := "DELETE FROM platforms WHERE id = $1 RETURNING *;"
	res := db.QueryRow(deletePlatformQuery, platformId)

	if err := res.Err(); err != nil {
		log.Fatal(err)
	}

	var deletedPlatformInfo platform

	if err := res.Scan(&deletedPlatformInfo.ID, &deletedPlatformInfo.Name); err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, deletedPlatformInfo)
}

func getGamesByPlatform(c *gin.Context) {
	platformId := c.Param("id")

	db := openSqlConnection()
	defer db.Close()

	checkPlatformQuery := "SELECT id FROM platforms WHERE id = $1;"
	row := db.QueryRow(checkPlatformQuery, platformId)
	if err := row.Err(); err != nil {
		log.Fatal(err)
	}

	getGamesBasedOnPlatformQuery := `
		SELECT g.id, g.name, g.price, g.genre_id, gr.id, gr.name FROM games_platforms gp 
		JOIN games g ON gp.game_id = g.id 
		JOIN platforms p ON gp.platform_id = p.id
		JOIN genres gr ON g.genre_id = gr.id
		WHERE platform_id = $1;
	`

	rows, err := db.Query(getGamesBasedOnPlatformQuery, platformId)
	if err != nil {
		log.Fatal(err)
	}

	var relevantGamesList []game

	for rows.Next() {
		var relevantGame game

		err := rows.Scan(&relevantGame.ID, &relevantGame.Name, &relevantGame.Price, &relevantGame.GenreID, &relevantGame.Genre.ID, &relevantGame.Genre.Name)
		if err != nil {
			log.Fatal(err)
		}

		relevantGamesList = append(relevantGamesList, relevantGame)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, relevantGamesList)
}

func postPlatform(c *gin.Context) {
	var newPlatform platform

	err := c.BindJSON(&newPlatform)
	if err != nil {
		log.Fatal(err)
	}

	db := openSqlConnection()
	defer db.Close()

	insertPlatformQuery := "INSERT INTO platforms (name) VALUES ($1);"
	_, err = db.Exec(insertPlatformQuery, newPlatform.Name)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfuly created new platform data."})
}

func updatePlatform(c *gin.Context) {
	var updatedPlatform platform

	err := c.BindJSON(&updatedPlatform)
	if err != nil {
		log.Fatal(err)
	}

	db := openSqlConnection()
	defer db.Close()

	checkPlatformValidityQuery := "SELECT id FROM platforms WHERE id = $1;"
	row := db.QueryRow(checkPlatformValidityQuery, updatedPlatform.ID)
	if err = row.Err(); err != nil {
		log.Fatal(err)
	}

	updatePlatformQuery := `
		UPDATE platforms
		SET name = $1
		WHERE id = $2
		RETURNING *;
	`

	var updatedPlatformInfo platform

	err = db.QueryRow(updatePlatformQuery, updatedPlatform.Name, updatedPlatform.ID).Scan(&updatedPlatformInfo.ID, &updatedPlatformInfo.Name)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, updatedPlatformInfo)
}
