package net

import (
	"io/ioutil"
	"net/http"
)

func IP() string {
	url := "https://api.ipify.org?format=text"
	resp, err := http.Get(url)
	if err != nil {
		return ""
	}

	defer resp.Body.Close()

	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	return string(ip)
}
