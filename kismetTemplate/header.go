package kismetTemplate

// Header
type header struct {
	Title     string
	Path      string
	Connected string
}

func tmplHeader(pages map[int]page) string {
	txt := `
<html>
<head>
<title>GoKismet - {{.Title}}</title>
<link rel="stylesheet" type="text/css" href="global.css">
{{if eq .Title "Discover"}}
<script src="/kismet.js"></script>
<script src="/discover.js"></script>
{{end}}
<script>
function kismetOnOff() {
   if (document.getElementById("cmn-toggle-7").checked == false) {
      conn.send(JSON.stringify({message: "kismetDISCONNECT"}));
   } else {
      conn.send(JSON.stringify({message: "kismetCONNECT"}));
   }
}

function kismetAddSource(nic) {
   conn.send(JSON.stringify({message: "nicADDSOURCE", nic: nic, name: document.getElementsByName("alias")[0].value}));
   setTimeout('conn.send(JSON.stringify({message: "statsNIC", nic: "' + nic + '"}))', 3000);
}

function kismetDelSource(nic) {
   conn.send(JSON.stringify({message: "nicDELSOURCE", nic: nic}));
   setTimeout('conn.send(JSON.stringify({message: "statsNIC", nic: "' + nic + '"}))', 3000);
}


function kismetParseInfo(obj) {
   document.getElementById("statsNcount").innerHTML = obj.networks;
   document.getElementById("statsRcount").innerHTML = obj.rogues;
   document.getElementById("statsCcount").innerHTML = obj.clients;
   document.getElementById("statsPcount").innerHTML = obj.total;
   document.getElementById("statsPrate").innerHTML = obj.rate;
   document.getElementById("statsCrypted").innerHTML = obj.crypt;
   document.getElementById("statsDropped").innerHTML = obj.dropped;
   document.getElementById("statsFiltered").innerHTML = obj.filtered;
   document.getElementById("statsManagement").innerHTML = obj.mgmt;
   document.getElementById("statsData").innerHTML = obj.data;
}

function kismetParseTerminate(obj) {
   document.getElementById("cmn-toggle-7").checked = false;
}

function ssidINFO(obj) {
   divtxt = "<form><table class='data'>\
<tr><th style='text-align: left' colspan='2' bgcolor='#000000'><font color='#ffffff'>" + obj.ssid + "</font></th></tr>\
<tr><th>Cloaked</th><td>" + obj.cloaked + "</td></tr>\
<tr><th>Channels</th><td>" + obj.keys + "</td></tr>\
<tr><th>First</th><td>" + obj.firstseen + "</td></tr>\
<tr><th>Last</th><td>" + obj.lastseen + "</td></tr>\
<tr><th>Max Rate</th><td>" + obj.maxrate + "</td></tr>\
<tr><th>Min DBM</th><td>" + obj.min + "</td></tr>\
<tr><th>Max DBM</th><td>" + obj.max + "</td></tr>\
<tr><th>Client Count</th><td>" + obj.clients + "</td></tr>\
<tr><th>BSSID Count</th><td>" + obj.aps + "</td></tr>\
<tr><th>Encryption</th><td>" + obj.crypt + "</td></tr>\
</table></form>";

   document.getElementById("wsOutput").innerHTML = divtxt;
}

function nicINFO(obj) {
   var active = "inactive";

   if (obj.active == 1) {
      active = "active";
   } else if (obj.active == 2) {
      active = obj.alias;
   }

   var hopping = "Yes";
   if (obj.hop == 0) {
      hopping = "No";
   }

   divtxt = "<form><table class='data'>\
<tr><th style='text-align: left' colspan='2' bgcolor='#000000'><font color='#ffffff'>" + obj.nic+ "</font></th></tr>\
<tr><td colspan='2'>Active<input type='radio' name='source' value='active'";

   if (active != 'inactive') {
      divtxt = divtxt + " checked";
   }  else {
      divtxt = divtxt + " onChange='kismetAddSource(\"" + obj.nic+ "\")'";
   }

   if ((active != 'active') && (active != 'inactive')) {
      divtxt = divtxt + " disabled";
   }

   divtxt = divtxt + "> Inactive<input type='radio' name='source' value='inactive'";
   
   if (active == 'inactive') {
      divtxt = divtxt + " checked";
   }  else {
      divtxt = divtxt + " onChange='kismetDelSource(\"" + obj.nic+ "\")'";
   } 

   if ((active != 'active') && (active != 'inactive')) {
      divtxt = divtxt + " disabled";
   } 

   divtxt = divtxt + ">"

   if ((active != 'active') && (active != 'inactive')) {
      divtxt = divtxt + " <font size='-2' color='#ff0000'>Via " + active + "</font>";
   }

   divtxt = divtxt + "</td></tr>";

   if ((active == 'inactive') || ((obj.physical != '') && (obj.nic != obj.physical))) {
      divtxt = divtxt + "<tr><th>Alias</th><td>";
      
      if (active == 'inactive') {
	 divtxt = divtxt + "<input type='text' name='alias' value=''></input>";
      } else {
	 divtxt = divtxt + obj.physical;
      }
      
      divtxt = divtxt + "</td></tr>";
   }

   if (obj.channel != 0) {
      divtxt = divtxt + "<tr><th>WiFi Channel</th><td>" + obj.channel + "</td></tr>\
<tr><th>Channel Hopping</th><td>" + hopping + "</td></tr>\
<tr><th>Rate</th><td>" + obj.velocity + " channels/sec</td></tr>\
<tr><th>Channel List</th><td>" + obj.cList + "</td></tr>"
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
   var json = JSON.parse(evt.data);
   window[json.message](json);
}
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
