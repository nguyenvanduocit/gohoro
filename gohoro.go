package main

import (
	"strings"
	"fmt"
	"net/http"
	"io/ioutil"
	"regexp"
	"flag"
	"sort"
)

var endpoint string = "http://www.horoscope.com/us/horoscopes/general/horoscope-general-daily-today.aspx";

var signMap = map[int]string{
	1:"ARIES",
	2:"TAURUS",
	3:"GEMINI",
	4:"CANCER",
	5:"LEO",
	6:"VIRGO",
	7:"LIBRA",
	8:"SCORPIO",
	9:"SAGITTARIUS",
	10:"CAPRICORN",
	11:"AQUARIUS",
	12:"PISCES",
}

func GetSignId(signName string) (int, error) {
	signName = strings.ToUpper(signName)
	for k, v := range signMap {
		if (v == signName) {
			return k, nil
		}
	}
	return 0, fmt.Errorf("Sign id not found")
}

func GetHoroscope(signName string) (string, error) {
	signId, err := GetSignId(signName);
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf("%s?sign=%d", endpoint, signId)
	resp, err := http.Get(url);
	if err != nil {
		return "", fmt.Errorf("Can not get data");
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Can not read body")
	}
	content := string(body)
	r, _ := regexp.Compile("(?s)<div class=\"block-horoscope-text[^>]*\">([^/]*)</div>")
	loginResultSubmatch := r.FindStringSubmatch(content)
	if loginResultSubmatch == nil {
		return "", fmt.Errorf("Can not read data.")
	}
	return strings.TrimSpace(loginResultSubmatch[1]), nil
}

func main() {
	var signName string
	flag.StringVar(&signName, "sign", "", "Sign name")
	flag.Parse()
	if signName == "" {
		var signId int = -1
		var keys []int
		for k := range signMap {
			keys = append(keys, k)
		}
		sort.Ints(keys)
		for _, k := range keys {
			fmt.Println(k, "\t", signMap[k])
		}
		fmt.Println(0, "\t", "Exit")
		isFirst:=true
		for signMap[signId] == ""{
			if !isFirst{
				fmt.Print("This sign is not exist, please type a number from the list below :")
			}else{
				isFirst = false;
				fmt.Print("Please choose a horoscope (number):")
			}
			fmt.Scanf("%d", &signId)
			if signId==0{
				return
			}
		}
		signName = signMap[signId]
	}
	fmt.Println("Please wait ...")
	content, err := GetHoroscope(signName)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(content)
}
