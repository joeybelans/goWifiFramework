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
<script src="global.js"></script>
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
