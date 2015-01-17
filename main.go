package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/joeybelans/gokismet/kdb"
	"github.com/joeybelans/gokismet/kismet"
	"github.com/joeybelans/gokismet/kismetTemplate"
)

// Prints the program usage
func Usage() {
	fmt.Println("Usage: gokismet -lport <web port> -khost <kismet IP> -kport <kismet port> [-ssid <essid>[,<essid>...]]\n")
	fmt.Println("Defaults:")
	flag.PrintDefaults()
}

// Command line arguments
type ssid []string

func (s *ssid) String() string {
	return fmt.Sprint(*s)
}

func (s *ssid) Set(value string) error {
	if len(*s) > 0 {
		return errors.New("SSID already set")
	}
	for _, ssid := range strings.Split(value, ",") {
		flag := false
		for _, ele := range *s {
			if ele == ssid {
				flag = true
				break
			}
		}
		if !flag {
			*s = append(*s, ssid)
		}
	}
	return nil
}

// Main function
func main() {
	// Parse arguments
	var (
		ssids  ssid
		lhost  = flag.String("lhost", "127.0.0.1", "HTTP service host")
		lport  = flag.Int("lport", 8080, "HTTP service port")
		khost  = flag.String("khost", "127.0.0.1", "Kismet service host")
		kport  = flag.Int("kport", 2501, "Kismet service port")
		dbfile = flag.String("sql", "gokismet.db", "SQLite3 filename")
		outdir = flag.String("outdir", ".", "Kismet output file directory")
		debug  = flag.Bool("debug", false, "Debug flag")
	)
	flag.Var(&ssids, "ssids", "List of in-scope SSIDs (comma separated)")
	flag.Parse()

	// Display SSIDs, if any
	if *debug && len(ssids) > 0 {
		fmt.Println("lhost ", *lhost, ": lport ", *lport, ": khost ", *khost, ": kport ", *kport, ": dbfile ", *dbfile, ": outdir ", *outdir, ": ssids ", ssids)
		fmt.Println("In-scope networks:")
		for i := 0; i < len(ssids); i++ {
			fmt.Printf("\t%s\n", ssids[i])
		}
	}
	fmt.Println()

	// Get the full output directory path
	*outdir, _ = filepath.Abs(*outdir)

	// Create/Open sqlite3 file
	db, err := sql.Open("sqlite3", *dbfile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	kdb.CreateDB(db, *dbfile)

	// Start kismet handler
	kismet.Run(*khost, *kport, db, *debug, ssids)

	// HTTP handlers
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		kismetTemplate.HttpHome(w, r, *outdir+"/"+*dbfile, ssids)
	})
	http.HandleFunc("/discover", func(w http.ResponseWriter, r *http.Request) {
		kismetTemplate.HttpDiscover(w, r, *outdir+"/"+*dbfile, ssids)
	})
	http.HandleFunc("/global.css", func(w http.ResponseWriter, r *http.Request) {
		kismetTemplate.HttpCSS(w, r)
	})
	http.HandleFunc("/kismet.js", func(w http.ResponseWriter, r *http.Request) {
		kismetTemplate.HttpKismetJS(w, r)
	})
	http.HandleFunc("/discover.js", func(w http.ResponseWriter, r *http.Request) {
		kismetTemplate.HttpDiscoverJS(w, r, ssids)
	})
	http.HandleFunc("/ws", kismet.ServeWS)

	// Start web service
	fmt.Printf("Browse to http://%s:%d to access the web interface\n", *lhost, *lport)
	fmt.Println("Press CTRL+C to stop the server")
	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", *lhost, *lport), nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
