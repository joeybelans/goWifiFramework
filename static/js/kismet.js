var sort = 'network';
var order = 'asc';
var filterScope = [];
var filterNetwork = [];
var filterChannel = [];
var filterBSSID = [];
var inScope = ['2turntables', '2turntablesG'];

var kismet = new function() {
   var accessPoints = new Object();
   var networks = new Object();
   var clients = new Object();
   var display = [];
   var sort = 'PACKETS';
   var order = 'asc';
   var filter = new Object();

   // Create filter array
   filter['scope'] = '';
   filter['ssid'] = [];
   filter['channel'] = [];
   filter['bssid'] = [];

   this.process = function(type, obj) {
      eval('parse' + type + '(obj)')
   }

   this.setSort = function(col) {
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

      eval('sortBy' + sort + '()');
   }

   this.displayDIV = function(element) {
      var id = element.id;
      var msg = '';
      eval('msg = display' + id.replace('filter', 'DIV') + '()');
      document.getElementById('divFilter').innerHTML = msg;
      document.getElementById('divFilter').style.display = 'block';
   }

   this.setDIVScope = function(val) {
      filter.scope = val;
      document.getElementById('divFilter').style.display = 'none';
   }

   function parseINFO(obj) {
      document.getElementById("networkCount").innerHTML = obj.networkCount;
      document.getElementById("rogueCount").innerHTML = obj.rogueCount;
      document.getElementById("clientCount").innerHTML = obj.clientCount;
      document.getElementById("apCount").innerHTML = obj.apCount;
      document.getElementById("packetCount").innerHTML = obj.packets;
      document.getElementById("packetRate").innerHTML = obj.rate;
      document.getElementById("cryptedCount").innerHTML = obj.crypt;
      document.getElementById("droppedCount").innerHTML = obj.dropped;
      document.getElementById("filteredCount").innerHTML = obj.filtered;
      document.getElementById("mgmtCount").innerHTML = obj.llcpackets;
      document.getElementById("dataCount").innerHTML = obj.datapackets;
      displayAccessPoints();
   }

   function parseBSSID(obj) {
      if (accessPoints[obj.bssid] === undefined) {
	 accessPoints[obj.bssid] = new Object();
         accessPoints[obj.bssid].ssid = "";
         accessPoints[obj.bssid].channel = obj.channel;
         accessPoints[obj.bssid].firsttime = obj.firsttime;
         accessPoints[obj.bssid].manuf= obj.manuf;
         accessPoints[obj.bssid].beaconInfo = "";
	 accessPoints[obj.bssid].clientCount = 0;
      }

      accessPoints[obj.bssid].lasttime = obj.lasttime;
      accessPoints[obj.bssid].signal = obj.signal_dbm;
      accessPoints[obj.bssid].maxsignal = obj.maxsignal_dbm;
      accessPoints[obj.bssid].minsignal = obj.minsignal_dbm;
      accessPoints[obj.bssid].packets = obj.numPackets;

      eval('sortBy' + sort + '()');
      displayAccessPoints();
   }

   function parseSSID(obj) {
      if (accessPoints[obj.mac] === undefined) {
	 accessPoints[obj.mac] = new Object();
         accessPoints[obj.mac].ssid = "";
         accessPoints[obj.mac].channel = "";
         accessPoints[obj.mac].firsttime = obj.firsttime;
         accessPoints[obj.mac].signal = "";
         accessPoints[obj.mac].minsignal = "";
         accessPoints[obj.mac].maxsignal = "";
         accessPoints[obj.mac].packets = 0;
         accessPoints[obj.mac].clientCount = 0;
      }

      if (accessPoints[obj.mac].ssid == "") {
         accessPoints[obj.mac].ssid = obj.ssid;
      }
      accessPoints[obj.mac].lasttime = obj.lasttime;
      accessPoints[obj.mac].beaconInfo = obj.beaconinfo;

      if (networks[obj.ssid] === undefined) {
	 networks[obj.ssid] = new Object;
         networks[obj.ssid].firsttime = obj.firsttime;
      }

      networks[obj.ssid].cloaked = obj.cloaked;
      networks[obj.ssid].lasttime = obj.lasttime;
   }

   function parseCLIENT(obj) {
      if (clients[obj.mac] === undefined) {
	 clients[obj.mac] = new Object();
         clients[obj.mac].firsttime = obj.firsttime;
         clients[obj.mac].manuf = obj.manuf;
	 clients[obj.mac].bssid = "";
      }

      if ((clients[obj.mac].bssid == "") && (obj.bssid != "")) {
         clients[obj.mac].bssid = obj.bssid;
	 accessPoints[obj.bssid].clientCount++;
      } else if ((clients[obj.mac].bssid != "") && (obj.bssid == "")) {
         clients[obj.mac].bssid = "";
	 accessPoints[obj.bssid].clientCount++;
      }

      clients[obj.mac].lasttime = obj.lasttime;
      clients[obj.mac].signal_dbm = obj.signal_dbm;
      clients[obj.mac].maxsignal_dbm = obj.maxsignal_dbm;
      clients[obj.mac].minsignal_dbm = obj.minsignal_dbm;
      clients[obj.mac].packets = obj.cryptpackets + obj.datapackets + obj.llcpackets;
   }

   function displayAccessPoints() {
      /*
      var msg = "<table onClick='getAccessPoint(event)' class='accessPoints'> \
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
*/
      var msg = "<table onClick='getAccessPoint(event)' class='accessPoints'> \
<thead>\n \
<tr><th><a class='apHeader' href='' onClick='kismet.setSort(\"SSID\"); return false;'>Network</a></th>\n \
<th><a class='apHeader' href='' onClick='kismet.setSort(\"BSSID\"); return false;'>BSSID</a></th>\n \
<th><a class='apHeader' href='' onClick='kismet.setSort(\"CHANNEL\"); return false;'>Channel</a></th>\n \
<th><a class='apHeader' href='' onClick='kismet.setSort(\"LASTTIME\"); return false;'>Last Seen</a></th>\n \
<th><a class='apHeader' href='' onClick='kismet.setSort(\"SIGNAL\"); return false;'>Power</a></th>\n \
<th><a class='apHeader' href='' onClick='kismet.setSort(\"MAXSIGNAL\"); return false;'>Max</a></th>\n \
<th><a class='apHeader' href='' onClick='kismet.setSort(\"CLIENTS\"); return false;'>Clients</a></th>\n \
<th><a class='apHeader' href='' onClick='kismet.setSort(\"PACKETS\"); return false;'>Packets</a></th></tr>\n \
</thead>\n \
<tbody>\n";

      var len = display.length;
      for (var i = 0; i < len; i++) {
         var ap = accessPoints[display[i]];
         if ((filterScope == '') || ((filterScope == 'in-scope') && (inScope.indexOf(ap.network) >= 0)) || ((filterScope == 'rogue') && (inScope.indexOf(ap.network) == -1))) {
	    /*
            msg += "<tr><td><div class='network'>" + ap.ssid+ "</div></td><td><div class='bssid'>" + display[i] + "</div></td><td><div class='channel'>" + ap.channel +
	        "</div></td><td><div class='last'>" + ap.lasttime + "</div></td><td><div class='power'>" + ap.signal + "</div></td><td><div class='max'>" + ap.maxsignal +
		"</div></td><td><div class='clients'></div></td><td><div class='packets'>" + ap.packets + "</div></td></tr>\n";
		//"</div></td><td><div class='clients'>" + ap.clients + "</div></td><td><div class='packets'>" + ap.packets + "</div></td></tr>\n";
		*/
            msg += "<tr><td>" + ap.ssid+ "</td><td>" + display[i] + "</td><td>" + ap.channel +
	        "</td><td>" + ap.lasttime + "</td><td>" + ap.signal + "</td><td>" + ap.maxsignal +
		"</td><td>" + ap.clientCount + "</td><td>" + ap.packets + "</td></tr>\n";
         }
      }

      msg = msg + "</tbody></table>\n";
      document.getElementById("apList").innerHTML = msg;
   }

   function sortBySSID() {
      var sorted = new Object();

      for (var bssid in accessPoints) {
         var ap = accessPoints[bssid];
         if (!(ap.ssid in sorted)) {
	    sorted[ap.ssid] = [];
	 }
         sorted[ap.ssid].push(bssid);
      }

      sortDisplay(sorted);
   }

   function sortByBSSID() {
      display = [];
      var aLen = accessPoints.length;
      var aKeys = Object.keys(accessPoints);

      aKeys.sort();
      if (order == 'des') {
         aKeys.reverse();
      }

      display = aKeys;
   }

   function sortByCHANNEL() {
      var sorted = new Object();

      for (var bssid in accessPoints) {
         var ap = accessPoints[bssid];
         if (!(ap.channel in sorted)) {
	    sorted[ap.channel] = [];
	 }
         sorted[ap.channel].push(bssid);
      }

      sortDisplay(sorted);
   }

   function sortByLASTTIME() {
      var sorted = new Object();

      for (var bssid in accessPoints) {
         var ap = accessPoints[bssid];
         if (!(ap.lasttime in sorted)) {
	    sorted[ap.lasttime] = [];
	 }
         sorted[ap.lasttime].push(bssid);
      }

      sortDisplay(sorted);
   }

   function sortBySIGNAL() {
      var sorted = new Object();

      for (var bssid in accessPoints) {
         var ap = accessPoints[bssid];
         if (!(ap.signal in sorted)) {
	    sorted[ap.signal] = [];
	 }
         sorted[ap.signal].push(bssid);
      }

      sortDisplay(sorted);
   }

   function sortByMAXSIGNAL() {
      var sorted = new Object();

      for (var bssid in accessPoints) {
         var ap = accessPoints[bssid];
         if (!(ap.maxsignal in sorted)) {
	    sorted[ap.maxsignal] = [];
	 }
         sorted[ap.maxsignal].push(bssid);
      }

      sortDisplay(sorted);
   }

   function sortByCLIENTS() {
      var sorted = new Object();

      for (var bssid in accessPoints) {
         var ap = accessPoints[bssid];
         if (!(ap.clientCount in sorted)) {
	    sorted[ap.clientCount] = [];
	 }
         sorted[ap.clientCount].push(bssid);
      }

      sortDisplay(sorted);
   }

   function sortByPACKETS() {
      var sorted = new Object();

      for (var bssid in accessPoints) {
         var ap = accessPoints[bssid];
         if (!(ap.packets in sorted)) {
	    sorted[ap.packets] = [];
	 }
         sorted[ap.packets].push(bssid);
      }

      sortDisplay(sorted);
   }

   function sortDisplay(sorted) {
      display = [];

      var sKeys = Object.keys(sorted);
      var sLen = sKeys.length;
      sKeys.sort();
      if (order == 'des') {
	 sKeys.reverse();
      }

      for (var i = 0; i < sLen; i++) {
	 sorted[sKeys[i]].sort();
	 for (var j = 0; j < sorted[sKeys[i]].length; j++) {
	    display.push(sorted[sKeys[i]][j]);
	 }
      }
   }

   function displayDIVScope() {
      var msg = "<b>Select the networks to display:</b>\n";
      msg += "<p>\n";
      msg += "<input type='radio' name='fscope' id='fscope' onClick=\"kismet.setDIVScope('')\"";
      if (filter.scope == '') { msg += " checked"; }
      msg += ">Show All</input>\n";
      msg += "<input type='radio' name='fscope' id='fscope' onClick=\"kismet.setDIVScope('in-scope')\"";
      if (filter.scope == 'in-scope') { msg += " checked"; }
      msg += ">In-scope</input>\n";
      msg += "<input type='radio' name='fscope' id='fscope' onClick=\"kismet.setDIVScope('rogue')\"";
      if (filter.scope == 'rogue') { msg += " checked"; }
      msg += ">Rogues</input>\n";

      return (msg);
   }

   function displayDIVSSID() {
      var msg = "<b>Select the SSIDs to display:</b>\n";
      msg += "<p>\n";

      var ssids = [];
      for (var bssid in accessPoints) {
	 var ap = accessPoints[bssid];
         if (ssids.indexOf(ap.ssid) >= 0) {
	    ssids.push(ap.ssid)
	 }
      }

      ssids.sort();
      sLen = ssids.length;

      for (var i = 0; i < sLen; i++) {
         msg += "<input type='checkbox' name='fssid' value='" + ssids[i] + "'";
         if (filter['ssid'].indexOf(ssids[i]) >= 0) {
            msg += ' checked';
         }
         msg += ">" + ssids[i] + "</input><br>\n";
      }

      return (msg);
   }
}

