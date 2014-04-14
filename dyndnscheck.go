// DynDNSCheck is a monitoring tool for your DynDNS host.
//
// Copyright (C) 2014 CODE2K:LABS
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this
// list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice,
// this list of conditions and the following disclaimer in the documentation
// and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its contributors
// may be used to endorse or promote products derived from this software without
// specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package main

import (
	"encoding/json"
	"flag"
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

const (
	VERSION string = "1.0"
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

var (
	config = &Config{}
)

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

// getDynIP returns the current IP of your DnyDNS host.
func getDynIP(host string) (ip string, err error) {
	ips, err := net.LookupIP(host)
	if err != nil {
		log.Printf("getDynIP: %v", err)
		return
	}
	ip = ips[0].String()
	log.Printf("DynDns IP: %v", ip)
	return
}

// getCurrentIP returns your external IP by using a web service
// like http://ifconfig.me/ip or http://icanhazip.com.
func getCurrentIP() (ip string, err error) {
	res, err := http.Get(config.CurrentIPURL)
	if err != nil {
		log.Printf("getCurrentIP: %v", err)
		return
	}
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("Can't get IP from %s: %s", config.CurrentIPURL, res.Status)
		return
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Printf("getCurrentIP: %v", err)
		return
	}
	ip = strings.TrimSpace(string(body))
	if net.ParseIP(ip) == nil {
		log.Printf("Result is not an IP address: %s", ip)
		ip = ""
		err = fmt.Errorf("%s does not return an IP address!", config.CurrentIPURL)
	}
	log.Printf("Current IP: %v", ip)
	return
}

func init() {
	config := flag.String("config", "config.json", "location of the configuration file")
	flag.Parse()

	log.Print("DynDNSCheck " + VERSION)
	log.Print("Configuration: " + *config)

	loadConfig(*config)
}

func main() {
	currentIP, cerr := getCurrentIP()
	dynIP, derr := getDynIP(config.DynDNSHost)

	if cerr != nil || derr != nil || currentIP != dynIP {
		// Something is wrong. Generate alert email:
		var message string

		if cerr != nil {
			message += fmt.Sprintf("Error getting current IP: %v\n", cerr)
		}
		if derr != nil {
			message += fmt.Sprintf("Error getting DynDNS IP: %v\n", derr)
		}
		if currentIP != dynIP && len(message) == 0 {
			message += fmt.Sprintf("IPs are different!\n\nDynDNS IP:%s\n", dynIP)
		}
		if cerr == nil {
			// if available always append the current IP
			message += fmt.Sprintf("Current IP: %s\n", currentIP)
		}

		message += fmt.Sprintf("\nChecked Host: %s\n", config.DynDNSHost)

		hostname, err := os.Hostname()
		if err != nil {
			hostname = fmt.Sprintf("Error: %v", err)
		}
		message += "DynDNSCheck " + VERSION + " is running on " + hostname

		log.Print("Check failed... sending alert")
		sendMail(config.EMailSubject, message)

	} else {
		log.Print("Check OK")
	}

}
