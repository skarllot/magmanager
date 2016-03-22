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

	"gopkg.in/mgo.v2"
	"gopkg.in/raiqub/web.v0"
)

// writeObjectIdError returns a not found object ID error when aplicable;
// otherwise returns a internal server error.
func writeObjectIdError(w http.ResponseWriter, id string, err error) {
	var jerr web.JSONError
	if err == mgo.ErrNotFound {
		jerr = web.NewJSONError().
			FromError(NotFoundObjectId(id)).
			Status(http.StatusNotFound).
			Build()
	} else {
		jerr = web.NewJSONError().
			FromError(err).
			Build()
	}
	web.JSONWrite(w, jerr.Status, jerr)
}

type InvalidObjectId string

func (e InvalidObjectId) Error() string {
	return fmt.Sprintf("Invalid object ID for '%s'", string(e))
}

type NotFoundObjectId string

func (e NotFoundObjectId) Error() string {
	return fmt.Sprintf("The object ID '%s' was not found", string(e))
}
