package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type RequestBuffer struct {
	buf         []*http.Request
	client      http.Client
	respHandler func(*http.Response) error
}

func NewRequestBuffer(client http.Client, sz int, respHandler func(*http.Response) error) *RequestBuffer {
	return &RequestBuffer{
		buf:         make([]*http.Request, 0, sz),
		client:      client,
		respHandler: respHandler,
	}
}

func (b *RequestBuffer) Len() int {
	return len(b.buf)
}

func (b *RequestBuffer) AddReq(req *http.Request) {
	b.buf = append(b.buf, req)
}

func (b *RequestBuffer) Exec() {
	var wg sync.WaitGroup
	wg.Add(len(b.buf))
	for _, r := range b.buf {
		go func(req *http.Request) {
			defer wg.Done()
			if resp, err := b.client.Do(req); err == nil {
				if err := b.respHandler(resp); err != nil {
					log.Fatal(err)
				}
			}
		}(r)
	}
	wg.Wait()
	b.buf = b.buf[:0]
}

const reqBufSize = 5

func main() {
	if err := os.Mkdir("../../web/flags", os.ModePerm); err != nil {
		log.Fatal(err)
	}
	tr := &http.Transport{
		MaxIdleConns:        reqBufSize,
		TLSHandshakeTimeout: 0 * time.Second,
		DisableCompression:  true,
	}
	reqBuf := NewRequestBuffer(http.Client{Transport: tr}, reqBufSize, func(r *http.Response) error {
		defer r.Body.Close()
		if r.StatusCode != http.StatusOK {
			io.Copy(io.Discard, r.Body)
			return nil
		}
		path := r.Request.URL.Path
		fileName := path[strings.LastIndex(path, "/")+1:]
		file, err := os.Create("../../web/flags/" + fileName)
		if err != nil {
			return err
		}
		_, err = io.Copy(file, r.Body)
		return err
	})
	for c1 := 'A'; c1 <= 'Z'; c1++ {
		for c2 := 'A'; c2 <= 'Z'; c2++ {
			if reqBuf.Len() < reqBufSize {
				req, err := http.NewRequest("GET", fmt.Sprintf("https://osu.ppy.sh//images/flags/%c%c.png", c1, c2), nil)
				if err != nil {
					log.Fatal(err)
				}
				reqBuf.AddReq(req)
			} else {
				reqBuf.Exec()
			}
		}
	}
}
