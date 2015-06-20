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
	"github.com/skarllot/magmanager/models"
	rqhttp "github.com/skarllot/raiqub/http"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type ProductController struct {
	dbname  string
	session *mgo.Session
}

func NewProductController(db string, s *mgo.Session) *ProductController {
	return &ProductController{db, s}
}

func (self *ProductController) getVendorId(
	w http.ResponseWriter,
	r *http.Request,
	id *bson.ObjectId,
) bool {
	vid := mux.Vars(r)["vid"]
	if !bson.IsObjectIdHex(vid) {
		rqhttp.JsonWrite(w, http.StatusNotFound, "Invalid vendor ID")
		return false
	}

	*id = bson.ObjectIdHex(vid)
	return true
}

func (self *ProductController) getProductId(
	w http.ResponseWriter,
	r *http.Request,
	id *bson.ObjectId,
) bool {
	pid := mux.Vars(r)["pid"]
	if !bson.IsObjectIdHex(pid) {
		rqhttp.JsonWrite(w, http.StatusNotFound, "Invalid product ID")
		return false
	}

	*id = bson.ObjectIdHex(pid)
	return true
}

func (self *ProductController) GetProduct(
	w http.ResponseWriter,
	r *http.Request,
) {
	var vid, pid bson.ObjectId
	if !self.getVendorId(w, r, &vid) ||
		!self.getProductId(w, r, &pid) {
		return
	}

	v := models.Vendor{}
	sel := bson.M{"products": bson.M{"$elemMatch": bson.M{"_id": pid}}}
	if err := self.session.
		DB(self.dbname).
		C(models.C_VENDORS_NAME).
		Find(bson.M{"_id": vid}).
		Select(sel).
		One(&v); err != nil {
		rqhttp.JsonWrite(w, http.StatusNotFound, "Product ID not found")
		return
	}
	if len(v.Products) != 1 {
		rqhttp.JsonWrite(w, http.StatusNotFound,
			"More than one product was found")
	}

	rqhttp.JsonWrite(w, http.StatusOK, v.Products[0])
}

func (self *ProductController) CreateProduct(
	w http.ResponseWriter,
	r *http.Request,
) {
	var vid bson.ObjectId
	if !self.getVendorId(w, r, &vid) {
		return
	}

	p := models.Product{}
	if !rqhttp.JsonRead(r.Body, &p, w) {
		return
	}
	p.Id = bson.NewObjectId()

	err := self.session.
		DB(self.dbname).
		C(models.C_VENDORS_NAME).
		UpdateId(vid, bson.M{"$push": bson.M{"products": p}})
	if err != nil {
		rqhttp.JsonWrite(w, http.StatusInternalServerError, err.Error())
		return
	}

	rqhttp.HttpHeader_Location().
		SetValue(fmt.Sprintf("/vendor/%s/product/%s", vid.Hex(), p.Id.Hex())).
		SetWriter(w.Header())
	rqhttp.JsonWrite(w, http.StatusCreated, p)
}

func (self *ProductController) RemoveProduct(
	w http.ResponseWriter,
	r *http.Request,
) {
	var vid, pid bson.ObjectId
	if !self.getVendorId(w, r, &vid) ||
		!self.getProductId(w, r, &pid) {
		return
	}

	err := self.session.
		DB(self.dbname).
		C(models.C_VENDORS_NAME).
		UpdateId(vid, bson.M{"$pull": bson.M{"products": bson.M{"_id": pid}}})
	if err != nil {
		rqhttp.JsonWrite(w, http.StatusNotFound, "")
		return
	}

	rqhttp.JsonWrite(w, http.StatusOK, "")
}

func (self *ProductController) Routes() rqhttp.Routes {
	return rqhttp.Routes{
		rqhttp.Route{
			"GetProduct",
			"GET",
			"/vendor/{vid}/product/{pid}",
			false,
			self.GetProduct,
		},
		rqhttp.Route{
			"CreateProduct",
			"POST",
			"/vendor/{vid}/product",
			false,
			self.CreateProduct,
		},
		rqhttp.Route{
			"RemoveProduct",
			"DELETE",
			"/vendor/{vid}/product/{pid}",
			false,
			self.RemoveProduct,
		},
	}
}
