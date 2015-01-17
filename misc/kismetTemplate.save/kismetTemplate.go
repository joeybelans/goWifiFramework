package kismetTemplate

import (
	"fmt"
	"html/template"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joeybelans/gokismet/kismet"
)

type page struct {
	URL   string
	Title string
	iface interface{}
}

type header struct {
	Title     string
	Path      string
	Connected string
}

type home struct {
	ServerVersion string
	ServerName    string
	StartTxt      string
	DBFile        string
	SSIDs         []string
	Stats         map[string]int
	Interfaces    []string
}

type discover struct {
	DBFile string
	SSIDs  []string
}

var templates map[string]*template.Template

func getInterfaces() []string {
	interfaces, _ := net.Interfaces()

	var iNames []string
	for _, iface := range interfaces {
		if _, err := os.Stat("/sys/class/net/" + iface.Name + "/wireless"); err == nil {
			iNames = append(iNames, iface.Name)
		}
	}

	return iNames
}

func init() {
	pages := map[int]page{
		0: page{URL: "/", Title: "Home", iface: tmplHome},
		1: page{URL: "/discover", Title: "Discover", iface: tmplDiscover},
		2: page{URL: "/profile", Title: "Profile", iface: nil},
		3: page{URL: "/networks", Title: "Networks", iface: nil},
		4: page{URL: "/aps", Title: "Access Points", iface: nil},
		5: page{URL: "/clients", Title: "Clients", iface: nil},
		6: page{URL: "/reports", Title: "Reports", iface: nil},
		7: page{URL: "/logs", Title: "Logs", iface: nil},
	}

	createPages(pages)
	templates["header"] = template.New("header")
	templates["header"], _ = templates["header"].Parse(tmplHeader(pages))
}

// Home
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

// Discover
func HttpDiscover(w http.ResponseWriter, req *http.Request, dbfile string, ssids []string) {
	connected := ""
	if kismet.Connected() {
		connected = "checked"
	}

	templates["header"].Execute(w, header{"Discover", req.URL.Path, connected})
	templates["/discover"].Execute(w, discover{dbfile, ssids})
}

func createPages(pages map[int]page) {
	templates = map[string]*template.Template{}

	for index := range pages {
		page := pages[index]
		if page.iface != nil {
			f := page.iface
			templates[page.URL] = template.New(page.URL)
			templates[page.URL], _ = templates[page.URL].Parse(f.(func() string)())
		}
	}
}

