package main

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"strings"
)

func GenURL(p string, params ...string) (string) {
	if len(params) == 0 {
		return _Config.Host + p
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
	return _Config.Host + p
}

func GenURLWithParam(p string, param map[string]string) string {
	idx := 0
	for key, val := range param {
		if idx == 0 {
			p += "?"
		} else {
			p += "&"
		}
		p += (key + "=" + val)
		idx ++
	}
	return _Config.Host + p
}

func GETParse(url string, resObj interface{}) (error) {
	req, err := http.NewRequest("GET", url, strings.NewReader(""))
	if err != nil {
		return err
	}
	if _Config.Token != "" {
		req.Header.Add("ct", _Config.Token)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		PrintThenExit(err.Error())
	}
	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(resBody, &resObj); err != nil {
		return err
	}

	return nil
}

func POSTParse(url string, param interface{}, resObj interface{}) error {
	var body string
	if param != nil {
		if b, err := json.Marshal(param); err != nil {
			return err
		} else {
			body = string(b)
		}
	}
	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		return err
	}
	if _Config.Token != "" {
		req.Header.Add("ct", _Config.Token)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		PrintThenExit(err.Error())
	}
	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(resBody, &resObj); err != nil {
		return err
	}

	return nil
}
