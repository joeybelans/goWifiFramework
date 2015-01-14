package kismetTemplate

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/joeybelans/gokismet/kismet"
)

// Home
type home struct {
	ServerVersion string
	ServerName    string
	StartTxt      string
	DBFile        string
	SSIDs         []string
	Stats         map[string]int
	Interfaces    []string
}

func HttpHome(w http.ResponseWriter, req *http.Request, dbfile string, ssids []string) {
	connected := ""
	if kismet.Connected() {
		connected = "checked"
	}

	templates["header"].Execute(w, header{"Home", req.URL.Path, connected})

	startInt, _ := strconv.ParseInt(kismet.ServerStart(), 10, 64)
	startTime := time.Unix(startInt, 0)
	hour, min, sec := startTime.Clock()
	startTxt := fmt.Sprintf("%02d:%02d:%02d", hour, min, sec)

	nCount, cCount, rCount, pCount, pRate, filtered := kismet.Stats()

	iNames := getInterfaces()

	templates["/"].Execute(w, home{kismet.ServerVersion(), kismet.ServerName(), startTxt, dbfile, ssids, map[string]int{"nCount": nCount, "cCount": cCount, "rCount": rCount, "pCount": pCount,
		"pRate": pRate, "filtered": filtered}, iNames})
}

func tmplHome() string {
	return `
<center>
<style>
span.stitle {
   font-weight: bold;
   font-size: medium;
   color: #004415;
}
table.stats {
   margin-left: 20px;
   border: 0;
}
table.stats th {
   font-weight: bold;
   font-size: small;
   color: #004415;
}
table.stats td {
   font-size: small;
   color: #004415;
}
table.data {
   margin-left: 20px;
   border: 0;
   border-spacing: 10px;
}
table.data th {
   font-weight: bold;
   font-size: medium;
   text-align: right;
}
table.data td {
   font-size: small;
}
ul.stats li {
   font-size: small;
   color: #004415;
}
a.stats:link { 
   color: #227755; 
} 
a.stats:visited { 
   color: #227755; 
} 
a.stats:active { 
   color: #227755; 
} 
a.stats:hover { 
	color: #116644; 
}
.verticalLine {
       border-right: 2px solid #ddeebb;
}
</style>
<table width="100%">
<tr><td width="30%">
<div class="verticalLine">
<span class="stitle">Kismet Server</span><br>
<table class="stats">
<tr><th align="left">Version</th><td>{{.ServerVersion}}</td></tr>
<tr><th align="left">Name</th><td>{{.ServerName}}</td></tr>
<tr><th align="left">Start Time</th><td>{{.StartTxt}}</td></tr>
</table>
<p>
<span class="stitle">gokismet Client</span><br>
<table class="stats">
<tr><th align="left" colspan="2">DB File</th></tr>
<tr><td align="left" colspan="2">{{.DBFile}}</td></tr>
</table>
<p>
<span class="stitle">In-scope Networks</span><br>
<blockquote><ul class="stats">{{range .SSIDs}}
<li><a class="stats" onClick='conn.send("statsSSID:{{.}}"); return false;' href="">{{.}}</a></li>{{end}}
</ul></blockquote>
<p>
<span class="stitle">Stats</span><br>
<table class="stats">
<tr><th align='left'>Network Count</th><td><div id="statsNcount">{{.Stats.nCount}}</div></td></tr>
<tr><th align='left'>Client Count</th><td><div id="statsCcount">{{.Stats.cCount}}</div></td></tr>
<tr><th align='left'>Rogue Count</th><td><div id="statsRcount">{{.Stats.rCount}}</div></td></tr>
<tr><th align='left'>Packet Count</th><td><div id="statsPcount">{{.Stats.pCount}}</div></td></tr>
<tr><th align='left'>Packet/Sec</th><td><div id="statsPrate">{{.Stats.pRate}}</div></td></tr>
<tr><th align='left'>Filtering</th><td><div id="statsFiltered">{{if eq .Stats.filtered 0}}No{{else}}Yes{{end}}</div></td></tr>
</table>
<p>
<span class="stitle">Network Interfaces</span><br>
<blockquote><ul class="stats">{{range .Interfaces}}
<li><a class="stats" onClick='conn.send("statsNIC:{{.}}"); return false;' href="">{{.}}</a></li>{{end}}
</ul></blockquote>
</div>
</td><td align="left" valign="top"><div name="wsOutput" id="wsOutput"></div></td></tr>
</table></center>
</body>
</html>
`
	//<li><a class="stats" href="/statsNIC/{{.}}">{{.}}</a></li>{{end}}
}
