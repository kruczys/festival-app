package main

import (
    "database/sql"
    "database/sql/driver"
    "errors"
    "fmt"
    "log"
    "net/http"
    "os"
    "strings"
    "time"
    "github.com/gorilla/mux"
    "github.com/joho/godotenv"
    _ "github.com/lib/pq"
)

type Location [2]float64

func (l Location) Value() (driver.Value, error) {
    return fmt.Sprintf("(%f,%f)", l[0], l[1]), nil
}

func (l *Location) Scan(value interface{}) error {
    if value == nil {
        return errors.New("location cannot be null")
    }

    switch v := value.(type) {
    case string:
        str := strings.Trim(v, "()")
        _, err := fmt.Sscanf(str, "%f,%f", &l[0], &l[1])
        return err
    case []byte:
        str := strings.Trim(string(v), "()")
        _, err := fmt.Sscanf(str, "%f,%f", &l[0], &l[1])
        return err
    default:
        return fmt.Errorf("unsupported location format: %T", value)
    }
}

type Festival struct {
    ID        int       `json:"id"`
    Name      string    `json:"name"`
    StartDate time.Time `json:"start_date"`
    EndDate   time.Time `json:"end_date"`
    Location  Location  `json:"location"`
}

type Performance struct {
    ID         int       `json:"id"`
    Name       string    `json:"name"`
    Genre      string    `json:"genre"`
    StartTime  time.Time `json:"start_time"`
    EndTime    time.Time `json:"end_time"`
    FestivalID int       `json:"festival_id"`
}

type App struct {
    DB     *sql.DB
    Router *mux.Router
}

func (app *App) Initialize() error {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using environment variables")
    }

    connStr := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
    )

    var err error
    app.DB, err = sql.Open("postgres", connStr)
    if err != nil {
        return fmt.Errorf("error connecting to the database: %v", err)
    }

    if err := app.createTables(); err != nil {
        return fmt.Errorf("error creating tables: %v", err)
    }

    app.Router = mux.NewRouter()
    app.setupRoutes()
    return nil
}

func (app *App) createTables() error {
    _, err := app.DB.Exec(`
        CREATE TABLE IF NOT EXISTS festivals (
            id SERIAL PRIMARY KEY,
            name VARCHAR(100) NOT NULL,
            start_date DATE NOT NULL,
            end_date DATE NOT NULL,
            location POINT NOT NULL
        )
    `)
    if err != nil {
        return err
    }

    _, err = app.DB.Exec(`
        CREATE TABLE IF NOT EXISTS performances (
            id SERIAL PRIMARY KEY,
            festival_id INTEGER REFERENCES festivals(id) ON DELETE CASCADE,
            name VARCHAR(100) NOT NULL,
            genre VARCHAR(50) NOT NULL,
            start_time TIMESTAMP NOT NULL,
            end_time TIMESTAMP NOT NULL
        )
    `)
    return err
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
