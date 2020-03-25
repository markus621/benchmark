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

type ConfigApp struct {
	Client   Client    `json:"client"`
	Requests []Request `json:"requests"`
}

type Client struct {
	Host     string `json:"host"`
	Count    int    `json:"count"`
	Timelife int    `json:"timelife"`
	Timeout  int    `json:"timeout"`
}

type Request struct {
	Url     string   `json:"url"`
	Method  string   `json:"method"`
	Headers []Header `json:"headers"`
	Body    string   `json:"body"`
}

type Header struct {
	Key string `json:"key"`
	Val string `json:"val"`
}
