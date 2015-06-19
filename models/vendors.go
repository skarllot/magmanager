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

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

const (
	C_VENDORS_NAME = "vendors"
)

type Vendor struct {
	Id       bson.ObjectId `bson:"_id" json:"id"`
	Name     string        `bson:"name" json:"name"`
	Products []Product     `bson:"products" json:"products"`
}

type Product struct {
	Id         bson.ObjectId `bson:"_id" json:"id"`
	Technology Technology    `bson:"technology" json:"technology"`
	Name       string        `bson:"name" json:"name"`
	Tapes      []Tape        `bson:"tapes" json:"tapes"`
}

type Tape struct {
	Id         bson.ObjectId `bson:"_id" json:"id"`
	Pool       bson.ObjectId `bson:"poolId" json:"poolId"`
	Container  bson.ObjectId `bson:"containerId" json:"containerId"`
	Serial     string        `bson:"serial" json:"serial"`
	Label      string        `bson:"label" json:"label"`
	LastWrite  time.Time     `bson:"lastWrite" json:"lastWrite"`
	LastUpdate time.Time     `bson:"lastUpdate" json:"lastUpdate"`
}

func PreInitVendors() []Vendor {
	return []Vendor{
		Vendor{bson.NewObjectId(), "None", []Product{
			Product{bson.NewObjectId(), TAPE_FILE, "Local File", []Tape{}},
			Product{bson.NewObjectId(), TAPE_FILE, "Remote File", []Tape{}},
		}},
		Vendor{bson.NewObjectId(), "Fujifilm", []Product{
			Product{bson.NewObjectId(), TAPE_LTO1, "LTO Ultrium 1", []Tape{}},
			Product{bson.NewObjectId(), TAPE_LTO2, "LTO Ultrium 2", []Tape{}},
			Product{bson.NewObjectId(), TAPE_LTO3, "LTO Ultrium 3", []Tape{}},
			Product{bson.NewObjectId(), TAPE_LTO4, "LTO Ultrium 4", []Tape{}},
			Product{bson.NewObjectId(), TAPE_LTO5, "LTO Ultrium 5", []Tape{}},
			Product{bson.NewObjectId(), TAPE_LTO6, "LTO Ultrium 6", []Tape{}},
		}},
		Vendor{bson.NewObjectId(), "HP", []Product{
			Product{bson.NewObjectId(), TAPE_LTO1, "C7971A", []Tape{}},
			Product{bson.NewObjectId(), TAPE_LTO2, "C7972A", []Tape{}},
			Product{bson.NewObjectId(), TAPE_LTO3, "C7973A", []Tape{}},
			Product{bson.NewObjectId(), TAPE_LTO4, "C7974A", []Tape{}},
			Product{bson.NewObjectId(), TAPE_LTO5, "C7975A", []Tape{}},
			Product{bson.NewObjectId(), TAPE_LTO6, "C7976A", []Tape{}},
		}},
		Vendor{bson.NewObjectId(), "IBM", []Product{}},
		Vendor{bson.NewObjectId(), "Imation", []Product{}},
		Vendor{bson.NewObjectId(), "Sony", []Product{
			Product{bson.NewObjectId(), TAPE_LTO1, "LTX100G", []Tape{}},
			Product{bson.NewObjectId(), TAPE_LTO2, "LTX200G", []Tape{}},
			Product{bson.NewObjectId(), TAPE_LTO3, "LTX400G", []Tape{}},
			Product{bson.NewObjectId(), TAPE_LTO4, "LTX800G", []Tape{}},
			Product{bson.NewObjectId(), TAPE_LTO5, "LTX1500G", []Tape{}},
			Product{bson.NewObjectId(), TAPE_LTO6, "LTX2500G", []Tape{}},
		}},
	}
}
