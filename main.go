package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

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
	bookmarked bool
}

// func searchMovieInDb(DB *sql.DB) {
// 	sqlQuery := `select title, overview FROM movies where Title="Avatar"`
// 	var title, overview string
// 	err := DB.QueryRow(sqlQuery).Scan(&title, &overview)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			log.Fatalf("No rows found with query: %s", sqlQuery)
// 		}
// 		fmt.Println("cant find", err)
// 	}

// 	fmt.Printf("Title: %s\n", title)
// 	fmt.Printf("overview: %s\n", overview)
// }

func updateBookmark(DB *sql.DB, movieID int, bookmarked bool) error {
    sqlQuery := `UPDATE movies SET bookmarked = ? WHERE id = ?`

    _, err := DB.Exec(sqlQuery, bookmarked, movieID)

    if err != nil {

        log.Printf("Something went wrong with bookmark update, status: %v", err)
        return err
    }

    return nil
}

func getMovieByIdFromDB(DB *sql.DB, movieId int) (Movie, error) {
	 sqlQuery := `SELECT * FROM movies WHERE id = ?;`

	var movie Movie

	err := DB.QueryRow(sqlQuery, movieId).Scan(&movie.id, &movie.original_title, &movie.budget, &movie.popularity, &movie.release_date, &movie.revenue, &movie.title, &movie.vote_average, &movie.vote_count, &movie.overview, &movie.tagline, &movie.uid, &movie.director_id, &movie.bookmarked);
	if err != nil {
		fmt.Println("database cant be open, err: %s", err)
		return Movie{}, err
	}
	return movie, nil
}

func getMovies(DB *sql.DB, sqlQuery string, args ...interface{}) ([]Movie, error) {
	data := []Movie{}
	movie := Movie {}
	rows, err := DB.Query(sqlQuery, args...)
	if err != nil {
		fmt.Println("(getMovies) database cant be open, err: %s", err)
		return nil, err
	}

	defer rows.Close()	
   
	for rows.Next() {
		err := rows.Scan(&movie.id, &movie.original_title, &movie.budget, &movie.popularity, &movie.release_date, &movie.revenue, &movie.title, &movie.vote_average, &movie.vote_count, &movie.overview, &movie.tagline, &movie.uid, &movie.director_id, &movie.bookmarked)
		if err != nil {
			log.Fatal("Error with query, status: ", err)
			return nil, err
		}
		
		data = append(data, movie)
	}
	return data, nil
}

func columnExist(DB *sql.DB, tableName string, columnName string)(bool, error) {
	sqlQuery := fmt.Sprintf("PRAGMA table_info(%s);", tableName)

	rows, err := DB.Query(sqlQuery, tableName)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	for rows.Next() {
		var cid int
		var name, ctype string
		var notnull, pk int
		var dflt_value *int

		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dflt_value, &pk); err != nil {
			return false, err
		}

		if name == columnName {
			return true, nil 
		}
	}
	return false, nil
}

func handleBookmarkChange(DB *sql.DB, movieID int, flag bool, w http.ResponseWriter, r *http.Request) {
	updateBookmark(DB, movieID, flag)
	movie, err := getMovieByIdFromDB(DB, movieID)
	if err != nil {
		log.Fatal("something went wrong with getting movie from DB, err: %s", err)
	}
	movieThumbnail(movie).Render(r.Context(), w)
}

func main() {
	DB, err:= sql.Open("sqlite3", "data/movie.sqlite")
	if err != nil {
		fmt.Println("database cant be open", err)
	}
	defer DB.Close()

	exists, err := columnExist(DB, "movies", "bookmarked")
	if err != nil {
		log.Fatal("error handling collumnexist", err)
	}

	if !exists {
		_, err = DB.Exec("ALTER TABLE movies ADD COLUMN bookmarked BOOLEAN DEFAULT 0")
		if err != nil {
			log.Println("Error adding new column:", err)
		} else {
			fmt.Println("Column 'bookmarked' added successfully.")
		}
	} else {
		log.Printf("Column '%s' already exists", "bookmarked")
	}

	
	fs := http.FileServer(http.Dir("static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))
	
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
	})

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		title := r.FormValue("search")

		sqlQuery := `SELECT * FROM movies WHERE title LIKE ? LIMIT 100`
		searchTitle := "%" + title + "%"
		
		data, err := getMovies(DB, sqlQuery, searchTitle)
		if(err != nil) {
			log.Fatal(err)
		}
		searchTabTemplate(data).Render(r.Context(), w)
	})

	http.HandleFunc("GET /saved/{id}", func(w http.ResponseWriter, r *http.Request) {
		movieID, err :=  strconv.Atoi(r.PathValue("id"))
		if err != nil {
			log.Fatal("something went wrong with integer to string conversion, status: %s", err)
		}
		handleBookmarkChange(DB, movieID, true, w, r)
	})

	http.HandleFunc("GET /unsaved/{id}", func(w http.ResponseWriter, r *http.Request) {
		movieID, err :=  strconv.Atoi(r.PathValue("id"))
		if err != nil {
			log.Fatal("something went wrong with integer to string conversion, status: %s", err)
		}
		handleBookmarkChange(DB, movieID, false, w, r)
	})

	http.HandleFunc("GET /favorites", func(w http.ResponseWriter, r *http.Request) {
		sqlQuery := `select * from movies where bookmarked = ? LIMIT 100`
		args := "1"
		data, err := getMovies(DB, sqlQuery, args)
		if err != nil {
			log.Fatal("something went wrong with getting favorites, error: %s", err)
		}

		//handle bookmarks
		favorites(data).Render(r.Context(), w)
	})

	http.HandleFunc("GET /movie/{id}", func(w http.ResponseWriter, r *http.Request) {
		movieID, err :=  strconv.Atoi(r.PathValue("id"))
		if err != nil {
			log.Fatal("something went wrong with integer to string conversion, status: %s", err)
		}

		movie, err := getMovieByIdFromDB(DB, movieID)
		if err != nil {
			log.Fatal("Something went wrong with getting the movie from db, err: %s", err)
		}

		movieDetails(movie).Render(r.Context(), w)
	})

	// Start the server.
	fmt.Println("listening on http://localhost:8000")
	if err := http.ListenAndServe("localhost:8000", nil); err != nil {
		log.Printf("error listening: %v", err)
	}
}

