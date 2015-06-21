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

	"github.com/skarllot/magmanager/models"
	rqhttp "github.com/skarllot/raiqub/http"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type VendorController struct {
	dbCollection *mgo.Collection
}

func NewVendorController(db *mgo.Database) *VendorController {
	return &VendorController{db.C(models.C_VENDORS_NAME)}
}

func (self *VendorController) GetVendor(
	w http.ResponseWriter,
	r *http.Request,
) {
	var id bson.ObjectId
	if !readObjectId(r, "id", &id) {
		jerr := rqhttp.NewJsonErrorFromError(
			http.StatusGone, InvalidObjectId("vendor"))
		rqhttp.JsonWrite(w, jerr.Status, jerr)
		return
	}

	v := models.Vendor{}
	err := self.dbCollection.FindId(id).One(&v)
	if err != nil {
		writeObjectIdError(w, id.Hex(), err)
		return
	}

	rqhttp.JsonWrite(w, http.StatusOK, v)
}

func (self *VendorController) GetVendorList(
	w http.ResponseWriter,
	r *http.Request,
) {
	list := make([]models.Vendor, 0)
	err := self.dbCollection.Find(nil).All(&list)
	if err != nil {
		jerr := rqhttp.NewJsonErrorFromError(http.StatusInternalServerError, err)
		rqhttp.JsonWrite(w, jerr.Status, jerr)
		return
	}

	rqhttp.JsonWrite(w, http.StatusOK, list)
}

func (self *VendorController) CreateVendor(
	w http.ResponseWriter,
	r *http.Request,
) {
	v := models.Vendor{}
	if !rqhttp.JsonRead(r.Body, &v, w) {
		return
	}
	v.Id = bson.NewObjectId()

	err := self.dbCollection.Insert(v)
	if err != nil {
		jerr := rqhttp.NewJsonErrorFromError(
			http.StatusInternalServerError, err)
		rqhttp.JsonWrite(w, jerr.Status, jerr)
		return
	}

	rqhttp.HttpHeader_Location().
		SetValue(fmt.Sprintf("/vendor/%s", v.Id.Hex())).
		SetWriter(w.Header())
	rqhttp.JsonWrite(w, http.StatusCreated, v)
}

func (self *VendorController) RemoveVendor(
	w http.ResponseWriter,
	r *http.Request,
) {
	var id bson.ObjectId
	if !readObjectId(r, "id", &id) {
		jerr := rqhttp.NewJsonErrorFromError(
			http.StatusGone, InvalidObjectId("vendor"))
		rqhttp.JsonWrite(w, jerr.Status, jerr)
		return
	}

	err := self.dbCollection.RemoveId(id)
	if err != nil {
		writeObjectIdError(w, id.Hex(), err)
		return
	}

	rqhttp.JsonWrite(w, http.StatusNoContent, nil)
}

func (self *VendorController) UpdateVendor(
	w http.ResponseWriter,
	r *http.Request,
) {
	var id bson.ObjectId
	if !readObjectId(r, "id", &id) {
		jerr := rqhttp.NewJsonErrorFromError(
			http.StatusGone, InvalidObjectId("vendor"))
		rqhttp.JsonWrite(w, jerr.Status, jerr)
		return
	}

	v := models.Vendor{}
	if !rqhttp.JsonRead(r.Body, &v, w) {
		return
	}

	err := self.dbCollection.UpdateId(v.Id, &v)
	if err != nil {
		writeObjectIdError(w, id.Hex(), err)
		return
	}

	rqhttp.JsonWrite(w, http.StatusOK, v)
}

func (self *VendorController) Routes() rqhttp.Routes {
	return rqhttp.Routes{
		rqhttp.Route{
			"GetVendor",
			"GET",
			"/vendor/{id}",
			false,
			self.GetVendor,
		},
		rqhttp.Route{
			"GetVendorList",
			"GET",
			"/vendor",
			false,
			self.GetVendorList,
		},
		rqhttp.Route{
			"CreateVendor",
			"POST",
			"/vendor",
			false,
			self.CreateVendor,
		},
		rqhttp.Route{
			"RemoveVendor",
			"DELETE",
			"/vendor/{id}",
			false,
			self.RemoveVendor,
		},
	}
}

// Ensure that VendorController implements Routable interface.
var _ rqhttp.Routable = (*VendorController)(nil)
