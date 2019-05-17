package serverLogic

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"Trabalho-XML-Cliente-Servidor/Servidor-Go/serverLog"

	//Manipulação de xml
	"github.com/lestrrat-go/libxml2"
	"github.com/lestrrat-go/libxml2/xpath"
	"github.com/lestrrat-go/libxml2/xsd"
)

const (
	XML_FORMAT_PATH               string = FILES_SOURCE_PATH + "format-response.xml"
	XPATH_METHODS_PARAMETERS      string = "/requisicao/metodo/parametros/parametro[0]"
	XML_FORMAT_FAILED_MSG         string = "XML está mal formatado..."
	XML_FORMAT_SUCCEDED_MSG       string = "XML está bem formatado!"
	XML_VALIDATION_FAILED_MSG     string = "XML não é válido..."
	XML_VALIDATION_SUCCEDED_MSG   string = "XML é válido!"
	MSG_FAILED_BUILD_XML_RESPONSE string = "O sistema falhou em construir um XML de resposta para o Cliente."
)

// Constroí a string do xml correspondente para a resposta ao Cliente.
func buildXMLResponse(value string) string {

	xml_format, err := os.Open(XML_FORMAT_PATH)
	error_sys := checkError(err, true)
	if error_sys {
		return MSG_FAILED_BUILD_XML_RESPONSE
	}
	defer xml_format.Close()

	buffer_xml, err := ioutil.ReadAll(xml_format)
	error_sys = checkError(err, true)
	if error_sys {
		return MSG_FAILED_BUILD_XML_RESPONSE
	}

	response_format := string(buffer_xml)
	response_str := strings.Replace(response_format, "{}", value, 1)

	return response_str
}

//Verificar sintaxe do xml - bem ou mal formatado.
func checkXML(xml string) bool {

	_, err := libxml2.ParseString(xml)
	poorly_formatted := checkError(err, true)

	if poorly_formatted {
		serverLog.PrintErrorMsg(XML_FORMAT_FAILED_MSG)
	} else {
		serverLog.PrintServerMsg(XML_FORMAT_SUCCEDED_MSG, false)
	}

	return !poorly_formatted
}

//Retorna o valor em string do element selecionado dado XPath.
func extractParameterValue(xml, xpath_str string) (string, bool) {

	var value string

	doc, err := libxml2.ParseString(xml)
	error_sys := checkError(err, true)

	if error_sys {
		return "", true
	}

	root_xml, err := doc.DocumentElement()
	error_sys = checkError(err, true)
	if error_sys {
		return "", true
	}

	ctx, err := xpath.NewContext(root_xml)
	error_sys = checkError(err, true)
	if error_sys {
		return "", true
	}
	defer ctx.Free()

	value = xpath.String(ctx.Find(xpath_str))

	return value, false
}

//Retorna os valores dos parâmetros recebidos pelo xml.
func extractParametersValues(xml string) (map[string]string, bool) {

	var (
		aux_name         string
		aux_value        string
		aux_format_xpath string = XPATH_METHODS_PARAMETERS
		error_sys        bool
	)

	params := make(map[string]string)

	for i := 1; ; i++ {
		aux_format_xpath = strings.Replace(aux_format_xpath, strconv.Itoa(i-1), strconv.Itoa(i), 1)

		aux_name, error_sys = extractParameterValue(xml, aux_format_xpath+"/nome")

		if aux_name == "" && !error_sys {
			break
		}

		if error_sys {
			return nil, true
		}

		aux_value, error_sys = extractParameterValue(xml, aux_format_xpath+"/valor")
		if error_sys {
			return nil, true
		}

		params[strings.ToLower(aux_name)] = aux_value
	}

	return params, false
}

//Valida um xml dado um xsd.
func validateXML(xml, xsd_path string) (bool, bool) {

	var error_concat string

	xsd_path_os := filepath.FromSlash(xsd_path)

	schema, err := xsd.ParseFromFile(xsd_path_os)
	error_sys := checkError(err, true)
	if error_sys {
		return false, true
	}

	defer schema.Free()

	doc, err := libxml2.ParseString(xml)
	error_sys = checkError(err, true)
	if error_sys {
		return false, true
	}

	if err := schema.Validate(doc); err != nil {
		for _, e := range err.(xsd.SchemaValidationError).Errors() {
			error_concat = error_concat + "Error: " + e.Error() + "\n"
		}
		serverLog.PrintServerMsgWithTitle("Erros de Validação. XSD: "+xsd_path, error_concat)
		serverLog.PrintErrorMsg(XML_VALIDATION_FAILED_MSG)
		return false, false
	}

	serverLog.PrintServerMsg(XML_VALIDATION_SUCCEDED_MSG, false)

	return true, false
}
