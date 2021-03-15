package main

import (
	"fmt"
	"time"

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
	queue      bool
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

func main() {
	fmt.Println(
		FGYellow +`
 ██████╗  ██████╗ ███████╗███╗   ██╗██╗██████╗ ███████╗ 
██╔════╝ ██╔═══██╗██╔════╝████╗  ██║██║██╔══██╗██╔════╝ 
██║  ███╗██║   ██║███████╗██╔██╗ ██║██║██████╔╝█████╗   
██║   ██║██║   ██║╚════██║██║╚██╗██║██║██╔═══╝ ██╔══╝   
╚██████╔╝╚██████╔╝███████║██║ ╚████║██║██║     ███████╗ 
 ╚═════╝  ╚═════╝ ╚══════╝╚═╝  ╚═══╝╚═╝╚═╝     ╚══════╝` +  Reset)
	initFlags()
	flag.Parse()
	//args := flag.Args()
	//TODO: add in frontend code
}
func initFlags() {
	flag.Float64VarP(&offset, "offset", "o", 0, "offset in milliseconds before snipe.")
	flag.IntVarP(&speedlimit, "auto-offset", "a", 3, "automatically set offset with X requests.")
	flag.IntVarP(&speedlimit, "speed-limit", "l", 0, "offset between requests.")
	flag.IntVarP(&snipereqs, "requests", "r", 2, "number of requests.")
	flag.BoolVarP(&msa, "microsoft", "m", false, "load a microsoft account.")
	flag.BoolVarP(&queue, "queue", "q", false, "enables STDIN name queueing")
	flag.StringVarP(&bearer, "bearer", "b", "", "load a microsoft account with this response. requires -m")
	flag.StringVarP(&name, "name", "n", "", "name to snipe.")

}

func timeSnipe(ch chan snipeRes, timestamp time.Time, lines []string, name1 ...string) {
	time.Sleep(time.Until(timestamp.Add(time.Second * -20)))
	bearers, labels, _ := sliceStrToBearers(lines)
	i := 0
	if isFlagPassed("microsoft") {
		i++
	}
	for _, bearer2 := range bearers {
		for j := 0; j < snipereqs; j++ {
			config := configuration{
				Bearer:    bearer2,
				Name:      name,
				Offset:    offset - float64(speedlimit*i*snipereqs+speedlimit*j),
				Timestamp: timestamp,
				Label:     &labels[i],
			}
			go snipe(config, ch)
		}
		i++
	}
}
