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
	if ip == "" {
		ip = CONN_HOST
	}
	if port == "" {
		port = CONN_PORT
	}

	l, err := net.Listen(CONN_TYPE, ip+":"+port)
	if err != nil {
		serverLog.PrintErrorMsg("Erro listening: " + err.Error())
		os.Exit(1)
	}
	defer l.Close()

	serverLog.PrintWaitingMsg("Listening em " + ip + ":" + port + "...")

	for {
		conn, err := l.Accept()
		if err != nil {
			serverLog.PrintErrorMsg("Error accepting: " + err.Error())
			os.Exit(1)
		}

		serverLog.PrintServerMsg("Conexão do cliente recebida.", false)
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

	if err != nil {

		serverLog.PrintErrorMsg("Erro durante a leitura: " + err.Error())
		return nil, err
	}
	if reqLen < 0 {
		return []byte("-1"), nil
	}
	return buf, nil
}

func readProtocolType(conn net.Conn) int {

	protocoloByte, err := readIncomingMsg(conn, 1)
	if err != nil {

		serverLog.PrintErrorMsg("Erro durante a leitura do tipo de protocolo, assumindo 1 como padrão: " + err.Error())
		return 1
	}

	protocolo, err := convertIncomingByteToNumber(protocoloByte)

	if err != nil || protocolo != 2 {
		return 1
	}
	return protocolo
}

func readCommunication(conn net.Conn) ([]byte, error) {

	serverLog.PrintWaitingMsg("Esperando ler tamanho...")
	reqLen, err := readIncomingMsg(conn, 5)
	if err != nil {
		return nil, err
	}
	tamMsg, err := convertIncomingByteToNumber(reqLen)
	if err != nil {

		serverLog.PrintErrorMsg("Erro durante a conversão do tamanho: " + err.Error())
		return nil, err
	}
	if tamMsg < 0 {
		return reqLen, nil
	}

	serverLog.PrintServerMsg("Tamanho lido: "+strconv.Itoa(tamMsg), false)

	serverLog.PrintWaitingMsg("Começar a ler mensagem...")
	buf, err := readBufferLimitedIncomingMsg(conn, tamMsg)
	if err != nil {

		serverLog.PrintErrorMsg("Erro durante a leitura da mensagem: " + err.Error())
		return nil, err
	}

	serverLog.PrintServerMsg("Número de bytes lidos do cliente: "+string(reqLen)+"\n\n Array de bytes lido convertido para string:\n\n "+truncateMsgToPrint(buf), false)

	return buf, nil
}

func readSize(conn net.Conn, tamBuffer int) (int, error) {
	tamByte, err := readIncomingMsg(conn, tamBuffer)
	if err != nil {
		return -1, err
	}

	tam, err := convertIncomingByteToNumber(tamByte)
	if err != nil {

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

	tamCabecalho, err := readSize(conn, 2)
	if err != nil || tamCabecalho < 1 {
		return []byte("-1"), err
	}

	serverLog.PrintServerMsg("Tamanho do cabeçalho lido: "+strconv.Itoa(tamCabecalho), false)
	serverLog.PrintWaitingMsg("Começar a ler tamanho do arquivo...")

	tamMsg, err := readSize(conn, tamCabecalho)

	serverLog.PrintServerMsg("Tamanho do arquivo: "+strconv.Itoa(tamMsg), false)
	if err != nil || tamMsg < 1 {
		return []byte("-1"), err
	}

	serverLog.PrintWaitingMsg("Começar a ler arquivo de " + strconv.Itoa(tamMsg) + " bytes...")
	buf, err := readBufferLimitedIncomingMsg(conn, tamMsg)
	if err != nil {

		serverLog.PrintErrorMsg(err.Error())
		return nil, err
	}

	serverLog.PrintServerMsg("Número de bytes lidos do cliente: "+strconv.Itoa(tamMsg)+"\n\nArray de bytes lido convertido para string:\n\n"+truncateMsgToPrint(buf), false)

	return buf, nil
}

func readFromClient(conn net.Conn, protocolo int) ([]byte, error) {
	if protocolo == 1 {
		return readCommunication(conn)
	} else {
		return readCommunicationWithHeader(conn)
	}

}

func handleRequest(conn net.Conn) {
	myid := len(clinetConnList)
	clinetConnList = append(clinetConnList, conn)

	serverLog.PrintWaitingMsg("Esperando o tipo de protocolo do cliente de número " + strconv.Itoa(myid) + "...")

	protocolo := readProtocolType(conn)

	serverLog.PrintServerMsg("Tipo de protocolo recebido: "+strconv.Itoa(protocolo), false)

	for {

		serverLog.PrintWaitingMsg("Esperando mensagem do cliente...")
		buf, err := readFromClient(conn, protocolo)
		if err != nil || string(buf) == "-1" {

			serverLog.PrintErrorMsg("Conexão com o cliente " + strconv.Itoa(myid) + " encerrada")
			break
		}

		serverLog.PrintWaitingMsg("Executando lógica do servidor...")
		notify(buf, myid, protocolo)
		wg.Wait()
	}

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

	serverLog.PrintServerMsg("Tamanho da mensagem a ser enviada: "+tamMsg, false)
	conn.Write([]byte(tamMsg))

	serverLog.PrintWaitingMsg("Enviando mensagem:\n\n" + truncateMsgToPrint(msg) + "\n\n...")
	conn.Write(msg)

	serverLog.PrintServerMsg("Mensagem enviada!", false)
}

func SendToClientWithHeader(msg []byte, clinetId int) {
	conn := clinetConnList[clinetId]

	tamMsg := strconv.Itoa(len(msg))

	tamHeader := strconv.Itoa(len(tamMsg))
	tamHeader = "0" + tamHeader
	tamHeader = tamHeader[len(tamHeader)-2:]

	serverLog.PrintServerMsg("Tamanho do cabeçalho mensagem a ser enviado: "+tamHeader, false)

	conn.Write([]byte(tamHeader))

	serverLog.PrintServerMsg("Cabeçalho a ser enviado: "+tamMsg, false)
	conn.Write([]byte(tamMsg))

	serverLog.PrintWaitingMsg("Enviando mensagem: \n\n" + truncateMsgToPrint(msg) + "\n\n...")
	conn.Write(msg)

	serverLog.PrintServerMsg("Mensagem enviada!", false)
}

func SendToClient(msg []byte, clinetId int, protocolo int) {
	defer wg.Done()

	serverLog.PrintWaitingMsg("Enviando mensagem ao cliente de número " + strconv.Itoa(clinetId) + " pelo protocolo " + strconv.Itoa(protocolo) + "...")
	if protocolo == 1 {
		SendToClientWithoutHeader(msg, clinetId)
	} else {
		SendToClientWithHeader(msg, clinetId)
	}

	return
}
