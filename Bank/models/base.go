package models

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

var db *gorm.DB //database

func init() {

	e := godotenv.Load() //Load .env file
	if e != nil {
		fmt.Print(e)
	}

	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")
	//fmt.Println("->", username, ":", password, ":", dbName, ":", dbHost)
	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, username, dbName, password) //Build connection string
	fmt.Println("URI:", dbUri)

	conn, err := gorm.Open("postgres", dbUri)
	if err != nil {
		fmt.Println("error:", err)
	}

	db = conn
	db.Debug().AutoMigrate(&User{}) //Database migration
}

//returns a handle to the DB object
func GetDB() *gorm.DB {
	return db
}