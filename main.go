package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/a-h/templ"
)


func main() {
	// Use a template that doesn't take parameters.
	http.Handle("/", templ.Handler(home()))


	http.HandleFunc("/clicked", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.FormValue("search"))
		resp, err := http.Get(`https://search.imdbot.workers.dev/?q=` + r.FormValue("search"))

		if(err != nil){
			fmt.Println("Error: ", err)
		}

		fmt.Println(resp)
		searchResultTemplate("test").Render(r.Context(), w)
	})

	// Start the server.
	fmt.Println("listening on http://localhost:8000")
	if err := http.ListenAndServe("localhost:8000", nil); err != nil {
		log.Printf("error listening: %v", err)
	}
}

