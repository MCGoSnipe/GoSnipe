package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/eiannone/keyboard"
)

func err(v ...interface{}) {
	fmt.Printf("[%vERROR%v] ", FGRed, Reset)
	for _, h := range v {
		fmt.Printf("%v ", h)
	}
	fmt.Print("\n")
}

func info(v ...interface{}) {
	fmt.Printf("[%vINFO%v] ", FGCyan, Reset)
	for _, h := range v {
		fmt.Printf("%v ", h)
	}
	fmt.Print("\n")
}

func snipeInfo(v snipeRes) {
	if v.Status == nil {
		err("No response was detected when trying to print status.")
		return
	}
	statusColor := (map[bool]string{true: FGGreen, false: FGRed})[*v.Status == 200]
	fmt.Printf("[%v%v%v] [%v%v%v] >> Sent at %v - Recv at %v\n", FGCyan, *v.Label, Reset, statusColor, *v.Status, Reset, v.Sent.Format("2006/01/02 15:04:05.0000000"), v.Recv.Format("2006/01/02 15:04:05.0000000"))
}

func yesNo(prompt string) bool {
	fmt.Printf("[%vY/N%v] %v\n", FGYellow, Reset, prompt)
	fmt.Println("Use arrow keys to move, enter to select.")
	fmt.Print(" * Yes      No\r * Yes")
	out := true
	keyboard.Open()
	defer keyboard.Close()
	// r/badcode moment
	for {
		_, code, _ := keyboard.GetSingleKey()
		if code == keyboard.KeyArrowDown || code == keyboard.KeyArrowUp || code == keyboard.KeyArrowLeft || code == keyboard.KeyArrowRight {
			out = !out
		} else if code == keyboard.KeyEnter {
			break
		}
		if out {
			fmt.Print("\r * Yes      No\r * Yes")
		} else {
			fmt.Print("\r   Yes    * No")
		}
	}
	fmt.Print("\n")
	return out
}

func inputString(prompt string) string {
	fmt.Printf("[%vIN%v] %v\n", FGBlue, Reset, prompt)
	fmt.Println("Press RETURN to finish.")
	fmt.Print(" >> ")
	read := bufio.NewReader(os.Stdin)
	readStr, _ := read.ReadString('\n')
	readSlice := []byte(readStr)
	readSlice = readSlice[:len(readSlice)-2]
	if '\n' == readSlice[len(readSlice)-1] {
		readSlice = readSlice[:len(readSlice)-2]
	}
	return readStr
}

/* uncomment this to test methods
func main() {
	err("Example error.")
	info("Example info.")
	info("You chose", yesNo("Choose an option."))
	info("You said", inputString("Type something."))
	label := "Success Info"
	status := 200
	timest := time.Now()
	snipeInfo(snipeRes{
		Label:  &label,
		Sent:   &timest,
		Recv:   &timest,
		Status: &status,
	})
	label = "Error Info"
	status = 403
	snipeInfo(snipeRes{
		Label:  &label,
		Sent:   &timest,
		Recv:   &timest,
		Status: &status,
	})
}

*/
