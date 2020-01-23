package main

import (
	"net/http"

	"github.com/gamegos/jsend"
)

type category struct {
	ID       int64  `json:"Id"`
	Category string `json:"Category"`
}

func getCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	db := DBconnection()
	query := "SELECT id,category FROM category WHERE status = 1"
	results, err := db.Query(query)

	if err != nil {
		jsend.Wrap(w).
			Status(400).
			Message("Error gettting categories").
			Data(err).
			Send()

		return
	}

	defer results.Close()

	var cat []category

	for results.Next() {
		var c category
		err = results.Scan(
			&c.ID,
			&c.Category)

		cat = append(cat, c)
		checkErr(err)
	}

	jsend.Wrap(w).
		Status(200).
		Message("Categories").
		Data(cat).
		Send()
}
