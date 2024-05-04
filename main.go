package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	_ "github.com/mattn/go-sqlite3"
)

type MovieData struct {
	id int
	title string
	voteAverage float64
	voteCount int
	status string
	releaseDate string
	revenue int
	length int
	adult bool
	backdropPath string
	budget int
	homepage string
	imdbId string
	originalLanguage string
	originalTitle string
	overview string
	popularity string
	posterPath string
	tagline string
	genres string
	productionCompanies string
	productionCountries string
	spokenLanguages string
	keywords string
}

func stringToFloat(s string) (data float64) {
	f, err := strconv.ParseFloat(s, 64)
	if(err != nil) {
		panic(err)
	}
	return f
}

func stringToInt(s string) (data int) {
	s = strings.TrimSuffix(s, ".0")
	i, err := strconv.Atoi(s)
	if(err != nil) {
		panic(err)
	}
	return i
}

func stringToBool(s string) (data bool) {
	b, err := strconv.ParseBool(s)
	if(err != nil) {
		panic(err)
	}
	return b
}

func parseCvs(filePath string) ([][]string, []string) {
	file, err := os.Open(filePath)

	if err != nil {
		log.Fatal("CSV file problem", err)
	}

	defer file.Close()

	cvsReader := csv.NewReader(file)

	var records [][]string
	
	for {
		record, err := cvsReader.Read()

		if err == io.EOF {
			break
		}
	
		if err != nil {
			log.Fatal("cvsReader.Read() file problem", err)
		}
		
		records = append(records, record)
	}
	return records[1:], records[0]
}

func parseMovies(records [][]string) []MovieData {
	var movies []MovieData

	for _, record := range records {
		movie := MovieData{
			id: stringToInt(record[0]),
			title: record[1],
			voteAverage: stringToFloat(record[2]),
			voteCount: stringToInt(record[3]), 
			status: record[4],
			releaseDate: record[5],
			revenue: stringToInt(record[6]), 
			length: stringToInt(record[7]), 
			adult: stringToBool(record[8]),
			backdropPath: record[9],
			budget: stringToInt(record[10]) ,
			homepage: record[11],
			imdbId: record[12],
			originalLanguage: record[13],
			originalTitle: record[14],
			overview: record[15],
			popularity: record[16],
			posterPath: record[17],
			tagline: record[18],
			genres: record[19],
			productionCompanies: record[20],
			productionCountries: record[21],
			spokenLanguages: record[22],
			keywords: record[23],
		}
		movies = append(movies, movie)
	}
	return movies
}

func insertData(db *sql.DB, movies []MovieData) {
	for _, movie := range movies {
			sql := `
				insert INTO movies(id,
				title,
				voteAverage,
				voteCount,
				status,
				releaseDate,
				revenue,
				length,
				adult,
				backdropPath,
				budget,
				homepage,
				imdbId,
				originalLanguage,
				originalTitle,
				overview,
				popularity,
				posterPath,
				tagline,
				genres,
				productionCompanies,
				productionCountries,
				spokenLanguages,
				keywords)

				VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
			`
		_, err := db.Exec(sql, movie.id, movie.title, movie.voteAverage, movie.voteCount, movie.status, movie.releaseDate, movie.revenue, movie.length, movie.adult, movie.backdropPath, movie.budget, movie.homepage, movie.imdbId, movie.originalLanguage, movie.originalTitle, movie.overview, movie.popularity, movie.posterPath, movie.tagline, movie.genres, movie.productionCompanies, movie.productionCountries, movie.spokenLanguages, movie.keywords)
		if err != nil {
			fmt.Println(err)
        }
	}
}

func createTable(db *sql.DB){
	sql := `
		  CREATE TABLE IF NOT EXISTS movies (
			id INT NOT NULL PRIMARY KEY,
			title TEXT,
			voteAverage FLOAT,
			voteCount INT,
			status TEXT,
			release_date TEXT,
			revenue INT,
			length INT,
			adult BOOL,
			backdropPath TEXT,
			budget INT,
			homepage TEXT,
			imdbId TEXT,
			originalLanguage TEXT,
			originalTitle TEXT,
			overview TEXT,
			popularity TEXT,
			posterPath TEXT,
			tagline TEXT,
			genres TEXT,
			productionCompanies TEXT,
			productionCountries TEXT,
			spokenLanguages TEXT,
			keywords TEXT
		  );
	`

	_, err := db.Exec(sql)
	if err != nil {
        fmt.Println(err)
    }
}

func openDataBase() *sql.DB {
	if _, err := os.Stat("movies.db"); err == nil {
		db, err := sql.Open("sqlite3", "/data/movies.db")
		
		if err != nil {
            fmt.Println(err)
        }
        return db
	}else {
		db, err := sql.Open("sqlite3", "data/movies.db")
		
		if err != nil {
            fmt.Println(err)
        }
		createTable(db)
        return db
	}
}

func main() {
	//db connection, csv to db 
	db := openDataBase()
	defer db.Close()

	data, _ := parseCvs("data/datasetscv.csv")
	insertData(db, parseMovies(data))

	// Use a template that doesn't take parameters.
	http.Handle("/", templ.Handler(home()))


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


		searchTabTemplate().Render(r.Context(), w)		
	})

	// Start the server.
	fmt.Println("listening on http://localhost:8000")
	if err := http.ListenAndServe("localhost:8000", nil); err != nil {
		log.Printf("error listening: %v", err)
	}
}