function processKismetUpdate(type) {
   if ((type == 'BSSID') || (type == 'SSID')) {
      var accessPoints = getAccessPoints();
      displayAccessPoints(accessPoints);
   } else if (type == 'SOURCE') {
      updateInterfaces();
   } else if (type == 'CLIENT') {
      displayClients(clients);
   }
}

function getAccessPoints() {
   var nKeys = getNetworkKeys();
   var accessPoints = getBSSIDs(nKeys);

   var fn = window['sortByS' + sort];
   var indexes = fn(accessPoints);
   return (sortedAccessPoints(indexes, accessPoints));
}

function getNetworkKeys() {
   if (filterNetwork.length > 0) {
      return (filterNetwork);
   }

   return (Object.keys(networks));
}

function getBSSIDs(nKeys) {
   var nLen = nKeys.length;
   var fLen = filterBSSID.length;
   var cLen = filterChannel.length;

   var accessPoints = [];
   for (var i = 0; i < nLen; i++) {
      var bKeys = Object.keys(networks[nKeys[i]].bssids);
      var bLen = bKeys.length;

      for (var j = 0; j < bLen; j++) {
         var bssid = null;
         if (fLen > 0) {
	    for (var k = 0; k < fLen; k++) {
	       if (((filterBSSID[k].length == 17) && (filterBSSID[k].toLowerCase() == bKeys[j].toLowerCase())) ||
	           ((filterBSSID[k].length < 17) && (bKeys[j].substring(0, filterBSSID[k].length).toLowerCase() == filterBSSID[k].toLowerCase()))) {
	          bssid = networks[nKeys[i]].bssids[bKeys[j]];
		  break;
	       }
	    }
	 } else {
	    bssid = networks[nKeys[i]].bssids[bKeys[j]];
	 }
	 if (bssid != null) {
            if ((cLen == 0) || (filterChannel.indexOf(bssid.channel) >= 0)) {
               accessPoints.push({
                  network: nKeys[i],
		  bssid: bKeys[j],
		  channel: bssid.channel,
		  last: bssid.lasttime,
		  power: bssid.power,
		  max: bssid.max,
		  clients: bssid.clients,
		  packets: bssid.packets});
	    }
	 }
      }
   }
   return (accessPoints);
}

