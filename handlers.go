package main

import (
    "database/sql"
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
    "strconv"
)

func (app *App) createFestival(w http.ResponseWriter, r *http.Request) {
    var festival Festival
    if err := json.NewDecoder(r.Body).Decode(&festival); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }

    err := app.DB.QueryRow(
        `INSERT INTO festivals (name, start_date, end_date, location) 
         VALUES ($1, $2, $3, $4) 
         RETURNING id`,
        festival.Name, festival.StartDate, festival.EndDate, festival.Location,
    ).Scan(&festival.ID)

    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusCreated, festival)
}

func (app *App) getFestivals(w http.ResponseWriter, r *http.Request) {
    rows, err := app.DB.Query(
        "SELECT id, name, start_date, end_date, location FROM festivals",
    )
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }
    defer rows.Close()

    var festivals []Festival
    for rows.Next() {
        var f Festival
        if err := rows.Scan(&f.ID, &f.Name, &f.StartDate, &f.EndDate, &f.Location); err != nil {
            respondWithError(w, http.StatusInternalServerError, err.Error())
            return
        }
        festivals = append(festivals, f)
    }

    respondWithJSON(w, http.StatusOK, festivals)
}

func (app *App) getFestival(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid festival ID")
        return
    }

    var festival Festival
    err = app.DB.QueryRow(
        "SELECT id, name, start_date, end_date, location FROM festivals WHERE id = $1",
        id,
    ).Scan(&festival.ID, &festival.Name, &festival.StartDate, &festival.EndDate, &festival.Location)

    if err == sql.ErrNoRows {
        respondWithError(w, http.StatusNotFound, "Festival not found")
        return
    }
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusOK, festival)
}

func (app *App) deleteFestival(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid festival ID")
        return
    }

    result, err := app.DB.Exec("DELETE FROM festivals WHERE id = $1", id)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    if rowsAffected == 0 {
        respondWithError(w, http.StatusNotFound, "Festival not found")
        return
    }

    respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (app *App) createPerformance(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    festivalID, err := strconv.Atoi(vars["id"])
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid festival ID")
        return
    }

    var performance Performance
    if err := json.NewDecoder(r.Body).Decode(&performance); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }

    performance.FestivalID = festivalID

    err = app.DB.QueryRow(
        `INSERT INTO performances (festival_id, name, genre, start_time, end_time) 
         VALUES ($1, $2, $3, $4, $5) 
         RETURNING id`,
        festivalID, performance.Name, performance.Genre,
        performance.StartTime, performance.EndTime,
    ).Scan(&performance.ID)

    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusCreated, performance)
}

func (app *App) updatePerformance(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    festivalID, err := strconv.Atoi(vars["id"])
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid festival ID")
        return
    }

    performanceID, err := strconv.Atoi(vars["performance_id"])
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid performance ID")
        return
    }

    var performance Performance
    if err := json.NewDecoder(r.Body).Decode(&performance); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }

    result, err := app.DB.Exec(
        `UPDATE performances 
         SET name = $1, genre = $2, start_time = $3, end_time = $4
         WHERE id = $5 AND festival_id = $6`,
        performance.Name, performance.Genre, performance.StartTime,
        performance.EndTime, performanceID, festivalID,
    )

    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    if rowsAffected == 0 {
        respondWithError(w, http.StatusNotFound, "Performance not found")
        return
    }

    respondWithJSON(w, http.StatusOK, performance)
}

func (app *App) getPerformances(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    festivalID, err := strconv.Atoi(vars["id"])
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid festival ID")
        return
    }

    rows, err := app.DB.Query(
        `SELECT id, festival_id, name, genre, start_time, end_time 
         FROM performances 
         WHERE festival_id = $1`,
        festivalID,
    )
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }
    defer rows.Close()

    var performances []Performance
    for rows.Next() {
        var p Performance
        if err := rows.Scan(
            &p.ID, &p.FestivalID, &p.Name, &p.Genre,
            &p.StartTime, &p.EndTime,
        ); err != nil {
            respondWithError(w, http.StatusInternalServerError, err.Error())
            return
        }
        performances = append(performances, p)
    }

    respondWithJSON(w, http.StatusOK, performances)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
    respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    response, err := json.Marshal(payload)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("Error encoding response"))
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(response)
}
