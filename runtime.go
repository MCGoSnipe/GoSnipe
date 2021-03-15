package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	minecraftServicesAPIHost = "api.minecraftservices.com"
	yggdrasilAuthURI         = "https://authserver.mojang.com"
	microsoftLoginAPI        = "https://login.live.com/oauth20_authorize.srf?client_id=9abe16f4-930f-4033-b593-6e934115122f&response_type=code&redirect_uri=https%3A%2F%2Fapi.gosnipe.tech%2Fapi%2Fauthenticate&scope=XboxLive.signin%20XboxLive.offline_access"
)

type configuration struct {
	Bearer    string
	Name      string
	Offset    float64
	Timestamp time.Time
	Label     *string
	Debug     bool
}

type snipeRes struct {
	Sent   *time.Time
	Recv   *time.Time
	Status *int
	Label  *string
}

type dropAPIRes struct {
	DropTime string `json:"time"`
}

type securityRes struct {
	Answer answerRes `json:"answer"`
}

type answerRes struct {
	ID int `json:"id"`
}

type accessTokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type accessTokenResponse struct {
	AccessToken *string `json:"accessToken"`
	YggError    *string `json:"error"`
}

func textToSliceStr(path string) ([]string, int) {
	file, err := os.Open(path)
	i := 0
	if err == nil {
		var txtSlice []string
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "\n") {
				line = line[:len(line)-1]
			}
			if strings.Contains(line, "\r") {
				line = line[:len(line)-1]
			}
			txtSlice = append(txtSlice, scanner.Text())
			i++
		}
		return txtSlice, i
	}
	return make([]string, 0), 0
}

func sliceStrToBearers(inputSlice []string) ([]string, []string, int) {
	outputSlice := make([]string, 0)
	outputSlice2 := make([]string, 0)
	i := 0
	client := &http.Client{}
	for _, input := range inputSlice {
		splitLogin := strings.Split(input, ":")
		if len(splitLogin) < 2 {
			continue
		}
		data := accessTokenRequest{
			Username: splitLogin[0],
			Password: splitLogin[1],
		}
		bytesToSend, err := json.Marshal(data)
		if err != nil {
			continue
		}
		req, err := http.NewRequest("POST", yggdrasilAuthURI+"/authenticate", bytes.NewBuffer(bytesToSend))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "GoSnipe/2.0")
		if err != nil {
			continue
		}
		res, err := client.Do(req)
		if err != nil {
			continue
		}
		if res.Status != "200 OK" {
			continue
		}
		respData, err := ioutil.ReadAll(res.Body)
		if err != nil {
			continue
		}
		var access accessTokenResponse
		err = json.Unmarshal(respData, &access)
		if err != nil || access.AccessToken == nil {
			continue
		}
		if len(splitLogin) != 5 {
			continue
		}
		req, err = http.NewRequest("GET", "https://api.mojang.com/user/security/challenges", nil)
		if err != nil {
			continue
		}
		req.Header.Set("Authorization", "Bearer "+*access.AccessToken)
		res, err = client.Do(req)
		if err != nil {
			continue
		}
		respData, err = ioutil.ReadAll(res.Body)
		if err != nil {
			continue
		}
		var security []securityRes
		err = json.Unmarshal(respData, &security)
		if err != nil {
			continue
		}
		if len(security) != 3 {
			continue
		}
		dataBytes := []byte(`[{"id": ` + strconv.Itoa(security[0].Answer.ID) + `, "answer": "` + splitLogin[2] + `"}, {"id": ` + strconv.Itoa(security[1].Answer.ID) + `, "answer": "` + splitLogin[3] + `"}, {"id": ` + strconv.Itoa(security[2].Answer.ID) + `, "answer": "` + splitLogin[4] + `"}]`)
		req, err = http.NewRequest("POST", "https://api.mojang.com/user/security/location", bytes.NewReader(dataBytes))
		if err != nil {
			continue
		}
		req.Header.Set("Authorization", "Bearer "+*access.AccessToken)
		client.Do(req)
		outputSlice = append(outputSlice, *access.AccessToken)
		outputSlice2 = append(outputSlice2, splitLogin[0])
		i++
	}
	return outputSlice, outputSlice2, i
}

func getDropTime(name string) *time.Time {
	res, err := http.Get("https://api.gosnipe.tech/api/status/name/" + name)
	if err != nil {
		return nil
	}
	apiRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil
	}
	var dropres dropAPIRes
	res.Body.Close()
	json.Unmarshal(apiRes, &dropres)
	timestamp, err := time.Parse(time.RFC3339, dropres.DropTime)
	if err != nil {
		return nil
	}
	return &timestamp
}

func autoOffset(count ...int) *float64 {
	c := 3
	if len(count) > 0 {
		c = count[0]
	}
	if c < 1 {
		c = 3
	}
	payload := "PUT /minecraft/profile/name/test HTTP/1.1\r\nHost: api.minecraftservices.com\r\nAuthorization: Bearer TestToken" + "\r\n"
	conn, err := tls.Dial("tcp", minecraftServicesAPIHost+":443", nil)
	if err != nil {
		return nil
	}
	sumNanos := int64(0)
	for i := 0; i < c; i++ {
		junk := make([]byte, 4096)
		conn.Write([]byte(payload))
		time1 := time.Now()
		conn.Write([]byte("\r\n"))
		conn.Read(junk)
		duration := time.Now().Sub(time1)
		sumNanos += duration.Nanoseconds()
	}
	conn.Close()
	sumNanos /= int64(c)
	avgMillis := float64(sumNanos)/float64(1000000) - float64(125)
	return &avgMillis
}

func snipe(config configuration, ch chan snipeRes) {
	time.Sleep(time.Until(config.Timestamp.Add(time.Millisecond * time.Duration(0-10000-config.Offset))))
	recvd := make([]byte, 4096)
	conn, err := tls.Dial("tcp", minecraftServicesAPIHost+":443", nil)
	if err != nil {
		if config.Debug {

			fmt.Print("\033[0;31m")
			fmt.Print(err)
			fmt.Print("\033[0m\n")
		}
		ch <- snipeRes{}
		return
	}
	payload := "PUT /minecraft/profile/name/" + config.Name + " HTTP/1.1\r\nHost: api.minecraftservices.com\r\nAuthorization: Bearer " + config.Bearer + "\r\n"
	conn.Write([]byte(payload))
	time.Sleep(time.Until(config.Timestamp.Add(time.Millisecond * time.Duration(0-config.Offset))))
	conn.Write([]byte("\r\n"))
	sent := time.Now()
	conn.Read(recvd)
	conn.Close()
	recv := time.Now()
	code, _ := strconv.Atoi(string(recvd[9:12]))
	label := config.Name
	if config.Label != nil {
		label += " @ " + *config.Label
	}

	ch <- snipeRes{
		Sent:   &sent,
		Recv:   &recv,
		Status: &code,
		Label:  &label,
	}
}