function sortBySnetwork(accessPoints) {
   var aLen = accessPoints.length;
   var indexes = new Object();

   for (var i = 0; i < aLen; i++) {
      var ap = accessPoints[i];
      if (!(ap.network in indexes)) {
	 indexes[ap.network] = new Object();
      }
      indexes[ap.network][ap.bssid] = new Object();
      indexes[ap.network][ap.bssid][ap.channel] = i;
   }

   return (indexes);
}

function sortBySbssid(accessPoints) {
   var aLen = accessPoints.length;
   var indexes = new Object();

   for (var i = 0; i < aLen; i++) {
      var ap = accessPoints[i];
      if (!(ap.bssid in indexes)) {
	 indexes[ap.bssid] = new Object();
      }
      indexes[ap.bssid][ap.network] = new Object();
      indexes[ap.bssid][ap.network][ap.channel] = i;
   }

   return (indexes);
}

function sortByOther(accessPoints, first) {
   var aLen = accessPoints.length;
   var indexes = new Object();

   for (var i = 0; i < aLen; i++) {
      var ap = accessPoints[i];
      var index = ap[first];
      if (!(index in indexes)) {
	 indexes[index] = new Object();
      }
      if (!(ap.network in indexes[index])) {
         indexes[index][ap.network] = new Object();
      }
      indexes[index][ap.network][ap.bssid] = i;
   }

   return (indexes);
}

