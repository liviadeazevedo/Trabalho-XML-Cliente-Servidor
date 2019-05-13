package serverLogic

import (
	"strconv"
	"strings"

	"Trabalho-XML-Cliente-Servidor/Servidor-Go/serverLog"
)

const (
	FILES_SOURCE_PATH               string = "../Arquivos/"
	XSD_REQUEST_PATH                string = FILES_SOURCE_PATH + "requisicao.xsd"
	XSD_HISTORICO_PATH              string = FILES_SOURCE_PATH + "historico.xsd"
	XPATH_METHOD_NAME               string = "/requisicao/metodo/nome"
	MSG_NUMBER_PARAMETERS_WRONG     string = "Número de parâmetros inválido. Verifique os parâmetros do método considerado."
	MSG_INVALID_METHOD_NAME         string = "Nome do método inválido. Selecione um válido."
	MSG_FAILED_PARSE_XML_RESPONSE   string = "Falha no processamento do XML da requisição. Por favor, verifique este XML."
	MSG_FAILED_EXTRACT_PARMS_VALUES string = "Falha em extrair os valores dos parâmetros no XML da requisição."
	MSG_INVALID_XSD                 string = "XML da requisição é inválido. Por favor, verifique este XML."
)

func methodHandler(xml string, method func(map[string]string) string, num_parms int) (ret_value string) {

	var (
		parms     map[string]string
		execute   bool
		error_sys bool
	)

	parms, error_sys = extractParametersValues(xml)

	if error_sys {
		return MSG_FAILED_EXTRACT_PARMS_VALUES
	}

	execute = checkParametersNumber(parms, num_parms)

	if execute {
		ret_value = method(parms)
		return
	} else {
		return MSG_NUMBER_PARAMETERS_WRONG
	}
}

func RequestXMLHandler(xml string) string {
	var (
		xml_resp    string
		resp        string
		method_name string
		error_sys   bool
		valid_req   bool
	)

	serverLog.PrintWaitingMsg("Verificando validade do XML de requisição...")

	valid_req, error_sys = validateXML(xml, XSD_REQUEST_PATH)
	if valid_req {
		resp = MSG_INVALID_XSD
	}

	serverLog.PrintWaitingMsg("Extraindo parêmetros da requisição...")

	method_name, _ = extractParameterValue(xml, XPATH_METHOD_NAME)

	method_name = strings.ToLower(method_name)

	serverLog.PrintWaitingMsg("Executando método solicitado pelo Cliente...")

	switch method_name {
	case "submeter":
		resp = methodHandler(xml, submeter, 1)
	case "consultastatus":
		resp = methodHandler(xml, consultaStatus, 1)
	default:
		if !error_sys {
			serverLog.PrintErrorMsg("Falha na execução do método! Nome inválido!")
			resp = MSG_INVALID_METHOD_NAME
		}
	}

	serverLog.PrintWaitingMsg("Construindo XML de resposta...")

	xml_resp = buildXMLResponse(resp)
	return xml_resp
}

/* envia um boletim como parâmetro e retorna um número inteiro (0 - sucesso, 1 - XML inválido, 2 - XML mal-formado, 3 - Erro Interno) */

func submeter(parms map[string]string) string {

	parms_names := []string{"boletim"}

	msg_error, valid_parms := haveAllParameters(parms_names, parms)
	if !valid_parms {
		return msg_error
	}

	correct_formated := checkXML(parms[parms_names[0]])
	if !correct_formated {
		return "2"
	}

	valid_xml, error_sys := validateXML(parms[parms_names[0]], XSD_HISTORICO_PATH)

	if error_sys {
		return "3"
	}

	if !valid_xml {
		return "1"
	}

	return "0"
}

/* consulta o status da inscrição do candidato com o CPF informado como parâmetro. Possíveis retornos: 0 - Candidato não encontrado, 1 - Em processamento,
2 - Candidato Aprovado e Selecionado, 3 - Candidato Aprovado e em Espera, 4 - Candidato Não Aprovado. */

func consultaStatus(parms map[string]string) string {

	parms_names := []string{"cpf"}

	msg_error, valid_parms := haveAllParameters(parms_names, parms)
	if !valid_parms {
		return msg_error
	}

	msg_error, is_number := checkIsANumber(parms[parms_names[0]])
	if !is_number {
		return msg_error
	}

	cpf_int, _ := strconv.Atoi(parms[parms_names[0]])

	if validsCodes(cpf_int) {
		return strconv.Itoa(cpf_int)
	} else {
		return "0"
	}
}
