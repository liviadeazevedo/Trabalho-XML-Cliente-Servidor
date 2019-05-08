package main

import (
	"fmt"

	"./serverConnection"
	"./serverLogic"
	//"strconv"
)

/*
func serverLogicTest(){
	xml := "<resposta><retorno>0</retorno></resposta>"
	xsd_path := "../Arquivos/resposta.xsd"

}
*/
func recieveNotify(msg []byte, clinetId int, protocolo int) {
	var (
		xml      string
		resposta string
	)
	xml = string(msg)
	resposta = serverLogic.RequestXMLHandler(xml)
	serverConnection.SendToClient([]byte(resposta), clinetId, protocolo)
}

func serverConnectionTest() {
	fmt.Println("se registrando no observer")
	serverConnection.RegisterObserver(recieveNotify)
	serverConnection.OpenListener()
}

func main() {
	serverConnectionTest()
	//a := []byte("02")
	//fmt.Println(strconv.Atoi(string(a)))

}
