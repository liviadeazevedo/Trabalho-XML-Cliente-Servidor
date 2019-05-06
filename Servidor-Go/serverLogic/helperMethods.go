package serverLogic

import (
	"fmt"
	"strings"
)

const (
	MSG_INVALID_PARAMETERS string = "Invalid parameters received. Missing parameters: "
)

func checkError(err error, print_error bool) bool {
	
	if err != nil {
			if print_error{
				fmt.Println(err)
			}
			return true
		}
		return false
}

func checkParametersNumber(list map[string]string, num_parms int) bool {
	
	if len(list) != num_parms {
		return false
	}else{
		return true
	}
}

func haveAllParameters(list_names []string, parms map[string]string) (string,bool) {
	
	var missing_parms []string
	var missing bool = false

	for _,name := range list_names {
		if parms[name] == ""{
			missing_parms = append(missing_parms,name)
			missing = true
		}
	}

	if missing {
		return MSG_INVALID_PARAMETERS + strings.Join(missing_parms, ", "),false
	}else{
		return "",true
	}
}

func validsCodes(c int) bool {
	
	CODES := []int{0, 1, 2, 3, 4}
	
    for _, v_code := range CODES {
        if v_code == c {
            return true
        }
    }
    return false
}