package main

import (
	"Trabalho-XML-Cliente-Servidor/Servidor-Go/serverConnection"
	"Trabalho-XML-Cliente-Servidor/Servidor-Go/serverLog"
	"Trabalho-XML-Cliente-Servidor/Servidor-Go/serverLogic"
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const CONFIG_FILE string = "serverConfig.txt"

//Função auxiliar do método "defineServerIPAndPort".
func createConfigPortIpFile() (ip string, port string, e error) {
	reader := bufio.NewReader(os.Stdin)

	config_path_file_os := filepath.FromSlash(CONFIG_FILE)

	new_file, err := os.Create(config_path_file_os)
	if err != nil {
		e = err
		return
	}
	defer new_file.Close()

	serverLog.PrintServerMsg("Digite o IP da máquina do Servidor:", false)
	aux_ip, _ := reader.ReadString('\n')
	ip = strings.Replace(aux_ip, "\n", "", -1)
	serverLog.PrintServerMsg("Digite a porta da aplicação do Servidor:", false)
	aux_port, _ := reader.ReadString('\n')
	port = strings.Replace(aux_port, "\n", "", -1)

	_, err = new_file.WriteString(aux_ip)
	if err != nil {
		e = err
		return
	}
	_, err = new_file.WriteString(aux_port)
	if err != nil {
		e = err
		return
	}

	//Impor as alterações de escrita no novo arquivo no disco.
	new_file.Sync()

	serverLog.PrintServerMsg("IP e Porta configurados com sucesso!", false)
	e = nil
	return
}

func defineServerIPAndPort() (ip string, port string, e error) {
	var config_exists bool

	serverLog.PrintServerMsgOnlyTitle("Configuração do IP e Porta do Servidor")

	config_path_file_os := filepath.FromSlash(CONFIG_FILE)

	//Verificar se o arquivo já existe
	_, err := os.Stat(config_path_file_os)
	if err == nil {
		config_exists = true
	}

	if config_exists {
		//Abrir o arquivo de configuração.
		config_file, err := os.Open(config_path_file_os)
		if err != nil {
			e = err
			return
		}

		ip_port_bytes, err := ioutil.ReadAll(config_file)
		if err != nil {
			e = err
			return
		}

		config_file.Close()

		ip_port_str := string(ip_port_bytes)
		ip_port_split_array := strings.Split(ip_port_str, "\n")

		reader := bufio.NewReader(os.Stdin)

		for {
			serverLog.PrintServerMsg("Já existe uma configuração com: ip = "+ip_port_split_array[0]+" e porta = "+ip_port_split_array[1]+". Deseja continuar com o mesmo?\n(Sim ou Nao)", false)
			resp, _ := reader.ReadString('\n')
			resp = strings.ToLower(resp)
			resp = strings.Replace(resp, "\n", "", -1)

			if resp == "sim" {
				ip = ip_port_split_array[0]
				port = ip_port_split_array[1]
				e = nil
				return

			} else if resp == "nao" {
				err := os.Remove(config_path_file_os)
				if err != nil {
					e = err
					return
				}

				ip, port, e = createConfigPortIpFile()
				return
			} else {
				serverLog.PrintServerMsg("Opção inválida. Digite \"sim\" ou \"nao\".", false)
			}
		}

	} else {
		ip, port, e = createConfigPortIpFile()
		return
	}
}

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
	ip, port, e := defineServerIPAndPort()
	if e != nil {
		serverLog.PrintServerMsg("Falha ao definir o ip e a porta por input. O padrão será utilizado")
		ip, port = "", ""
	}
	serverLog.PrintWaitingMsg("Se registrando no observer...")
	serverConnection.RegisterObserver(recieveNotification)
	serverLog.PrintServerMsg("Registrado no observer!", false)
	serverLog.PrintWaitingMsg("Abrindo servidor...")
	serverConnection.OpenListener(ip, port)
}
