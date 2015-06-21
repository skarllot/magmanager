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

package controllers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

func readObjectId(r *http.Request, key string, id *bson.ObjectId) bool {
	strId := mux.Vars(r)[key]
	if !bson.IsObjectIdHex(strId) {
		return false
	}

	*id = bson.ObjectIdHex(strId)
	return true
}

type InvalidObjectId string

func (e InvalidObjectId) Error() string {
	return fmt.Sprintf("Invalid object ID for '%s'", string(e))
}
