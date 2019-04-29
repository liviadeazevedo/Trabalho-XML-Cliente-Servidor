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

func checkError(err error) bool {
    if err != nil {
			//fmt.Println(err)
			return true
		}
		return false
}

//Verificar sintaxe do xml - bem ou mal formatado.
func CheckXML(xml string) bool {

	_, err := libxml2.ParseString(xml)
	poorly_formatted := checkError(err)

	if poorly_formatted {
		fmt.Println("\n===XML is poorly formatted...===")
	}else{
		fmt.Println("===XML is well formatted!===")
	}

	return poorly_formatted
}

//Valida um xml dado um xsd.
func ValidateXML(xml,xsd_path string) bool {

	schema, err := xsd.ParseFromFile(xsd_path)
	checkError(err)
	defer schema.Free()
  
	doc, err := libxml2.ParseString(xml)
	checkError(err)
  
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

// Constroí a string do xml correspondente para a resposta ao Cliente.
func BuildXMLResponse(value string) string {
	
	xml_format_path := "Arquivos/format-response.xml"
  xml_format, err := os.Open(xml_format_path)
  checkError(err)
	defer xml_format.Close()
		
	buffer_xml, err := ioutil.ReadAll(xml_format)
	checkError(err)
	
	response_format := string(buffer_xml)
	response_str := strings.Replace(response_format, "{}", value, 1)

	return response_str
}