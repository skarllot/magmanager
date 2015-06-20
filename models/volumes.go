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
	C_VOLUMES_NAME = "volumes"
)

type Volume struct {
	Id        bson.ObjectId `bson:"_id" json:"id"`
	ProductId bson.ObjectId `bson:"product_id" json:"productId"`
	PoolId    bson.ObjectId `bson:"pool_id" json:"poolId"`
	Serial    string        `bson:"serial" json:"serial"`
	Label     string        `bson:"label" json:"label"`
	LastWrite time.Time     `bson:"last_write" json:"lastWrite"`
	History   []VolHistory  `bson:"history" json:"history"`
}

type VolHistory struct {
	Id          bson.ObjectId `bson:"_id" json:"id"`
	Date        time.Time     `bson:"date" json:"date"`
	Status      VolStatus     `bson:"status" json:"status"`
	ContainerId bson.ObjectId `bson:"container_id" json:"containerId"`
	Details     string        `bson:"details" json:"details"`
}
