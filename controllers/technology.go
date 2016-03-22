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
	"net/http"

	"github.com/skarllot/magmanager/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/raiqub/rest.v0"
	"gopkg.in/raiqub/web.v0"
)

type TechnologyController struct {
	dbCollection *mgo.Collection
}

func NewTechnologyController(db *mgo.Database) *TechnologyController {
	return &TechnologyController{db.C(models.C_VENDORS_NAME)}
}

func (self *TechnologyController) GetTechnologyList(
	w http.ResponseWriter,
	r *http.Request,
) {
	var list []models.Technology
	err := self.dbCollection.
		Find(bson.M{}).
		Distinct("products.technology", &list)
	if err != nil {
		jerr := web.NewJSONError().
			FromError(err).
			Status(http.StatusGone).
			Build()
		web.JSONWrite(w, jerr.Status, jerr)
		return
	}

	web.JSONWrite(w, http.StatusOK, list)
}

func (self *TechnologyController) Routes() rest.Routes {
	return rest.Routes{
		rest.Route{
			"GetTechnologyList",
			"GET",
			"/technology",
			false,
			self.GetTechnologyList,
		},
	}
}

// Ensure that TechnologyController implements Routable interface.
var _ rest.Routable = (*TechnologyController)(nil)
