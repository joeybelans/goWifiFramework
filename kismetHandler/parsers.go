// Manages server connection and data related to kismet server
package kismetHandler

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

func processERROR(obj map[string]string) {
	fmt.Println(obj["message"] + ": " + obj["text"])
}

func processSTATUS(obj map[string]string) {
	fmt.Println(obj["title"] + ": " + obj["text"])
}

func processSOURCE(obj map[string]string) {
	if obj["warning"] != "" {
		fmt.Println("NIC-WARNING: (" + obj["error"] + ") " + obj["warning"])
	}
	//fmt.Println(obj["message"] + ": " + obj["interface"] + ", " + obj["channel"])
}

func processINFO(obj map[string]string) {
	// Update the interface
	if okToSend("INFO") {
		// Update the object with the additional counts
		obj["networkCount"] = fmt.Sprintf("%d", len(networks))
		obj["rogueCount"] = fmt.Sprintf("%d", int(math.Max(0, float64(len(networks)-len(ssids)))))
		obj["apCount"] = fmt.Sprintf("%d", len(accessPoints))
		obj["clientCount"] = fmt.Sprintf("%d", len(clients))

		// Send the data to the socket
		sendString("INFO", obj)
	}
}

func processALERT(obj map[string]string) {
	fmt.Println(obj["message"] + ": " + obj["text"])
}

func processBSSIDSRC(obj map[string]string) {
	//fmt.Println(obj["message"] + ": " + obj["bssid"] + ", " + obj["uuid"] + ", " + obj["lastttime"] + ", " + obj["numpackets"])
}

func processBSSID(obj map[string]string) {
	// Define access point
	var ap accessPoint
	_, exists := accessPoints[obj["bssid"]]
	if !exists {
		ap = accessPoint{"", "", obj["manuf"], 0, 0, 0, 0, obj["rangeip"], obj["netmaskip"], obj["gatewayip"], 0, 0, 0, 0}
		ap.firsttime, _ = strconv.Atoi(obj["firsttime"])
	} else {
		ap = accessPoints[obj["bssid"]]
		ap.rangeip = obj["rangeip"]
		ap.netmaskip = obj["netmaskip"]
		ap.gatewayip = obj["gatewayip"]
	}

	// Update integer fields
	ap.channel, _ = strconv.Atoi(obj["channel"])
	ap.signalDBM, _ = strconv.Atoi(obj["signal_dbm"])
	ap.minSignalDBM, _ = strconv.Atoi(obj["minsignal_dbm"])
	ap.maxSignalDBM, _ = strconv.Atoi(obj["maxsignal_dbm"])

	llcpackets, _ := strconv.Atoi(obj["llcpackets"])
	datapackets, _ := strconv.Atoi(obj["datapackets"])
	ap.numPackets = llcpackets + datapackets

	lasttime, _ := strconv.Atoi(obj["lasttime"])
	if ap.lasttime < lasttime {
		ap.lasttime = lasttime
	}

	// Update accessPoints array
	accessPoints[obj["bssid"]] = ap

	// Send to socket
	if okToSend("BSSID") {
		obj["firsttime"] = time.Unix(int64(ap.firsttime), 0).Local().String()
		obj["lasttime"] = time.Unix(int64(ap.lasttime), 0).Local().String()
		obj["numPackets"] = strconv.Itoa(ap.numPackets)
		sendString("BSSID", obj)
	}
}

func processSSID(obj map[string]string) {
	if obj["ssid"] != "" {
		var ssid network
		_, exists := networks[obj["ssid"]]
		if !exists {
			ssid = network{false, false, 0, 0, 0, "", map[string]int{}}
			ssid.firsttime, _ = strconv.Atoi(obj["firsttime"])
		} else {
			ssid = networks[obj["ssid"]]
		}

		// Update fields
		ssid.cloaked = false
		if obj["cloaked"] != "0" {
			ssid.cloaked = true
		}

		lasttime, _ := strconv.Atoi(obj["lasttime"])
		if ssid.lasttime < lasttime {
			ssid.lasttime = lasttime
		}

		_, exists = ssid.bssids[obj["mac"]]
		if !exists || ssid.bssids[obj["mac"]] < lasttime {
			ssid.bssids[obj["mac"]] = lasttime
		}

		networks[obj["ssid"]] = ssid

		// Send to socket
		if okToSend("SSID") {
			obj["firsttime"] = time.Unix(int64(ssid.firsttime), 0).Local().String()
			obj["lasttime"] = time.Unix(int64(ssid.lasttime), 0).Local().String()
			sendString("SSID", obj)
		}

		// Update the access point, if exists
		_, exists = accessPoints[obj["mac"]]
		if exists {
			ap := accessPoints[obj["mac"]]
			ap.ssid = obj["ssid"]
			accessPoints[obj["mac"]] = ap
		}

	}
}

func processCLISRC(obj map[string]string) {
	//fmt.Println(obj["message"] + ": " + obj["bssid"] + ": " + obj["lastseen"])
}

func processNETTAG(obj map[string]string) {
	//fmt.Println(obj["message"] + ": " + obj["bssid"] + ": " + obj["tag"] + ": " + obj["value"])
}

func processCLITAG(obj map[string]string) {
	//fmt.Println(obj["message"] + ": " + obj["bssid"] + ": " + obj["tag"] + ": " + obj["value"])
}

func processCLIENT(obj map[string]string) {
	if obj["bssid"] != obj["mac"] {
		var clnt client
		_, exists := clients[obj["mac"]]
		if !exists {
			clnt = client{obj["bssid"], 0, 0, 0, 0, 0, 0}
		} else {
			clnt = clients[obj["mac"]]
			clnt.bssid = obj["bssid"]
		}

		if clnt.firsttime == 0 {
			clnt.firsttime, _ = strconv.Atoi(obj["firsttime"])
		}

		lasttime, _ := strconv.Atoi(obj["lasttime"])
		if clnt.lasttime < lasttime {
			clnt.lasttime = lasttime
		}

		clnt.signalDBM, _ = strconv.Atoi(obj["signal_dbm"])
		clnt.minSignalDBM, _ = strconv.Atoi(obj["minsignal_dbm"])
		clnt.maxSignalDBM, _ = strconv.Atoi(obj["maxsignal_dbm"])

		llcpackets, _ := strconv.Atoi(obj["llcpackets"])
		datapackets, _ := strconv.Atoi(obj["datapackets"])
		clnt.numPackets = llcpackets + datapackets

		clients[obj["mac"]] = clnt

		// Send to socket
		if okToSend("CLIENT") {
			sendString("CLIENT", obj)
		}
	}
}

func processTERMINATE(obj map[string]string) {
	fmt.Println(obj["message"] + ": " + obj["text"])
}
