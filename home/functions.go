// Manages home interface
// Miscellaneous functions
package home

import (
	"fmt"
	"net"
	"os"

	_ "github.com/joeybelans/gokismet/statik"
)

// Get the list of available interfaces
func getInterfaces() []string {
	interfaces, _ := net.Interfaces()

	var iNames []string
	for _, iface := range interfaces {
		if _, err := os.Stat("/sys/class/net/" + iface.Name + "/wireless"); err == nil {
			iNames = append(iNames, iface.Name)
		}
	}

	return iNames
}

func ProcessCommand() {
	fmt.Println("HOME")
}
