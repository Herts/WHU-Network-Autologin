package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

type Credential struct {
	Id       string `json:"id"`
	Password string `json:"password"`
}

func main() {
	// init a new http client
	cli := &http.Client{}
	// check if the router is connected to WHU-Network
	req, _ := http.NewRequest("GET", "http://captive.apple.com", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 11_2_6 like Mac OS X) AppleWebKit/604.5.6 (KHTML, like Gecko) Mobile/15D100")
	resp, err := cli.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	byteResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	//log.Println(string(byteResp))
	// if the redirect script is not found, it can be considered success
	regex1 := regexp.MustCompile("<script>top.self.location.href='(.*)'</script>")
	result := regex1.FindSubmatch(byteResp)
	if len(result) < 2 {
		return
	}
	log.Printf("%q", result[1])
	// get the query of the redirect url
	rawUrl, err := url.Parse(string(result[1]))
	if err != nil {
		log.Panic(err)
	}
	// load id and password
	var cre Credential
	f, err := os.Open("config.json")
	if err != nil {
		log.Panic("os.Open ->", err)
	}

	err = json.NewDecoder(f).Decode(&cre)
	if err != nil {
		log.Panic("json ->", err)
	}

	//construct the request
	data := url.Values{}
	data.Set("userId", cre.Id)
	data.Set("password", cre.Password)
	data.Set("service", "Internet")
	data.Set("queryString", rawUrl.RawQuery)
	data.Set("operatorPwd", "")
	data.Set("operatorUserId", "")
	data.Set("validcode", "")
	data.Set("passwordEncrypt", "false")

	loginReq, err := http.NewRequest("POST", "http://172.19.1.9:8080/eportal/InterFace.do?method=login", strings.NewReader(data.Encode()))
	if err != nil {
		log.Println("newRequest ->", err)
	}

	header := http.Header{}
	header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	loginReq.Header = header

	loginResp, err := cli.Do(loginReq)
	if err != nil {
		log.Println(err)
	}
	defer loginResp.Body.Close()
	byteResp, err = ioutil.ReadAll(loginResp.Body)
	log.Println(string(byteResp))
}
