package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func testNetwork() bool {
	_, err := http.Get("https://www.qq.com")
	if err != nil {
		return false
	}
	return true
}

func getIPv6() string {
	resp, err := http.Get("http://cippv6.ustb.edu.cn/get_ip.php")

	if err != nil {
		return ""
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	ipv6 := strings.Split(string(body), "'")[1]

	return ipv6
}

func doLogin(id, password string) {
	ipv6 := getIPv6()
	params := url.Values{}
	params.Set("DDDDD", id)
	params.Set("upass", password)
	params.Set("0MKKey", "123456789")
	params.Set("v6ip", ipv6)

	resp, err := http.PostForm("http://202.204.48.82", params)

	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	if strings.Contains(string(body), "successfully logged") {
		if ipv6 != "" {
			log.Println("Logged with IPv4 & IPv6")
		} else {
			log.Println("Logged with IPv4 only")
		}
	} else {
		log.Println("Error")
	}
}

func main() {
	id := flag.String("id", "", "Student's ID")
	pwd := flag.String("pwd", "", "Student's Password")
	flag.Parse()

	if *id == "" || *pwd == "" {
		log.Fatalln("Please input your id and password")
	}

	if !testNetwork() {
		doLogin(*id, *pwd)
	} else {
		log.Println("Already logged.")
	}
}
