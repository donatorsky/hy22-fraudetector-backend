package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os/exec"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type Tweet struct {
	ID         int      `db:"id" json:"id"`
	Link       string   `db:"link" json:"link"`
	Content    string   `db:"content" json:"content"`
	Label      *int     `db:"label" json:"label"`
	Prediction *float64 `db:"prediction" json:"prediction"`
}

type CreateTweetsReq struct {
	Tweets []Tweet `json:"tweets"`
}

type TweetActionReq struct {
	Action string `json:"action"`
}

func main() {
	db, err := sql.Open("sqlite3", "fraudector.db")
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()

	r.HandleFunc("/unverified-tweets", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT * FROM tweets ORDER BY id DESC LIMIT 10")
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		tweets := make([]Tweet, 0)

		for rows.Next() {
			tweet := Tweet{}

			err = rows.Scan(&tweet.Link, &tweet.Content, &tweet.ID, &tweet.Label, &tweet.Prediction)
			if err != nil {
				panic(err)
			}

			tweets = append(tweets, tweet)
		}
		if err := rows.Err(); err != nil {
			panic(err)
		}

		j, err := json.Marshal(tweets)
		if err != nil {
			panic(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(j)
	}).Methods("GET")

	r.HandleFunc("/tweets", func(w http.ResponseWriter, r *http.Request) {
		var tweets CreateTweetsReq
		if err := json.NewDecoder(r.Body).Decode(&tweets); err != nil {
			panic(err)
		}

		for _, tweet := range tweets.Tweets {
			args := []string{tweet.Content}
			out, err := exec.Command("python3 prediction/get_prediction.py", args...).Output()
			if err != nil {
				panic(err)
			}

			pred, err := strconv.ParseFloat(string(out), 64)
			if err == nil {
				panic(err)
			}

			_, err = db.Exec("INSERT INTO tweets (link, content, label, prediction) VALUES ($1, $2, $3, $4)", tweet.Link, tweet.Content, nil, pred)
			if err != nil {
				panic(err)
			}
		}
	}).Methods("POST")

	r.HandleFunc("/tweets/{id}/action", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		tweetID, err := strconv.Atoi(vars["id"])
		if err != nil {
			panic(err)
		}

		var action TweetActionReq
		if err := json.NewDecoder(r.Body).Decode(&action); err != nil {
			panic(err)
		}

		switch action.Action {
		case "decline":
			_, err := db.Exec("UPDATE tweets SET label=1 WHERE id=$1", tweetID)
			if err != nil {
				panic(err)
			}

		case "accept":
			_, err := db.Exec("UPDATE tweets SET label=0 WHERE id=$1", tweetID)
			if err != nil {
				panic(err)
			}
		}
	}).Methods("POST")

	http.ListenAndServe(":80", r)
}
