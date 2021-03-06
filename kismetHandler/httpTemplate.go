package kismetHandler

import (
	"net/http"

	"github.com/joeybelans/gokismet/header"
)

// Discover
type kismetData struct {
	SSIDs []string
}

//func HttpKismet(w http.ResponseWriter, req *http.Request, dbfile string, ssids []string) {
func HttpKismet(w http.ResponseWriter, req *http.Request, ssids []string) {
	header.Display(w, "Kismet", "kismet")
	httpTemplate.Execute(w, kismetData{ssids})
}

func httpSource() string {
	return `
<center>
<form>
<table width='100%' height='100%'>
<tr><td width='20%' height='100%' valign='top'>
<div class='verticalLine'>
<table width='100%' height='100%'>
<tr><td>
<div id='interfaces'></div>
<span class="ntitle">Network Interfaces</span><br>
<blockquote><ul class="nav">
<li><a class="nav" onClick='conn.send(JSON.stringify({type: "kismet", cmd: "GetNicStats", nic: "test"})); return false;' href="">test</a></li>
</ul></blockquote>
<p>
<span class="ntitle">Stats</span><br>
<table class="nav">
<tr><th align='left'>Network Count</th><td><div id="networkCount">0</div></td></tr>
<tr><th align='left'>Rogue Count</th><td><div id="rogueCount">0</div></td></tr>
<tr><th align='left'>Client Count</th><td><div id="clientCount">0</div></td></tr>
<tr><th align='left'>AP Count</th><td><div id="apCount">0</div></td></tr>
<tr><th align='left'>Packet/Sec</th><td><div id="packetRate">0</div></td></tr>
<tr><th align='left'>Crypted Packets</th><td><div id="cryptedCount">0</div></td></tr>
<tr><th align='left'>Packet Count</th><td><div id="packetCount">0</div></td></tr>
<tr><th align='left'>Management</th><td><div id="mgmtCount">0</div></td></tr>
<tr><th align='left'>Data</th><td><div id="dataCount">0</div></td></tr>
<tr><th align='left'>Dropped Packets</th><td><div id="droppedCount">0</div></td></tr>
<tr><th align='left'>Filtered</th><td><div id="filteredCount">0</div></td></tr>
</table>
<p>
<span class='ntitle'>Filters <a href='' onClick='clearFilters(); return (false);'><font size='-2'>(clear filters)</font></a></span><br>
<div id='divFilter' class='divFilter'></div>
<table class='nav'>
<tr><th align='right'><a id='filterScope' class='filter' href='' onClick="kismet.displayDIV(this); return false;">Scope</a></th><td><div id='divScopeStatus'>Show All</div></td></tr>
<tr><th align='right'><a id='filterSSID' class='filter' href='' onClick="kismet.displayDIV(this); return false;">SSID</a></th><td><div id='divSSIDStatus'>Show All</div></td></tr>
<tr><th align='right'><a id='filterChannel' class='filter' href='' onClick="kismet.displayDIV(this); return false;">Channel</a></th><td><div id='divChannelStatus'>Show All</div></td></tr>
<tr><th align='right'><a id='filterBSSID' class='filter' href='' onClick="kismet.displayDIV(this); return false;">BSSID</a></th><td><div id='divBSSIDStatus'>Show All</div></td></tr>
</table>
<p>
<span class='ntitle'>Notes</span><br>
<table class='nav'>
<tr><th>Location</th><td><input type='text' name='location' size='30'></td></tr>
<tr><td colspan='2'><textarea name='note' rows='10' cols='42'></textarea></td></tr>
<tr><td colspan='2'><center><input type='submit' name='add' value='Add Note'></td></tr>
</table>
</td></tr></table>
</div>
</td><td align='left' valign='top'>
<div class='apDetails' name='apDetails' id='apDetails' onClick="this.style.display = 'none';"></div>
<p>
<div name='apList' id='apList'></div>
<p>
<div name='clientOutput' id='clientOutput'></div>
</td></tr>
</table></center>
</body>
</html>
`
}
