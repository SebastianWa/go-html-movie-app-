package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/a-h/templ"
	_ "github.com/mattn/go-sqlite3"
)

type Movie struct {
	id int
	original_title string
	budget int
	popularity int
	release_date string
	revenue int
	title string
	vote_average string
	vote_count int
	overview string
	tagline string
	uid int
	director_id int
}

func searchMovieInDb(db *sql.DB) {
	sqlQuery := `select title, overview FROM movies where Title="Avatar"`
	var title, overview string
	err := db.QueryRow(sqlQuery).Scan(&title, &overview)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Fatalf("No rows found with query: %s", sqlQuery)
		}
		fmt.Println("cant find", err)
	}

	fmt.Printf("Title: %s\n", title)
	fmt.Printf("overview: %s\n", overview)
}

func getMoviesFromDB(db *sql.DB, searchTitle string) ([]Movie, error) {
	sqlQuery := `select * from movies where Title LIKE $1 LIMIT 100`
	fmt.Println(sqlQuery)
	data := []Movie{}
	rows, err := db.Query(sqlQuery, "%" + searchTitle + "%")
	if err != nil {
		fmt.Println("database cant be open", err)
		return nil, err
	}
	
	var id int
	var original_title string
	var budget int
	var popularity int
	var release_date string
	var revenue int
	var title string
	var vote_average string
	var vote_count int
	var overview string
	var tagline string
	var uid int
	var director_id int

	defer rows.Close()	
   
	for rows.Next() {
		err := rows.Scan(&id, &original_title, &budget, &popularity, &release_date, &revenue, &title, &vote_average, &vote_count, &overview, &tagline, &uid, &director_id)
		if err != nil {
			log.Fatal("Eror with query: ", err)
			return nil, err
		}
		
		data = append(data, Movie{
			id, original_title, budget, popularity, release_date, revenue, title, vote_average, vote_count, overview, tagline, uid, director_id})
	}
	fmt.Println(data)
	return data, nil
}

func main() {
	//db connection, csv to db 
	// db := openDataBase()
	// defer db.Close()

	// data, _ := parseCvs("data/datasetscv.csv")
	// insertData(db, parseMovies(data))

	db, err:= sql.Open("sqlite3", "data/movie.sqlite")
	if err != nil {
		fmt.Println("database cant be open", err)
	}
	defer db.Close()

	fs := http.FileServer(http.Dir("static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))
	
	// sqlQuery := `select * FROM movies where Title="Avatar"`
	
	// Use a template that doesn't take parameters.
	http.Handle("/", templ.Handler(Home()))


	http.HandleFunc("/clicked", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.FormValue("search"))
		resp, err := http.Get(`https://search.imdbot.workers.dev/?q=` + r.FormValue("search"))

		if(err != nil){
			fmt.Println("Error: ", err)
		}

		fmt.Println(resp)
	})

	http.HandleFunc("/searchTab", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("searchTab")

		// searchTabTemplate(nil).Render(r.Context(), w)		
	})

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		title := r.FormValue("search")
		data, err := getMoviesFromDB(db, title)
		if(err != nil) {
			log.Fatal(err)
		}
		fmt.Println(data)
		searchTabTemplate(data).Render(r.Context(), w)
	})

	// Start the server.
	fmt.Println("listening on http://localhost:8000")
	if err := http.ListenAndServe("localhost:8000", nil); err != nil {
		log.Printf("error listening: %v", err)
	}
}

