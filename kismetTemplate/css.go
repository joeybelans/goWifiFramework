package kismetTemplate

import (
	"net/http"
	"net/url"
)

// GlobalCSS
type css struct {
	Referer string
}

func HttpCSS(w http.ResponseWriter, req *http.Request) {
	referer, _ := url.Parse(req.Referer())
	w.Header().Set("Content-Type", "text/css")
	templates["/global.css"].Execute(w, css{referer.Path})
}

func tmplGlobalCSS() string {
	return `
label {
	background: transparent;
	border: 0;
	margin: 0;
	padding: 0;
	vertical-align: baseline;
}
#wrapper {
	min-width: 600px;
}
.settings {
	display: table;
	height: 45px;
	width: 100%;
}
.settings .switch {
	display: table-cell;
	vertical-align: middle;
	padding: 0;
}
.cmn-toggle {
	position: absolute;
	margin-left: -9999px;
	visibility: hidden;
}
.cmn-toggle + label {
	display: block;
	position: relative;
	cursor: pointer;
	outline: none;
	-webkit-user-select: none;
	-moz-user-select: none;
	-ms-user-select: none;
	user-select: none;
}
input.cmn-toggle-on-off + label {
	padding: 0;
	width: 150px;
	height: 45px;
}
input.cmn-toggle-on-off + label:before, input.cmn-toggle-on-off + label:after {
	display: block;
	position: absolute;
	top: 0;
	left: 0;
	bottom: 0;
	right: 0;
	font-family: "Roboto Slab", serif;
	font-size: 16px;
	text-align: center;
	font-weight: bold;
	line-height: 45px;
}
input.cmn-toggle-on-off + label:before {
	color: #ff0000;
	background-color: #dddddd;
	content: attr(data-off);
	-webkit-transition: -webkit-transform 0.5s;
	-moz-transition: -moz-transform 0.5s;
	-o-transition: -o-transform 0.5s;
	transition: transform 0.5s;
	-webkit-backface-visibility: hidden;
	-moz-backface-visibility: hidden;
	-ms-backface-visibility: hidden;
	-o-backface-visibility: hidden;
	backface-visibility: hidden;
}
input.cmn-toggle-on-off + label:after {
	color: #004415;
	background-color: #ddeebb;
	content: attr(data-on);
	-webkit-transition: -webkit-transform 0.5s;
	-moz-transition: -moz-transform 0.5s;
	-o-transition: -o-transform 0.5s;
	transition: transform 0.5s;
	-webkit-transform: rotateY(180deg);
	-moz-transform: rotateY(180deg);
	-ms-transform: rotateY(180deg);
	-o-transform: rotateY(180deg);
	transform: rotateY(180deg);
	-webkit-backface-visibility: hidden;
	-moz-backface-visibility: hidden;
	-ms-backface-visibility: hidden;
	-o-backface-visibility: hidden;
	backface-visibility: hidden;
}
input.cmn-toggle-on-off:checked + label:before {
	-webkit-transform: rotateY(180deg);
	-moz-transform: rotateY(180deg);
	-ms-transform: rotateY(180deg);
	-o-transform: rotateY(180deg);
	transform: rotateY(180deg);
}
input.cmn-toggle-on-off:checked + label:after {
	-webkit-transform: rotateY(0);
	-moz-transform: rotateY(0);
	-ms-transform: rotateY(0);
	-o-transform: rotateY(0);
	transform: rotateY(0);
}
div#hmenu {
	margin: 0;
	padding: .3em 0 .3em 0;
	background: #ddeebb;
	width: 100%;
	height: 35px;
	text-align: center;
}
div#hmenu ul {
	list-style: none;
	margin: 0;
	padding: 0;
}
div#hmenu ul li {
	margin: 0;
	padding: 0;
	display: inline;
}
div#hmenu ul li.current-page {
	margin: 0;
	padding: 0;
	display: inline;
	text-decoration: none;
	font-weight: bold;
	font-size: x-large;
	color: #004415;
}
div#hmenu ul a:link{
	margin: 0;
	padding: .3em .4em .3em .4em;
	text-decoration: none;
	font-weight: bold;
	font-size: medium;
	color: #004415;
}
div#hmenu ul a:visited{
	margin: 0;
	padding: .3em .4em .3em .4em;
	text-decoration: none;
	font-weight: bold;
	font-size: medium;
	color: #004415;
}
div#hmenu ul a:active{
	margin: 0;
	padding: .3em .4em .3em .4em;
	text-decoration: none;
	font-weight: bold;
	font-size: medium;
	color: #227755;
}
div#hmenu ul a:hover{
	margin: 0;
	padding: .3em .4em .3em .4em;
	text-decoration: none;
	font-weight: bold;
	font-size: medium;
	color: #f6f0cc;
	background-color: #227755;
}
{{if eq .Referer "/discover"}}
span.stitle {
   font-weight: bold;
   font-size: medium;
   color: #004415;
}
table.stats {
   margin-left: 20px;
   border: 0;
}
table.stats th {
   font-weight: bold;
   font-size: small;
   color: #004415;
   vertical-align: top;
}
table.stats td {
   font-size: small;
   color: #004415;
   vertical-align: top;
}
table.data {
   margin-left: 20px;
   border-collapse:collapse;
   height: 100px;
}
table.data tbody,
table.data thead { display: block; }
table.data th {
   font-weight: bold;
   font-size: medium;
   text-align: left;
   padding-left: 10px;
   padding-right: 10px;
   padding-top: 5px;
   padding-bottom: 5px;
}
table.data thead{
   height: 30px;
}
table.data tbody {
   height: 500px;
   overflow-y: auto;
   overflow-x: hidden;
}
table.data td {
   padding-left: 10px;
   padding-right: 10px;
   padding-top: 5px;
   padding-bottom: 5px;
}
table.data tr:nth-child(odd) td {
   font-size: small;
   background: #ffffff;
}
table.data tr:nth-child(even) td {
   font-size: small;
   background: #f5f9fa;
}
table.data tr:hover td {
   font-size: small;
   background: #fffbae;
}
table.data .network { width: 200px; }
table.data .bssid { width: 125px; }
table.data .channel { width: 50px; }
table.data .last { width: 200px; }
table.data .power { width: 50px; }
table.data .max { width: 50px; }
table.data .clients { width: 50px; }
table.data .packets { width: 50px; }
ul.stats li {
   font-size: small;
   color: #004415;
}
a.stats:link { 
   color: #227755; 
} 
a.stats:visited { 
   color: #227755; 
} 
a.stats:active { 
   color: #227755; 
} 
a.stats:hover { 
	color: #116644; 
}
a.data:link { 
   color: #000000
} 
a.data:visited { 
   color: #000000; 
} 
a.data:active { 
   color: #000000; 
} 
a.data:hover { 
   color: #000000; 
}
.verticalLine {
   border-right: 2px solid #ddeebb;
}
.selected {
   background: #ff0000;
}
.divFilter {
   display: none;
   position: absolute;
   text-align: left;
   border-style: solid;
   background-color: #ffffff;
   padding: 3px;
}{{end}}`
}
