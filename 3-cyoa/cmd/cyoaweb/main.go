package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

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

	// Create custom template with custom path function "/story/{chapter}" instead of "/{chapter}"
	tpl := template.Must(template.New("").Parse(storyTemplate))
	// To use functional options, call NewHandler and pass in option functions
	// for whichever option you want to pass.
	// We set options to the return values of functions defined in the API
	// so that no misunderstandings are created for the user using the API or code
	//
	// Pass in options by calling functional options within API
	h := cyoa.NewHandler(story, cyoa.WithTemplate(tpl), cyoa.WithPathFunc(pathFn))

	mux := http.NewServeMux()
	mux.Handle("/story/", h)
	fmt.Printf("Starting the server on port: %d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), h))
}

// Pass custom path function
func pathFn(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "/story" || path == "/story/" {
		path = "/story/intro"
	}
	return path[len("/story/"):]
}

var storyTemplate = `
    <!DOCTYPE html>
    <html>
    <head>
        <title>Choose Your Own Adventure</title>
        <meta charset="UTF-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
    </head>
    <body>
        <section class="page">
            <h1>{{.Title}}</h1>
            {{range .Paragraphs}}
            <p>{{.}}</p>
            {{end}}
            <ul>
                {{range .Options}}
                    <li><a href="/story/{{.Chapter}}">{{.Text}}</a></li>
                {{end}}
            </ul>
            <style>
                body {
                    font-family: helvetica, arial;
                }
                h1 {
                    text-align: center;
                    position: relative;
                }
                .page {
                    width: 80%;
                    max-width: 500px;
                    margin: auto;
                    margin-top: 40px;
                    margin-bottom: 40px;
                    padding: 80px;
                    background: #FFFCF6;
                    border: 1px solid #eee;
                    box-shadow: 0 10px 6px -6px #777;
                }
                ul {
                    border-top: 1px dotted #ccc;
                    padding: 10px 0 0 0;
                    -webkit-padding-start: 0;
                }
                li {
                    padding-top: 10px;
                }
                a,
                a:visited {
                    text-decoration: none;
                    color:#6295b5;
                }
                a:active,
                a:hover {
                    color: #7792a2;
                }
                p {
                    text-indent: 1em;
                }
            </style>
        </section>
    </body>
    </html>`
