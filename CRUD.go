package main

import (
	"database/sql"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Planta struct {
	ID             int    `json:"id_planta"`
	NomeCientifico string `json:"nome_cientifico"`
	NomePopular    string `json:"nome_popular"`
}

var db *sql.DB

func main() {
	var err error

	// Adjust username:password as needed
	db, err = sql.Open("mysql", "root:ceub123456@tcp(localhost:3306)/bd_planta")
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	r.Use(cors.Default())

	// CREATE
	r.POST("/plantas", func(c *gin.Context) {
		var p Planta
		if err := c.BindJSON(&p); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		result, err := db.Exec(
			"INSERT INTO PLANTAS (nome_cientifico, nome_popular) VALUES (?, ?)",
			p.NomeCientifico, p.NomePopular,
		)

		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		id, _ := result.LastInsertId()
		p.ID = int(id)

		c.JSON(201, p)
	})

	// READ ALL
	r.GET("/plantas", func(c *gin.Context) {
		rows, err := db.Query("SELECT id_planta, nome_cientifico, nome_popular FROM PLANTAS")
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var plantas []Planta
		for rows.Next() {
			var p Planta
			rows.Scan(&p.ID, &p.NomeCientifico, &p.NomePopular)
			plantas = append(plantas, p)
		}

		c.JSON(200, plantas)
	})

	// UPDATE
	r.PUT("/plantas/:id", func(c *gin.Context) {
		id := c.Param("id")
		var p Planta

		if err := c.BindJSON(&p); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		_, err := db.Exec(
			"UPDATE PLANTAS SET nome_cientifico=?, nome_popular=? WHERE id_planta=?",
			p.NomeCientifico, p.NomePopular, id,
		)

		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"updated": id})
	})

	// DELETE
	r.DELETE("/plantas/:id", func(c *gin.Context) {
		id := c.Param("id")

		_, err := db.Exec("DELETE FROM PLANTAS WHERE id_planta=?", id)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.Status(204)
	})

	r.Run(":8080")
}
