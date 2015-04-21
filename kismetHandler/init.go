// Manages server connection and data related to kismet server
package kismetHandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"

	"github.com/joeybelans/gokismet/header"
	"github.com/joeybelans/golibkismet"
)

// Initialize the package
func init() {
	//db = nil
	debug = false
	networks = map[string]network{}
	accessPoints = map[string]accessPoint{}
	clients = map[string]client{}
	ssids = []string{}

	// Create page
	httpTemplate = template.New("/kismet")
	httpTemplate, _ = httpTemplate.Parse(httpSource())
}

// Initialize the package
func Init(host string, port int, ldebug bool, lssids []string) {
	// Set package variables
	//db = ldb
	debug = ldebug
	ssids = lssids

	header.AddPage("/kismet", "Kismet")
	http.HandleFunc("/kismet", func(w http.ResponseWriter, r *http.Request) {
		HttpKismet(w, r, ssids)
	})

	// Establish connection to kismet server
	golibkismet.Connect(host, port, debug)

	// Listen for responses
	go listen()

	// Setup the client
	setupClient()
}

// Setup the client
func setupClient() {
	golibkismet.Enable("INFO", []string{"packets", "rate", "crypt", "dropped", "filtered", "llcpackets", "datapackets"})
	golibkismet.Enable("STATUS", []string{"text", "flags"})
	golibkismet.Enable("ERROR", []string{"cmdid", "text"})
	golibkismet.Enable("ACK", []string{"cmdid", "text"})
	golibkismet.Enable("TERMINATE", []string{"text"})
	golibkismet.Enable("TIME", []string{"timesec"})
	golibkismet.Enable("STATUS", []string{"text", "flags"})
	golibkismet.Enable("SOURCE", []string{"interface", "type", "username", "channel", "uuid", "packets", "hop", "velocity", "dwell", "hop_time_sec", "hop_time_usec", "channellist", "error", "warning"})
	golibkismet.Enable("ALERT", []string{"sec", "usec", "header", "bssid", "source", "dest", "other", "channel", "text"})
	golibkismet.Enable("BSSID", []string{"bssid", "type", "llcpackets", "datapackets", "cryptpackets", "manuf", "channel", "firsttime", "lasttime", "atype", "rangeip", "netmaskip", "gatewayip", "signal_dbm", "minsignal_dbm", "maxsignal_dbm"})
	golibkismet.Enable("SSID", []string{"mac", "checksum", "type", "ssid", "beaconinfo", "cryptset", "cloaked", "firsttime", "lasttime", "maxrate", "beaconrate", "packets", "beacons", "dot11d"})
	golibkismet.Enable("CLIENT", []string{"bssid", "mac", "type", "firsttime", "lasttime", "manuf", "llcpackets", "datapackets", "cryptpackets", "gpsfixed", "minlat", "minlon", "minalt", "maxlat", "maxlon", "maxalt", "agglat", "agglon", "aggalt", "signal_dbm", "noise_dbm", "minsignal_dbm", "minnoise_dbm", "maxsignal_dbm", "maxnoise_dbm", "signal_rssi", "noise_rssi", "minsignal_rssi", "minnoise_rssi", "maxsignal_rssi", "maxnoise_rssi", "bestlat", "bestlon", "bestalt", "atype", "ip", "gatewayip", "datasize", "maxseenrate", "encodingset", "carrierset", "decrypted", "channel", "fragments", "retries", "newpackets", "freqmhz", "cdpdevice", "cdpport", "dhcphost", "dhcpvendor", "datacryptset"})
	golibkismet.Enable("BSSIDSRC", []string{"bssid", "uuid", "lasttime", "numpackets"})
	golibkismet.Enable("CLISRC", []string{"bssid", "mac", "uuid", "lasttime", "numpackets", "signal_dbm", "minsignal_dbm", "maxsignal_dbm"})
	golibkismet.Enable("NETTAG", []string{"bssid", "tag", "value"})
	golibkismet.Enable("CLITAG", []string{"bssid", "mac", "tag", "value"})
	//golibkismet.AddSource("wlan0", "test")
}

// Listen for responses from kismet client interface
func listen() {
	for {
		var (
			msg []byte
			obj map[string]string
		)

		msg = <-golibkismet.Responses
		json.Unmarshal(msg, &obj)

		c, found := processors[obj["message"]]
		if found {
			c.(func(map[string]string))(obj)
		} else if debug {
			fmt.Println("UNKNOWN RESPONSE: " + obj["message"])
		}
	}
}
