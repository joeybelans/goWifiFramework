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
<table width='100%' height='100%'>
<tr><td width='20%' height='100%' valign='top'>
<div class='verticalLine'>
<table width='100%' height='100%'>
<tr><td>
<div id='interfaces'></div>
<span class='stitle'>Filters <a href='' onClick='clearFilters(); return (false);'><font size='-2'>(clear filters)</font></a></span><br>
<table class='stats'>
<tr><th align='right'>
<div onclick="displayDIV(this, event, 'divScope', 'block');" onmouseout="displayDIV(this, event, 'divScope', 'none');">Scope</div>
<div id="divScope" class="divFilter" onmouseover="displayDIV(this, event, 'divScope', 'block');" onmouseout="displayDIV(this, event, 'divScope', 'none');">
<b>Select the networks to display:</b>
<p>
<input type='radio' name='fscope' id='fscope' value='' checked>All</input> 
<input type='radio' name='fscope' id='fscope' value='in-scope'>In-scope</input> 
<input type='radio' name='fscope' id='fscope' value='rogue'>Rogues</input>
</div></th>
<td><div id='divScopeStatus' onclick="displayDIV(this, event, 'divScope', 'block');" onmouseout="displayDIV(this, event, 'divScope', 'none');">Show All</div>
</td></tr>
<tr><th align='right'>
<div onclick="displayDIV(this, event, 'divNetwork', 'block');" onmouseout="displayDIV(this, event, 'divNetwork', 'none');">Network</div>
<div id='divNetwork' class='divFilter' onmouseover="displayDIV(this, event, 'divNetwork', 'block');" onmouseout="displayDIV(this, event, 'divNetwork', 'none');"></div>
</th>
<td><div id='divNetworkStatus' onclick="displayDIV(this, event, 'divNetwork', 'block');" onmouseout="displayDIV(this, event, 'divNetwork', 'none');">Show All</div>
</td></tr>
<tr><th align='right'>
<div onclick="displayDIV(this, event, 'divChannel', 'block');" onmouseout="displayDIV(this, event, 'divChannel', 'none');">Channel</div>
<div id='divChannel' class='divFilter' onmouseover="displayDIV(this, event, 'divChannel', 'block');" onmouseout="displayDIV(this, event, 'divChannel', 'none');"></div>
</th>
<td><div id='divChannelStatus' onclick="displayDIV(this, event, 'divChannel', 'block');" onmouseout="displayDIV(this, event, 'divChannel', 'none');">Show All</div>
</td></tr>
<tr><th align='right'>
<div onclick="displayDIV(this, event, 'divBSSID', 'block');" onmouseout="displayDIV(this, event, 'divBSSID', 'none');">BSSID</div>
<div id='divBSSID' class='divFilter' onmouseover="displayDIV(this, event, 'divBSSID', 'block');" onmouseout="displayDIV(this, event, 'divBSSID', 'none');">
<b>Enter a single BSSID, or portion of, on each line:</b><br>
<font size='-2'>Ex. 11:22:33:44:55:66 or 11:22:33:</font>
<p>
<textarea name='fbssid' id='fbssid' rows='10' cols='30'></textarea>
</div></th>
<td><div id='divBSSIDStatus' onclick="displayDIV(this, event, 'divBSSID', 'block');" onmouseout="displayDIV(this, event, 'divBSSID', 'none');">Show All</div>
</td></tr>
</table>
<p>
<span class='stitle'>Notes</span><br>
<table class='stats'>
<tr><th>Location</th><td><input type='text' name='location' size='30'></td></tr>
<tr><td colspan='2'><textarea name='note' rows='10' cols='42'></textarea></td></tr>
<tr><td colspan='2'><center><input type='submit' name='add' value='Add Note'></td></tr>
</table>
</td></tr></table>
</div>
</td><td align='left' valign='top'>
<div class='apDetails' name='apDetails' id='apDetails' onClick="this.style.display = 'none';"></div>
<p>
<div name='wsOutput' id='wsOutput'></div>
<p>
<div name='clientOutput' id='clientOutput'></div>
</td></tr>
</table></center>
</body>
</html>
`
}
