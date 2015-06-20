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
	C_TAPES_NAME = "tapes"
)

type Tape struct {
	Id         bson.ObjectId `bson:"_id" json:"id"`
	Product    bson.ObjectId `bson:"product_id" json:"productId"`
	Pool       bson.ObjectId `bson:"pool_id" json:"poolId"`
	Container  bson.ObjectId `bson:"container_id" json:"containerId"`
	Serial     string        `bson:"serial" json:"serial"`
	Label      string        `bson:"label" json:"label"`
	LastWrite  time.Time     `bson:"last_write" json:"lastWrite"`
	LastUpdate time.Time     `bson:"last_update" json:"lastUpdate"`
}
