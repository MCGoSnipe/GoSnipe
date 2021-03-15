package main

import (
	"os"
	"runtime"
)

func configExists() bool {
	appdata := getOSAppDataLocation() + "/gosnipe"
	if _, err := os.Stat(appdata); os.IsNotExist(err) {
		return false
	}
	return true
}

func initializeConfig() {

}

func wipeConfig() {

}

func getOSAppDataLocation() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("LOCALAPPDATA")
	} else {
		return os.Getenv("$HOME") + "/.config"
	}
}
