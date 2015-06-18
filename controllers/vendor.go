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
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/skarllot/magmanager/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type VendorController struct {
	session *mgo.Session
}

func NewVendorController(s *mgo.Session) *VendorController {
	return &VendorController{s}
}

func (self *VendorController) GetVendor(
	w http.ResponseWriter,
	r *http.Request,
) {
	id := mux.Vars(r)["id"]
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	oid := bson.ObjectIdHex(id)
	v := models.Vendor{}
	if err := self.session.
		DB(models.DATABASE_NAME).
		C(models.COLLECTION_VENDORS_NAME).
		FindId(oid).
		One(&v); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	vj, _ := json.Marshal(v)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", vj)
}

func (self *VendorController) CreateVendor(
	w http.ResponseWriter,
	r *http.Request,
) {
	v := models.Vendor{}
	json.NewDecoder(r.Body).Decode(&v)
	v.Id = bson.NewObjectId()

	self.session.
		DB(models.DATABASE_NAME).
		C(models.COLLECTION_VENDORS_NAME).
		Insert(v)

	vj, _ := json.Marshal(v)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s", vj)
}

func (self *VendorController) RemoveVendor(
	w http.ResponseWriter,
	r *http.Request,
) {
	id := mux.Vars(r)["id"]
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	oid := bson.ObjectIdHex(id)
	if err := self.session.
		DB(models.DATABASE_NAME).
		C(models.COLLECTION_VENDORS_NAME).
		RemoveId(oid); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}
