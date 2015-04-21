// Manages home interface
// Initializes package
package home

import (
	"net/http"
	"text/template"

	"github.com/joeybelans/gokismet/header"
	_ "github.com/joeybelans/gokismet/statik"
	"github.com/joeybelans/gokismet/webSocketHandler"
)

var home *template.Template

// Called when package is loaded
func init() {
	// Create root page
	header.AddPage("/", "Home")
	webSocketHandler.AddPage("/", "home")

	// Create the template
	home = template.New("/")
	home, _ = home.Parse(templateSource())
}

// Initialize the package
func Init(lhost string, lport int, khost string, kport int, ssids []string) interface{} {
	// Create home HTTP handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		HttpHome(w, r, lhost, lport, khost, kport, ssids)
	})
	return (ProcessCommand)
}
