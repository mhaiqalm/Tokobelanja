package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"Tokobelanja/app/helper"
	"Tokobelanja/domain"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func StartDB() {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	sslMode := "disable"
	if os.Getenv("APP_ENV") == "production" {
		sslMode = "require"
	}

	config := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta", host, user, password, dbName, dbPort, sslMode)
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connection to database", err)
	}

	log.Println("Success connect to Database")
	db.AutoMigrate(domain.User{}, domain.Category{}, domain.Product{}, domain.Transaction{})
	Seeders(db)
	SetUpDBConnection(db)
}

func SetUpDBConnection(DB *gorm.DB) {
	db = DB
}

func GetDBConnection() *gorm.DB {
	return db
}

func Seeders(db *gorm.DB) {
	var user domain.User = domain.User{
		FullName:  "admin",
		Email:     "admin@gmail.com",
		Balance:   0,
		Role:      "admin",
		Password:  helper.HassPass("admin123"),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := db.First(&user, "email = ?", user.Email).Error
	if err != nil {
		db.Create(&user)
	}
}
