package kismetTemplate

import (
	"net/http"
	"strings"
)

// DiscoverJS
func HttpDiscoverJS(w http.ResponseWriter, req *http.Request, ssids []string) {
	w.Header().Set("Content-Type", "application/javascript")
	w.Write([]byte(`
var sort = 'network';
var order = 'asc';
var filterScope = [];
var filterNetwork = [];
var filterChannel = [];
var filterBSSID = [];
var inScope = ['` + strings.Join(ssids, "', '") + `'];

function processKismetUpdate(type) {
   if ((type == 'BSSID') || (type == 'SSID')) {
      var accessPoints = getAccessPoints();
      displayAccessPoints(accessPoints);
   } else if (type == 'SOURCE') {
      updateInterfaces();
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
      setFilter(name, status);
      updateFilterStatus();
      var accessPoints = getAccessPoints();
      displayAccessPoints(accessPoints);
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

      msg += "<span class='stitle'>" + name + "</span><br>\n \
<table class='stats'>\n \
<tr><td colspan='2' align='center'><input type='radio' name='kmethod_" + name + "' value='locked'>Locked \n \
<input type='radio' name='kmethod_" + name + "' value='hopping' checked>Hopping</td></tr>\n \
<tr><th align='right'>Channel List</th><td>" + interfaces[nic].channels.split(",").length + " channels @ XXX channels/sec</td></tr>\n \
<tr><th align='righ'>Current Channel</th><td>" + interfaces[nic].current + "</td></tr>\n \
</table><p>\n";
   }
   document.getElementById('interfaces').innerHTML = msg;
}

/*
function updateInterfaceDIV() {
<tr><td colspan='2'><input type='radio' name='kmethod' value='locked'>Locked <input type='radio' name='kmethod' value='hopping' checked>Hopping</td></tr>
<tr><th align='right'><div onclick="displayDIV(this, event, 'divKChannel', 'block');" onmouseout="displayDIV(this, event, 'divKChannel', 'none');">Channel List</div>
<div id="divKChannel" class="divFilter" onmouseover="displayDIV(this, event, 'divKChannel', 'block');" onmouseout="displayDIV(this, event, 'divKChannel', 'none');">
<table><tr><th align='right'>Channel List:</th><td><input type='text' name='clist' value='1,11,153' size='50'></td></tr>
<tr><th align='right'>Rate:</th><td><input type='text' name='crate' value='3' size='3'>channels/sec</td></tr>
<tr><td colspan='2' align='center'><input type='submit' value='Update Kismet'></td></tr></table></div></th>
<td><div onclick="displayDIV(this, event, 'divKChannel', 'block');" onmouseout="displayDIV(this, event, 'divKChannel', 'none');">10 channels @ 3 channels/sec</div></td></tr>
<tr><th align='right'>Current Channel</th><td><div id='interfaces'></div></td></tr>
</table>
*/
`))
}
