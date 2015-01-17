package kismetTemplate

import (
	"net/http"

	"github.com/joeybelans/gokismet/kismet"
)

// Discover
type discover struct {
	DBFile string
	SSIDs  []string
}

func HttpDiscover(w http.ResponseWriter, req *http.Request, dbfile string, ssids []string) {
	connected := ""
	if kismet.Connected() {
		connected = "checked"
	}

	templates["header"].Execute(w, header{"Discover", req.URL.Path, connected})
	templates["/discover"].Execute(w, discover{dbfile, ssids})
}

func tmplDiscover() string {
	return `
<center>
<form>
<table width="100%">
<tr><td width="15%">
<div class="verticalLine">
<span class="stitle">Filters</span><br>
<table class="stats">
<tr><th>Scope</th><td><select name='fscope'>
<option value=''>All</option>
<option value='in'>In-scope</option>
<option value='out'>Out-of-scope</option>
</select></td></tr>
<tr><th>Network</th><td><select name='fnetwork' id='fnetwork' onChange='filterNetwork=this.value; displayNetworks();'></select></td></tr>
<tr><th>Channel</th><td><select name='fchannel' id='fchannel' onChange='filterChannel=this.value; displayNetworks();'></select></td></tr>
<tr><th>BSSID</th><td><input type='text' name='fbssid' size='18'></td></tr>
</table>
<p>
<span class="stitle">Notes</span><br>
<table class="stats">
<tr><th>Location</th><td><input type='text' name='location' size='30'></td></tr>
<tr><td colspan='2'><textarea rows='10' cols='42'></textarea></td></tr>
<tr><td colspan='2'><center><input type='submit' name='add' value='Add Note'></td></tr>
</table>
</div>
</td><td align="left" valign="top"><div name="wsOutput" id="wsOutput"></div></td></tr>
</table></center>
</body>
</html>
`
	//<li><a class="stats" href="/statsNIC/{{.}}">{{.}}</a></li>{{end}}
}