func tmplHeader(pages map[int]page) string {
	txt := `
<html>
<head>
<title>GoKismet - {{.Title}}</title>
<style type='text/css'>
label {
	background: transparent;
	border: 0;
	margin: 0;
	padding: 0;
	vertical-align: baseline;
}
#wrapper {
	min-width: 600px;
}
.settings {
	display: table;
	height: 45px;
	width: 100%;
}
.settings .switch {
	display: table-cell;
	vertical-align: middle;
	padding: 0;
}
.cmn-toggle {
	position: absolute;
	margin-left: -9999px;
	visibility: hidden;
}
.cmn-toggle + label {
	display: block;
	position: relative;
	cursor: pointer;
	outline: none;
	-webkit-user-select: none;
	-moz-user-select: none;
	-ms-user-select: none;
	user-select: none;
}
input.cmn-toggle-on-off + label {
	padding: 0;
	width: 150px;
	height: 45px;
}
input.cmn-toggle-on-off + label:before, input.cmn-toggle-on-off + label:after {
	display: block;
	position: absolute;
	top: 0;
	left: 0;
	bottom: 0;
	right: 0;
	font-family: "Roboto Slab", serif;
	font-size: 16px;
	text-align: center;
	font-weight: bold; 
	line-height: 45px;
}
input.cmn-toggle-on-off + label:before {
	color: #ff0000;
	background-color: #dddddd;
	content: attr(data-off);
	-webkit-transition: -webkit-transform 0.5s;
	-moz-transition: -moz-transform 0.5s;
	-o-transition: -o-transform 0.5s;
	transition: transform 0.5s;
	-webkit-backface-visibility: hidden;
	-moz-backface-visibility: hidden;
	-ms-backface-visibility: hidden;
	-o-backface-visibility: hidden;
	backface-visibility: hidden;
}
input.cmn-toggle-on-off + label:after {
	color: #004415; 
	background-color: #ddeebb; 
	content: attr(data-on);
	-webkit-transition: -webkit-transform 0.5s;
	-moz-transition: -moz-transform 0.5s;
	-o-transition: -o-transform 0.5s;
	transition: transform 0.5s;
	-webkit-transform: rotateY(180deg);
	-moz-transform: rotateY(180deg);
	-ms-transform: rotateY(180deg);
	-o-transform: rotateY(180deg);
	transform: rotateY(180deg);
	-webkit-backface-visibility: hidden;
	-moz-backface-visibility: hidden;
	-ms-backface-visibility: hidden;
	-o-backface-visibility: hidden;
	backface-visibility: hidden;
}
input.cmn-toggle-on-off:checked + label:before {
	-webkit-transform: rotateY(180deg);
	-moz-transform: rotateY(180deg);
	-ms-transform: rotateY(180deg);
	-o-transform: rotateY(180deg);
	transform: rotateY(180deg);
}
input.cmn-toggle-on-off:checked + label:after {
	-webkit-transform: rotateY(0);
	-moz-transform: rotateY(0);
	-ms-transform: rotateY(0);
	-o-transform: rotateY(0);
	transform: rotateY(0);
}
div#hmenu { 
	margin: 0; 
	padding: .3em 0 .3em 0; 
	background: #ddeebb; 
	width: 100%; 
	height: 35px;
	text-align: center; 
} 
div#hmenu ul { 
	list-style: none; 
	margin: 0; 
	padding: 0; 
} 
div#hmenu ul li { 
	margin: 0; 
	padding: 0; 
	display: inline; 
} 
div#hmenu ul li.current-page {
	margin: 0; 
	padding: 0; 
	display: inline; 
	text-decoration: none; 
	font-weight: bold; 
	font-size: x-large; 
	color: #004415; 
} 
div#hmenu ul a:link{ 
	margin: 0; 
	padding: .3em .4em .3em .4em; 
	text-decoration: none; 
	font-weight: bold; 
	font-size: medium; 
	color: #004415; 
} 
div#hmenu ul a:visited{ 
	margin: 0; 
	padding: .3em .4em .3em .4em; 
	text-decoration: none; 
	font-weight: bold; 
	font-size: medium; 
	color: #004415; 
} 
div#hmenu ul a:active{ 
	margin: 0; 
	padding: .3em .4em .3em .4em; 
	text-decoration: none; 
	font-weight: bold; 
	font-size: medium; 
	color: #227755; 
} 
div#hmenu ul a:hover{ 
	margin: 0; 
	padding: .3em .4em .3em .4em; 
	text-decoration: none; 
	font-weight: bold; 
	font-size: medium; 
	color: #f6f0cc; 
	background-color: #227755; 
}
</style>
<script>
function kismetOnOff() {
   if (document.getElementById("cmn-toggle-7").checked == false) {
      conn.send("kismetDISCONNECT")
   } else {
      conn.send("kismetCONNECT")
   }
}

function kismetAddSource(nic) {
   conn.send("nicADDSOURCE:" + nic + ":" + document.getElementsByName("alias")[0].value);
   setTimeout('conn.send("statsNIC:' + nic + '")', 3000);
}

function kismetDelSource(nic) {
   conn.send("nicDELSOURCE:" + nic);
   setTimeout('conn.send("statsNIC:' + nic + '")', 3000);
}

function kismetParseInfo(msg) {
   var fields = msg.split(":");

   document.getElementById("statsNcount").innerHTML = fields[0];
   document.getElementById("statsCcount").innerHTML = fields[1];
   document.getElementById("statsRcount").innerHTML = fields[2];
   document.getElementById("statsPcount").innerHTML = fields[3];
   document.getElementById("statsPrate").innerHTML = fields[4];
   document.getElementById("statsFiltered").innerHTML = fields[5];
}

function kismetParseTerminate(msg) {
   document.getElementById("cmn-toggle-7").checked = false;
}

function ssidINFO(msg) {
   var fields = msg.split(";");

   divtxt = "<form><table class='data'>\
<tr><th style='text-align: left' colspan='2' bgcolor='#000000'><font color='#ffffff'>" + fields[0] + "</font></th></tr>\
<tr><th>Cloaked</th><td>" + fields[1] + "</td></tr>\
<tr><th>Channels</th><td>" + fields[10] + "</td></tr>\
<tr><th>First</th><td>" + fields[2] + "</td></tr>\
<tr><th>Last</th><td>" + fields[3] + "</td></tr>\
<tr><th>Max Rate</th><td>" + fields[4] + "</td></tr>\
<tr><th>Min DBM</th><td>" + fields[5] + "</td></tr>\
<tr><th>Max DBM</th><td>" + fields[6] + "</td></tr>\
<tr><th>Client Count</th><td>" + fields[7] + "</td></tr>\
<tr><th>BSSID Count</th><td>" + fields[8] + "</td></tr>\
<tr><th>Encryption</th><td>" + fields[9] + "</td></tr>\
</table></form>";

   document.getElementById("wsOutput").innerHTML = divtxt;
}

function nicINFO(msg) {
   var fields = msg.split(";");

   if (fields[1] == 0) {
      fields[1] = "inactive";
   } else if (fields[1] == 1) {
      fields[1] = "active";
   } else {
      fields[1] = fields[2];
   }

   if (fields[5] == 0) {
      fields[5] = "No";
   } else if ((fields[5] == 1) || (fields[5] == 3)) {
      fields[5] = "Yes";
   }

   divtxt = "<form><table class='data'>\
<tr><th style='text-align: left' colspan='2' bgcolor='#000000'><font color='#ffffff'>" + fields[0] + "</font></th></tr>\
<tr><td colspan='2'>Active<input type='radio' name='source' value='active'";

   if (fields[1] != 'inactive') {
      divtxt = divtxt + " checked";
   }  else {
      divtxt = divtxt + " onChange='kismetAddSource(\"" + fields[0] + "\")'";
   }

   if ((fields[1] != 'active') && (fields[1] != 'inactive')) {
      divtxt = divtxt + " disabled";
   }

   divtxt = divtxt + "> Inactive<input type='radio' name='source' value='inactive'";
   
   if (fields[1] == 'inactive') {
      divtxt = divtxt + " checked";
   }  else {
      divtxt = divtxt + " onChange='kismetDelSource(\"" + fields[0] + "\")'";
   } 

   if ((fields[1] != 'active') && (fields[1] != 'inactive')) {
      divtxt = divtxt + " disabled";
   } 

   divtxt = divtxt + ">"

   if ((fields[1] != 'active') && (fields[1] != 'inactive')) {
      divtxt = divtxt + " <font size='-2' color='#ff0000'>Via " + fields[1] + "</font>";
   }

   divtxt = divtxt + "</td></tr>";

   if ((fields[1] == 'inactive') || ((fields[3] != '') && (fields[2] != fields[3]))) {
      divtxt = divtxt + "<tr><th>Alias</th><td>";
      
      if (fields[1] == 'inactive') {
	 divtxt = divtxt + "<input type='text' name='alias' value=''></input>";
      } else {
	 divtxt = divtxt + fields[3];
      }
      
      divtxt = divtxt + "</td></tr>";
   }

   if (fields[4] != null) {
      divtxt = divtxt + "<tr><th>WiFi Channel</th><td>" + fields[4] + "</td></tr>\
<tr><th>Channel Hopping</th><td>" + fields[5] + "</td></tr>\
<tr><th>Rate</th><td>" + fields[6] + " channels/sec</td></tr>\
<tr><th>Channel List</th><td>" + fields[7] + "</td></tr>"
   }

   divtxt = divtxt + "</table></form>"

   document.getElementById("wsOutput").innerHTML = divtxt;
}

wsServer = window.location.hostname
if (window.location.port != "") {
   wsServer = wsServer + ":" + window.location.port;
}
conn = new WebSocket("ws://" + wsServer + "/ws");

conn.onopen = function (event) {
   conn.send("{{.Path}}");
};

conn.onclose = function (event) {
   document.getElementById("cmn-toggle-7").checked = false;
}

conn.onmessage = function(evt) {
   var i = evt.data.indexOf(':');
   window[evt.data.slice(0,i)](evt.data.slice(i+1));
};
</script>
</head>
<body>
<table border="0" width="100%" cellspacing="0" cellpadding="0"><tr><td width="5%">
<div class="switch">
<input id="cmn-toggle-7" class="cmn-toggle cmn-toggle-on-off" type="checkbox" onClick="kismetOnOff()" {{.Connected}}>
<label for="cmn-toggle-7" data-on="Connected" data-off="Disconnected"></label>
</div>
</td><td>
<div id="hmenu"> 
<ul>`
	for i := 0; i < len(pages); i++ {
		txt = txt + "<li{{ if eq .Title \"" + pages[i].Title + "\" }} class=\"current-page\">" + pages[i].Title + "{{ else }}><a href='" + pages[i].URL + "'>" + pages[i].Title + "</a>{{ end }}</li>\n"
	}
	txt = txt + "</ul></div></td></tr></table><p>\n"
	return txt
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

func tmplDiscover() string {
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
   border-collapse:collapse;
}
table.data th {
   font-weight: bold;
   font-size: medium;
   text-align: left;
   padding-left: 10px;
   padding-right: 10px;
   padding-top: 5px;
   padding-bottom: 5px;
}
table.data td {
   padding-left: 10px;
   padding-right: 10px;
   padding-top: 5px;
   padding-bottom: 5px;
}
table.data tr:nth-child(odd) td {
   font-size: small;
   background: #ffffff;
}
table.data tr:nth-child(even) td {
   font-size: small;
   background: #f5f9fa;
}
table.data tr:hover td {
   font-size: small;
   background: #fffbae;
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
a.data:link { 
   color: #000000
} 
a.data:visited { 
   color: #000000; 
} 
a.data:active { 
   color: #000000; 
} 
a.data:hover { 
   color: #000000; 
}
.verticalLine {
   border-right: 2px solid #ddeebb;
}
.selected {
   background: #ff0000;
}
</style>
<script>
var networks = new Object();
var sort = 'network';
var order = 'asc';
var filterScope = '';
var filterNetwork = '';
var filterChannel = '';
var filterBSSID = '';

function setSort(col) {
   if (sort == col) {
      if (order == 'asc') {
	 order = 'des';
      } else {
	 order = 'asc';
      }
   } else {
      sort = col;
      order = 'asc';
   } 

   displayNetworks();
}

function displaySortSSID() {
   var msg = ""

   var nKeys = Object.keys(networks);
   var nLen = nKeys.length;
   nKeys.sort();
   if (order == 'des') {
      nKeys.reverse();
   }
   for (var i = 0; i < nLen; i++) {
      var bKeys = Object.keys(networks[nKeys[i]].bssids);
      var bLen = bKeys.length;
      bKeys.sort();
      for (var j = 0; j < bLen; j++) {
	 var b = networks[nKeys[i]].bssids[bKeys[j]]

	 if (((filterNetwork == '') || (filterNetwork == nKeys[i])) &&
	     ((filterChannel == '') || (filterChannel == b.channel))) {
            msg = msg + "<tr><td>" + nKeys[i] + "</td><td>" + bKeys[j] + "</td><td>" + b.channel + "</td><td>" + b.lasttime + "</td><td>" + b.power +
	 	"</td><td>" + b.max + "</td><td>" + b.clients + "</td><td>" + b.packets + "</td></tr>\n";
	 }
      }
   }
   return (msg);
}

function displaySortBSSID() {
   var msg = "";
   var sBSSIDS = new Object();

   for (var ssid in networks) {
      for (var bssid in networks[ssid].bssids) {
	 sBSSIDS[bssid] = ssid;
      }
   }

   var sKeys = Object.keys(sBSSIDS);
   var sLen = sKeys.length;
   sKeys.sort();
   if (order == 'des') {
      sKeys.reverse();
   }
   for (var i = 0; i < sLen; i++) {
      var b = networks[sBSSIDS[sKeys[i]]].bssids[sKeys[i]];
      msg = msg + "<tr><td>" + sBSSIDS[sKeys[i]] + "</td><td>" + sKeys[i] + "</td><td>" + b.channel + "</td><td>" + b.lasttime + "</td><td>" + b.power +
	 	"</td><td>" + b.max + "</td><td>" + b.clients + "</td><td>" + b.packets + "</td></tr>\n";
   }
   return (msg);
}

function displaySortCHANNEL() {
   var msg = "";
   var pad = "000000";
   var sCHANNEL = new Object();

   for (var ssid in networks) {
      for (var bssid in networks[ssid].bssids) {
         cKey = (pad+networks[ssid].bssids[bssid].channel).slice(-pad.length)
         if (!(cKey in sCHANNEL)) {
            sCHANNEL[cKey] = new Object();
            sCHANNEL[cKey].networks = new Object();
         }
         sCHANNEL[cKey].channel = networks[ssid].bssids[bssid].channel;
         sCHANNEL[cKey].networks[ssid] = '';
      }
   }

   var cKeys = Object.keys(sCHANNEL);
   var cLen = cKeys.length;
   cKeys.sort();
   if (order == 'des') {
      cKeys.reverse();
   }
   for (var i = 0; i < cLen; i++) {
      var nKeys = Object.keys(sCHANNEL[cKeys[i]].networks);
      var nLen = nKeys.length;
      nKeys.sort();

      for (var j = 0; j < nLen; j++) {
         var bKeys = Object.keys(networks[nKeys[j]].bssids);
         var bLen = bKeys.length;
         bKeys.sort();

	 for (var k = 0; k < bLen; k++) {
            var b = networks[nKeys[j]].bssids[bKeys[k]];
	    if (b.channel == sCHANNEL[cKeys[i]].channel) {
               msg = msg + "<tr><td>" + nKeys[j] + "</td><td>" + bKeys[k] + "</td><td>" + b.channel + "</td><td>" + b.lasttime + "</td><td>" + b.power +
		  "</td><td>" + b.max + "</td><td>" + b.clients + "</td><td>" + b.packets + "</td></tr>\n";
	    }
	 }
      }
   }
   return (msg);
}

function displaySortLASTSEEN() {
   var msg = "";
   var sTIME = new Object();

   for (var ssid in networks) {
      for (var bssid in networks[ssid].bssids) {
         if (!(networks[ssid].bssids[bssid].lasttime in sTIME)) {
            sTIME[networks[ssid].bssids[bssid].lasttime] = new Object();
         }
         sTIME[networks[ssid].bssids[bssid].lasttime][ssid] = '';
      }
   }

   var tKeys = Object.keys(sTIME);
   var tLen = tKeys.length;
   tKeys.sort();
   if (order == 'asc') {
      tKeys.reverse();
   }
   for (var i = 0; i < tLen; i++) {
      var nKeys = Object.keys(sTIME[tKeys[i]]);
      var nLen = nKeys.length;
      nKeys.sort();

      for (var j = 0; j < nLen; j++) {
         var bKeys = Object.keys(networks[nKeys[j]].bssids);
         var bLen = bKeys.length;
         bKeys.sort();

	 for (var k = 0; k < bLen; k++) {
            var b = networks[nKeys[j]].bssids[bKeys[k]];
	    if (b.lasttime == tKeys[i]) {
               msg = msg + "<tr><td>" + nKeys[j] + "</td><td>" + bKeys[k] + "</td><td>" + b.channel + "</td><td>" + b.lasttime + "</td><td>" + b.power +
		  "</td><td>" + b.max + "</td><td>" + b.clients + "</td><td>" + b.packets + "</td></tr>\n";
	    }
	 }
      }
   }
   return (msg);
}

function displaySortPOWER() {
   var msg = "";
   var pad = "000000";
   var sPOWER = new Object();

   for (var ssid in networks) {
      for (var bssid in networks[ssid].bssids) {
         pKey = (pad+networks[ssid].bssids[bssid].power).slice(-pad.length)
         if (!(pKey in sPOWER)) {
            sPOWER[pKey] = new Object();
            sPOWER[pKey].networks = new Object();
         }
         sPOWER[pKey].power = networks[ssid].bssids[bssid].power;
         sPOWER[pKey].networks[ssid] = '';
      }
   }

   var pKeys = Object.keys(sPOWER);
   var pLen = pKeys.length;
   pKeys.sort();
   if (order == 'des') {
      pKeys.reverse();
   }
   for (var i = 0; i < pLen; i++) {
      var nKeys = Object.keys(sPOWER[pKeys[i]].networks);
      var nLen = nKeys.length;
      nKeys.sort();

      for (var j = 0; j < nLen; j++) {
         var bKeys = Object.keys(networks[nKeys[j]].bssids);
         var bLen = bKeys.length;
         bKeys.sort();

	 for (var k = 0; k < bLen; k++) {
            var b = networks[nKeys[j]].bssids[bKeys[k]];
	    if (b.power == sPOWER[pKeys[i]].power) {
               msg = msg + "<tr><td>" + nKeys[j] + "</td><td>" + bKeys[k] + "</td><td>" + b.channel + "</td><td>" + b.lasttime + "</td><td>" + b.power +
		  "</td><td>" + b.max + "</td><td>" + b.clients + "</td><td>" + b.packets + "</td></tr>\n";
	    }
	 }
      }
   }
   return (msg);
}

function displaySortMAX() {
   var msg = "";
   var pad = "000000";
   var sMAX = new Object();

   for (var ssid in networks) {
      for (var bssid in networks[ssid].bssids) {
         mKey = (pad+networks[ssid].bssids[bssid].max).slice(-pad.length)
         if (!(mKey in sMAX)) {
            sMAX[mKey] = new Object();
            sMAX[mKey].networks = new Object();
         }
         sMAX[mKey].max = networks[ssid].bssids[bssid].max;
         sMAX[mKey].networks[ssid] = '';
      }
   }

   var mKeys = Object.keys(sMAX);
   var mLen = mKeys.length;
   mKeys.sort();
   if (order == 'des') {
      mKeys.reverse();
   }
   for (var i = 0; i < mLen; i++) {
      var nKeys = Object.keys(sMAX[mKeys[i]].networks);
      var nLen = nKeys.length;
      nKeys.sort();

      for (var j = 0; j < nLen; j++) {
         var bKeys = Object.keys(networks[nKeys[j]].bssids);
         var bLen = bKeys.length;
         bKeys.sort();

	 for (var k = 0; k < bLen; k++) {
            var b = networks[nKeys[j]].bssids[bKeys[k]];
	    if (b.max == sMAX[mKeys[i]].max) {
               msg = msg + "<tr><td>" + nKeys[j] + "</td><td>" + bKeys[k] + "</td><td>" + b.channel + "</td><td>" + b.lasttime + "</td><td>" + b.power +
		  "</td><td>" + b.max + "</td><td>" + b.clients + "</td><td>" + b.packets + "</td></tr>\n";
	    }
	 }
      }
   }
   return (msg);
}

function displaySortCLIENTS() {
   var msg = "";
   var pad = "000000";
   var sCLIENTS = new Object();

   for (var ssid in networks) {
      for (var bssid in networks[ssid].bssids) {
         cKey = (pad+networks[ssid].bssids[bssid].clients).slice(-pad.length)
         if (!(cKey in sCLIENTS)) {
            sCLIENTS[cKey] = new Object();
            sCLIENTS[cKey].networks = new Object();
         }
         sCLIENTS[cKey].clients = networks[ssid].bssids[bssid].clients;
         sCLIENTS[cKey].networks[ssid] = '';
      }
   }

   var cKeys = Object.keys(sCLIENTS);
   var cLen = cKeys.length;
   cKeys.sort();
   if (order == 'des') {
      cKeys.reverse();
   }
   for (var i = 0; i < cLen; i++) {
      var nKeys = Object.keys(sCLIENTS[cKeys[i]].networks);
      var nLen = nKeys.length;
      nKeys.sort();

      for (var j = 0; j < nLen; j++) {
         var bKeys = Object.keys(networks[nKeys[j]].bssids);
         var bLen = bKeys.length;
         bKeys.sort();

	 for (var k = 0; k < bLen; k++) {
            var b = networks[nKeys[j]].bssids[bKeys[k]];
	    if (b.clients == sCLIENTS[cKeys[i]].clients) {
               msg = msg + "<tr><td>" + nKeys[j] + "</td><td>" + bKeys[k] + "</td><td>" + b.channel + "</td><td>" + b.lasttime + "</td><td>" + b.power +
		  "</td><td>" + b.max + "</td><td>" + b.clients + "</td><td>" + b.packets + "</td></tr>\n";
	    }
	 }
      }
   }
   return (msg);
}

function displaySortPACKETS() {
   var msg = "";
   var pad = "000000";
   var sPACKETS = new Object();

   for (var ssid in networks) {
      for (var bssid in networks[ssid].bssids) {
         pKey = (pad+networks[ssid].bssids[bssid].packets).slice(-pad.length)
         if (!(pKey in sPACKETS)) {
            sPACKETS[pKey] = new Object();
            sPACKETS[pKey].networks = new Object();
         }
         sPACKETS[pKey].packets = networks[ssid].bssids[bssid].packets;
         sPACKETS[pKey].networks[ssid] = '';
      }
   }

   var pKeys = Object.keys(sPACKETS);
   var pLen = pKeys.length;
   pKeys.sort();
   if (order == 'des') {
      pKeys.reverse();
   }
   for (var i = 0; i < pLen; i++) {
      var nKeys = Object.keys(sPACKETS[pKeys[i]].networks);
      var nLen = nKeys.length;
      nKeys.sort();

      for (var j = 0; j < nLen; j++) {
         var bKeys = Object.keys(networks[nKeys[j]].bssids);
         var bLen = bKeys.length;
         bKeys.sort();

	 for (var k = 0; k < bLen; k++) {
            var b = networks[nKeys[j]].bssids[bKeys[k]];
	    if (b.packets == sPACKETS[pKeys[i]].packets) {
               msg = msg + "<tr><td>" + nKeys[j] + "</td><td>" + bKeys[k] + "</td><td>" + b.channel + "</td><td>" + b.lasttime + "</td><td>" + b.power +
		  "</td><td>" + b.max + "</td><td>" + b.clients + "</td><td>" + b.packets + "</td></tr>\n";
	    }
	 }
      }
   }
   return (msg);
}

function displayNetworks() {
   var msg = "<table onClick='alert(this.rows)' class='data'> \
<tr><th><a class='data' href='' onClick='setSort(\"network\"); return false;'>Network</a></th>\n \
<th><a class='data' href='' onClick='setSort(\"bssid\"); return false;'>BSSID</a></th>\n \
<th><a class='data' href='' onClick='setSort(\"channel\"); return false;'>Channel</a></th>\n \
<th><a class='data' href='' onClick='setSort(\"lastseen\"); return false;'>Last Seen</a></th>\n \
<th><a class='data' href='' onClick='setSort(\"power\"); return false;'>Power</a></th>\n \
<th><a class='data' href='' onClick='setSort(\"max\"); return false;'>Max</a></th>\n \
<th><a class='data' href='' onClick='setSort(\"clients\"); return false;'>Clients</a></th>\n \
<th><a class='data' href='' onClick='setSort(\"packets\"); return false;'>Packets</a></th></tr>\n";

   if (sort == 'network') {
      msg = msg + displaySortSSID();
   } else if (sort == 'bssid') {
      msg = msg + displaySortBSSID();
   } else if (sort == 'channel') {
      msg = msg + displaySortCHANNEL();
   } else if (sort == 'lastseen') {
      msg = msg + displaySortLASTSEEN();
   } else if (sort == 'power') {
      msg = msg + displaySortPOWER();
   } else if (sort == 'max') {
      msg = msg + displaySortMAX();
   } else if (sort == 'clients') {
      msg = msg + displaySortCLIENTS();
   } else if (sort == 'packets') {
      msg = msg + displaySortPACKETS();
   }

   msg = msg + "</table>\n";
   document.getElementById("wsOutput").innerHTML = msg;

   var opt = document.createElement('option');
   opt.value = '';
   opt.innerHTML = 'All';
   var fnetwork  = document.getElementById("fnetwork");
   while(fnetwork.options.length > 0){
      fnetwork.remove(0);
   }
   fnetwork.appendChild(opt);

   opt = document.createElement('option');
   opt.value = '';
   opt.innerHTML = 'All';
   var fchannel = document.getElementById("fchannel");
   while(fchannel.options.length > 0){
      fchannel.remove(0);
   }
   fchannel .appendChild(opt);

   var nKeys = Object.keys(networks);
   var nLen = nKeys.length;
   nKeys.sort();
   for (var i = 0; i < nLen; i++) {
      opt = document.createElement('option');
      opt.value = nKeys[i];
      opt.innerHTML = nKeys[i];

      if (filterNetwork == nKeys[i]) {
	 opt.selected = 'selected';
      }

      fnetwork.appendChild(opt);
   }

   var pad = "000000";
   var channels = new Object();
   for (var ssid in networks) {
      for (var bssid in networks[ssid].bssids) {
         cKey = (pad+networks[ssid].bssids[bssid].channel).slice(-pad.length)
         if (!(cKey in channels)) {
            channels[cKey] = networks[ssid].bssids[bssid].channel;
         }
      }
   }
   var cKeys = Object.keys(channels);
   var cLen = cKeys.length;
   cKeys.sort();
   for (var i = 0; i < cLen; i++) {
      if (channels[cKeys[i]] != '') {
         opt = document.createElement('option');
         opt.value = channels[cKeys[i]];
         opt.innerHTML = channels[cKeys[i]];

         if (filterChannel == channels[cKeys[i]]) {
	    opt.selected = 'selected';
         }

         fchannel.appendChild(opt);
      }
   }
}

function discoverNetworkBSSID(msg) {
   var fields = msg.split(";");

   for (var ssid in networks) {
      for (var bssid in networks[ssid].bssids) {
         if (bssid == fields[0]) {
	    networks[ssid].bssids[bssid].channel = fields[1];
	    networks[ssid].bssids[bssid].lasttime = fields[2];
	    networks[ssid].bssids[bssid].power = fields[3];
	    networks[ssid].bssids[bssid].clients = fields[4];
	    networks[ssid].bssids[bssid].max = fields[5];
	    networks[ssid].bssids[bssid].packets = fields[6];
            break;
	 }
      }
   }

   displayNetworks();
}

function discoverNetworkSSID(msg) {
   var fields = msg.split(";");

   if (!(fields[0] in networks)) {
      networks[fields[0]] = {
         lastseen: fields[2],
         bssids: {}
      };
   } else if (!(fields[1] in networks[fields[0]].bssids)) {
      networks[fields[0]].bssids[fields[1]] = {
         channel: '',
         lasttime: '',
         power: '',
	 max: '',
	 clients: '',
	 packets: ''
      };
   }

   displayNetworks();
}
</script>
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

/*
// Networks
func HttpNetworks(w http.ResponseWriter, req *http.Request) {
	HttpHeader(w, "Networks")
	io.WriteString(w, `
<form>
<center><table border="1">
<tr><th colspan="4">Kismet Options</th></tr>
<tr><th>Channel(s)</th><td><input type="text" name="channels" size="25"></td>
<th>Hop</th><td>Yes<input type="radio" name="hop" value="yes" checked> No<input type="radio" name="hop" value="no"></td></tr>
<tr><th>Networks</th><td colspan="2"><select name="networks" size="1">
<option></option>
<option>test</option>
</select></td><td rowspan="3"> </td></tr>
<tr><th>BSSIDs</th><td colspan="2"><select name="bssids" size="1">
<option></option>
<option>aa:bb:cc:dd:ee:ff</option>
</select></td></tr>
<tr><th>Clients</th><td colspan="2"><select name="clients" size="1">
<option></option>
<option>Client1</option>
</select></td></tr>
</table></center>
</form>
`)
}

// Clients
func HttpClients(w http.ResponseWriter, req *http.Request) {
	HttpHeader(w, "Clients")
	io.WriteString(w, `
<form>
<center><table border="1">
<tr><th>Client</th><td><select name="client" size="1">
<option></option>
<option>Client1</option>
</select></td></tr>
<tr><th>MAC</th><td>aa:bb:cc:dd:ee:ff</td></tr>
<tr><th>Network</th><td>test</td></tr>
<tr><th>Probes</th><td>Probe1</td></tr>
</table></center>
</form>
`)
}

// Rogues
func HttpRogues(w http.ResponseWriter, req *http.Request) {
	HttpHeader(w, "Rogues")
	io.WriteString(w, `
<form>
<center><table border="1">
<tr><th>Network</th><td><select name="network" size="1">
<option></option>
<option>test</option>
</select></td></tr>
<tr><th>BSSIDs</th><td colspan="2"><select name="bssids" size="1">
<option></option>
<option>aa:bb:cc:dd:ee:ff</option>
</select></td></tr>
<tr><th>Clients</th><td colspan="2"><select name="clients" size="1">
<option></option>
<option>Client1</option>
</select></td></tr>
</table></center>
</form>
`)
}

func HttpWS(w http.ResponseWriter, req *http.Request, message chan string) {
	msg := <-message
	io.WriteString(w, msg)
}
*/
