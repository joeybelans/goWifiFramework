wsServer = window.location.hostname
if (window.location.port != "") {
   wsServer = wsServer + ":" + window.location.port;
}
conn = new WebSocket("ws://" + wsServer + "/ws");

conn.onopen = function (event) {
   conn.send(window.location.pathname);
};

conn.onmessage = function(evt) {
   var json = JSON.parse(evt.data);
   //alert(json.message)
   //window[json.component].process(json.message, json.obj);
   //kismet.process(json.message, json.obj);
}

function numSort(a,b) {
   return(a-b)
}
