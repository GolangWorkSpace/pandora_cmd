package main

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
)

func GenURL(p string, params ...string) (string) {
	if len(params) == 0 {
		return _Host + p
	}
	for idx, param := range params {
		i := idx / 2
		im := idx % 2
		if im == 0 {
			if i == 0 {
				p += "?"
			} else {
				p += "&"
			}
			p += ( param + "=" )
		} else {
			p += param
		}
	}
	return _Host + p
}

func GETParse(url string, res interface{}) (error) {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, res); err != nil {
		return err
	}
	return nil
}
