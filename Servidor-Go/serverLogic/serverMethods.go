package serverLogic

import (
	"strings"
	"strconv"
)

const (
	FILES_SOURCE_PATH string = "../Arquivos/"
	XSD_REQUEST_PATH string = FILES_SOURCE_PATH + "requisicao.xsd"
	XSD_RESPONSE_PATH string = FILES_SOURCE_PATH + "resposta.xsd"
	XSD_HISTORICO_PATH string = FILES_SOURCE_PATH + "historico.xsd"
	XPATH_METHOD_NAME string = "/requisicao/metodo/nome"
	MSG_NUMBER_PARAMETERS_WRONG string = "Wrong parameters number. Check the method parameters." 
	MSG_INVALID_METHOD_NAME string = "Invalid method name. Select a valid one."
	MSG_FAILED_PARSE_XML_RESPONSE string = "Failed to parser request XML. Please, check the request XML."
	MSG_FAILED_EXTRACT_PARMS_VALUES string = "Failed in extract parameters values from XML request."
)

func methodHandler(xml string, method func(map[string]string)string, num_parms int) (ret_value string) {

	var (
		parms map[string]string
		execute bool
		error_sys bool	
	)

	parms, error_sys = extractParametersValues(xml, num_parms)

	if error_sys {
		return MSG_FAILED_EXTRACT_PARMS_VALUES
	}

	execute = checkParametersNumber(parms,num_parms)

	if execute {
		ret_value = method(parms)
		return	
	}else{
		return MSG_NUMBER_PARAMETERS_WRONG
	}
}

func RequestXMLHandler(xml string) string {
	var (
		xml_resp string
		resp string
		method_name string
		error_sys bool
	)

	method_name,error_sys = extractParameterValue(xml,XPATH_METHOD_NAME)

	if error_sys {
		resp = MSG_FAILED_PARSE_XML_RESPONSE
	}

	method_name = strings.ToLower(method_name)

	switch  method_name {
		case "submeter":
			resp = methodHandler(xml,submeter,1)
		case "consultastatus":
			resp = methodHandler(xml,consultaStatus,1)
		default:
			if !error_sys {
				resp = MSG_INVALID_METHOD_NAME	
			}
	}

	xml_resp = buildXMLResponse(resp)
	return xml_resp
}

/* envia um boletim como parâmetro e retorna um número inteiro (0 - sucesso, 1 - XML inválido, 2 - XML mal-formado, 3 - Erro Interno) */

//func submeter(Boletim string) int
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

	valid_xml, error_sys := validateXML(parms[parms_names[0]],XSD_HISTORICO_PATH)
	
	if error_sys{
		return "3"
	}
	
	if !valid_xml {
		return "1"
	}

	return "0"
}

/* consulta o status da inscrição do candidato com o CPF informado como parâmetro. Possíveis retornos: 0 - Candidato não encontrado, 1 - Em processamento, 
2 - Candidato Aprovado e Selecionado, 3 - Candidato Aprovado e em Espera, 4 - Candidato Não Aprovado. */

//func consultaStatus(cpf string) int
func consultaStatus(parms map[string]string) string {

	parms_names := []string{"cpf"}

	msg_error, valid_parms := haveAllParameters(parms_names, parms)
	if !valid_parms {
		return msg_error
	}

	cpf_int,_ := strconv.Atoi(parms[parms_names[0]])

	if validsCodes(cpf_int) {
		return strconv.Itoa(cpf_int)
	}else{
		return "0"
	}
}
