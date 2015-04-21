package httpTemplate

import "net/http"

// Home
type home struct {
	Kismet     bool
	DBFile     string
	SSIDs      []string
	Interfaces []string
}

func HttpHome(w http.ResponseWriter, req *http.Request, webHost string, webPort int, kismetHost string, kismetPort int, dbfile string, outdir string, ssids []string) {
	templates["header"].Execute(w, header{"Home"})

	iNames := getInterfaces()

	templates["/"].Execute(w, home{true, dbfile, ssids, iNames})
}

func tmplHome() string {
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
<span class="ntitle">gokismet Client</span><br>
<table class="nav">
<tr><th align="left" colspan="2">DB File</th></tr>
<tr><td align="left" colspan="2">{{.DBFile}}</td></tr>
</table>
<p>
<span class="ntitle">In-scope Networks</span><br>
<blockquote><ul class="nav">{{range .SSIDs}}
<li><a class="nav" onClick='conn.send(JSON.stringify({message: "statsSSID", ssid: {{.}}})); return false;' href="">{{.}}</a></li>{{end}}
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
