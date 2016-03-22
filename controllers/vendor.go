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
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/raiqub/rest.v0"
	"gopkg.in/raiqub/web.v0"
)

const (
	vendorInputMaxLength = 8192
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
	id, ok := rest.Vars(r).GetObjectId("id")
	if !ok {
		jerr := web.NewJSONError().
			FromError(InvalidObjectId("vendor")).
			Status(http.StatusGone).
			Build()
		web.JSONWrite(w, jerr.Status, jerr)
		return
	}

	v := models.Vendor{}
	// db.vendors.findOne({_id: id})
	err := self.dbCollection.FindId(id).One(&v)
	if err != nil {
		writeObjectIdError(w, id.Hex(), err)
		return
	}

	web.JSONWrite(w, http.StatusOK, v)
}

func (self *VendorController) GetVendorList(
	w http.ResponseWriter,
	r *http.Request,
) {
	list := make([]models.Vendor, 0)
	// db.vendors.find()
	err := self.dbCollection.Find(nil).All(&list)
	if err != nil {
		jerr := web.NewJSONError().
			FromError(err).
			Build()
		web.JSONWrite(w, jerr.Status, jerr)
		return
	}

	web.JSONWrite(w, http.StatusOK, list)
}

func (self *VendorController) CreateVendor(
	w http.ResponseWriter,
	r *http.Request,
) {
	v := models.Vendor{}
	if !web.JSONRead(r.Body, vendorInputMaxLength, &v, w) {
		return
	}
	v.Id = bson.NewObjectId()

	// db.vendors.insert(v)
	err := self.dbCollection.Insert(v)
	if err != nil {
		jerr := web.NewJSONError().
			FromError(err).
			Build()
		web.JSONWrite(w, jerr.Status, jerr)
		return
	}

	web.NewHeader().
		Location().
		SetValue(fmt.Sprintf("/vendor/%s", v.Id.Hex())).
		Write(w.Header())
	web.JSONWrite(w, http.StatusCreated, v)
}

func (self *VendorController) RemoveVendor(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, ok := rest.Vars(r).GetObjectId("id")
	if !ok {
		jerr := web.NewJSONError().
			FromError(InvalidObjectId("vendor")).
			Status(http.StatusGone).
			Build()
		web.JSONWrite(w, jerr.Status, jerr)
		return
	}

	// db.vendors.remove({_id: id})
	err := self.dbCollection.RemoveId(id)
	if err != nil {
		writeObjectIdError(w, id.Hex(), err)
		return
	}

	web.JSONWrite(w, http.StatusNoContent, nil)
}

func (self *VendorController) UpdateVendor(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, ok := rest.Vars(r).GetObjectId("id")
	if !ok {
		jerr := web.NewJSONError().
			FromError(InvalidObjectId("vendor")).
			Status(http.StatusGone).
			Build()
		web.JSONWrite(w, jerr.Status, jerr)
		return
	}

	v := models.Vendor{}
	if !web.JSONRead(r.Body, vendorInputMaxLength, &v, w) {
		return
	}

	v.Id = bson.ObjectId("")
	// db.vendors.update({_id: id}, v)
	err := self.dbCollection.UpdateId(id, v)
	if err != nil {
		writeObjectIdError(w, id.Hex(), err)
		return
	}

	web.JSONWrite(w, http.StatusNoContent, nil)
}

func (self *VendorController) Routes() rest.Routes {
	return rest.Routes{
		rest.Route{
			"GetVendorList",
			"GET",
			"/vendor",
			false,
			self.GetVendorList,
		},
		rest.Route{
			"GetVendor",
			"GET",
			"/vendor/{id}",
			false,
			self.GetVendor,
		},
		rest.Route{
			"CreateVendor",
			"POST",
			"/vendor",
			false,
			self.CreateVendor,
		},
		rest.Route{
			"UpdateVendor",
			"PUT",
			"/vendor/{id}",
			false,
			self.UpdateVendor,
		},
		rest.Route{
			"RemoveVendor",
			"DELETE",
			"/vendor/{id}",
			false,
			self.RemoveVendor,
		},
	}
}

// Ensure that VendorController implements Routable interface.
var _ rest.Routable = (*VendorController)(nil)
