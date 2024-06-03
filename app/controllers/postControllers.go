package controllers

import (
	"CRUD/app/models"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func CreatePostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(*models.UserClaims)
		if user.Role != "admin" && user.Role != "user" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var post models.Post
		err := json.NewDecoder(r.Body).Decode(&post)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		query := `INSERT INTO posts (title, content, status, publish_date) VALUES ($1, $2, $3, $4) RETURNING id`
		err = db.QueryRow(query, post.Title, post.Content, post.Status, post.PublishDate).Scan(&post.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, tag := range post.Tags {
			tagID := 0
			err = db.QueryRow(`INSERT INTO tags (label) VALUES ($1) ON CONFLICT (label) DO UPDATE SET label = EXCLUDED.label RETURNING id`, tag.Label).Scan(&tagID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			_, err = db.Exec(`INSERT INTO post_tags (post_id, tag_id) VALUES ($1, $2)`, post.ID, tagID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(post)
	}
}

func GetPostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		postID, err := strconv.Atoi(vars["post_id"])
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		var post models.Post
		query := `SELECT id, title, content, status, publish_date FROM posts WHERE id = $1`
		err = db.QueryRow(query, postID).Scan(&post.ID, &post.Title, &post.Content, &post.Status, &post.PublishDate)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Post not found", http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		rows, err := db.Query(`SELECT t.id, t.label FROM tags t INNER JOIN post_tags pt ON t.id = pt.tag_id WHERE pt.post_id = $1`, postID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var tag models.Tag
			err := rows.Scan(&tag.ID, &tag.Label)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			post.Tags = append(post.Tags, tag)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(post)
	}
}

func UpdatePostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(*models.UserClaims)
		if user.Role != "admin" && user.Role != "user" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		vars := mux.Vars(r)
		postID, err := strconv.Atoi(vars["post_id"])
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		var post models.Post
		err = json.NewDecoder(r.Body).Decode(&post)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if user.Role == "user" && post.Status == "publish" {
			http.Error(w, "Forbidden: users cannot publish posts", http.StatusForbidden)
			return
		}

		query := `UPDATE posts SET title = $1, content = $2, status = $3, publish_date = $4 WHERE id = $5`
		_, err = db.Exec(query, post.Title, post.Content, post.Status, post.PublishDate, postID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = db.Exec(`DELETE FROM post_tags WHERE post_id = $1`, postID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, tag := range post.Tags {
			tagID := 0
			err = db.QueryRow(`INSERT INTO tags (label) VALUES ($1) ON CONFLICT (label) DO UPDATE SET label = EXCLUDED.label RETURNING id`, tag.Label).Scan(&tagID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			_, err = db.Exec(`INSERT INTO post_tags (post_id, tag_id) VALUES ($1, $2)`, postID, tagID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(post)
	}
}

func DeletePostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(*models.UserClaims)
		if user.Role != "admin" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		vars := mux.Vars(r)
		postID, err := strconv.Atoi(vars["post_id"])
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		_, err = db.Exec(`DELETE FROM posts WHERE id = $1`, postID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = db.Exec(`DELETE FROM post_tags WHERE post_id = $1`, postID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
