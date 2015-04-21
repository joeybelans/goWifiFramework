// Manages home interface
// HTTP Template
package home

import (
	"net/http"

	"github.com/joeybelans/gokismet/header"
)

// Template data
type templateData struct {
	Kismet     bool
	SSIDs      []string
	Interfaces []string
}

// Home HTTP handler function
func HttpHome(w http.ResponseWriter, req *http.Request, webHost string, webPort int, kismetHost string, kismetPort int, ssids []string) {
	header.Display(w, "Home", "home")

	iNames := getInterfaces()

	home.Execute(w, templateData{true, ssids, iNames})
}

// Home source code
func templateSource() string {
	return `
<center>
<table width="100%">
<tr><td width="20%">
<div class="verticalLine">
<span class="ntitle">Dependencies</span><br>
<table class="nav">
<tr><th align="left">Kismet</th><td>{{if .Kismet}}Found{{else}}<font color="#ff0000">Not Found</font>{{end}}</td></tr>
<tr><th align="left">Kismograph</th><td></td></tr>
<tr><th align="left">EvilAP</th><td></td></tr>
<tr><th align="left">Airodump</th><td></td></tr>
</table>
<p>
<span class="ntitle">In-scope Networks</span><br>
<blockquote><ul class="nav">{{range .SSIDs}}
<!-- <li><a class="nav" onClick='conn.send(JSON.stringify({message: "statsSSID", ssid: {{.}}})); return false;' href="">{{.}}</a></li> -->
<li><a class="nav" onClick='getSSID("{{.}}"); return false;' href="">{{.}}</a></li>{{end}}
</ul></blockquote>
<p>
<span class="ntitle">Network Interfaces</span><br>
<blockquote><ul class="nav">{{range .Interfaces}}
<li><a class="nav" onClick='conn.send(JSON.stringify({type: "kismet", cmd: "GetNicStats", nic: {{.}}})); return false;' href="">{{.}}</a></li>{{end}}
</ul></blockquote>
</div>
</td><td align="left" valign="top"><div name="wsOutput" id="wsOutput"></div></td></tr>
</table></center>
</body>
</html>
`
}
