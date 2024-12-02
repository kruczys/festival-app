package main

import (
    "database/sql"
    "encoding/json"
    "net/http"
    "log"
    "os"
    "time"
)
/* POST /festivals/
GET /festivals/
GET /festivals/{id}
DELETE /festivals/{id}
POST /festivals/{id}/performances
PUT /festivals/{id}/performances/{performance_id}
GET /festivals/{id}/performances */
type Festival struct {
    ID          int         `json:"id"`
    Name        string      `json:"name"`
    StartDate   time.Time   `json:"start_date"`
    EndDate     time.Time   `json:"end_date"`
    Location    [2]float64  `json:"location"`
}

type Performance struct {
    ID          int         `json:"id"`
    Name        string      `json:"name"`
    Genre       string      `json:"genre"`
    StartTime   time.Time   `json:"start_time"`
    EndTime     time.Time   `json:"end_time"`
}

func main() {
    return;
}
