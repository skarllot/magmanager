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

package models

type Technology string

var Technologies []Technology

func addTechnology(s string) Technology {
	Technologies = append(Technologies, Technology(s))
	return Technology(s)
}

var (
	TAPE_FILE = addTechnology("File")
	TAPE_LTO1 = addTechnology("LTO-1")
	TAPE_LTO2 = addTechnology("LTO-2")
	TAPE_LTO3 = addTechnology("LTO-3")
	TAPE_LTO4 = addTechnology("LTO-4")
	TAPE_LTO5 = addTechnology("LTO-5")
	TAPE_LTO6 = addTechnology("LTO-6")
)
