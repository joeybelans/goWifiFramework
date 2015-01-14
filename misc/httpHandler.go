package main

import (
	"fmt"
	"html/template"
	"io"
	"net"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/joeybelans/gokismet/kismet"
	"github.com/joeybelans/gokismet/kismetTemplate"
)

type iface struct {
	mac     string
	monitor []string
}

// Header function
func HttpHeader(w http.ResponseWriter, title string) {
	// Define headers
	headers := []string{"Home", "Discover", "Profile", "Networks", "AccessPoints", "Clients", "Reports", "Logs"}
	headerMap := map[string]string{
		"Home":         "/",
		"Discover":     "discover/",
		"Profile":      "profile/",
		"Networks":     "networks/",
		"AccessPoints": "accesspoints/",
		"Clients":      "clients/",
		"Reports":      "reports/",
		"Logs":         "logs/",
	}

	// Create HTML header
	io.WriteString(w, `
<html>
<head>
<title>GoKismet - `+title+`</title>
<style type='text/css'>
table.navtable th {
	border-width: 1px;
	padding: 8px;
	border-style: solid;
	border-color: #666666;
	background-color: #ADD8E6;
}
table.altrowstable {
	font-family: verdana,arial,sans-serif;
	font-size:11px;
	color:#333333;
	border-width: 1px;
	border-color: #a9c6c9;
	border-collapse: collapse;
}
table.altrowstable th {
	border-width: 1px;
	padding: 8px;
	border-style: solid;
	border-color: #a9c6c9;
}
table.altrowstable td {
	border-width: 1px;
	padding: 8px;
	border-style: solid;
	border-color: #a9c6c9;
}
.oddrowcolor{
	background-color:#d4e3e5;
}
.evenrowcolor{
	background-color: #ADD8E6;
}
.interfaces{
	width: 300px;
	height: 200px;
	border: 1px solid;
	padding-top: 10px;
	padding-right: 10px;
	padding-bottom: 10px;
	padding-left: 10px;
}
</style>
<script>
function getNetwork() {
   alert('hello');
   document.getElementById("wsOutput").value="It works";
}
</script>
</head>
<body>
<center><table class='navtable'>
<tr>`)
	for i := range headers {
		io.WriteString(w, "<th>")
		if headers[i] == title {
			io.WriteString(w, "<font size='+2' color='#ffffff'>")
		} else {
			io.WriteString(w, "<a href='"+headerMap[headers[i]]+"'>")
		}
		io.WriteString(w, headers[i])
		if headers[i] == title {
			io.WriteString(w, "</font>")
		} else {
			io.WriteString(w, "</a>")
		}
		io.WriteString(w, "</th>\n")
	}
	io.WriteString(w, `
		</tr>
		</table></center>
		<p>
		<center><table border='0'>
		<tr><td>
		`)
}

// Get physical and monitor interfaces
func getInterfaces() []string {
	ifaces, _ := net.Interfaces()

	wireless := []string{}
	for i := range ifaces {
		// Get the interface attributes
		out, _ := exec.Command("iwconfig", ifaces[i].Name).Output()
		if string(out) != "" {
			outstr := strings.Split(string(out), "\n")[0]
			matched, _ := regexp.MatchString(`Mode:Monitor`, outstr)
			if !matched {
				wireless = append(wireless, ifaces[i].Name)
			}
		}
	}
	return wireless
}

// Home
//func HttpHome(w http.ResponseWriter, req *http.Request, message chan string) {
func HttpHome(w http.ResponseWriter, req *http.Request, outdir string, ssids []string, header *template.Template, page *template.Template) {
	header.Execute(w, kismetTemplate.Header{"Home"})

	startInt, _ := strconv.ParseInt(kismet.ServerStart(), 10, 64)
	startTime := time.Unix(startInt, 0)
	hour, min, sec := startTime.Clock()
	startTxt := fmt.Sprintf("%02d:%02d:%02d", hour, min, sec)

	page.Execute(w, kismetTemplate.Home{kismet.ServerVersion(), kismet.ServerName(), startTxt, outdir, ssids, map[string]int{"nCount": 1, "cCount": 2, "rCount": 3, "pCount": 4, "pSec": 5, "filtering": 0}, []string{"wlan0", "wlan1"}})
}

// Networks
func HttpNetworks(w http.ResponseWriter, req *http.Request) {
	HttpHeader(w, "Networks")
	io.WriteString(w, `
<form>
<center><table border="1">
<tr><th colspan="4">Kismet Options</th></tr>
<tr><th>Channel(s)</th><td><input type="text" name="channels" size="25"></td>
<th>Hop</th><td>Yes<input type="radio" name="hop" value="yes" checked> No<input type="radio" name="hop" value="no"></td></tr>
<tr><th>Networks</th><td colspan="2"><select name="networks" size="1">
<option></option>
<option>test</option>
</select></td><td rowspan="3"> </td></tr>
<tr><th>BSSIDs</th><td colspan="2"><select name="bssids" size="1">
<option></option>
<option>aa:bb:cc:dd:ee:ff</option>
</select></td></tr>
<tr><th>Clients</th><td colspan="2"><select name="clients" size="1">
<option></option>
<option>Client1</option>
</select></td></tr>
</table></center>
</form>
`)
}

// Clients
func HttpClients(w http.ResponseWriter, req *http.Request) {
	HttpHeader(w, "Clients")
	io.WriteString(w, `
<form>
<center><table border="1">
<tr><th>Client</th><td><select name="client" size="1">
<option></option>
<option>Client1</option>
</select></td></tr>
<tr><th>MAC</th><td>aa:bb:cc:dd:ee:ff</td></tr>
<tr><th>Network</th><td>test</td></tr>
<tr><th>Probes</th><td>Probe1</td></tr>
</table></center>
</form>
`)
}

// Rogues
func HttpRogues(w http.ResponseWriter, req *http.Request) {
	HttpHeader(w, "Rogues")
	io.WriteString(w, `
<form>
<center><table border="1">
<tr><th>Network</th><td><select name="network" size="1">
<option></option>
<option>test</option>
</select></td></tr>
<tr><th>BSSIDs</th><td colspan="2"><select name="bssids" size="1">
<option></option>
<option>aa:bb:cc:dd:ee:ff</option>
</select></td></tr>
<tr><th>Clients</th><td colspan="2"><select name="clients" size="1">
<option></option>
<option>Client1</option>
</select></td></tr>
</table></center>
</form>
`)
}

func HttpWS(w http.ResponseWriter, req *http.Request, message chan string) {
	msg := <-message
	io.WriteString(w, msg)
}
