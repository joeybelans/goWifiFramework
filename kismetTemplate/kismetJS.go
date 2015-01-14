package kismetTemplate

import "net/http"

// KismetJS
func HttpKismetJS(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")
	w.Write([]byte(`
var networks = new Object();
var interfaces = new Object();

function kismetParseBSSID(msg) {
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

   processKismetUpdate('BSSID');
}

function kismetParseSSID(msg) {
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

   processKismetUpdate('SSID');
}

function kismetParseSOURCE(msg) {
   var fields = msg.split(";");

   if (!(fields[0] in interfaces)) {
      interfaces[fields[0]] = {
	 name: fields[1],
	 current: fields[2],
	 hop: fields[3],
	 velociy: fields[4],
	 channels: fields[5]
      };
   } else {
      interfaces[fields[0]].current = fields[2];
      interfaces[fields[0]].hop = fields[3];
      interfaces[fields[0]].velociy = fields[4];
      interfaces[fields[0]].channels = fields[5];
   }

   processKismetUpdate('SOURCE');
}
`))
}
