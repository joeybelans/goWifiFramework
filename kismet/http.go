package kismet

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

func wsSend(obj map[string]interface{}) {
	if wsconn != nil {
		switch obj["message"] {
		case "kismetParseInfo":
			if curPage == "/" {
				wsconn.WriteJSON(obj)
			}
		case "kismetParseSSID":
			if curPage == "/discover" {
				wsconn.WriteJSON(obj)
			}
		case "kismetParseBSSID":
			if curPage == "/discover" {
				wsconn.WriteJSON(obj)
			}
		case "kismetParseSOURCE":
			if curPage == "/discover" {
				wsconn.WriteJSON(obj)
			}
		case "kismetParseCLIENT":
			if curPage == "/discover" {
				wsconn.WriteJSON(obj)
			}
		default:
			wsconn.WriteJSON(obj)
		}
	}
}

func processWSCommand(data map[string]interface{}) {
	switch data["message"].(string) {
	case "kismetDISCONNECT":
		kismet.conn.Close()

	case "kismetCONNECT":
		Run(kismet.host, kismet.port, kismet.db, kismet.debug, ssids)

	case "statsNIC":
		nic := data["nic"].(string)
		stats := map[string]interface{}{
			"message":  "nicINFO",
			"nic":      nic,
			"active":   0,
			"physical": "",
			"alias":    "",
			"channel":  0,
			"hop":      0,
			"velocity": 0,
			"cList":    "",
		}
		iface, _ := net.InterfaceByName(nic)
		fmt.Println(iface)
		for _, ele := range kismet.interfaces {
			fmt.Println(stats["active"], ele.hwaddr, iface.HardwareAddr.String())
			if stats["active"] == 0 && ele.hwaddr == iface.HardwareAddr.String() {
				fmt.Println("found")
				stats["active"] = "2"
				stats["physical"] = ele.pname
				stats["alias"] = ele.lname
				stats["channel"] = ele.channel
				stats["hop"] = ele.hop
				stats["velocity"] = ele.velocity
				stats["cList"] = ele.channellist
			}
			if ele.pname == nic {
				stats["active"] = "1"
				break
			}
		}
		fmt.Println(stats)
		wsSend(stats)

	case "statsSSID":
		ssid := data["ssid"].(string)
		wlan := networks[ssid]

		channels := make(map[int]string)
		clientCount := 0
		min := 0
		max := -99
		for bssid := range wlan.bssids {
			ap := accessPoints[bssid]
			if ap.minSignalDBM < min {
				min = ap.minSignalDBM
			}
			if ap.maxSignalDBM > max {
				max = ap.maxSignalDBM
			}
			channels[ap.channel] = string(ap.channel)

			for clientMAC := range clients {
				if clientMAC == bssid {
					clientCount++
				}
			}
		}

		cloaked := "No"
		if wlan.cloaked == 1 {
			cloaked = "Yes"
		}

		keys := []string{}
		for k := range channels {
			keys = append(keys, strconv.Itoa(k))
		}

		i, _ := strconv.ParseInt(wlan.firsttime, 10, 64)
		firstTime := time.Unix(i, 0)

		i, _ = strconv.ParseInt(wlan.lasttime, 10, 64)
		lastTime := time.Unix(i, 0)

		wsSend(map[string]interface{}{
			"message":   "ssidINFO",
			"ssid":      ssid,
			"cloaked":   cloaked,
			"firstseen": firstTime.String(),
			"lastseen":  lastTime.String(),
			"maxrate":   wlan.maxrate,
			"min":       min,
			"max":       max,
			"clients":   clientCount,
			"aps":       len(wlan.bssids),
			"crypt":     wlan.encryption,
			"keys":      strings.Join(keys, ","),
		})

	case "nicADDSOURCE":
		txt := data["nic"].(string)
		if data["name"].(string) != "" {
			txt = txt + ":name=" + data["name"].(string)
		}
		Send("ADDSOURCE", txt)

	case "nicDELSOURCE":
		for uid, ele := range kismet.interfaces {
			if ele.pname == data["nic"].(string) {
				Send("DELSOURCE", uid)
				delete(kismet.interfaces, uid)
				break
			}
		}
	case "nicLOCK":
		for uid, ele := range kismet.interfaces {
			if ele.pname == data["nic"].(string) {
				Send("HOPSOURCE", uid+" LOCK "+data["channel"].(string))
				break
			}
		}
	case "nicHOP":
		for uid, ele := range kismet.interfaces {
			if ele.pname == data["nic"].(string) {
				Send("HOPSOURCE", uid+" HOP "+data["velocity"].(string))
				break
			}
		}
	case "nicCHANSOURCE":
		for uid, ele := range kismet.interfaces {
			if ele.pname == data["nic"].(string) {
				Send("CHANSOURCE", uid+" "+data["cList"].(string))
				break
			}
		}
	case "getAPDetails":
		ap := accessPoints[data["bssid"].(string)]
		wlan := networks[ap.ssid]

		i, _ := strconv.ParseInt(ap.firsttime, 10, 64)
		apFirstTime := time.Unix(i, 0)

		i, _ = strconv.ParseInt(ap.lasttime, 10, 64)
		apLastTime := time.Unix(i, 0)

		i, _ = strconv.ParseInt(wlan.firsttime, 10, 64)
		wlanFirstTime := time.Unix(i, 0)

		i, _ = strconv.ParseInt(wlan.lasttime, 10, 64)
		wlanLastTime := time.Unix(i, 0)

		wsSend(map[string]interface{}{
			"message":     "apDetails",
			"ssid":        ap.ssid,
			"cloaked":     wlan.cloaked,
			"firstseen":   wlanFirstTime.String(),
			"lastseen":    wlanLastTime.String(),
			"maxrate":     wlan.maxrate,
			"crypt":       wlan.encryption,
			"aps":         len(wlan.bssids),
			"bssid":       data["bssid"].(string),
			"aptype":      ap.apType,
			"manuf":       ap.manuf,
			"channel":     ap.channel,
			"apfirstseen": apFirstTime.String(),
			"aplastseen":  apLastTime.String(),
			"atype":       ap.atype,
			"ip":          ap.rangeip,
			"netmask":     ap.netmaskip,
			"gateway":     ap.gatewayip,
			"power":       ap.signalDBM,
			"min":         ap.minSignalDBM,
			"max":         ap.maxSignalDBM,
			"packets":     ap.numPackets,
		})
	}
}
