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
)

// Constroí a string do xml correspondente para a resposta ao Cliente.
func buildXMLResponse(value string) string {

  xml_format, err := os.Open(XML_FORMAT_PATH)
  checkError(err,true)
	defer xml_format.Close()
		
	buffer_xml, err := ioutil.ReadAll(xml_format)
	checkError(err,true)
	
	response_format := string(buffer_xml)
	response_str := strings.Replace(response_format, "{}", value, 1)

	return response_str
}

func checkError(err error, print_error bool) bool {
    if err != nil {
			if print_error{
				fmt.Println(err)
			}
			return true
		}
		return false
}

//Verificar sintaxe do xml - bem ou mal formatado.
func checkXML(xml string) bool {

	_, err := libxml2.ParseString(xml)
	poorly_formatted := checkError(err,false)

	if poorly_formatted {
		fmt.Println("\n===XML is poorly formatted...===")
	}else{
		fmt.Println("===XML is well formatted!===")
	}

	return poorly_formatted
}

//Retorna o valor em string do element selecionado dado XPath.
func extractParameterValue(xml,xpath_str string) string {

	var value string

	doc, err := libxml2.ParseString(xml)
	checkError(err,false)

	root_xml, err := doc.DocumentElement()
	checkError(err,false)

	ctx, err := xpath.NewContext(root_xml)
	checkError(err,false)
	defer ctx.Free()

	value = xpath.String(ctx.Find(xpath_str))
	
	return value
}

//Retorna os valores dos parâmetros recebidos pelo xml.
func extractParametersValues(xml string, qtd_params int) (map[string]string) {

	var (
		aux_name string
		aux_value string
		aux_format_xpath string = XPATH_METHODS_PARAMETERS
	) 

	params := make(map[string]string)

	for i := 1; i <= qtd_params; i++ {
		fmt.Println(i)
		aux_format_xpath = strings.Replace(aux_format_xpath,strconv.Itoa(i-1),strconv.Itoa(i),1)
		fmt.Println(aux_format_xpath)

		aux_name = extractParameterValue(xml,aux_format_xpath + "/nome")
		aux_value = extractParameterValue(xml,aux_format_xpath + "/valor")
		params[aux_name] = aux_value

	}

	return params
}

//Valida um xml dado um xsd.
func validateXML(xml,xsd_path string) bool {

	schema, err := xsd.ParseFromFile(xsd_path)
	checkError(err,false)
	defer schema.Free()
  
	doc, err := libxml2.ParseString(xml)
	checkError(err,false)
  
	if err := schema.Validate(doc); err != nil {
	  for _, e := range err.(xsd.SchemaValidationError).Errors() {
			fmt.Println("Error: ", e.Error())
		}
		fmt.Println("===XML is not valid...===")
	  return false
	}
  
	fmt.Println("===Validation Successful!===")

	return true
}