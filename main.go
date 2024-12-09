package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Location [2]float64

type Festival struct {
	ID           uint          `json:"id" gorm:"primaryKey"`
	Name         string        `json:"name" gorm:"not null"`
	StartDate    time.Time     `json:"start_date" gorm:"not null"`
	EndDate      time.Time     `json:"end_date" gorm:"not null"`
	Location     Location      `json:"location" gorm:"type:point;not null"`
	Performances []Performance `json:"performances,omitempty" gorm:"foreignKey:FestivalID"`
}

type Performance struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	FestivalID uint      `json:"festival_id"`
	Name       string    `json:"name" gorm:"not null"`
	Genre      string    `json:"genre" gorm:"not null"`
	StartTime  time.Time `json:"start_time" gorm:"not null"`
	EndTime    time.Time `json:"end_time" gorm:"not null"`
}

type App struct {
	DB     *gorm.DB
	Router *mux.Router
}

func (app *App) Initialize() error {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	var err error
	app.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("error connecting to the database: %v", err)
	}

	err = app.DB.AutoMigrate(&Festival{}, &Performance{})
	if err != nil {
		return fmt.Errorf("error migrating database: %v", err)
	}

	app.Router = mux.NewRouter()
	app.setupRoutes()
	return nil
}

func (app *App) setupRoutes() {
	app.Router.HandleFunc("/festivals", app.createFestival).Methods("POST")
	app.Router.HandleFunc("/festivals", app.getFestivals).Methods("GET")
	app.Router.HandleFunc("/festivals/{id:[0-9]+}", app.getFestival).Methods("GET")
	app.Router.HandleFunc("/festivals/{id:[0-9]+}", app.deleteFestival).Methods("DELETE")
	app.Router.HandleFunc("/festivals/{id:[0-9]+}/performances", app.createPerformance).Methods("POST")
	app.Router.HandleFunc("/festivals/{id:[0-9]+}/performances/{performance_id:[0-9]+}", app.updatePerformance).Methods("PUT")
	app.Router.HandleFunc("/festivals/{id:[0-9]+}/performances", app.getPerformances).Methods("GET")
}

func main() {
	app := &App{}
	if err := app.Initialize(); err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, app.Router))
}
