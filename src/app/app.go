package app

import (
	"log"
	"sync"

	"github.com/markus621/benchmark/src/managers/client"
	"github.com/markus621/benchmark/src/managers/config"
)

type stat struct {
	sync.Mutex
	ttl float64
	ok  int64
	er  int64
	c   int
}

// Run ...
func Run(conf *config.Config) {
	ch := make(chan client.CHAN)

	for i := 0; i < conf.Clients; i++ {
		go client.New(conf).Run(ch)
	}

	s := stat{}

	for {
		select {
		case t := <-ch:
			s.Lock()
			s.c++

			s.ttl += t.T
			s.ok += t.OK
			s.er += t.ER

			s.Unlock()
		default:
			s.Lock()
			if s.c >= conf.Clients {
				s.Unlock()

				result := 60000.0 * float64(s.ok+s.er) / s.ttl

				log.Printf("\n\nUrl: %v %v\nClients: %v\nRequests: %v#/sec\nGood req: %v\nBad req:%v\n\n", conf.Method, conf.URL, conf.Clients, int64(result), s.ok, s.er)
				return
			}
			s.Unlock()
		}
	}
}
