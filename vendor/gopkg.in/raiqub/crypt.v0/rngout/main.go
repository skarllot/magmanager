/*
 * Copyright (C) 2015 Fabrício Godoy <skarllot@gmail.com>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 2
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 59 Temple Place - Suite 330, Boston, MA  02111-1307, USA.
 */

// This test program outputs to STDOUT SSTDEG random generator.
// Its output randomness can be checked against FIPS 140-2 tests provided by
// 'rngtest' program.
package main

import (
	"bufio"
	"os"

	"github.com/raiqub/crypt"
)

func main() {
	rng := crypt.NewSSTDEG()
	defer rng.Close()

	scanner := bufio.NewScanner(rng)
	sout := bufio.NewWriter(os.Stdout)

	for scanner.Scan() {
		sout.Write(scanner.Bytes())
	}
}
