// Generates the header source code
// Exported functions
package header

import (
	"fmt"
	"net/http"
	"text/template"

	_ "github.com/joeybelans/gokismet/statik"
)

// Add a page
func AddPage(url string, title string) {
	pages[len(pages)] = page{URL: url, Title: title}
}

// Display header
func Display(w http.ResponseWriter, title string, file string) {
	header.Execute(w, templateData{title, file})
}

// Create the header template
//func Create() {
func Create(packages map[string]interface{}) {
	fmt.Println(packages)
	header = template.New("header")
	header, _ = header.Parse(templateSource(pages))
}
