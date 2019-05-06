package serverLogic

import (
	"io/ioutil"
	"fmt"
	"os"
	"strings"
	"strconv"

	//Manipulação de xml
	"github.com/lestrrat-go/libxml2"
	"github.com/lestrrat-go/libxml2/xpath"
	"github.com/lestrrat-go/libxml2/xsd"
)

const (
	XML_FORMAT_PATH string = FILES_SOURCE_PATH + "format-response.xml"
	XPATH_METHODS_PARAMETERS string = "/requisicao/metodo/parametros/parametro[0]"
	XML_FORMAT_FAILED_MSG string = "==========XML is poorly formatted...=========="
	XML_FORMAT_SUCCEDED_MSG string = "==========XML is well formatted!=========="
	XML_VALIDATION_FAILED_MSG string = "==========XML is not valid...=========="
	XML_VALIDATION_SUCCEDED_MSG string = "==========XML is valid!=========="
	MSG_FAILED_BUILD_XML_RESPONSE string = "System failed to build XML response to Client."
)

// Constroí a string do xml correspondente para a resposta ao Cliente.
func buildXMLResponse(value string) string {

  xml_format, err := os.Open(XML_FORMAT_PATH)
	error_sys := checkError(err,true)
	if error_sys{
		return MSG_FAILED_BUILD_XML_RESPONSE
	}
	defer xml_format.Close()
		
	buffer_xml, err := ioutil.ReadAll(xml_format)
	error_sys = checkError(err,true)
	if error_sys{
		return MSG_FAILED_BUILD_XML_RESPONSE
	}

	response_format := string(buffer_xml)
	response_str := strings.Replace(response_format, "{}", value, 1)

	return response_str
}

//Verificar sintaxe do xml - bem ou mal formatado.
func checkXML(xml string) bool {

	_, err := libxml2.ParseString(xml)
	poorly_formatted := checkError(err,false)

	if poorly_formatted {
		fmt.Println("\n" + XML_FORMAT_FAILED_MSG)
	}else{
		fmt.Println(XML_FORMAT_SUCCEDED_MSG)
	}

	return !poorly_formatted
}

//Retorna o valor em string do element selecionado dado XPath.
func extractParameterValue(xml,xpath_str string) (string,bool) {

	var value string

	doc, err := libxml2.ParseString(xml)
	error_sys := checkError(err,false)
	if error_sys {
		return "",true
	}

	root_xml, err := doc.DocumentElement()
	error_sys = checkError(err,false)
	if error_sys {
		return "",true
	}

	ctx, err := xpath.NewContext(root_xml)
	error_sys = checkError(err,false)
	if error_sys {
		return "",true
	}
	defer ctx.Free()

	value = xpath.String(ctx.Find(xpath_str))
	
	return value,true
}

//Retorna os valores dos parâmetros recebidos pelo xml.
func extractParametersValues(xml string, qtd_params int) (map[string]string,bool) {

	var (
		aux_name string
		aux_value string
		aux_format_xpath string = XPATH_METHODS_PARAMETERS
		error_sys bool 
	) 

	params := make(map[string]string)

	for i := 1; i <= qtd_params; i++ {
		aux_format_xpath = strings.Replace(aux_format_xpath,strconv.Itoa(i-1),strconv.Itoa(i),1)

		aux_name, error_sys = extractParameterValue(xml,aux_format_xpath + "/nome")
		if error_sys {
			return nil,true
		}

		aux_value, error_sys = extractParameterValue(xml,aux_format_xpath + "/valor")
		if error_sys {
			return nil,true
		}

		params[strings.ToLower(aux_name)] = aux_value
	}

	return params,false
}

//Valida um xml dado um xsd.
func validateXML(xml,xsd_path string) (bool,bool) {

	schema, err := xsd.ParseFromFile(xsd_path)
	error_sys := checkError(err,false)
	if error_sys {
		return false,true
	}

	defer schema.Free()
  
	doc, err := libxml2.ParseString(xml)
	error_sys = checkError(err,false)
	if error_sys {
		return false,true
	}
  
	if err := schema.Validate(doc); err != nil {
	  for _, e := range err.(xsd.SchemaValidationError).Errors() {
			fmt.Println("Error: ", e.Error())
		}
		fmt.Println(XML_VALIDATION_FAILED_MSG)
	  return false,false
	}
  
	fmt.Println(XML_VALIDATION_SUCCEDED_MSG)

	return true,false
}