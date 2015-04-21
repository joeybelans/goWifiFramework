package newPackage

import (
	"net/http"

	"github.com/joeybelans/gokismet/httpHandler"
)

// Discover
type templateData struct {
}

func HttpFunction(w http.ResponseWriter, req *http.Request) {
	// Add the header
	httpHandler.Header(w, "Title")
	httpTemplate.Execute(w, templateData{})
}

func httpSource() string {
	return `
	Body of html
</body>
</html>
`
}
