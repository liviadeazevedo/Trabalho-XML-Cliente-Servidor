package serverLogic

import (
	"strings"

	"../serverLog"
)

const (
	MSG_INVALID_PARAMETERS string = "Parâmetros recebidos inválidos. Parâmetros faltantes: "
	MSG_NOT_A_NUMBER       string = "Parâmetro 'cpf' inválido. Por favor, envie um válido."
)

func checkError(err error, print_error bool) bool {

	if err != nil {
		if print_error {
			serverLog.PrintServerMsg(err.Error(), false)
		}
		return true
	}
	return false
}

func checkIsANumber(txt string) (string, bool) {

	NUMBERS_ZERO_TO_NINE := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	var char_is_a_num bool = false

	for _, char := range txt {

		for _, num := range NUMBERS_ZERO_TO_NINE {
			if string(char) == num {
				char_is_a_num = true
				break
			}
		}

		if !char_is_a_num {
			return MSG_NOT_A_NUMBER, false
		} else {
			char_is_a_num = false
		}
	}

	return "", true
}

func checkParametersNumber(list map[string]string, num_parms int) bool {

	if len(list) != num_parms {
		return false
	} else {
		return true
	}
}

func haveAllParameters(list_names []string, parms map[string]string) (string, bool) {

	var missing_parms []string
	var missing bool = false

	for _, name := range list_names {
		if parms[name] == "" {
			missing_parms = append(missing_parms, name)
			missing = true
		}
	}

	if missing {
		return MSG_INVALID_PARAMETERS + strings.Join(missing_parms, ", "), false
	} else {
		return "", true
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
