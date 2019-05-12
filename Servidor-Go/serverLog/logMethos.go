package serverLog

import "fmt"

func PrintServerMsg(msg string, simpleMsg bool) {
	if !simpleMsg {
		fmt.Println()
		fmt.Println("---------------------------------------------------------------------------")
		fmt.Println(msg)
		fmt.Println("---------------------------------------------------------------------------")
	} else {
		fmt.Println(msg)
	}
}

func PrintServerMsgWithTitle(title string, msg string) {
	fmt.Println()
	fmt.Println("===========================================================================")
	fmt.Println(title)
	fmt.Println("===========================================================================")
	fmt.Println()
	fmt.Println(msg)
	fmt.Println("---------------------------------------------------------------------------")
}

func PrintServerMsgOnlyTitle(title string) {
	fmt.Println()
	fmt.Println("===========================================================================")
	fmt.Println(title)
	fmt.Println("===========================================================================")
}

func PrintErrorMsg(err string) {
	fmt.Println()
	fmt.Println("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
	fmt.Println(err)
	fmt.Println("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
}

func PrintWaitingMsg(msg string) {
	fmt.Println()
	fmt.Println("###########################################################################")
	fmt.Println(msg)
	fmt.Println("###########################################################################")
}
