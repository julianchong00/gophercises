package main

import (
	"fmt"
	"strings"

	"github.com/julianchong00/link"
)

var exampleHTML = `
<html>

<body>
    <a href="/dog-cat">dog cat
        <!-- commented text SHOULD NOT be included! -->
    </a>
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
