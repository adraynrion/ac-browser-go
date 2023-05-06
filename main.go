package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

type Assassin struct {
	ID        uint64 `json:"id"`
	Name      string `json:"name"`
	Age       uint16 `json:"age"`
	Gender    string `json:"gender"`
	Speed     uint16 `json:"speed"`
	Strength  uint16 `json:"damage"`
	Busy      bool   `json:"busy"`
	CreatedAt uint64 `json:"created_at"`
	UpdatedAt uint64 `json:"updated_at"`
}

var db *sql.DB

func initEnv() {
	// Open the .env file
	file, err := os.Open(".env")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Read the .env file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Split each line on the equal sign
		parts := strings.Split(scanner.Text(), "=")
		if len(parts) != 2 {
			// Skip lines that don't have an equal sign
			continue
		}

		// Set the environment variable
		os.Setenv(parts[0], parts[1])
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Check that the environment variables are set
	if os.Getenv("MYSQL_USER") == "" {
		log.Fatal("MYSQL_USER must be set")
	}
	if os.Getenv("MYSQL_PASSWORD") == "" {
		log.Fatal("MYSQL_PASSWORD must be set")
	}
	if os.Getenv("MYSQL_DATABASE") == "" {
		log.Fatal("MYSQL_DATABASE must be set")
	}
}

func initDriver() {
	// connect to database
	cfg := mysql.Config{
		User:   os.Getenv("MYSQL_USER"),
		Passwd: os.Getenv("MYSQL_PASSWORD"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: os.Getenv("MYSQL_DATABASE"),
	}
	// get a database handle
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	// check the connection
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")
}

func main() {
	log.SetPrefix("main: ")
	// disable printing the time, source file, and line number
	// log.SetFlags(0)

	initEnv()
	initDriver()

	// create router
	router := gin.Default()
	router.GET("/assassins", getAssassins)
	router.GET("/assassins/:id", findAssassin)
	router.POST("/assassins", createAssassin)

	// start server
	router.Run(":8080")
}

func getAllAssassins() ([]Assassin, error) {
	var assassins []Assassin

	rows, err := db.Query("SELECT * FROM assassins")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var a Assassin
		if err := rows.Scan(
			&a.ID,
			&a.Name,
			&a.Age,
			&a.Gender,
			&a.Speed,
			&a.Strength,
			&a.Busy,
			&a.CreatedAt,
			&a.UpdatedAt,
		); err != nil {
			return nil, err
		}
		assassins = append(assassins, a)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return assassins, nil
}
func getAssassins(c *gin.Context) {
	var assassinList []Assassin
	var err error
	if assassinList, err = getAllAssassins(); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, assassinList)
}

func findAssassin(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	var a Assassin
	if err := db.QueryRow("SELECT * FROM assassins WHERE id = ?", id).Scan(
		&a.ID,
		&a.Name,
		&a.Age,
		&a.Gender,
		&a.Speed,
		&a.Strength,
		&a.Busy,
		&a.CreatedAt,
		&a.UpdatedAt,
	); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, a)
}

func createAssassin(c *gin.Context) {
	var isDef bool
	var assassinName string
	if assassinName, isDef = c.GetPostForm("name"); !isDef {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}
	var assassinGender string
	if assassinGender, isDef = c.GetPostForm("gender"); !isDef {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "gender is required"})
		return
	}

	var newAssassin Assassin
	newAssassin.Name = assassinName
	newAssassin.Age = uint16(rand.Intn(80-20+1) + 20) // min = 20, max = 80
	newAssassin.Gender = assassinGender
	newAssassin.Speed = uint16(rand.Intn(10))
	newAssassin.Strength = uint16(rand.Intn(10))
	newAssassin.Busy = false
	newAssassin.CreatedAt = uint64(time.Now().Unix())
	newAssassin.UpdatedAt = uint64(time.Now().Unix())

	result, err := db.Exec("INSERT INTO assassins SET ?", newAssassin)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newAssassin.ID = uint64(id)

	c.IndentedJSON(http.StatusCreated, newAssassin)
}
