package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	DynDNSHost    string
	CurrentIPURL  string
	EMailFrom     string
	EMailTo       string
	EMailSubject  string
	EMailServer   string
	EMailPort     int
	EMailPassword string
}

var config = &Config{}

func loadConfig(filename string) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("open config: ", err)
	}
	if err = json.Unmarshal(file, config); err != nil {
		log.Fatalf("parse config: ", err)
	}
}

func sendMail(subject string, message string) {
	auth := smtp.PlainAuth(
		"",
		config.EMailFrom,
		config.EMailPassword,
		config.EMailServer)

	err := smtp.SendMail(
		config.EMailServer+":"+strconv.Itoa(config.EMailPort),
		auth,
		config.EMailFrom,
		[]string{config.EMailTo},
		[]byte("Subject: "+subject+"\n\n"+message))

	if err != nil {
		log.Printf("sendMail: %v", err)
	}

}

func getDynIP(host string) (ip string, err error) {
	ips, err := net.LookupIP(host)
	if err != nil {
		log.Printf("getDynIP: %v", err)
	} else {
		ip = ips[0].String()
	}
	return
}

func getCurrentIP() (ip string, err error) {
	res, err := http.Get(config.CurrentIPURL)
	if err != nil {
		log.Printf("Error: %v", err)
	} else if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("Can't get IP from %s: %s", config.CurrentIPURL, res.Status)
	} else {
		body, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			log.Printf("Error: %v", err)
		} else {
			ip = strings.TrimSpace(string(body))
		}
	}
	return
}

func main() {
	loadConfig("config.json")

	currentIP, cerr := getCurrentIP()
	dynIP, derr := getDynIP(config.DynDNSHost)

	var message string

	if cerr != nil {
		message += fmt.Sprintf("Error getting current IP: %v\n", cerr)
	} else {
		log.Printf("Current IP: %v", currentIP)
	}
	if derr != nil {
		message += fmt.Sprintf("Error getting DynDNS IP: %v\n", derr)
	} else {
		log.Printf("DynDns IP: %v", dynIP)
	}
	if currentIP != dynIP && len(message) == 0 {
		message += fmt.Sprintf("IPs are different!\n\nCurrent IP: %s\nDynDNS IP:%s\n", currentIP, dynIP)
	}

	if len(message) > 0 {
		message += fmt.Sprintf("\nChecked Host: %s\n", config.DynDNSHost)

		hostname, err := os.Hostname()
		if err != nil {
			hostname = fmt.Sprintf("Error: %v", err)
		}
		message += "DynDNSCheck is running on host " + hostname

		log.Print("Check failed... sending alert")
		sendMail(config.EMailSubject, message)
	} else {
		log.Print("Check OK")
	}

}
