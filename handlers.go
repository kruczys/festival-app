package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (app *App) createFestival(w http.ResponseWriter, r *http.Request) {
	var festival Festival
	if err := json.NewDecoder(r.Body).Decode(&festival); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	result := app.DB.Create(&festival)
	if result.Error != nil {
		respondWithError(w, http.StatusInternalServerError, result.Error.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, festival)
}

func (app *App) getFestivals(w http.ResponseWriter, r *http.Request) {
	var festivals []Festival
	result := app.DB.Find(&festivals)
	if result.Error != nil {
		respondWithError(w, http.StatusInternalServerError, result.Error.Error())
		return
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
	result := app.DB.Preload("Performances").First(&festival, id)
	if result.Error != nil {
		respondWithError(w, http.StatusNotFound, "Festival not found")
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

	result := app.DB.Delete(&Festival{}, id)
	if result.Error != nil {
		respondWithError(w, http.StatusInternalServerError, result.Error.Error())
		return
	}

	if result.RowsAffected == 0 {
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

	performance.FestivalID = uint(festivalID)

	result := app.DB.Create(&performance)
	if result.Error != nil {
		respondWithError(w, http.StatusInternalServerError, result.Error.Error())
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

	result := app.DB.Model(&Performance{}).Where("id = ? AND festival_id = ?", performanceID, festivalID).Updates(map[string]interface{}{
		"name":       performance.Name,
		"genre":      performance.Genre,
		"start_time": performance.StartTime,
		"end_time":   performance.EndTime,
	})

	if result.Error != nil {
		respondWithError(w, http.StatusInternalServerError, result.Error.Error())
		return
	}

	if result.RowsAffected == 0 {
		respondWithError(w, http.StatusNotFound, "Performance not found")
		return
	}

	performance.ID = uint(performanceID)
	performance.FestivalID = uint(festivalID)
	respondWithJSON(w, http.StatusOK, performance)
}

func (app *App) getPerformances(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	festivalID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid festival ID")
		return
	}

	var performances []Performance
	result := app.DB.Where("festival_id = ?", festivalID).Find(&performances)
	if result.Error != nil {
		respondWithError(w, http.StatusInternalServerError, result.Error.Error())
		return
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
