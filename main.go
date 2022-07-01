package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
)

/*
read login, password, url, sleeptime after success from file
create go routine that fires login request and sleeps after success
*/
var config = koanf.New(":")
var logger = log.Logger{}

func main() {
	logfile, err := os.OpenFile("cp_log.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Could not create/append file: %v", err)
	}
	logger = *log.New(logfile, "CP_", log.LstdFlags)
	err = config.Load(file.Provider("config.yml"), yaml.Parser())
	if err != nil {
		logger.Fatalf("read config error; %v", err)
	}
	login(config.String("url"), config.String("login:user"), config.String("login:password"), time.Duration(config.Int("timeout")))
}

func login(url string, user string, password string, duration time.Duration) {
	logger.Print("CP Login started")
	retry := 0
	client := http.Client{}

	res, err := client.Get(url)
	if err != nil {
		logger.Printf("error while sending base request %v", err)
	}

	if res.StatusCode == 200 {
		retry = 0
		time.Sleep(duration * time.Minute)
	} else if res.StatusCode == 302 {
		authurl := res.Header.Get("Location")
		if len(authurl) > 0 {
			loginreq, err := http.NewRequest("GET", "authurl", nil)
			if err != nil {
				logger.Printf("error while creating login request: %v", err)
			}
			loginreq.Header.Add("Authorization", "Basic")
			client.Do(loginreq)
		}

		retry++
		time.Sleep(time.Duration(retry) * time.Minute)
	} else {
		//todo
	}

	//loginreq.SetBasicAuth("hess", "Kopenha§1")
	res, err := client.Do(loginreq)
	if err != nil {
		logger.Printf("error while sending simple request %v", err)
	}

	for _, element := range res.Header {
		logger.Printf("header-> %v", element)
	}
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Printf("client: could not read response body: %s\n", err)
	}
	logger.Printf("client: response body: %s\n", resBody)
	if true {
		return
	}
	if res.StatusCode == 200 {
		retry = 0
		time.Sleep(duration * time.Minute)
	} else {
		client.Do(loginreq)
		retry++
		time.Sleep(time.Duration(retry) * time.Minute)
	}

	logger.Printf("client: status code: %d\n", res.StatusCode)

}
