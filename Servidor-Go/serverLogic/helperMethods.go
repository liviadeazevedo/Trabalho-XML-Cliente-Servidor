package serverLogic

func checkParametersNumber(list map[string]string, num_parms int) bool {
	if len(list) != num_parms {
		return false
	}else{
		return true
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