package main

import "fmt"

const (
	FGBlack         string = "\033[30;30m"
	FGRed           string = "\033[30;31m"
	FGGreen         string = "\033[30;32m"
	FGYellow        string = "\033[30;33m"
	FGBlue          string = "\033[30;34m"
	FGMagenta       string = "\033[30;35m"
	FGCyan          string = "\033[30;36m"
	FGWhite         string = "\033[30;37m"
	BGBlack         string = "\033[30;40m"
	BGRed           string = "\033[30;41m"
	BGGreen         string = "\033[30;42m"
	BGYellow        string = "\033[30;43m"
	BGBlue          string = "\033[30;44m"
	BGMagenta       string = "\033[30;45m"
	BGCyan          string = "\033[30;46m"
	BGWhite         string = "\033[30;47m"
	BrightFGBlack   string = "\033[30;90m"
	BrightFGRed     string = "\033[30;91m"
	BrightFGGreen   string = "\033[30;92m"
	BrightFGYellow  string = "\033[30;93m"
	BrightFGBlue    string = "\033[30;94m"
	BrightFGMagenta string = "\033[30;95m"
	BrightFGCyan    string = "\033[30;96m"
	BrightFGWhite   string = "\033[30;97m"
	BrightBGBlack   string = "\033[30;100m"
	BrightBGRed     string = "\033[30;101m"
	BrightBGGreen   string = "\033[30;102m"
	BrightBGYellow  string = "\033[30;103m"
	BrightBGBlue    string = "\033[30;104m"
	BrightBGMagenta string = "\033[30;105m"
	BrightBGCyan    string = "\033[30;106m"
	BrightBGWhite   string = "\033[30;107m"
	Reset           string = "\033[0m"
)

func setFgCol(r, g, b int) {
	fmt.Printf("\033[38;2;%d;%d;%dm", r, g, b)
}

func setBgCol(r, g, b int) {
	fmt.Printf("\033[48;2;%d;%d;%dm", r, g, b)
}
