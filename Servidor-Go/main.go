package main

import (
	"fmt"
	"./serverLogic"
	"./serverConnection"
)

/*
func serverLogicTest(){
	xml := "<resposta><retorno>0</retorno></resposta>"
	xsd_path := "../Arquivos/resposta.xsd"

}
*/
func recieveNotify(msg []byte, clinetId int) []byte {
	var (
		xml string
		resposta string
	)
	xml = string(msg)
	resposta = serverLogic.RequestXMLHandler(xml)
	return []byte(resposta)
}

func serverConnectionTest(){
	fmt.Println("se registrando no observer")
	serverConnection.RegisterObserver(recieveNotify)
	serverConnection.OpenListener()
}

func main() {
	serverConnectionTest()
}