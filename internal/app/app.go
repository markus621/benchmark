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
	"math"
	"os"
	"time"

	"github.com/deweppro/core/pkg/app"
	"github.com/olekukonko/tablewriter"
)

type App struct {
	cfg  *ConfigApp
	fc   *app.ForceClose
	c    chan CHAN
	sctx context.Context
	cncl context.CancelFunc

	ct float64
	tt float64
	cf float64
	tf float64
	ce float64
	te float64
}

type CHAN struct {
	T  float64
	OK bool
	ER error
}

func NewApp(c *ConfigApp, fc *app.ForceClose) *App {
	ctx, cncl := context.WithTimeout(fc.C, time.Duration(c.Client.Timelife)*time.Second)

	return &App{
		cfg:  c,
		fc:   fc,
		c:    make(chan CHAN, c.Client.Count*100),
		sctx: ctx,
		cncl: cncl,
	}
}

func (a *App) Up() error {
	start, cncl := context.WithCancel(a.sctx)
	defer cncl()

	for i := 0; i < a.cfg.Client.Count; i++ {
		if err := NewClient(a.cfg, a.c, start, a.sctx); err != nil {
			return err
		}
	}

	go a.run()
	return nil
}

func (a *App) Down() error {
	a.cncl()
	return nil
}

func (a *App) run() {
	t := time.Now()
	a.Tiker()
	a.result(time.Since(t).Seconds())
	a.fc.Close()
}

func (a *App) result(timelife float64) {
	data := [][]interface{}{
		{"Success", a.ct, timeop(a.ct, a.tt), opsec(a.ct, timelife)},
		{"Fail", a.cf, timeop(a.cf, a.tf), opsec(a.cf, timelife)},
		{"Error", a.ce, timeop(a.ce, a.te), opsec(a.ce, timelife)},
	}

	fmt.Println("+---------------------------------------")
	fmt.Println("| Domain   : ", a.cfg.Client.Host)
	fmt.Println("| Clients  : ", a.cfg.Client.Count)
	fmt.Println("| Timelife : ", math.Floor(timelife), "s")

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Type", "Total op", "time/op", "op/sec"})
	for _, inf := range data {
		table.Append([]string{str(inf[0]), str(inf[1]), str(inf[2]), str(inf[3])})
	}
	table.Render()
}

func (a *App) Tiker() {
	for {
		select {
		case <-a.sctx.Done():
			return

		case data := <-a.c:
			if data.ER != nil {
				a.ce++
				a.te += data.T
			} else if data.OK {
				a.ct++
				a.tt += data.T
			} else {
				a.cf++
				a.tf += data.T
			}
		}
	}
}
