package main

import (
   "fmt"
   "flag"
   "net/http"
   "log"
   "database/sql"
   "os"
   _ "github.com/mattn/go-sqlite3"
)

// Create networks structure
type networks []string
 
func (n *networks) String() string {
    return fmt.Sprintf("%s", *n)
}

func (n *networks) Set(network string) error {
   *n = append(*n, network)
   return nil
}
 
// Prints the program usage
func Usage() {
   fmt.Println("Usage: gokismet -lport <web port> -khost <kismet IP> -kport <kismet port> [-ssid <essid>]...\n")
   fmt.Println("Defaults:")
   flag.PrintDefaults()
}

// Command line arguments
var (
   mynetworks networks
   lport = flag.Int("lport", 8080, "Local HTTP Service Port")
   khost = flag.String("khost", "127.0.0.1", "Kismet Service Host")
   kport = flag.Int("kport", 2501, "Kismet Service Port")
   dbfile = flag.String("sql", "gokismet.db", "SQLite3 Filename")
)

// Main function
func main() {
   // Parse arguments
   flag.Var(&mynetworks, "ssid", "List of In-Scope SSIDs")
   flag.Parse()

   // Display SSIDs, if any
   if len(mynetworks) > 0 {
      fmt.Println("In-scope networks:")
      for i := 0; i < len(mynetworks); i++ {
         fmt.Printf("\t%s\n", mynetworks[i])
      }
   }
   fmt.Println()

   // Create/Open sqlite3 file
   db, err := sql.Open("sqlite3", *dbfile)
   if err != nil {
      log.Fatal(err)
   }
   defer db.Close()

   // Create the schema, if necessary
   _, err = os.Stat(*dbfile)
   if os.IsNotExist(err) {
      CreateDB(db)
   }

   // Start kismet handler
   message := make(chan string)
   go kismet.run(*khost, *kport, message)

   // HTTP handlers
   http.HandleFunc("/", HttpHome)
   http.HandleFunc("/networks/", HttpNetworks)
   http.HandleFunc("/clients/", HttpClients)
   http.HandleFunc("/rogues/", HttpRogues)
   http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
      HttpWS(w, r, message)
   })

   // Start web service
   fmt.Printf("Browse to http://127.0.0.1:%d to access the web interface\n", *lport)
   fmt.Println("Press CTRL+C to stop the server")
   if err := http.ListenAndServe(fmt.Sprintf(":%d", *lport), nil); err != nil {
      log.Fatal("ListenAndServe:", err)
   }
}
