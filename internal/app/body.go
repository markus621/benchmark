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

import "io"

type Body struct {
	body string
	done bool
}

func NewBody(m, t string) io.Reader {
	switch m {
	case "GET":
	case "OPTION":
	default:
		return &Body{body: t}
	}
	return nil
}

func (b *Body) Read(p []byte) (n int, err error) {
	if b.done {
		return 0, io.EOF
	}
	p = []byte(b.body)
	b.done = true
	return len(p), nil
}