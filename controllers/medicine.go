package controllers

import (
	"lily-med/repository"
	"net/http"
	"strconv"
)

func CreateMedicineHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	hourStr := r.FormValue("timeToTakeHour")
	minStr := r.FormValue("timeToTakeMin")

	hour, err := strconv.Atoi(hourStr)
	if err != nil {
		http.Error(w, "Invalid hour format", http.StatusBadRequest)
		return
	}

	min, err := strconv.Atoi(minStr)
	if err != nil {
		http.Error(w, "Invalid min format", http.StatusBadRequest)
		return
	}
	timeToTake := repository.TimeOfDay{Hour: hour, Min: min}
	medicine := repository.Medicine{
		Name:       r.FormValue("name"),
		Taken:      false,
		Disabled:   false,
		TimeToTake: timeToTake,
	}
	if _, err = medicine.AddMedicine(r.Context()); err != nil {
		http.Error(w, "failed to add medicine", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

}
