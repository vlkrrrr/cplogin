package main

import (
	"fmt"
	"io/ioutil"
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

func main() {
	config.Load(file.Provider("config.yml"), yaml.Parser())
	go login(config.String("url"), config.String("login:user"), config.String("login:password"), time.Duration(config.Int("timeout")))
}

func login(url string, user string, password string, duration time.Duration) {
	client := http.Client{}
	loginreq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("error while creating request: %v", err)
	}
	loginreq.SetBasicAuth(user, password)
	res, err := client.Get(url)
	if err != nil {
		fmt.Printf("error while sending simple request %v", err)
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

	fmt.Printf("client: status code: %d\n", res.StatusCode)
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("client: response body: %s\n", resBody)
}
