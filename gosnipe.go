//
//
// Used for getting legacy code only, will be deleted in a later commit.
//
//
// package main

// import (
// 	"bufio"
// 	"encoding/json"
// 	"fmt"
// 	"os"
// 	"strconv"
// 	"time"

// 	flag "github.com/spf13/pflag"
// )

// func timeSnipeQueue(ch chan snipeRes, timestamp time.Time, lines []string, name string) {
// 	time.Sleep(time.Until(timestamp.Add(time.Second * -20)))
// 	bearers, labels, _ := sliceStrToBearers(lines)
// 	i := 0
// 	if isFlagPassed("microsoft") {
// 		i++
// 	}
// 	for _, bearer2 := range bearers {
// 		for j := 0; j < snipereqs; j++ {
// 			config := configuration{
// 				Bearer:    bearer2,
// 				Name:      name,
// 				Offset:    offset - float64(speedlimit*i*snipereqs+speedlimit*j),
// 				Timestamp: timestamp,
// 				Label:     &labels[i],
// 			}
// 			go snipe(config, ch)
// 		}
// 		i++
// 	}
// }

// func timeSnipe(ch chan snipeRes, timestamp time.Time, lines []string) {
// 	time.Sleep(time.Until(timestamp.Add(time.Second * -20)))
// 	bearers, labels, _ := sliceStrToBearers(lines)
// 	i := 0
// 	if isFlagPassed("microsoft") {
// 		i++
// 	}
// 	for _, bearer2 := range bearers {
// 		for j := 0; j < snipereqs; j++ {
// 			config := configuration{
// 				Bearer:    bearer2,
// 				Name:      name,
// 				Offset:    offset - float64(speedlimit*i*snipereqs+speedlimit*j),
// 				Timestamp: timestamp,
// 				Label:     &labels[i],
// 			}
// 			go snipe(config, ch)
// 		}
// 		i++
// 	}
// }

// func getResp(ch chan snipeRes) {
// 	for true {
// 		snipeRes := <-ch
// 		if snipeRes.Status != nil {
// 			fmt.Println(*snipeRes.Label + " Status: " + strconv.Itoa(*snipeRes.Status) + " Sent:" + snipeRes.Sent.Format("2006/01/02 15:04:05.0000000") + " Recv: " + snipeRes.Recv.Format("2006/01/02 15:04:05.0000000"))
// 		} else {
// 			fmt.Println("Status was nil.")
// 		}
// 	}
// }

// func main() {
// 	fmt.Println(
// 		FGYellow +
// 			" ██████╗  ██████╗ ███████╗███╗   ██╗██╗██████╗ ███████╗\n" +
// 			"██╔════╝ ██╔═══██╗██╔════╝████╗  ██║██║██╔══██╗██╔════╝\n" +
// 			"██║  ███╗██║   ██║███████╗██╔██╗ ██║██║██████╔╝█████╗  \n" +
// 			"██║   ██║██║   ██║╚════██║██║╚██╗██║██║██╔═══╝ ██╔══╝  \n" +
// 			"╚██████╔╝╚██████╔╝███████║██║ ╚████║██║██║     ███████╗\n" +
// 			" ╚═════╝  ╚═════╝ ╚══════╝╚═╝  ╚═══╝╚═╝╚═╝     ╚══════╝" + Reset)
// 	initFlags()
// 	flag.Parse()
// 	if !isFlagPassed("name") || name == "" {
// 		if !isFlagPassed("queue") {
// 			fmt.Println("No name specified. Exiting.")
// 			os.Exit(1)
// 		}
// 	}
// 	if !isFlagPassed("path") && !isFlagPassed("microsoft") {
// 		fmt.Println("No accounts file was loaded and no MS accounts were loaded. Exiting.")
// 		os.Exit(1)
// 	}
// 	read := bufio.NewReader(os.Stdin)
// 	if !isFlagPassed("offset") {
// 		if isFlagPassed("auto-offset") {
// 			offsetTemp := autoOffset(autooffset)
// 			if offsetTemp == nil {
// 				offset = 0
// 			} else {
// 				offset = *offsetTemp
// 				fmt.Println("Auto-offset used: " + strconv.FormatFloat(offset, 'f', 4, 64))
// 			}
// 		}
// 	}
// 	var lines []string
// 	if isFlagPassed("path") {
// 		lines, _ = textToSliceStr(path)
// 	}

