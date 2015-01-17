package kismetTemplate

import "net/http"

// GlobalCSS
type js struct {
	Referer string
}

func HttpJS(w http.ResponseWriter, req *http.Request) {
	//referer, _ := url.Parse(req.Referer())
	w.Header().Set("Content-Type", "application/javascript")
	//templates["/global.js"].Execute(w, js{referer.Path})
	w.Write([]byte(tmplGlobalJS()))
}

func tmplGlobalJS() string {
	return `
var networks = new Object();
var sort = 'network';
var order = 'asc';
var filterScope = '';
var filterNetwork = [];
var filterChannel = [];
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

	 if (((filterNetwork.length == 0) || (filterNetwork.indexOf(nKeys[i]) >= 0)) && ((filterChannel.length == 0) || (filterChannel.indexOf(b.channel) >= 0))) {
            msg = msg + "<tr><td><div class='network'>" + nKeys[i] + "</div></td><td><div class='bssid'>" + bKeys[j] + "</div></td><td><div class='channel'>" + b.channel +
	        "</div></td><td><div class='last'>" + b.lasttime + "</div></td><td><div class='power'>" + b.power + "</div></td><td><div class='max'>" + b.max +
		"</div></td><td><div class='clients'>" + b.clients + "</div></td><td><div class='packets'>" + b.packets + "</div></td></tr>\n";
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
<thead>\n \
<tr><th><div class='network'><a class='data' href='' onClick='setSort(\"network\"); return false;'>Network</a></div></th>\n \
<th><div class='bssid'><a class='data' href='' onClick='setSort(\"bssid\"); return false;'>BSSID</a></div></th>\n \
<th><div class='channel'><a class='data' href='' onClick='setSort(\"channel\"); return false;'>Channel</a></div></th>\n \
<th><div class='last'><a class='data' href='' onClick='setSort(\"lastseen\"); return false;'>Last Seen</a></div></th>\n \
<th><div class='power'><a class='data' href='' onClick='setSort(\"power\"); return false;'>Power</a></div></th>\n \
<th><div class='max'><a class='data' href='' onClick='setSort(\"max\"); return false;'>Max</a></div></th>\n \
<th><div class='clients'><a class='data' href='' onClick='setSort(\"clients\"); return false;'>Clients</a></div></th>\n \
<th><div class='packets'><a class='data' href='' onClick='setSort(\"packets\"); return false;'>Packets</a></div></th></tr>\n \
</thead>\n \
<tbody>\n";

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

   msg = msg + "</tbody></table>\n";
   document.getElementById("wsOutput").innerHTML = msg;
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

function is_child_of(parent, child) {
   if( child != null ) {			
      while( child.parentNode ) {
         if( (child = child.parentNode) == parent ) {
            return true;
         }
      }
   }
   return false;
}

function displayDIV(element, event, name, status) {
   var current_mouse_target = null;
   if ( event.toElement ) {				
      current_mouse_target = event.toElement;
   } else if( event.relatedTarget ) {				
      current_mouse_target = event.relatedTarget;
   }
   if( !is_child_of(element, current_mouse_target) && element != current_mouse_target && element.id == name && document.getElementById(name).style.display == 'block' && status == 'none') {
      setFilter(name, status);
   }
   document.getElementById(name).style.display = status;
}

function setFilter(name, status) {
   if ((document.getElementById(name).style.display == 'block') && (status == 'none')) {
      if (name == 'divScope') {
	 var scope = document.getElementsByName('fscope');
         for (i = 0; i < scope.length; i++) {
	    if (scope[i].checked) {
	       filterScope = scope[i].value;
	       break;
	    }
	 }
      } else if ((name == 'divNetwork') || (name == 'divChannel')) {
	 var checkboxes = null;
         if (name == 'divNetwork') {
            checkboxes = document.getElementsByName('fnetwork');
	    filterNetwork = [];
	 } else {
            checkboxes = document.getElementsByName('fchannel');
	    filterChannel = [];
	 }
         for (var i=0; i < checkboxes.length; i++) {
            if (checkboxes[i].checked) {
               if (name == 'divNetwork') {
                  filterNetwork.push(checkboxes[i].value);
               } else {
                  filterChannel.push(checkboxes[i].value);
               }
            }
         }
      } else if (name == 'divBSSID') {
	 filterBSSID = document.getElementById('fbssid').value.split("\n").join();
      }
   }
}
`
}
