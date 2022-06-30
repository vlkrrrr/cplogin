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
	test := duration * time.Minute
	logger.Printf("CP Login started, Duration %v * minutes %v", duration, test)
	client := http.Client{}
	loginreq, err := http.NewRequest("GET", url, nil)
	if true {
		return
	}
	if err != nil {
		logger.Printf("error while creating request: %v", err)
	}
	loginreq.SetBasicAuth(user, password)
	res, err := client.Get(url)
	if err != nil {
		logger.Printf("error while sending simple request %v", err)
	}
	retry := 0
	if res.StatusCode == 200 {
		retry = 0
		time.Sleep(duration * time.Minute)
	} else {
		client.Do(loginreq)
		retry++
		time.Sleep(time.Duration(retry) * time.Minute)
	}

	logger.Printf("client: status code: %d\n", res.StatusCode)
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Printf("client: could not read response body: %s\n", err)
	}
	logger.Printf("client: response body: %s\n", resBody)
}
