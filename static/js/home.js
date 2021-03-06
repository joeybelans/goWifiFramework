// Request the SSID information of the provided ssid
function getSSID(ssid) {
   conn.send(JSON.stringify({message: "statsSSID", ssid: ssid}));
   alert(ssid);
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

//<tr><th style='text-align: left' colspan='2' bgcolor='#000000'><font color='#ffffff'>" + obj.ssid + "</font></th></tr>\
function ssidINFO(obj) {
   divtxt = "<form><table class='data'>\
<tr><th class='header' colspan='2'>" + obj.ssid + "</th></tr>\
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

//<tr><th style='text-align: left' colspan='2' bgcolor='#000000'><font color='#ffffff'>" + obj.nic+ "</font></th></tr>\
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
<tr><th class='header' colspan='2'>" + obj.nic+ "</th></tr>\
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
