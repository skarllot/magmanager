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
	rqhttp "github.com/skarllot/raiqub/http"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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
		jerr := rqhttp.NewJsonErrorFromError(
			http.StatusGone, err)
		rqhttp.JsonWrite(w, jerr.Status, jerr)
		return
	}

	rqhttp.JsonWrite(w, http.StatusOK, list)
}

func (self *TechnologyController) Routes() rqhttp.Routes {
	return rqhttp.Routes{
		rqhttp.Route{
			"GetTechnologyList",
			"GET",
			"/technology",
			false,
			self.GetTechnologyList,
		},
	}
}

// Ensure that TechnologyController implements Routable interface.
var _ rqhttp.Routable = (*TechnologyController)(nil)