function sortBySchannel(accessPoints) {
   return (sortByOther(accessPoints, 'channel'));
}

function sortBySlastseen(accessPoints) {
   return (sortByOther(accessPoints, 'last'));
}

function sortBySpower(accessPoints) {
   return (sortByOther(accessPoints, 'power'));
}

function sortBySmax(accessPoints) {
   return (sortByOther(accessPoints, 'max'));
}

function sortBySclients(accessPoints) {
   return (sortByOther(accessPoints, 'clients'));
}

function sortBySpackets(accessPoints) {
   return (sortByOther(accessPoints, 'packets'));
}

function numSort(a,b) {
   return(a-b)
}

function sortedAccessPoints(indexes, accessPoints) {
   var sorted = [];

   var aKeys = Object.keys(indexes);
   var aLen = aKeys.length;

   if ((sort != 'network') && (sort != 'bssid') && (sort != 'lastseen')) {
      aKeys.sort(numSort);
   } else {
      aKeys.sort();
   }
   if (order == 'des') {
      aKeys.reverse();
   }

   for (var i = 0; i < aLen; i++) {
      var bKeys = Object.keys(indexes[aKeys[i]]);
      var bLen = bKeys.length;
      bKeys.sort();

      for (var j = 0; j < bLen; j++) {
         var cKeys = Object.keys(indexes[aKeys[i]][bKeys[j]]);
         var cLen = cKeys.length;
         cKeys.sort();

	 for (var k = 0; k < cLen; k++) {
	    sorted.push(accessPoints[indexes[aKeys[i]][bKeys[j]][cKeys[k]]]);
	 }
      }
   }
   return (sorted);
}

