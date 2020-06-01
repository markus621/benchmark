/*
 * Copyright (c) 2020.  Mikhail Knyazhev <markus621@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/gpl-3.0.html>.
 */

package app

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	cUserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.146 Safari/537.36"
)

type HttpClient struct {
	cfg    *ConfigApp
	client *http.Client
	c      chan<- CHAN
	start  context.Context
	stop   context.Context
	list   []*http.Request
}

func NewClient(cfg *ConfigApp, c chan<- CHAN, start, stop context.Context) error {
	timelife := time.Duration(cfg.Client.Timelife*2) * time.Second
	dial := &net.Dialer{Timeout: timelife, KeepAlive: timelife}
	hc := &HttpClient{
		cfg:   cfg,
		c:     c,
		start: start,
		stop:  stop,
		client: &http.Client{
			Timeout: time.Duration(cfg.Client.Timeout) * time.Second,
			Transport: &http.Transport{
				DialContext:         dial.DialContext,
				TLSHandshakeTimeout: time.Duration(0),
			},
		},
		list: make([]*http.Request, 0, len(cfg.Requests)),
	}

	for _, r := range cfg.Requests {
		request, err := http.NewRequestWithContext(stop, r.Method,
			fmt.Sprintf("%s%s", cfg.Client.Host, r.Url), NewBody(r.Method, r.Body))
		if err != nil {
			return err
		}

		request.Header.Set("User-Agent", cUserAgent)
		request.Header.Set("Connection", "keep-alive")
		request.Header.Set("Accept", "*/*")
		for _, h := range r.Headers {
			request.Header.Set(h.Key, h.Val)
		}

		hc.list = append(hc.list, request)
	}

	go hc.run()

	return nil
}

func (hc *HttpClient) run() {
	<-hc.start.Done()

	for {
		select {
		case <-hc.stop.Done():
			return
		default:
			hc.make()
		}
	}
}

func (hc *HttpClient) dorequest(r *http.Request) {
	t := time.Now()
	response, err := hc.client.Do(r)
	ttl := time.Since(t).Seconds()

	if err != nil {
		hc.c <- CHAN{
			T:  ttl,
			OK: false,
			ER: err,
		}

		logrus.WithError(err).Error("connection error")
		return
	}

	defer response.Body.Close()

	ok := 200 <= response.StatusCode && response.StatusCode < 400
	if ok {
		io.Copy(ioutil.Discard, response.Body)
	} else {
		b, _ := ioutil.ReadAll(response.Body)
		logrus.WithField("body", string(b)).Errorf("error code: %d", response.StatusCode)
	}

	hc.c <- CHAN{
		T:  ttl,
		OK: ok,
		ER: nil,
	}
}

func (hc *HttpClient) make() {
	for _, r := range hc.list {
		hc.dorequest(r)
	}
}
