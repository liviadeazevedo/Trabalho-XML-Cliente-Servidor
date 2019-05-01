package serverLogic

import (
	"io/ioutil"
	"fmt"
	"os"
	//"path/filepath"
	"strings"

	//Manipulação de xml
	"github.com/lestrrat-go/libxml2"
	//"github.com/lestrrat-go/libxml2/parser"
	//"github.com/lestrrat-go/libxml2/types"
	//"github.com/lestrrat-go/libxml2/xpath"
	"github.com/lestrrat-go/libxml2/xsd"

	
)

const XML_FORMAT_PATH = FILES_SOURCE_PATH + "format-response.xml"

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

//Obs: Uso do XPath sujeito a alteração
func extractParameterValue(xml,xpath string) string {

	var value string

	//...(code)...

	return value
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