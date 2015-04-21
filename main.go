package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/joeybelans/gokismet/header"
	"github.com/joeybelans/gokismet/home"
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

type pkgMap struct {
	Title     string
	Processor interface{}
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

	// Initialize packages
	// Each package will create a link on the navigation banner in the order the packages are loaded
	// Each Init function should return the address of its ProcessSocket function
	//packages := map[string]pkgMap{}
	packages := []interface{}
	header.Init()
	//packages["/"] = pkgMap{"Home", home.Init(*lhost, *lport, *khost, *kport, ssids)}
	packages = Append(interface{"/", "Home", home.Init(*lhost, *lport, *khost, *kport, ssids)})
	fmt.Println(packages)

	// Create the HTML header
	header.Create(packages)

	/*
		kismetHandler.Init(*khost, *kport, *debug, ssids)
	*/

	// Start web service
	fmt.Printf("Browse to http://%s:%d to access the web interface\n", *lhost, *lport)
	fmt.Println("Press CTRL+C to stop the server")
	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", *lhost, *lport), nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
