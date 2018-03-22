package client

import (
	"bytes"
	"log"
	"net/http"
	"time"

	"github.com/markus621/benchmark/src/managers/config"
)

const (
	// CMethodPost ...
	CMethodPost = "POST"
	// CMethodGet ...
	CMethodGet = "GET"
	// CUserAgent ...
	CUserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.146 Safari/537.36"
)

// Agent ...
type Agent struct {
	conf   *config.Config
	client *http.Client

	stat    []float64
	okCount int64
	erCount int64
}

// CHAN ...
type CHAN struct {
	T  float64
	OK int64
	ER int64
}

// New ...
func New(conf *config.Config) *Agent {
	return &Agent{
		conf: conf,
		client: &http.Client{
			Timeout: time.Duration(conf.Timeout) * time.Second,
			//Transport: nil,
			//CheckRedirect:nil,
			//Jar:nil,
		},
		stat:    make([]float64, 0),
		okCount: 0,
		erCount: 0,
	}
}

// Run ...
func (c *Agent) Run(ch chan CHAN) {

	timeTick := time.Tick(1 * time.Nanosecond)
	timeAfter := time.After(time.Duration(c.conf.Timelife) * time.Second)

	for {
		select {
		case <-timeTick:
			c.request()

		case <-timeAfter:
			sum := 0.0

			for _, v := range c.stat {
				sum += v
			}

			ch <- CHAN{
				ER: c.erCount,
				OK: c.okCount,
				T:  sum,
			}

			return
		}
	}
}

func (c *Agent) request() {
	switch c.conf.Method {
	case CMethodGet:
		c.exec()
	case CMethodPost:
		c.exec()
	default:
		log.Fatalf("Method <%s> not support", c.conf.Method)
	}
}

func (c *Agent) body() *bytes.Buffer {
	return bytes.NewBuffer([]byte(c.conf.Body))
}

func (c *Agent) exec() {

	if request, err := http.NewRequest(c.conf.Method, c.conf.URL, c.body()); err == nil {
		t := time.Now()

		request.Header.Set("Content-Type", c.conf.ContentType)
		request.Header.Set("User-Agent", CUserAgent)

		response, err := c.client.Do(request)
		len := float64(time.Since(t) / time.Millisecond)

		if response != nil {
			defer response.Body.Close()
		}

		if err != nil {
			c.erCount++
		} else if response.StatusCode >= 200 || response.StatusCode < 400 {
			c.okCount++
		} else {
			c.erCount++
		}

		c.stat = append(c.stat, len)
	} else {
		println(err.Error())
	}
}
