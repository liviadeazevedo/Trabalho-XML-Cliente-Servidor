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

import  "strings"

const (
	FILES_SOURCE_PATH string = "../Arquivos/"
	XSD_REQUEST_PATH string = FILES_SOURCE_PATH + "requisicao.xsd"
	XSD_RESPONSE_PATH string = FILES_SOURCE_PATH + "resposta.xsd"
	XSD_HISTORICO_PATH string = FILES_SOURCE_PATH + "historico.xsd"
)

func submeterHandler(xml string) int {
	var (
		cod int
		boletim string
	) 

	boletim = extractParameterValue(xml,"<xpath-do-elemento-parametro-historico>")
	cod = submeter(boletim)

	return cod
}

func consultaStatusHandler(xml string) int {
	var (
		cod int
		cpf string
	) 

	cpf = extractParameterValue(xml,"<xpath-do-elemento-parametro-cpf>")
	cod = consultaStatus(cpf)
	
	return cod
}


func RequestXMLHandler(xml string) string {
	var (
		xml_resp string
		cod_int int
		method_name string
	)

	//Obs: Uso do XPath sujeito a alteração
	method_name = strings.ToLower(extractParameterValue(xml,"<xpath-do-elemento-nome>"))

	switch  method_name {
		case "submeter":
			cod_int = submeterHandler(xml)
		case "consultastatus":
			cod_int = consultaStatusHandler(xml)
	}

	xml_resp = buildXMLResponse(string(cod_int))
	return xml_resp
}

/* envia um boletim como parâmetro e retorna um número inteiro (0 - sucesso, 1 - XML inválido, 2 - XML mal-formado, 3 - Erro Interno) */
func submeter(Boletim string) int {
	
	//...(code)...

	return 0
}

/* consulta o status da inscrição do candidato com o CPF informado como parâmetro. Possíveis retornos: 0 - Candidato não encontrado, 1 - Em processamento, 
2 - Candidato Aprovado e Selecionado, 3 - Candidato Aprovado e em Espera, 4 - Candidato Não Aprovado. */
func consultaStatus(cpf string) int {
	
	//...(code)...
	
	return 0
}
