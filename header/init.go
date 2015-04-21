// Generates header source code
// Initialization file
package header

import (
	"log"
	"net/http"
	"text/template"

	_ "github.com/joeybelans/gokismet/statik"
	"github.com/joeybelans/gokismet/webSocketHandler"
	"github.com/rakyll/statik/fs"
)

type page struct {
	URL   string
	Title string
}

var (
	pages  map[int]page
	header *template.Template
)

// Called when package is loaded
func init() {
	// Initialize pages map
	pages = make(map[int]page)
}

// Initialize the package
func Init() {
	/*
		fmt.Println(packages)
			        // Create the header template
				header = template.New("header")
				header, _ = header.Parse(templateSource(pages))
	*/

	// Static files
	statikFS, err := fs.New()
	if err != nil {
		log.Fatalf(err.Error())
	}

	// Create additional HTTP handlers
	http.HandleFunc("/ws", webSocketHandler.ServeWS)
	http.Handle("/css/", http.FileServer(statikFS))
	http.Handle("/js/", http.FileServer(statikFS))
}
