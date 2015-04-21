var networks = new Object();
var interfaces = new Object();
var clients = new Object();

conn.onclose = function (event) {
   document.getElementById("cmn-toggle-7").checked = false;
}

function kismetOnOff() {
   if (document.getElementById("cmn-toggle-7").checked == false) {
      conn.send(JSON.stringify({message: "kismetDISCONNECT"}));
   } else {
      conn.send(JSON.stringify({message: "kismetCONNECT"}));
   }
}

function kismetParseTerminate(obj) {
   document.getElementById("cmn-toggle-7").checked = false;
}

function kismetParseBSSID(obj) {
   for (var ssid in networks) {
      for (var bssid in networks[ssid].bssids) {
         if (bssid == obj.bssid) {
	    networks[ssid].bssids[bssid].channel = obj.channel;
	    networks[ssid].bssids[bssid].lasttime = obj.lastseen;
	    networks[ssid].bssids[bssid].power = obj.power;
	    networks[ssid].bssids[bssid].clients = obj.power;
	    networks[ssid].bssids[bssid].max = obj.max;
	    networks[ssid].bssids[bssid].packets = obj.packets;
            break;
	 }
      }
   }

   processKismetUpdate('BSSID');
}

function kismetParseSSID(obj) {
   if (!(obj.ssid in networks)) {
      networks[obj.ssid] = {
         lastseen: obj.lastseen,
         bssids: {}
      };
   } else if (!(obj.bssid in networks[obj.ssid].bssids)) {
      networks[obj.ssid].bssids[obj.bssid] = {
         channel: '',
         lasttime: '',
         power: '',
	 max: '',
	 clients: '',
	 packets: ''
      };
   }

   processKismetUpdate('SSID');
}

function kismetParseSOURCE(obj) {
   if (!(obj.nic in interfaces)) {
      interfaces[obj.nic] = {
	 name: obj.name,
	 current: obj.channel,
	 hop: obj.hop,
	 velociy: obj.velocity,
	 channels: obj.chList
      };
   } else {
      interfaces[obj.nic].current = obj.channel;
      interfaces[obj.nic].hop = obj.hop;
      interfaces[obj.nic].velocity = obj.velocity;
      interfaces[obj.nic].channels = obj.chList;
   }

   processKismetUpdate('SOURCE');
}

function kismetParseCLIENT(obj) {
   if (!(obj.mac in clients)) {
      clients[obj.mac] = {
	 last: obj.lastseen,
	 power: obj.power,
	 min: obj.min,
	 max: obj.max,
	 packets: obj.packets
      };
   } else {
      clients[obj.mac].last = obj.lastseen;
      clients[obj.mac].power = obj.power;
      clients[obj.mac].min = obj.min;
      clients[obj.mac].max = obj.max;
      clients[obj.mac].packets = obj.packets;
   }

   processKismetUpdate('CLIENT');
}
