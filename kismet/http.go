package kismet

import (
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

func wsSend(msgType string, msg string) {
	if wsconn != nil {
		switch msgType {
		case "kismetParseInfo":
			if curPage == "/" {
				wsconn.WriteMessage(websocket.TextMessage, []byte(msgType+":"+msg))
			}
		case "kismetParseSSID":
			if curPage == "/discover" {
				wsconn.WriteMessage(websocket.TextMessage, []byte(msgType+":"+msg))
			}
		case "kismetParseBSSID":
			if curPage == "/discover" {
				wsconn.WriteMessage(websocket.TextMessage, []byte(msgType+":"+msg))
			}
		case "kismetParseSOURCE":
			if curPage == "/discover" {
				wsconn.WriteMessage(websocket.TextMessage, []byte(msgType+":"+msg))
			}
		case "kismetParseCLIENT":
			if curPage == "/discover" {
				wsconn.WriteMessage(websocket.TextMessage, []byte(msgType+":"+msg))
			}
		default:
			wsconn.WriteMessage(websocket.TextMessage, []byte(msgType+":"+msg))
		}
	}
}

func processWSCommand(msg string) {
	s := strings.SplitN(msg, ":", 2)

	switch s[0] {
	case "kismetDISCONNECT":
		kismet.conn.Close()

	case "kismetCONNECT":
		Run(kismet.host, kismet.port, kismet.db, kismet.debug, ssids)

	case "statsNIC":
		iface, _ := net.InterfaceByName(s[1])
		active := "0"
		stats := ""
		for _, ele := range kismet.interfaces {
			if stats == "" && ele.hwaddr == iface.HardwareAddr.String() {
				active = "2"
				stats = ele.pname + ";" + ele.lname + ";" + ele.channel + ";" + ele.hop + ";" + ele.velocity + ";" + ele.channellist
			}
			if ele.pname == s[1] {
				active = "1"
				break
			}
		}
		stats = s[1] + ";" + active + ";" + stats
		wsSend("nicINFO", stats)

	case "statsSSID":
		wlan := networks[s[1]]

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

		wsSend("ssidINFO", s[1]+";"+cloaked+";"+firstTime.String()+";"+lastTime.String()+";"+strconv.Itoa(wlan.maxrate)+";"+strconv.Itoa(min)+";"+
			strconv.Itoa(max)+";"+strconv.Itoa(clientCount)+";"+strconv.Itoa(len(wlan.bssids))+";"+wlan.encryption+";"+strings.Join(keys, ","))

	case "nicADDSOURCE":
		t := strings.Split(s[1], ":")
		txt := t[0]
		if t[1] != "" {
			txt = txt + ":name=" + t[1]
		}
		Send("ADDSOURCE", txt)

	case "nicDELSOURCE":
		for uid, ele := range kismet.interfaces {
			if ele.pname == s[1] {
				Send("DELSOURCE", uid)
				delete(kismet.interfaces, uid)
				break
			}
		}
	case "nicLOCK":
		t := strings.Split(s[1], ":")
		for uid, ele := range kismet.interfaces {
			if ele.pname == t[0] {
				Send("HOPSOURCE", uid+" LOCK "+t[1])
				break
			}
		}
	case "nicHOP":
		t := strings.Split(s[1], ":")
		for uid, ele := range kismet.interfaces {
			if ele.pname == t[0] {
				Send("HOPSOURCE", uid+" HOP "+t[1])
				break
			}
		}
	case "nicCHANSOURCE":
		t := strings.Split(s[1], ":")
		for uid, ele := range kismet.interfaces {
			if ele.pname == t[0] {
				Send("CHANSOURCE", uid+" "+t[1])
				break
			}
		}
	case "getAPDetails":
		ap := accessPoints[s[1]]
		wlan := networks[ap.ssid]

		i, _ := strconv.ParseInt(ap.firsttime, 10, 64)
		apFirstTime := time.Unix(i, 0)

		i, _ = strconv.ParseInt(ap.lasttime, 10, 64)
		apLastTime := time.Unix(i, 0)

		i, _ = strconv.ParseInt(wlan.firsttime, 10, 64)
		wlanFirstTime := time.Unix(i, 0)

		i, _ = strconv.ParseInt(wlan.lasttime, 10, 64)
		wlanLastTime := time.Unix(i, 0)

		wsSend("apDetails", ap.ssid+";"+strconv.Itoa(wlan.cloaked)+";"+wlanFirstTime.String()+";"+wlanLastTime.String()+";"+strconv.Itoa(wlan.maxrate)+";"+wlan.encryption+";"+
			strconv.Itoa(len(wlan.bssids))+";"+s[1]+";"+ap.apType+";"+ap.manuf+";"+strconv.Itoa(ap.channel)+";"+apFirstTime.String()+";"+apLastTime.String()+";"+
			strconv.Itoa(ap.atype)+";"+ap.rangeip+";"+ap.netmaskip+";"+ap.gatewayip+";"+strconv.Itoa(ap.signalDBM)+";"+strconv.Itoa(ap.minSignalDBM)+";"+
			strconv.Itoa(ap.maxSignalDBM)+";"+strconv.Itoa(ap.numPackets))
	}
}
