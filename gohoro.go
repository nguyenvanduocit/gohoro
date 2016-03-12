package main

import (
	"strings"
	"fmt"
	"net/http"
	"io/ioutil"
	"regexp"
	"flag"
)

var endpoint string = "http://www.horoscope.com/us/horoscopes/general/horoscope-general-daily-today.aspx";

var signMap = map[string]int{
	"ARIES":1,
	"TAURUS":2,
	"GEMINI":3,
	"CANCER":4,
	"LEO":5,
	"VIRGO":6,
	"LIBRA":7,
	"SCORPIO":8,
	"SAGITTARIUS":9,
	"CAPRICORN":10,
	"AQUARIUS":11,
	"PISCES":12,
}

func GetSignNameById(signId int) string{
	for name,id := range signMap{
		if(id == signId){
			return name
		}
	}
	return ""
}

func GetHoroscope(signName string) (string, error) {
	signName = strings.ToUpper(signName)
	signId, ok := signMap[signName];
	if !ok {
		return "", fmt.Errorf("I don't know anything about this sign.")
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
		for k, v := range signMap {
			fmt.Println(k, "--", v)
		}
		fmt.Println(0, "\t", "Exit")
		isFirst:=true
		for signName == ""{
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
			signName = GetSignNameById(signId)
		}
	}
	fmt.Println("Please wait ...")
	content, err := GetHoroscope(signName)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(content)
}
