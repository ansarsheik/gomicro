package main

import (
	"net/http"

	"github.com/gamegos/jsend"
)

type Post struct {
	ID      int64  `json:"postid"`
	Title   string `json:"title"`
	Details string `json:"details"`
}

func addNewPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var postdata Post
	postdata.Title = r.FormValue("title")
	postdata.Details = r.FormValue("details")

	if len(postdata.Title) < 1 || len(postdata.Details) < 1 {
		jsend.Wrap(w).
			Status(400).
			Message("Post Title & Details are required to create post").
			Data("").
			Send()

		return
	}

	db := DBconnection()
	query := "INSERT INTO post (title, details, category_id) values(?,?, 1)"
	result, err := db.Exec(query, postdata.Title, postdata.Details)

	if err != nil {
		jsend.Wrap(w).
			Status(400).
			Message("Error adding post").
			Data("").
			Send()

		return
	}

	postdata.ID, _ = result.LastInsertId()

	jsend.Wrap(w).
		Status(201).
		Message("Post added").
		Data(postdata).
		Send()
}
