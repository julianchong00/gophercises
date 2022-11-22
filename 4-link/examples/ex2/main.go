package main

import (
	"fmt"
	"strings"

	"github.com/julianchong00/link"
)

var exampleHTML = `
<html>

<head>
    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/font-awesome/4.7.0/css/font-awesome.min.css" />
</head>

<body>
    <h1>Social stuffs</h1>
    <div>
        <a href="https://www.twitter.com/joncalhoun">
            Check me out on twitter
            <i class="fa fa-twitter" aria-hidden="true"></i>
        </a>
        <a href="https://github.com/gophercises">
            Gophercises is on <strong>Github</strong>!
        </a>
    </div>
</body>

</html>
`

func main() {
	// strings.NewReader returns a reader from a string input
	r := strings.NewReader(exampleHTML)
	links, err := link.Parse(r)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", links)
}
