package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	gosnipe "github.com/MCGoSnipe/Runtime"
	flag "github.com/spf13/pflag"
)

var (
	autooffset int
	speedlimit int
	msa        bool
	offset     float64
	name       string
	path       string
	bearer     string
	snipereqs  int
)

type msaRes struct {
	AccessToken *string `json:"access_token"`
	MSAError    *string `json:"error"`
}

// https://stackoverflow.com/a/54747682
func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func initFlags() {
	flag.Float64VarP(&offset, "offset", "o", 0, "offset in milliseconds before snipe.")
	flag.IntVarP(&speedlimit, "auto-offset", "a", 3, "automatically set offset with X requests.")
	flag.IntVarP(&speedlimit, "speed-limit", "l", 0, "offset between requests.")
	flag.IntVarP(&snipereqs, "requests", "r", 2, "number of requests.")
	flag.BoolVarP(&msa, "microsoft", "m", false, "load a microsoft account.")
	flag.StringVarP(&bearer, "bearer", "b", "", "load a microsoft account with this response. requires -m")
	flag.StringVarP(&name, "name", "n", "", "name to snipe.")
	flag.StringVarP(&path, "path", "p", "", "path to accounts text file.")

}

func timeSnipe(timestampTemp *time.Time, lines []string) {
	timestamp := *timestampTemp
	ch := make(chan gosnipe.SnipeRes)
	time.Sleep(time.Until(timestamp.Add(time.Second * time.Duration(20))))
	bearers, labels, _ := gosnipe.SliceStrToBearers(lines)
	for i, bearer := range bearers {
		config := gosnipe.Configuration{
			Bearer:    bearer,
			Name:      name,
			Offset:    offset + float64(speedlimit*i),
			Timestamp: timestamp,
			Label:     &labels[i],
		}
		for j := 0; j < snipereqs; j++ {
			go gosnipe.Snipe(config, ch)
			go getResp(ch)
		}
	}
}

func getResp(ch chan gosnipe.SnipeRes) {
	snipeRes := <-ch
	if snipeRes.Status != nil {
		fmt.Println(*snipeRes.Label + "Status: " + strconv.Itoa(*snipeRes.Status) + " Sent:" + snipeRes.Sent.Format(time.RFC3339) + " Recv: " + snipeRes.Recv.Format(time.RFC3339))
	}
}

func main() {
	fmt.Println( //frick you linter, ruining how this should look
		" ██████╗  ██████╗ ███████╗███╗   ██╗██╗██████╗ ███████╗\n" +
			"██╔════╝ ██╔═══██╗██╔════╝████╗  ██║██║██╔══██╗██╔════╝\n" +
			"██║  ███╗██║   ██║███████╗██╔██╗ ██║██║██████╔╝█████╗  \n" +
			"██║   ██║██║   ██║╚════██║██║╚██╗██║██║██╔═══╝ ██╔══╝  \n" +
			"╚██████╔╝╚██████╔╝███████║██║ ╚████║██║██║     ███████╗\n" +
			" ╚═════╝  ╚═════╝ ╚══════╝╚═╝  ╚═══╝╚═╝╚═╝     ╚══════╝")
	initFlags()
	flag.Parse()
	if !isFlagPassed("name") || name == "" {
		fmt.Println("No name specified. Exiting.")
		os.Exit(1)
	}
	if !isFlagPassed("path") && !isFlagPassed("microsoft") {
		fmt.Println("No accounts file was loaded and no MS accounts were loaded. Exiting.")
		os.Exit(1)
	}
	read := bufio.NewReader(os.Stdin)
	if !isFlagPassed("offset") {
		if isFlagPassed("auto-offset") {
			offsetTemp := gosnipe.AutoOffset(autooffset)
			if offsetTemp == nil {
				offset = 0
			} else {
				offset = *offsetTemp
			}
		}
	}
	bearers := make([]string, 0)
	labels := make([]string, 0)
	var lines []string
	if isFlagPassed("path") {
		lines, _ = gosnipe.TextToSliceStr(path)
	}
	if isFlagPassed("microsoft") {
		if msa {
			fmt.Println("Head to the link below, authorize the app, and paste the page shown afterwards here.")
			fmt.Println(gosnipe.MicrosoftLoginAPI)
			if !isFlagPassed("bearer") {
				res, err := read.ReadString('\n')
				if err != nil {
					fmt.Println("Failed to read from STDIN.")
					os.Exit(1)
				}
				var resp msaRes
				json.Unmarshal([]byte(res), &resp)
				bearers = append(bearers, *resp.AccessToken)
				labels = append(labels, "Microsoft Account")
			} else {
				var resp msaRes
				json.Unmarshal([]byte(bearer), &resp)
				bearers = append(bearers, *resp.AccessToken)
				labels = append(labels, "Microsoft Account")
			}
		} else {
			if !isFlagPassed("path") {
				fmt.Println("No accounts file was loaded and no MS accounts were loaded. Exiting.")
				os.Exit(1)
			}
		}
	}
	timestampTemp := gosnipe.GetDropTime(name)
	if timestampTemp == nil {
		fmt.Println("Failed to fetch droptime.")
		os.Exit(1)
	}
	go timeSnipe(timestampTemp, lines)
	fmt.Println("Snipe running. Press enter to stop.")
	read.ReadString('\n')
}
