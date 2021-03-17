package main

import "fmt"

func switchCommand(cmd string) {
	switch cmd {
	case "snipe":
		info(fmt.Sprintf("attempting snipe with speedlimit of %v\n", speedlimit))
	default:
		err(fmt.Sprintf("command \"%s\" is not a valid command.\n", cmd))
	}
}
