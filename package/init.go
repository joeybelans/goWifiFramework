// New package template
package newPackage

import (
	"net/http"
	"text/template"

	"github.com/joeybelans/gokismet/httpHandler"
)

// This is called when the package is loaded
func init() {
	// Define the template to be loaded when the page is requested
	httpTemplate = template.New("/path")
	httpTemplate, _ = httpTemplate.Parse(httpSource())
}

// Initialize the package
// This function must be defined and included in main.go
// Update with the necessary arguments
func Init() {
	// Add the page to the list of pages
	httpHandler.AddPage("/path", "Title")
	http.HandleFunc("/path", func(w http.ResponseWriter, r *http.Request) {
		HttpFunction(w, r)
	})
}
