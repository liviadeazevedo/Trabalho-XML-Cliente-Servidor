package serverConnection

import (
	"net"
	"os"
	"strconv"
	"sync"

	"Trabalho-XML-Cliente-Servidor/Servidor-Go/serverLog"
)

const (
	CONN_HOST            = "localhost"
	CONN_PORT            = "4446"
	CONN_TYPE            = "tcp4"
	MAX_TAM_MSG_TO_PRINT = 400
	MAX_PKG_SIZE         = 1400
)

var (
	observerList   []func([]byte, int, int)
	clinetConnList []net.Conn
	wg             sync.WaitGroup
)

func Init() {
	if observerList == nil {
		observerList = make([]func([]byte, int, int), 0)
	}
	if clinetConnList == nil {
		clinetConnList = make([]net.Conn, 0)
	}
}

func Restart() {
	for _, clienteSocket := range clinetConnList {
		clienteSocket.Close()
	}
	clinetConnList = nil
	Init()
}

func OpenListener(ip string, port string) {
	Init()
	// Listen for incoming connections.
	if ip == "" {
		ip = CONN_HOST
	}
	if port == "" {
		port = CONN_PORT
	}

	l, err := net.Listen(CONN_TYPE, ip+":"+port)
	if err != nil {
		//fmt.Println("Erro listening:", err.Error())
		serverLog.PrintErrorMsg("Erro listening: " + err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()

	//fmt.Println("Listening em " + CONN_HOST + ":" + CONN_PORT)
	serverLog.PrintWaitingMsg("Listening em " + ip + ":" + port + "...")

	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			//fmt.Println("Error accepting: ", err.Error())
			serverLog.PrintErrorMsg("Error accepting: " + err.Error())
			os.Exit(1)
		}
		//fmt.Println("Conexão do cliente recebida.")
		serverLog.PrintServerMsg("Conexão do cliente recebida.", false)
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

func convertIncomingByteToNumber(byteNumber []byte) (int, error) {
	return strconv.Atoi(string(byteNumber))
}

func truncateMsgToPrint(msg []byte) string {
	msgToPrint := string(msg)
	if len(msgToPrint) > MAX_TAM_MSG_TO_PRINT {
		msgToPrint = msgToPrint[:MAX_TAM_MSG_TO_PRINT] + " [...]"
	}
	return msgToPrint
}

func readBufferLimitedIncomingMsg(conn net.Conn, tamBuffer int) ([]byte, error) {
	auxBuff := make([]byte, 0)
	truncated := false
	if tamBuffer > MAX_PKG_SIZE {
		truncated = true
		count := 0
		for true {
			buf, err := readIncomingMsg(conn, MAX_PKG_SIZE)
			if err != nil || string(buf) == "-1" {
				return buf, err
			}
			auxBuff = append(auxBuff, buf...)
			count += 1
			valorLido := count * MAX_PKG_SIZE
			restaLer := tamBuffer - valorLido
			if restaLer < MAX_PKG_SIZE {
				tamBuffer = restaLer
				break
			}
		}
	}

	buf, err := readIncomingMsg(conn, tamBuffer)
	if err != nil || string(buf) == "-1" {
		return buf, err
	}

	if truncated {
		return append(auxBuff, buf...), nil
	}

	return buf, nil
}

func readIncomingMsg(conn net.Conn, tamBuffer int) ([]byte, error) {

	buf := make([]byte, tamBuffer)
	reqLen, err := conn.Read(buf)
	//fmt.Println("\n bytes a ler", tamBuffer, "\n bytes lidos", reqLen, "\n msg lida", string(buf), "\n")
	if err != nil {
		//fmt.Println("Erro durante a leitura:", err.Error())
		serverLog.PrintErrorMsg("Erro durante a leitura: " + err.Error())
		return nil, err
	}
	if reqLen < 0 {
		return []byte("-1"), nil
	}
	//fmt.Printf("qtd de bytes lidos %d\n", reqLen)
	return buf, nil
}

func readProtocolType(conn net.Conn) int {
	//fmt.Println("\nesperando receber o tipo de protocolo")
	protocoloByte, err := readIncomingMsg(conn, 1)
	if err != nil {
		//fmt.Println("Erro durante a leitura do tipo de protocolo, assumindo 1 como padrão: ", err.Error())
		serverLog.PrintErrorMsg("Erro durante a leitura do tipo de protocolo, assumindo 1 como padrão: " + err.Error())
		return 1
	}

	protocolo, err := convertIncomingByteToNumber(protocoloByte)
	//fmt.Println("número do protocolo", protocolo)
	if err != nil || protocolo != 2 {
		return 1
	}
	return protocolo
}

func readCommunication(conn net.Conn) ([]byte, error) {
	//fmt.Println("\nEsperando ler tamanho")
	serverLog.PrintWaitingMsg("Esperando ler tamanho...")
	reqLen, err := readIncomingMsg(conn, 5)
	if err != nil {
		return nil, err
	}
	tamMsg, err := convertIncomingByteToNumber(reqLen)
	if err != nil {
		//fmt.Println("Erro durante a conversão do tamanho:", err.Error())
		serverLog.PrintErrorMsg("Erro durante a conversão do tamanho: " + err.Error())
		return nil, err
	}
	if tamMsg < 0 {
		return reqLen, nil
	}

	//fmt.Printf("tamanho lido %d\n", tamMsg)
	serverLog.PrintServerMsg("Tamanho lido: "+strconv.Itoa(tamMsg), false)
	//fmt.Println("começar a ler mensagem")
	serverLog.PrintWaitingMsg("Começar a ler mensagem...")
	buf, err := readBufferLimitedIncomingMsg(conn, tamMsg)
	if err != nil {
		//fmt.Println("Erro durante a leitura da mensagem:", err.Error())
		serverLog.PrintErrorMsg("Erro durante a leitura da mensagem: " + err.Error())
		return nil, err
	}

	serverLog.PrintServerMsg("Número de bytes lidos do cliente: "+string(reqLen)+"\n\n Array de bytes lido convertido para string:\n\n "+truncateMsgToPrint(buf), false)
	//fmt.Println("Número de bytes lidos do cliente:", string(reqLen))
	//fmt.Println("Array de bytes lido convertido para string:\n\n", string(buf))

	return buf, nil
}

func readSize(conn net.Conn, tamBuffer int) (int, error) {
	tamByte, err := readIncomingMsg(conn, tamBuffer)
	if err != nil {
		return -1, err
	}
	//fmt.Println("=====tambyte", string(tamByte))
	tam, err := convertIncomingByteToNumber(tamByte)
	if err != nil {
		//fmt.Println("Erro durante a conversão do tamanho :", err.Error())
		serverLog.PrintErrorMsg("Erro durante a conversão do tamanho: " + err.Error())
		return -1, err
	}
	if tam < 0 {
		return -1, nil
	}
	return tam, nil
}

func readCommunicationWithHeader(conn net.Conn) ([]byte, error) {
	serverLog.PrintWaitingMsg("Esperando ler tamanho do cabeçalho...")
	//fmt.Println("\nesperando ler tamanho do cabeçalho")
	tamCabecalho, err := readSize(conn, 2)
	if err != nil || tamCabecalho < 1 {
		return []byte("-1"), err
	}
	//PAREI AQUI A FORMATAÇÃO DA MENSAGEM DO SERVIDOR!
	//fmt.Println("tamanho do cabeçalho lido", tamCabecalho)

	//fmt.Println("\ncomeçar a ler tamanho do arquivo")

	serverLog.PrintServerMsg("Tamanho do cabeçalho lido: "+strconv.Itoa(tamCabecalho), false)
	serverLog.PrintWaitingMsg("Começar a ler tamanho do arquivo...")

	tamMsg, err := readSize(conn, tamCabecalho)
	//fmt.Println("tamanho do arquivo", tamMsg)
	serverLog.PrintServerMsg("Tamanho do arquivo: "+strconv.Itoa(tamMsg), false)
	if err != nil || tamMsg < 1 {
		return []byte("-1"), err
	}

	//fmt.Println("\ncomeçar a ler arquivo de", tamMsg, "bytes")
	serverLog.PrintWaitingMsg("Começar a ler arquivo de " + strconv.Itoa(tamMsg) + " bytes...")
	buf, err := readBufferLimitedIncomingMsg(conn, tamMsg)
	if err != nil {
		//fmt.Println("Erro durante a leitura do arquivo:", err.Error())
		serverLog.PrintErrorMsg(err.Error())
		return nil, err
	}

	//fmt.Println("\nNúmero de bytes lidos do cliente:", tamMsg)
	//fmt.Println("Array de bytes lido convertido para string:\n\n", truncateMsgToPrint(buf))
	serverLog.PrintServerMsg("Número de bytes lidos do cliente: "+strconv.Itoa(tamMsg)+"\n\nArray de bytes lido convertido para string:\n\n"+truncateMsgToPrint(buf), false)

	return buf, nil
}

//return arquivo, erro
func readFromClient(conn net.Conn, protocolo int) ([]byte, error) {
	if protocolo == 1 {
		return readCommunication(conn)
	} else {
		return readCommunicationWithHeader(conn)
	}

}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	myid := len(clinetConnList)
	clinetConnList = append(clinetConnList, conn)
	//fmt.Println("\n Esperando o tipo de protocolo do cliente de número", myid)
	serverLog.PrintWaitingMsg("Esperando o tipo de protocolo do cliente de número " + strconv.Itoa(myid) + "...")

	protocolo := readProtocolType(conn)
	//fmt.Println("\n Tipo de protocolo recebido :", protocolo)
	serverLog.PrintServerMsg("Tipo de protocolo recebido: "+strconv.Itoa(protocolo), false)

	for {
		//fmt.Println("\n Esperando mensagem do cliente")
		serverLog.PrintWaitingMsg("Esperando mensagem do cliente...")
		buf, err := readFromClient(conn, protocolo)
		if err != nil || string(buf) == "-1" {
			//fmt.Println("\n Conexão com o cliente", myid, "encerrada")
			serverLog.PrintErrorMsg("Conexão com o cliente " + strconv.Itoa(myid) + "encerrada")
			break
		}
		//fmt.Println("\n\n Executando lógica do servidor... \n")
		serverLog.PrintWaitingMsg("Executando lógica do servidor...")
		notify(buf, myid, protocolo)
		wg.Wait()
	}
	// Close the connection when you're done with it.
	clinetConnList[myid] = nil
	conn.Close()
}

func RegisterObserver(recieveNotify func([]byte, int, int)) {
	if observerList == nil {
		observerList = make([]func([]byte, int, int), 0)
	}
	observerList = append(observerList, recieveNotify)
}

func notify(msg []byte, id int, protocolo int) {
	for _, observerCallback := range observerList {
		wg.Add(1)
		observerCallback(msg, id, protocolo)
	}
}

func SendToClientWithoutHeader(msg []byte, clinetId int) {
	conn := clinetConnList[clinetId]

	tamMsg := strconv.Itoa(len(msg))
	tamMsg = "00000" + tamMsg
	tamMsg = tamMsg[len(tamMsg)-5:]

	//fmt.Println("\nTamanho da mensagem a ser enviada", tamMsg)
	serverLog.PrintServerMsg("Tamanho da mensagem a ser enviada "+tamMsg, false)
	conn.Write([]byte(tamMsg))
	//fmt.Println("Enviando mensagem \n", truncateMsgToPrint(msg), "\n")
	serverLog.PrintWaitingMsg("Enviando mensagem:\n\n" + truncateMsgToPrint(msg) + "\n\n...")
	conn.Write(msg)
	//fmt.Println("Mensagem enviada...\n")
	serverLog.PrintServerMsg("Mensagem enviada!", false)
}

func SendToClientWithHeader(msg []byte, clinetId int) {
	conn := clinetConnList[clinetId]

	tamMsg := strconv.Itoa(len(msg))

	tamHeader := strconv.Itoa(len(tamMsg))
	tamHeader = "0" + tamHeader
	tamHeader = tamHeader[len(tamHeader)-2:]

	//fmt.Println("\nTamanho do cabeçalho mensagem a ser enviado", tamHeader)
	serverLog.PrintServerMsg("Tamanho do cabeçalho mensagem a ser enviado: "+tamHeader, false)

	conn.Write([]byte(tamHeader))

	//fmt.Println("Cabeçalho a ser enviado", tamMsg)
	serverLog.PrintServerMsg("Cabeçalho a ser enviado"+tamMsg, false)
	conn.Write([]byte(tamMsg))
	//fmt.Println("Enviando mensagem \n\n'", truncateMsgToPrint(msg), "'\n")
	serverLog.PrintWaitingMsg("Enviando mensagem: \n\n" + truncateMsgToPrint(msg) + "\n\n...")
	conn.Write(msg)
	//fmt.Println("Mensagem enviada...\n")
	serverLog.PrintServerMsg("Mensagem enviada!", false)
}

func SendToClient(msg []byte, clinetId int, protocolo int) {
	defer wg.Done()
	//fmt.Println("\n\n Enviando mensagem ao cliente de número", clinetId, "pelo protocolo", protocolo)
	serverLog.PrintWaitingMsg("Enviando mensagem ao cliente de número " + strconv.Itoa(clinetId) + " pelo protocolo " + strconv.Itoa(protocolo) + "...")
	if protocolo == 1 {
		SendToClientWithoutHeader(msg, clinetId)
	} else {
		SendToClientWithHeader(msg, clinetId)
	}

	return
}
