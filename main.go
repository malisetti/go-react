package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "madras"
	DB_NAME     = "banks"
)

type Bank struct {
	Id   int64  `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

type Branch struct {
	Ifsc     string `db:"ifsc" json:"ifsc"`
	BankId   int64  `db:"bank_id" json:"bank_id"`
	Branch   string `db:"branch" json:"branch"`
	Address  string `db:"address" json:"address"`
	City     string `db:"city" json:"city"`
	District string `db:"district" json:"district"`
	State    string `db:"state" json:"state"`
}

func main() {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	checkErr(err)
	defer db.Close()

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	router.GET("/bank", func(c *gin.Context) {
		banks := []Bank{}
		rows, err := db.Query("SELECT * FROM banks")
		checkErr(err)
		defer rows.Close()

		for rows.Next() {
			bank := Bank{}
			err = rows.Scan(&bank.Name, &bank.Id)
			checkErr(err)
			banks = append(banks, bank)
		}

		c.JSON(200, banks)
	})

	router.GET("/city/:bank", func(c *gin.Context) {
		bank := c.Params.ByName("bank")
		//find the cities of a bank
		cities := []string{}
		rows, err := db.Query("SELECT DISTINCT city FROM branches WHERE bank_id=$1", bank)
		checkErr(err)
		defer rows.Close()

		for rows.Next() {
			var city string
			err = rows.Scan(&city)
			checkErr(err)
			cities = append(cities, city)
		}

		c.JSON(200, cities)
	})

	router.GET("/branch/:bank/:city", func(c *gin.Context) {
		bank := c.Params.ByName("bank")
		city := c.Params.ByName("city")
		branches := []Branch{}

		rows, err := db.Query("SELECT ifsc, bank_id, branch, address, city, district, state FROM bank_branches WHERE bank_id=$1 AND city=$2", bank, city)
		checkErr(err)
		defer rows.Close()

		for rows.Next() {
			branch := Branch{}
			err = rows.Scan(&branch.Ifsc, &branch.BankId, &branch.Branch, &branch.Address, &branch.City, &branch.District, &branch.State)
			checkErr(err)
			branches = append(branches, branch)
		}

		c.JSON(200, branches)
	})

	router.Run(":" + port)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