function updateFilters(accessPoints) {
   var aLen = accessPoints.length;
   var networks = [];
   var channels = [];

   for (var i = 0; i < aLen; i++) {
      var ap = accessPoints[i];
      if (networks.indexOf(ap.network) == -1) {
	 networks.push(ap.network);
      }
      if (channels.indexOf(ap.channel) == -1) {
	 channels.push(ap.channel);
      }
   }

   updateNetworkFilters(networks);
   updateChannelFilters(channels);
}

function updateNetworkFilters(networks) {
   if (document.getElementById('divNetwork').style.display != 'block') {
      updateFilter(networks, filterNetwork, 'networks', 'fnetwork', 'divNetwork');
   }
}

function updateChannelFilters(channels) {
   if (document.getElementById('divChannel').style.display != 'block') {
      updateFilter(channels, filterChannel, 'channels', 'fchannel', 'divChannel');
   }
}

function updateFilter(list, filter, type, cname, dname) {
   var lLen = list.length;

   if (type == 'channels') {
      list.sort(numSort);
   } else {
      list.sort();
   }

   var msg = "<b>Select the " + type + " to display:</b>\n<p>\n";

   for (var i = 0; i < lLen; i++) {
      msg += "<input type='checkbox' name='" + cname + "' value='" + list[i] + "'";
      if (filter.indexOf(list[i]) >= 0) {
	 msg += ' checked';
      }
      msg += ">" + list[i] + "</input><br>\n";
   }

   document.getElementById(dname).innerHTML = msg;
}

