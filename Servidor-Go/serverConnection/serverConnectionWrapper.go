package serverConnection

import (
	"fmt"
	"net"
	"os"
	"bufio"
	"strconv"
//	"sync"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "4446"
	CONN_TYPE = "tcp4"
)

var (
	observerList []func([]byte, int)[]byte
	clinetConnList []net.Conn
)

func Init(){
	if observerList == nil{
		observerList = make([]func([]byte, int)[]byte, 0)
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
	//defer wg.Done()
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

func readIncomingMsg(conn net.Conn, tamBuffer int) ([]byte, error){
	// Make a buffer to hold incoming data.
	buf := make([]byte, tamBuffer)
	// Read the incoming connection into the buffer.
	reqLen, err := bufio.NewReader(conn).Read(buf) //conn.Read(buf)

	if err != nil {
		fmt.Println("Erro durante a leitura:", err.Error())
		return nil, err
	}
	if reqLen < 0{
		return []byte("-1"), nil
	}
	fmt.Printf("qtd de bytes lidos %d\n", reqLen)
	return buf, nil
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	myid := len(clinetConnList)
	clinetConnList = append(clinetConnList, conn)
	
	for {
		fmt.Println("\nesperando ler tamanho")
		reqLen, err := readIncomingMsg(conn, 5)
		if err != nil {
			return
		}
		tamMsg, err := strconv.Atoi(string(reqLen))
		if err != nil {
			fmt.Println("Erro durante a conversão do tamanho:", err.Error())
			return
		}
		if tamMsg < 0 {
			break
		}
		fmt.Printf("tamanho lido %d\n", tamMsg)
		fmt.Println("começar a ler mensagem")
		buf, err := readIncomingMsg(conn, tamMsg)
		if err != nil {
			fmt.Println("Erro durante a leitura da mensagem:", err.Error())
			return
		}

		fmt.Println("Número de bytes lidos do cliente: ",string(reqLen))
		fmt.Println("Array de bytes lido convertido para string:\n\n ",string(buf))

		notify(buf, myid)
		// Send a response back to person contacting us.
		//conn.Write([]byte("Message received."))
	}
	// Close the connection when you're done with it.
	clinetConnList[myid] = nil
	conn.Close()
}
func RegisterObserver(recieveNotify func([]byte, int)[]byte){
	if observerList == nil{
		observerList = make([]func([]byte, int)[]byte, 0)
	}
	observerList = append(observerList, recieveNotify)
}

func notify(msg[]byte, id int){
	for _, observerCallback := range observerList {
		SendToClient(observerCallback(msg, id), id)
	}
}

func SendToClient(msg []byte, clinetId int){
	conn := clinetConnList[clinetId]

	tamMsg := strconv.Itoa(len(msg))
	tamMsg = "00000" + tamMsg
	tamMsg = tamMsg[len(tamMsg) - 5:]

	fmt.Println("\ntamanho da mensagem enviado ", tamMsg)
	conn.Write([]byte(tamMsg))
	conn.Write(msg)
	
	return 
}