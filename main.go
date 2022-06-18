// 参考：https://golang.hateblo.jp/entry/golang-http-cache

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Cache struct {
	Resp *http.Response
	Body []byte
}

var cacheMap = make(map[string]*Cache)

func sendRequest(rawURL string) ([]byte, *http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, nil, err
	}

	cache := cacheMap[rawURL]
	if cache != nil {
		etag := cache.Resp.Header.Get("etag")
		req.Header.Set("if-none-match", etag)
		lastModified := cache.Resp.Header.Get("last-modified")
		req.Header.Set("if-modified-since", lastModified)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	fmt.Println("----------")
	fmt.Println("Status code:", resp.StatusCode)
	// fmt.Printf("printing fetched data...\n%s\n", body[:20])

	if resp.StatusCode == http.StatusNotModified {
		// fmt.Printf("printing cached data...\n%s\n", cache.Body[:20])
		return cache.Body, cache.Resp, nil
	}
	cacheMap[rawURL] = &Cache{
		Resp: resp,
		Body: body,
	}
	return body, resp, nil
}

func main() {
	rawURL := "https://cdn-ak.f.st-hatena.com/images/fotolife/g/golang/20181009/20181009042416.png"

	for i := 0; i < 2; i++ {
		now := time.Now()
		_, _, err := sendRequest(rawURL)
		if err != nil {
			log.Fatal(err)
		}
		timeElapsed := time.Since(now).Milliseconds()
		fmt.Println("time elapsed(ms):", timeElapsed)
		time.Sleep(1000 * time.Millisecond)
	}

	fmt.Println("done.")
}