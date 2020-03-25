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
	"fmt"
	"math"
	"time"
)

func str(d interface{}) string {
	return fmt.Sprint(d)
}

func opsec(c, t float64) float64 {
	if c == 0 || t == 0 {
		return 0.00
	}

	return math.Floor(100*c/t) / 100
}

func timeop(c, t float64) time.Duration {
	if c == 0 || t == 0 {
		return time.Duration(0)
	}

	return time.Duration(math.Floor(float64(time.Second.Nanoseconds()) * c / t))
}