// 	ch := make(chan snipeRes)
// 	var timestamp time.Time
// 	if !isFlagPassed("queue") {
// 		timestampTemp := getDropTime(name)
// 		if timestampTemp == nil {
// 			fmt.Println("Failed to fetch droptime.")
// 			os.Exit(1)
// 		}
// 		timestamp = *timestampTemp
// 	}
// 	var resp msaRes
// 	if isFlagPassed("microsoft") {
// 		if msa {
// 			label := "Microsoft Account"
// 			if !isFlagPassed("bearer") {
// 				fmt.Println("Head to the link below, authorize the app, and paste the page shown afterwards here.")
// 				fmt.Println(microsoftLoginAPI)
// 				res, err := read.ReadString('\n')
// 				if err != nil {
// 					fmt.Println("Failed to read from STDIN.")
// 					os.Exit(1)
// 				}
// 				json.Unmarshal([]byte(res), &resp)
// 				if resp.AccessToken == nil {
// 					fmt.Println("Failed to authenticate, exiting...")
// 					os.Exit(1)
// 				}
// 				if !isFlagPassed("queue") {
// 					for j := 0; j < snipereqs; j++ {

// 						config := configuration{
// 							Bearer:    *resp.AccessToken,
// 							Name:      name,
// 							Offset:    offset + float64(speedlimit*j),
// 							Timestamp: timestamp,
// 							Label:     &label,
// 						}
// 						go snipe(config, ch)
// 					}
// 				}
// 			} else {
// 				json.Unmarshal([]byte(bearer), &resp)
// 				if resp.AccessToken == nil {
// 					fmt.Println("Failed to authenticate, exiting...")
// 					os.Exit(1)
// 				}
// 				if !isFlagPassed("queue") {
// 					for j := 0; j < snipereqs; j++ {
// 						config := configuration{
// 							Bearer:    *resp.AccessToken,
// 							Name:      name,
// 							Offset:    offset + float64(speedlimit*j),
// 							Timestamp: timestamp,
// 							Label:     &label,
// 						}
// 						go snipe(config, ch)
// 					}
// 				}
// 			}
// 		} else {
// 			if !isFlagPassed("path") {
// 				fmt.Println("No accounts file was loaded and no MS accounts were loaded. Exiting.")
// 				os.Exit(1)
// 			}
// 		}
// 	}
// 	if !isFlagPassed("queue") {
// 		fmt.Println("Snipe running. Press enter to close.")
// 		go getResp(ch)
// 		go timeSnipe(ch, timestamp, lines)
// 		read.ReadString('\n')
// 		os.Exit(0)
// 	} else {
// 		go getResp(ch)
// 		queueFunc(read, resp, lines, ch)
// 	}
// }

// func queueFunc(reader *bufio.Reader, msaResp msaRes, lines []string, ch chan snipeRes) {
// 	name := " "
// 	for name != "" {
// 		fmt.Print("Enter name to queue or leave blank to finalize: ")
// 		name, err := reader.ReadString('\n')
// 		if err != nil {
// 			fmt.Println("Error occurred reading name, finalizing...")
// 			break
// 		}
// 		namebytes := []byte(name)
// 		if namebytes[len(namebytes)-1] == '\n' {
// 			namebytes = namebytes[:len(namebytes)-1]
// 		}
// 		name = string(namebytes)
// 		if name == "" {
// 			break
// 		}
// 		if namebytes[len(namebytes)-1] == '\r' {
// 			namebytes = namebytes[:len(namebytes)-1]
// 		}
// 		name = string(namebytes)
// 		if name == "" {
// 			break
// 		}
// 		timestampTemp := getDropTime(name)
// 		if timestampTemp == nil {
// 			fmt.Println("Failed to fetch droptime for " + name + ".")
// 			continue
// 		}
// 		timestamp := *timestampTemp
// 		go timeSnipeQueue(ch, timestamp, lines, name)
// 		if msaResp.AccessToken != nil {
// 			label := "Microsoft Account"
// 			for j := 0; j < snipereqs; j++ {

// 				config := configuration{
// 					Bearer:    *msaResp.AccessToken,
// 					Name:      name,
// 					Offset:    offset + float64(speedlimit*j),
// 					Timestamp: timestamp,
// 					Label:     &label,
// 				}
// 				go snipe(config, ch)
// 			}
// 		}
// 		fmt.Println("Snipe started for " + name + ".")
// 	}
// 	fmt.Println("Finalized. Press enter to quit.")
// 	reader.ReadString('\n')
// 	os.Exit(0)
// }
