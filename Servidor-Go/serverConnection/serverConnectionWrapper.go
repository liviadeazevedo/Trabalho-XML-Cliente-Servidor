package serverConnection

import (
	"fmt"
	"net"
	"os"
	"bufio"
	"strconv"
	"sync"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "4446"
	CONN_TYPE = "tcp4"
)

var (
	observerList []func([]byte, int, int)
	clinetConnList []net.Conn
	wg sync.WaitGroup
)

func Init(){
	if observerList == nil{
		observerList = make([]func([]byte, int, int), 0)
	}
	if clinetConnList == nil{
		clinetConnList = make([]net.Conn, 0)
	}
}

func Restart(){
	observerList = nil
	for _, clienteSocket := range clinetConnList {
		clienteSocket.Close()
	}
	clinetConnList = nil
	Init()
}

func OpenListener(){
	Init()
	// Listen for incoming connections.
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Erro listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()

	fmt.Println("Listening em " + CONN_HOST + ":" + CONN_PORT)

	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		fmt.Println("Conexão do cliente recebida.")
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

func convertIncomingByteToNumber(byteNumber []byte) (int, error){
	return strconv.Atoi(string(byteNumber))
}

func readIncomingMsg(conn net.Conn, tamBuffer int) ([]byte, error){
	// Make a buffer to hold incoming data.
	buf := make([]byte, tamBuffer)
	// Read the incoming connection into the buffer.
	reqLen, err := bufio.NewReader(conn).Read(buf) //conn.Read(buf)

	fmt.Println("\n bytes a ler ", tamBuffer, "\n bytes lidos ", reqLen, "\n msg lida ", string(buf), "\n")
	if err != nil {
		fmt.Println("Erro durante a leitura:", err.Error())
		return nil, err
	}
	if reqLen < 0{
		return []byte("-1"), nil
	}
	//fmt.Printf("qtd de bytes lidos %d\n", reqLen)
	return buf, nil
}

func readProtocolType(conn net.Conn) int{
	fmt.Println("\nesperando receber o tipo de protocolo")
	protocoloByte, err := readIncomingMsg(conn, 1)
	if err != nil {
		fmt.Println("Erro durante a leitura do tipo de protocolo, assumindo 1 como padrão: ", err.Error())
		return 1;
	}

	protocolo, err := convertIncomingByteToNumber(protocoloByte)
	fmt.Println("número do protocolo ", protocolo)
	if err != nil || protocolo != 2{
		return 1
	}
	return protocolo
}

func readCommunication(conn net.Conn) ([]byte, error){
	fmt.Println("\nesperando ler tamanho")
	reqLen, err := readIncomingMsg(conn, 5)
	if err != nil {
		return nil, err
	}
	tamMsg, err := convertIncomingByteToNumber(reqLen)
	if err != nil {
		fmt.Println("Erro durante a conversão do tamanho:", err.Error())
		return nil, err
	}
	if tamMsg < 0 {
		return reqLen, nil
	}

	fmt.Printf("tamanho lido %d\n", tamMsg)
	fmt.Println("começar a ler mensagem")
	buf, err := readIncomingMsg(conn, tamMsg)
	if err != nil {
		fmt.Println("Erro durante a leitura da mensagem:", err.Error())
		return nil, err
	}

	fmt.Println("Número de bytes lidos do cliente: ",string(reqLen))
	fmt.Println("Array de bytes lido convertido para string:\n\n ",string(buf))

	return buf, nil
}

func readSize(conn net.Conn, tamBuffer int) (int, error){
	tamByte, err := readIncomingMsg(conn, tamBuffer)
	if err != nil {
		return -1, err
	}
	fmt.Println("=====tambyte ", string(tamByte))
	tam, err := convertIncomingByteToNumber(tamByte)
	if err != nil {
		fmt.Println("Erro durante a conversão do tamanho :", err.Error())
		return -1, err
	}
	if tam < 0 {
		return -1, nil
	}
	return tam, nil
}

func readCommunicationWithHeader(conn net.Conn) ([]byte, error){
	fmt.Println("\nesperando ler tamanho do cabeçalho")
	tamCabecalho, err := readSize(conn, 2)
	if err != nil || tamCabecalho < 1{
		return []byte("-1"), err
	}

	fmt.Printf("tamanho do cabeçalho lido %d\n", tamCabecalho)
	fmt.Println("começar a ler tamanho do arquivo")
	tamMsg, err := readSize(conn, tamCabecalho)
	fmt.Println("tamanho da mensagem ", tamMsg)
	if err != nil || tamMsg < 1{
		return []byte("-1"), err
	}

	buf, err := readIncomingMsg(conn, tamMsg)
	if err != nil {
		fmt.Println("Erro durante a leitura da mensagem:", err.Error())
		return nil, err
	}

	fmt.Println("Número de bytes lidos do cliente: ",string(tamMsg))
	fmt.Println("Array de bytes lido convertido para string:\n\n ",string(buf))

	return buf, nil
}

//return arquivo, erro
func readFromClient(conn net.Conn, protocolo int) ([]byte, error){
	if protocolo == 1{
		return readCommunication(conn)
	}else {
		return readCommunicationWithHeader(conn)
	}

}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	myid := len(clinetConnList)
	clinetConnList = append(clinetConnList, conn)
	protocolo := readProtocolType(conn)
	fmt.Println("\n tipo de protocolo recebido : ", protocolo)
	for {
		buf, err := readFromClient(conn, protocolo)
		if err != nil || string(buf) == "-1"{
			break
		}
		notify(buf, myid, protocolo)
		wg.Wait()
	}
	// Close the connection when you're done with it.
	clinetConnList[myid] = nil
	conn.Close()
}

func RegisterObserver(recieveNotify func([]byte, int, int)){
	if observerList == nil{
		observerList = make([]func([]byte, int, int), 0)
	}
	observerList = append(observerList, recieveNotify)
}

func notify(msg[]byte, id int, protocolo int){
	for _, observerCallback := range observerList {
		wg.Add(1)
		observerCallback(msg, id, protocolo)
	}
}

func SendToClientWithoutHeader(msg []byte, clinetId int){
	conn := clinetConnList[clinetId]

	tamMsg := strconv.Itoa(len(msg))
	tamMsg = "00000" + tamMsg
	tamMsg = tamMsg[len(tamMsg) - 5:]

	fmt.Println("\ntamanho da mensagem enviado ", tamMsg)
	conn.Write([]byte(tamMsg))
	conn.Write(msg)
}

func SendToClientWithHeader(msg []byte, clinetId int){
	conn := clinetConnList[clinetId]

	tamMsg := strconv.Itoa(len(msg))
	
	tamHeader := strconv.Itoa(len(tamMsg))
	tamHeader = "0" + tamHeader
	tamHeader = tamHeader[len(tamHeader) - 2:]

	fmt.Println("\ntamanho do cabeçalho mensagem a ser enviado ", tamHeader)
	fmt.Println("\ncabeçalho a ser enviado ", tamMsg)
	conn.Write([]byte(tamHeader))
	conn.Write([]byte(tamMsg))
	conn.Write(msg)
}

func SendToClient(msg []byte, clinetId int, protocolo int){
	defer wg.Done()
	if protocolo == 1{
		SendToClientWithoutHeader(msg, clinetId)
	}else{
		SendToClientWithHeader(msg, clinetId)
	}
	
	return 
}