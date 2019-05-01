package main

import (
	"fmt"
	//"./serverLogic"
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
	fmt.Println("\n\n======== mensagem recebida do cliente: \n", xml)
	//resposta = serverLogic.RequestXMLHandler(xml)
	resposta = "vc enviou yeah!"
	fmt.Println("======== resposta a ser enviada para o cliente: \n", resposta)
	return []byte(resposta)
	//serverConnection.SendToClient([]byte(resposta), clinetId)
	
	//fmt.Println("fui notificado com : '", string(msg), "'")
}

func serverConnectionTest(){
	fmt.Println("se registrando no observer")
	serverConnection.RegisterObserver(recieveNotify)
	serverConnection.OpenListener()
}

func main() {
	serverConnectionTest()
/*	
	a := 123

	b := strconv.Itoa(a)
	b = "00000"+b

	fmt.Println(b[len(b)-5:])
*/	
	/*
	serverConnection.Init()
	serverConnection.RegisterObserver(recieveNotify)

	serverConnection.Notify()
	*/

}