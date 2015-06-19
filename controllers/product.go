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

func (self *ProductController) Getproduct(
	w http.ResponseWriter,
	r *http.Request,
) {
	vid := mux.Vars(r)["vid"]
	if !bson.IsObjectIdHex(vid) {
		rqhttp.JsonWrite(w, http.StatusNotFound, "Invalid vendor ID")
		return
	}
	pid := mux.Vars(r)["pid"]
	if !bson.IsObjectIdHex(pid) {
		rqhttp.JsonWrite(w, http.StatusNotFound, "Invalid product ID")
		return
	}

	v := models.Vendor{}
	query := bson.M{
		"_id": bson.ObjectIdHex(vid),
	}
	sel := bson.M{"products": bson.M{"$elemMatch": bson.M{
		"_id": bson.ObjectIdHex(pid)},
	}}
	if err := self.session.
		DB(self.dbname).
		C(models.C_VENDORS_NAME).
		Find(query).
		Select(sel).
		One(&v); err != nil {
		rqhttp.JsonWrite(w, http.StatusNotFound, "Product ID not found")
		return
	}
	if len(v.Products) != 1 {
		rqhttp.JsonWrite(w, http.StatusNotFound, "More than one product was found")
	}

	rqhttp.JsonWrite(w, http.StatusOK, v.Products[0])
}

func (self *ProductController) CreateProduct(
	w http.ResponseWriter,
	r *http.Request,
) {
	vid := mux.Vars(r)["vid"]
	if !bson.IsObjectIdHex(vid) {
		rqhttp.JsonWrite(w, http.StatusNotFound, "Invalid vendor ID")
		return
	}

	v := models.Vendor{}
	if err := self.session.
		DB(self.dbname).
		C(models.C_VENDORS_NAME).
		FindId(bson.ObjectIdHex(vid)).
		One(&v); err != nil {
		rqhttp.JsonWrite(w, http.StatusNotFound, "")
		return
	}

	p := models.Product{}
	if !rqhttp.JsonRead(r.Body, &p, w) {
		return
	}
	p.Id = bson.NewObjectId()
	v.Products = append(v.Products, p)

	if err := self.session.
		DB(self.dbname).
		C(models.C_VENDORS_NAME).
		UpdateId(v.Id, v); err != nil {
		rqhttp.JsonWrite(w, http.StatusInternalServerError, err.Error())
		return
	}

	rqhttp.HttpHeader_Location().
		SetValue(fmt.Sprintf("/vendor/%s/product/%s", v.Id.Hex(), p.Id.Hex())).
		SetWriter(w.Header())
	rqhttp.JsonWrite(w, http.StatusCreated, p)
}

func (self *ProductController) RemoveProduct(
	w http.ResponseWriter,
	r *http.Request,
) {
	id := mux.Vars(r)["id"]
	if !bson.IsObjectIdHex(id) {
		rqhttp.JsonWrite(w, http.StatusNotFound, "")
		return
	}

	oid := bson.ObjectIdHex(id)
	if err := self.session.
		DB(self.dbname).
		C(models.C_VENDORS_NAME).
		RemoveId(oid); err != nil {
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
			self.Getproduct,
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
