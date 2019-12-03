package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
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

func parseUsage() {
	resp, err := http.Get("http://202.204.48.82")
	if err != nil {
		log.Println("Failed to get usage")
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	bodyStr := string(body)

	regx := regexp.MustCompile(`time='(\d+)\s+';flow='(\d+)\s+';fsele=`)
	match := regx.FindStringSubmatch(bodyStr)
	useTime, _ := strconv.Atoi(match[1])
	flowV4, _ := strconv.Atoi(match[2])

	flow0 := flowV4 % 1024
	flow1 := flowV4 - flow0
	flow0 = flow0 * 1000
	flow0 = flow0 - flow0%1024
	useFlowV4 := float64(flow1)/1024 + float64(flow0)/1024000

	regx = regexp.MustCompile(`v6af=(\d+);v6df=`)
	match = regx.FindStringSubmatch(bodyStr)
	flowV6, _ := strconv.Atoi(match[1])
	useFlowV6 := float64(flowV6) / 4096

	regx = regexp.MustCompile(`fee='(\d+)\s+';xsele=`)
	match = regx.FindStringSubmatch(bodyStr)
	fee, _ := strconv.Atoi(match[1])
	fee1 := fee - fee%100
	money := float64(fee1) / 10000

	log.Printf("Time: %d min, FlowV4: %.2f MB, FlowV6: %.2f MB, Money: CNY %.2f\n", useTime, useFlowV4, useFlowV6, money)
}

func main() {
	id := flag.String("id", "", "Student's ID")
	pwd := flag.String("pwd", "", "Student's Password")
	flag.Parse()

	if !testNetwork() {
		if *id == "" || *pwd == "" {
			log.Fatalln("Please input your id and password")
		}

		doLogin(*id, *pwd)
	} else {
		log.Println("Already logged.")
	}

	parseUsage()
}
