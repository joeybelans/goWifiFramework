package kismetHandler

import (
	"database/sql"
	"text/template"
)

// Packet statistics
type packetStats struct {
	total    int
	rate     int
	crypt    int
	dropped  int
	filtered int
	mgmt     int
	data     int
}

// Access point object
type accessPoint struct {
	apType       string
	ssid         string
	manuf        string
	channel      int
	firsttime    int
	lasttime     int
	atype        int
	rangeip      string
	netmaskip    string
	gatewayip    string
	signalDBM    int
	minSignalDBM int
	maxSignalDBM int
	numPackets   int
}

// Wireless network object
type network struct {
	inscope    bool
	cloaked    bool
	firsttime  int
	lasttime   int
	maxrate    int
	encryption string
	bssids     map[string]int
}

// Wireless client object
type client struct {
	bssid        string
	firsttime    int
	lasttime     int
	signalDBM    int
	minSignalDBM int
	maxSignalDBM int
	numPackets   int
}

// Package variables
var (
	db           *sql.DB
	debug        bool
	networks     map[string]network
	accessPoints map[string]accessPoint
	clients      map[string]client
	ssids        []string
	httpTemplate *template.Template
)

/*
Capabilities
        ERROR:          cmdid, text
        STATUS:         text, flags
        SOURCE:         interface, type, username, channel, uuid, packets, hop, velocity, dwell, hop_time_sec, hop_time_usec, channellist, error, warning
        INFO:           networks, packets, rate, numsources, numerrorsources, crypt, dropped, filtered, llcpackets, datapackets
        ALERT:          sec, usec, header, bssid, source, dest, other, channel, text
        BSSIDSRC:       bssid, uuid, lasttime, numpackets
        BSSID:          bssid, type, llcpackets, datapackets, cryptpackets, manuf, channel, firsttime, lasttime, atype, rangeip, netmaskip, gatewayip, gpsfixed, minlat, minlon, minalt,
                        minspd, maxlat, maxlon, maxalt, maxspd, signal_dbm, noise_dbm, minsignal_dbm, minnoise_dbm, maxsignal_dbm, maxnoise_dbm, signal_rssi, noise_rssi, minsignal_rssi,
                        minnoise_rssi, maxsignal_rssi, maxnoise_rssi, bestlat, bestlon, bestalt, agglat, agglon, aggalt, aggpoints, datasize, turbocellnid, turbocellmode, turbocellsat,
                        carrierset, maxseenrate, encodingset, decrypted, dupeivpackets, bsstimestamp, cdpdevice, cdpport, fragments, retries, newpackets, freqmhz, datacryptset
        SSID:           mac, checksum, type, ssid, beaconinfo, cryptset, cloaked, firsttime, lasttime, maxrate, beaconrate, packets, beacons, dot11d
        CLISRC:         bssid, mac, uuid, lasttime, numpackets, signal_dbm, minsignal_dbm, maxsignal_dbm
        NETTAG:         bssid, tag, value
        CLITAG:         bssid, mac, tag, value
        CLIENT:         bssid, mac, type, firsttime, lasttime, manuf, llcpackets, datapackets, cryptpackets, gpsfixed, minlat, minlon, minalt, maxlat, maxlon, maxalt, agglat, agglon, aggalt,
                        signal_dbm, noise_dbm, minsignal_dbm, minnoise_dbm, maxsignal_dbm, maxnoise_dbm, signal_rssi, noise_rssi, minsignal_rssi, minnoise_rssi, maxsignal_rssi, maxnoise_rssi,
                        bestlat, bestlon, bestalt, atype, ip, gatewayip, datasize, maxseenrate, encodingset, carrierset, decrypted, channel, fragments, retries, newpackets, freqmhz, cdpdevice,
                        cdpport, dhcphost, dhcpvendor, datacryptset
        TERMINATE:      text
*/
// Kismet response processors
var processors = map[string]interface{}{
	"ERROR":     processERROR,
	"STATUS":    processSTATUS,
	"SOURCE":    processSOURCE,
	"INFO":      processINFO,
	"ALERT":     processALERT,
	"BSSIDSRC":  processBSSIDSRC,
	"BSSID":     processBSSID,
	"SSID":      processSSID,
	"CLISRC":    processCLISRC,
	"NETTAG":    processNETTAG,
	"CLITAG":    processCLITAG,
	"CLIENT":    processCLIENT,
	"TERMINATE": processTERMINATE,
}
