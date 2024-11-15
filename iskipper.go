package main

import (
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
)

func main() {
	// get config
	_, err := get_config()
	if err != nil {
		log.Fatal(err)
	}
	// setup http client
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{
		Jar: jar,
	}
	// sign in to iCringe
	iCringeUrl := "https://student.iclicker.com/#/courses"

	cookies, err := get_cookies()
	if err != nil {
		log.Fatal(err)
	}

	client.Jar.SetCookies(&url.URL{Scheme: "http", Host: iCringeUrl}, cookies)

	request, err := http.NewRequest("GET", iCringeUrl, nil)
	if err != nil {
		log.Fatalln(err)
	}

	request.Header.Add("Referer", "https://passport.identity.ucsb.edu/")
	request.Header.Add("Priority", "u=0, i")
	request.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")

	resp, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	stringBody := string(body)
	log.Printf(resp.Request.URL.String())
	log.Printf(resp.Status)
	log.Printf(stringBody)
}

func get_cookies() ([]*http.Cookie, error) {
	cookies := []*http.Cookie{}
	data, err := os.ReadFile("./cookies.txt")
	if err != nil {
		return nil, err
	}
	stringData := string(data)
	variables := strings.Split(stringData, ";")
	for _, variable := range variables {
		pairing := strings.SplitN(variable, "=", 2)
		key := pairing[0]
		value, err := url.QueryUnescape(pairing[1])
		if err != nil {
			return nil, err
		}
		cookies = append(cookies, &http.Cookie{Name: key, Value: value})
	}
	return cookies, nil
}
