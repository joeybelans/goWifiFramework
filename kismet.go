package main

import (
   "net"
   "fmt"
)

type kismetclient struct {
   test string
}

var kismet = kismetclient {"test string"}

func (kismet *kismetclient) run(host string, port int, message (chan string)) {
   conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
   if err != nil {
      fmt.Println("Can't connect to kismet server")
   }
   fmt.Fprintf(conn, "test")
}

/*
status, err := bufio.NewReader(conn).ReadString('\n')
*/
