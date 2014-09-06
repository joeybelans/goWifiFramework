package main

import (
	"io"
	"net/http"
)

// Header function
func HttpHeader(w http.ResponseWriter) {
   io.WriteString(w, `
<html>
<head></head>
<body>
<center><table>
<tr><td><a href="/">Home/Stats</a></td><td><a href="/networks/">Networks</a></td><td><a href="/clients/">Clients</a></td><td><a href="/rogues/">Rogues</a></td></tr>
</table></center>
<p>
`)
}

// Home/Stats
func HttpHome(w http.ResponseWriter, req *http.Request) {
   HttpHeader(w)
   io.WriteString(w, `
<center><table>
<tr><th>In-scope Networks</th><td>test</td></tr>
<tr><th>Wireless Interfaces</th><td>wlan0</td></tr>
<tr><th>Monitor Interfaces</th><td>mon0</td></tr>
<tr><th>Kismet (host:port)</th><td>127.0.0.1:2501</td></tr>
<tr><th>Discovered Clients</th><td>10</td></tr>
<tr><th>Discovered BSSIDS</th><td>323</td></tr>
</table></center>
`)
}

// Networks
func HttpNetworks(w http.ResponseWriter, req *http.Request) {
   HttpHeader(w)
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
   HttpHeader(w)
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
   HttpHeader(w)
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

func HttpWS(w http.ResponseWriter, req *http.Request, message (chan string)) {
   msg := <-message
   io.WriteString(w, msg)
}
