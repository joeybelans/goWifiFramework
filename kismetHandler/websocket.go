package kismetHandler

import "github.com/joeybelans/gokismet/webSocketHandler"

func okToSend(message string) bool {
	return webSocketHandler.OKToSend("kismet", message)
}

func sendString(message string, obj map[string]string) {
	webSocketHandler.SendString("kismet", message, obj)
}

func sendInterface(message string, obj map[string]interface{}) {
	webSocketHandler.SendInterface("kismet", message, obj)
}

/*
func processKismet(data map[string]interface{}) {
	switch data["cmd"].(string) {
	case "GetNicStats":
		fmt.Println("GetNicStats")
	}
}
*/