function displayAccessPoints(accessPoints) {
   var msg = "<table onClick='getAccessPoint(event)' class='data'> \
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

   var aLen = accessPoints.length;
   for (var i = 0; i < aLen; i++) {
      var ap = accessPoints[i];
      if ((filterScope == '') || ((filterScope == 'in-scope') && (inScope.indexOf(ap.network) >= 0)) || ((filterScope == 'rogue') && (inScope.indexOf(ap.network) == -1))) {
         msg += "<tr><td><div class='network'>" + ap.network + "</div></td><td><div class='bssid'>" + ap.bssid + "</div></td><td><div class='channel'>" + ap.channel +
	        "</div></td><td><div class='last'>" + ap.last + "</div></td><td><div class='power'>" + ap.power + "</div></td><td><div class='max'>" + ap.max +
		"</div></td><td><div class='clients'>" + ap.clients + "</div></td><td><div class='packets'>" + ap.packets + "</div></td></tr>\n";
      }
   }

   msg = msg + "</tbody></table>\n";
   document.getElementById("wsOutput").innerHTML = msg;
}

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

   var accessPoints = getAccessPoints();
   displayAccessPoints(accessPoints);
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

   if ((event.type == 'click') && ((name == 'divNetwork') || (name == 'divChannel'))) {
      var accessPoints = getAccessPoints();
      updateFilters(accessPoints);
   } else if ((!is_child_of(element, current_mouse_target)) && (element != current_mouse_target) &&
      (element.id == name) && (document.getElementById(name).style.display == 'block' && status == 'none')) {
      if ((name == 'divNetwork') || (name == 'divChannel') || (name == 'divScope') || (name == 'divBSSID')) {
         setFilter(name, status);
         updateFilterStatus();
         var accessPoints = getAccessPoints();
         displayAccessPoints(accessPoints);
      } else {
	 var nic = name.substring(3);
         //conn.send('nicCHANSOURCE:' + nic + ':' + document.getElementById('div' + nic + 'Clist').value);
         //conn.send('nicHOP:' + nic + ":" + document.getElementById('div' + nic + 'Crate').value);
	 conn.send(JSON.stringify({message: "nicCHANSOURCE", nic: nic, cList: document.getElementById('div' + nic + 'Clist').value}));
	 conn.send(JSON.stringify({message: "nicHOP", nic: nic, velocity: document.getElementById('div' + nic + 'Crate').value}));
      }
   }
   document.getElementById(name).style.display = status;
}

