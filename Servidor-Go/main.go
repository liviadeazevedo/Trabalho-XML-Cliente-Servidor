package main

import (
	"./serverConnection"
	"./serverLog"
	"./serverLogic"
)

func recieveNotification(msg []byte, clinetId int, protocolo int) {
	var (
		xml      string
		resposta string
	)
	xml = string(msg)
	resposta = serverLogic.RequestXMLHandler(xml)
	serverConnection.SendToClient([]byte(resposta), clinetId, protocolo)
}

func main() {
	//fmt.Println("se registrando no observer")
	serverLog.PrintServerMsgOnlyTitle("Se registrando no observer...")
	serverConnection.RegisterObserver(recieveNotification)
	serverConnection.OpenListener()
}
