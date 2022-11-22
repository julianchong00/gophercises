package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/julianchong00/cyoa"
)

// Use this file to run web application
func main() {
	port := flag.Int("port", 3000, "the port to start the CYOA web application on")
	file := flag.String("file", "gopher.json", "the JSON file with the CYOA story")
	flag.Parse()
	fmt.Printf("Using the story in %v.\n", *file)

	f, err := os.Open(*file)
	if err != nil {
		panic(err)
	}

	story, err := cyoa.JsonStory(f)
	if err != nil {
		panic(err)
	}

    tpl := template.Must(template.New("").Parse("Hello!"))
    // To use functional options, call NewHandler and pass in option functions
    // for whichever option you want to pass.
    // We set options to the return values of functions defined in the API
    // so that no misunderstandings are created for the user using the API or code
	h := cyoa.NewHandler(story, cyoa.WithTemplate(tpl))
	fmt.Printf("Starting the server on port: %d\n", *port)
    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), h))
}