function setFilter(name, status) {
   if ((document.getElementById(name).style.display == 'block') && (status == 'none')) {
      if (name == 'divScope') {
	 var scope = document.getElementsByName('fscope');
	 var sLen = scope.length;
	 var i;
         for (i = 0; i < sLen; i++) {
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
	 var fbssid = document.getElementById('fbssid').value
	 if (fbssid.length == 0) {
	    filterBSSID = [];
	 } else {
	    filterBSSID = fbssid.split("\n");
	    var fLen = filterBSSID.length;
	    for (var i = 0; i < fLen; i++) {
	       if (filterBSSID[i].indexOf(":") == -1) {
		  var sLen = filterBSSID[i].length;
		  var tmp = '';
		  for (var j = 0; j < sLen; j++) {
		     tmp += filterBSSID[i][j];
		     if ((j+1)%2 == 0) {
			tmp += ':';
		     }
		  }
		  filterBSSID[i] = tmp;
	       }
	    }
	 }
      }
   }
}

function updateFilterStatus() {
   if (filterScope == 'in-scope') {
      document.getElementById('divScopeStatus').innerHTML = 'In-scope';
   } else if (filterScope == 'rogue') {
      document.getElementById('divScopeStatus').innerHTML = 'Rogues';
   } else {
      document.getElementById('divScopeStatus').innerHTML = 'Show All';
   }

   if (filterNetwork.length == 0) {
      document.getElementById('divNetworkStatus').innerHTML = 'Show All';
   } else {
      document.getElementById('divNetworkStatus').innerHTML = filterNetwork.join("<br>\n");
   }

   if (filterChannel.length == 0) {
      document.getElementById('divChannelStatus').innerHTML = 'Show All';
   } else {
      document.getElementById('divChannelStatus').innerHTML = filterChannel.join("<br>\n");
   }

   if (filterBSSID.length == 0) {
      document.getElementById('divBSSIDStatus').innerHTML = 'Show All';
   } else {
      document.getElementById('divBSSIDStatus').innerHTML = filterBSSID.join("<br>\n");
   }
}

function clearFilters() {
   var scope = document.getElementsByName('fscope');
   scope[0].checked = true;
   filterScope = '';

   document.getElementById('fbssid').value = '';
   filterNetwork = [];
   filterChannel = [];
   filterBSSID = [];

   updateFilterStatus();
   var accessPoints = getAccessPoints();
   displayAccessPoints(accessPoints);
}

function updateInterfaces() {
   var msg = '';

   for (nic in interfaces) {
      var name;
      if (nic == interfaces[nic].name) {
	 name = nic;
      } else {
	 name = interfaces[nic].name + ' (' + nic + ')';
      }

      if ((document.getElementById('div' + name) != null) && (document.getElementById('div' + name).style.display == 'block')) {
	 return;
      }

      msg += "<span class='stitle'>" + name + "</span><br>\n \
</span><br>\n \
<table class='stats'>\n \
<tr><td colspan='2' align='center'><input type='radio' name='kmethod_" + name + "' value='locked' onClick=\"channelLock('" + name + "');\"";

      if (interfaces[nic].hop == 0) {
	 msg += ' checked';
      }
      
      msg += ">Locked \n \
<input type='radio' name='kmethod_" + name + "' value='hopping' onChange=\"channelHop('" + name + "');\"";

      if ((interfaces[nic].hop == 3) || (interfaces[nic].hop == 1)) {
	 msg += ' checked';
      }
      
      msg += ">Hopping</td></tr>\n";

      if ((interfaces[nic].hop == 3) || (interfaces[nic].hop == 1)) {
         msg += "<tr><th align='right'>\n \
<div onclick=\"displayDIV(this, event, 'div" + name + "', 'block');\" onmouseout=\"displayDIV(this, event, 'div" + name + "', 'none');\">Channel List</div>\n \
<div id='div" + name + "' class='divFilter' onmouseover=\"displayDIV(this, event, 'div" + name +"', 'block');\" onmouseout=\"displayDIV(this, event, 'div" + name +"', 'none');\">\n \
<table><tr><th align='right'>Channel List:</th><td><input type='text' id='div" + name + "Clist' value='" + interfaces[nic].channels + "' size='50'></td></tr>\n \
<tr><th align='right'>Rate:</th><td><input type='text' id='div" + name + "Crate' value='" + interfaces[nic].velocity + "' size='3'>channels/sec</td></tr></table></div></th>\n \
<td>" + interfaces[nic].channels.split(",").length + " channels @ " + interfaces[nic].velocity + " channels/sec</td></tr>\n";
      }

      msg += "<tr><th align='righ'>Current Channel</th><td>" + interfaces[nic].current + "</td></tr>\n \
</table><p>\n";
   }
   document.getElementById('interfaces').innerHTML = msg;
}

function channelLock(name) {
   var channel = prompt('Enter channel', '');
   if (isNaN(channel)) {
      alert (channel + ' is not a number');
   } else {
      //conn.send('nicLOCK:' + name + ':' + channel);
      conn.send(JSON.stringify({message: "nicLOCK", nic: name, channel: channel}));
   }
}

function channelHop(name) {
   var vel;
   if (document.getElementById('div' + name + 'Crate') != null) {
      vel = document.getElementById('div' + name + 'Crate').value;
   } else {
      vel = 3;
   }
   //conn.send('nicHOP:' + name + ":" + vel);
   conn.send(JSON.stringify({message: "nicHOP", nic: name, velocity: vel}));
}

function getAccessPoint(e) {
   var t = e.target;
   var p;
   if (t.nodeName == 'DIV') {
      p = t.parentNode.parentNode;
   } else if (t.nodeName == 'TD') {
      p = t.parentNode;
   } else {
      return;
   }

   //conn.send('getAPDetails:' + p.childNodes[1].childNodes[0].innerHTML);
   conn.send(JSON.stringify({message: "getAPDetails", bssid: p.childNodes[1].childNodes[0].innerHTML}));
   //alert (p.childNodes[0].childNodes[0].innerHTML);
   //alert (p.childNodes[1].childNodes[0].innerHTML);
}

function apDetails(obj) {
   msg = "<span class='stitle'>" + obj.ssid + "</span>\n \
<blockquote>\n \
<table class='stats'>\n \
<tr><th>Cloaked</th><td>" + obj.cloaked + "</td></tr>\n \
<tr><th>First Seen</th><td>" + obj.firstseen + "</td></tr>\n \
<tr><th>Last Seen</th><td>" + obj.lastseen + "</td></tr>\n \
<tr><th>Max Rate</th><td>" + obj.maxrate + "</td></tr>\n \
<tr><th>Encryption</th><td>" + obj.crypt + "</td></tr>\n \
<tr><th># of APs</th><td>" + obj.aps + "</td></tr>\n \
</table></blockquote>\n \
<p>\n \
<span class='stitle'>" + obj.bssid + "</span>\n \
<blockquote>\n \
<table class='stats'>\n \
<tr><th>Type</th><td>" + obj.aptype + "</td></tr>\n \
<tr><th>Manufacturer</th><td>" + obj.manuf + "</td></tr>\n \
<tr><th>Channel</th><td>" + obj.channel + "</td></tr>\n \
<tr><th>First Seen</th><td>" + obj.apfirstseen + "</td></tr>\n \
<tr><th>Last Seen</th><td>" + obj.aplastseen + "</td></tr>\n \
<tr><th>Network</th><td>" + obj.ip + "</td></tr>\n \
<tr><th>Netmask</th><td>" + obj.netmask + "</td></tr>\n \
<tr><th>Gateway</th><td>" + obj.gateway + "</td></tr>\n \
<tr><th>Power</th><td>" + obj.power + "</td></tr>\n \
<tr><th>Min Power</th><td>" + obj.min + "</td></tr>\n \
<tr><th>Max Power</th><td>" + obj.max + "</td></tr>\n \
<tr><th># of Packets</th><td>" + obj.packets + "</td></tr>\n \
</table></blockquote>";

   document.getElementById('apDetails').innerHTML = msg;
   document.getElementById('apDetails').style.display = 'block';
}

function displayClients(clients) {
   var msg = "<table class='data'> \
<thead>\n \
<tr><th><div class='network'><a class='data' href='' onClick='setSort(\"network\"); return false;'>Client</a></div></th>\n \
<th><div class='last'><a class='data' href='' onClick='setSort(\"lastseen\"); return false;'>Last Seen</a></div></th>\n \
<th><div class='power'><a class='data' href='' onClick='setSort(\"power\"); return false;'>Power</a></div></th>\n \
<th><div class='max'><a class='data' href='' onClick='setSort(\"max\"); return false;'>Min</a></div></th>\n \
<th><div class='max'><a class='data' href='' onClick='setSort(\"max\"); return false;'>Max</a></div></th>\n \
<th><div class='packets'><a class='data' href='' onClick='setSort(\"packets\"); return false;'>Packets</a></div></th></tr>\n \
</thead>\n \
<tbody>\n";

   for (var mac in clients) {
      var client = clients[mac];
      msg += "<tr><td><div class='network'>" + mac + "</div></td><td><div class='last'>" + client.last + "</div></td><td><div class='power'>" + client.power + "</div></td><td><div class='max'>" +
                client.min + "</div></td><td><div class='max'>" + client.max + "</div></td><td><div class='packets'>" + client.packets + "</div></td></tr>\n";
   }

   msg = msg + "</tbody></table>\n";
   document.getElementById("clientOutput").innerHTML = msg;
}
