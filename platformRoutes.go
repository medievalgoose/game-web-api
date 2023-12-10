package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
