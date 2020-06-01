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
	"io"
)

type Body struct {
	body []byte
	off  int
	done bool
}

func NewBody(m, t string) io.Reader {
	switch m {
	case "GET":
	case "OPTION":
	default:
		return &Body{body: []byte(t), off: 0, done: false}
	}
	return nil
}

func (b *Body) Read(p []byte) (int, error) {
	if b.done {
		b.done = false
		b.off = 0
		return 0, io.EOF
	}
	n := copy(p, b.body[b.off:])
	b.off += n
	if b.off == len(b.body) {
		b.done = true
	}
	return n, nil
}
