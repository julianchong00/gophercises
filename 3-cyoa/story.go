package cyoa

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
)

func init() {
	tpl = template.Must(template.New("").Parse(defaultHandlerTemplate))
}

var tpl *template.Template

var defaultHandlerTemplate = `
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
                    <li><a href="/{{.Chapter}}">{{.Text}}</a></li>
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

func NewHandler(s Story) http.Handler {
	return handler{s}
}

type handler struct {
	s Story
}

func (h handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	path := strings.TrimSpace(r.URL.Path)
	if path == "" || path == "/" {
		// When no path, assume starting story from the intro
		path = "/intro"
	}
	// "/intro" => "intro" - remove slashes
	path = path[1:]

	// path = ["intro"]
	// this if statement checks if the map has something in it with the
	// "ok" variable, otherwise it might still enter conditional without
	// anything inside the map

    // checks if chapter link exists in map, then writes template to rw
    // with template.Execute
	if chapter, ok := h.s[path]; ok {
		err := tpl.Execute(rw, chapter)
		if err != nil {
			log.Printf("%v", err)
			http.Error(rw, "Something went wrong...", http.StatusInternalServerError)
		}
		return
	}
	http.Error(rw, "Chapter not found.", http.StatusNotFound)
}

func JsonStory(r io.Reader) (Story, error) {
	d := json.NewDecoder(r)
	var story Story
	if err := d.Decode(&story); err != nil {
		return nil, err
	}

	return story, nil
}

// define Story type which will store our story data
type Story map[string]Chapter

// define chapter struct
type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

// define options struct
type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}
