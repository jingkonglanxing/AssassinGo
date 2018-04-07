package intruder

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
)

// Intruder intrudes the target.
type Intruder struct {
	target      string
	header      string
	intrudeType string
	re          *regexp.Regexp
	payload     []string
}

// NewIntruder returns a new intruder.
func NewIntruder(target, header string) *Intruder {
	return &Intruder{
		target: target,
		header: header,
		re:     regexp.MustCompile(`\$\$(.*?)\$\$`),
	}
}

// Run the inturder.
func (i *Intruder) Run(conn *websocket.Conn) {
	// 	i.header = `GET /$$1$$ HTTP/1.1
	// Host: sweety-birdsong-recg.cc
	// User-Agent: Mozilla/5.0 (X11; Linux x86_64; rv:61.0) Gecko/20100101 Firefox/61.0
	// Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8
	// Accept-Language: en-US,en;q=0.5
	// Accept-Encoding: gzip, deflate
	// Connection: close
	// Upgrade-Insecure-Requests: 1`
	payload := []string{"1", "2", "3"}
	for _, p := range payload {
		resp := i.fetch(p)
		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		ret := map[string]string{
			"payload":     p,
			"resp_status": strconv.Itoa(resp.StatusCode),
			"resp_len":    strconv.Itoa(len(string(body))),
		}
		conn.WriteJSON(ret)
	}
}

func (i *Intruder) fetch(payload string) *http.Response {
	client := &http.Client{}
	req := i.parse(payload)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	return resp
}

func (i *Intruder) parse(payload string) *http.Request {
	i.header = i.re.ReplaceAllString(i.header, payload)
	hr := strings.Split(i.header, "\n")

	method := strings.Split(hr[0], " ")[0]
	path := strings.Split(hr[0], " ")[1]
	req, _ := http.NewRequest(method, "http://"+i.target+path, nil)

	for _, row := range hr[1:] {
		hh := strings.Split(row, ": ")
		k := hh[0]
		v := hh[1]
		req.Header.Add(k, v)
	}
	return req
}