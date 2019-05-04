/*

//(...codigo servidor...)

//msg seria em []byte, provavelmente, e gostaria também que aqui só tenha as informações do xml, sem a parte do tamanho.
xml = string(msg)

//"resposta" seria uma string do xml de reposta a ser enviado ao Cliente.
//Ex de string retornada da resposta: 	<resposta>
//    										<retorno>0</retorno>
//										</resposta>
resposta = serverLogic.ResquestXMLHandler(xml) //Interface da lógica do XML, que vou implmentar.

//Converte a string em []byte para poder enviar ao Cliente de volta
msg_enviar = []byte(resposta)

//(...codigo servidor...)

*/

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
)

func methodHandler(xml string, method func(map[string]string)string, num_parms int) (ret_value string) {

	var (
		parms map[string]string
		execute bool	
	)

	parms = extractParametersValues(xml, num_parms)
	execute = checkParametersNumber(parms,num_parms)

	if execute {
		ret_value = method(parms)
		return	
	}else{
		return "Wrong parameters number."
	}
}

func RequestXMLHandler(xml string) string {
	var (
		xml_resp string
		resp string
		method_name string
	)

	method_name = strings.ToLower(extractParameterValue(xml,XPATH_METHOD_NAME))

	switch  method_name {
		case "submeter":
			resp = methodHandler(xml,submeter,1)
		case "consultastatus":
			resp = methodHandler(xml,consultaStatus,1)
	}

	xml_resp = buildXMLResponse(resp)
	return xml_resp
}

/* envia um boletim como parâmetro e retorna um número inteiro (0 - sucesso, 1 - XML inválido, 2 - XML mal-formado, 3 - Erro Interno) */

//func submeter(Boletim string) int
func submeter(parms map[string]string) string {
	
	//...(code)...

	return "0"
}

/* consulta o status da inscrição do candidato com o CPF informado como parâmetro. Possíveis retornos: 0 - Candidato não encontrado, 1 - Em processamento, 
2 - Candidato Aprovado e Selecionado, 3 - Candidato Aprovado e em Espera, 4 - Candidato Não Aprovado. */

//func consultaStatus(cpf string) int
func consultaStatus(parms map[string]string) string {

	cpf_int,_ := strconv.Atoi(parms["cpf"])

	if validsCodes(cpf_int) {
		return strconv.Itoa(cpf_int)
	}else{
		return "0"
	}
}
