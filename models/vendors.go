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
	"gopkg.in/mgo.v2/bson"
)

const (
	C_VENDORS_NAME = "vendors"
)

type Vendor struct {
	Id       bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string        `bson:"name" json:"name"`
	Products []Product     `bson:"products" json:"products"`
}

type Product struct {
	Id         bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Technology Technology    `bson:"technology" json:"technology"`
	Name       string        `bson:"name" json:"name"`
}

func PreInitVendors() []Vendor {
	return []Vendor{
		Vendor{bson.NewObjectId(), "None", []Product{
			Product{bson.NewObjectId(), TECH_FILE, "Local File"},
			Product{bson.NewObjectId(), TECH_FILE, "Remote File"},
		}},
		Vendor{bson.NewObjectId(), "Fujifilm", []Product{
			Product{bson.NewObjectId(), TECH_LTO1, "LTO Ultrium 1"},
			Product{bson.NewObjectId(), TECH_LTO2, "LTO Ultrium 2"},
			Product{bson.NewObjectId(), TECH_LTO3, "LTO Ultrium 3"},
			Product{bson.NewObjectId(), TECH_LTO4, "LTO Ultrium 4"},
			Product{bson.NewObjectId(), TECH_LTO5, "LTO Ultrium 5"},
			Product{bson.NewObjectId(), TECH_LTO6, "LTO Ultrium 6"},
		}},
		Vendor{bson.NewObjectId(), "HP", []Product{
			Product{bson.NewObjectId(), TECH_LTO1, "C7971A"},
			Product{bson.NewObjectId(), TECH_LTO2, "C7972A"},
			Product{bson.NewObjectId(), TECH_LTO3, "C7973A"},
			Product{bson.NewObjectId(), TECH_LTO4, "C7974A"},
			Product{bson.NewObjectId(), TECH_LTO5, "C7975A"},
			Product{bson.NewObjectId(), TECH_LTO6, "C7976A"},
		}},
		Vendor{bson.NewObjectId(), "IBM", []Product{}},
		Vendor{bson.NewObjectId(), "Imation", []Product{}},
		Vendor{bson.NewObjectId(), "Sony", []Product{
			Product{bson.NewObjectId(), TECH_LTO1, "LTX100G"},
			Product{bson.NewObjectId(), TECH_LTO2, "LTX200G"},
			Product{bson.NewObjectId(), TECH_LTO3, "LTX400G"},
			Product{bson.NewObjectId(), TECH_LTO4, "LTX800G"},
			Product{bson.NewObjectId(), TECH_LTO5, "LTX1500G"},
			Product{bson.NewObjectId(), TECH_LTO6, "LTX2500G"},
		}},
	}
}
