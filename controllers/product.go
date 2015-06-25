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

type ProductController struct {
	dbCollection *mgo.Collection
}

func NewProductController(db *mgo.Database) *ProductController {
	return &ProductController{db.C(models.C_VENDORS_NAME)}
}

func (self *ProductController) GetProduct(
	w http.ResponseWriter,
	r *http.Request,
) {
	var vid, pid bson.ObjectId
	if !readObjectId(r, "vid", &vid) {
		jerr := rqhttp.NewJsonErrorFromError(
			http.StatusGone, InvalidObjectId("vendor"))
		rqhttp.JsonWrite(w, jerr.Status, jerr)
		return
	}
	if !readObjectId(r, "pid", &pid) {
		jerr := rqhttp.NewJsonErrorFromError(
			http.StatusGone, InvalidObjectId("product"))
		rqhttp.JsonWrite(w, jerr.Status, jerr)
		return
	}

	v := models.Vendor{}
	sel := bson.M{"products": bson.M{"$elemMatch": bson.M{"_id": pid}}}
	err := self.dbCollection.Find(bson.M{"_id": vid}).Select(sel).One(&v)
	if err != nil {
		writeObjectIdError(w,
			fmt.Sprintf("%s/%s", vid.Hex(), pid.Hex()), err)
		return
	}
	if len(v.Products) != 1 {
		jerr := rqhttp.NewJsonErrorFromError(
			http.StatusNotFound, fmt.Errorf("More than one product was found"))
		rqhttp.JsonWrite(w, jerr.Status, jerr)
		return
	}

	rqhttp.JsonWrite(w, http.StatusOK, v.Products[0])
}

func (self *ProductController) CreateProduct(
	w http.ResponseWriter,
	r *http.Request,
) {
	var vid bson.ObjectId
	if !readObjectId(r, "vid", &vid) {
		jerr := rqhttp.NewJsonErrorFromError(
			http.StatusGone, InvalidObjectId("vendor"))
		rqhttp.JsonWrite(w, jerr.Status, jerr)
		return
	}

	p := models.Product{}
	if !rqhttp.JsonRead(r.Body, &p, w) {
		return
	}
	p.Id = bson.NewObjectId()

	err := self.dbCollection.
		UpdateId(vid, bson.M{"$push": bson.M{"products": p}})
	if err != nil {
		writeObjectIdError(w, vid.Hex(), err)
		return
	}

	rqhttp.HttpHeader_Location().
		SetValue(fmt.Sprintf("/vendor/%s/product/%s", vid.Hex(), p.Id.Hex())).
		SetWriter(w.Header())
	rqhttp.JsonWrite(w, http.StatusCreated, p)
}

func (self *ProductController) UpdateProduct(
	w http.ResponseWriter,
	r *http.Request,
) {
	var vid, pid bson.ObjectId
	if !readObjectId(r, "vid", &vid) {
		jerr := rqhttp.NewJsonErrorFromError(
			http.StatusGone, InvalidObjectId("vendor"))
		rqhttp.JsonWrite(w, jerr.Status, jerr)
		return
	}
	if !readObjectId(r, "pid", &pid) {
		jerr := rqhttp.NewJsonErrorFromError(
			http.StatusGone, InvalidObjectId("product"))
		rqhttp.JsonWrite(w, jerr.Status, jerr)
		return
	}

	p := models.Product{}
	if !rqhttp.JsonRead(r.Body, &p, w) {
		return
	}

	p.Id = bson.ObjectId("")
	//err := self.dbCollection
}

func (self *ProductController) RemoveProduct(
	w http.ResponseWriter,
	r *http.Request,
) {
	var vid, pid bson.ObjectId
	if !readObjectId(r, "vid", &vid) {
		jerr := rqhttp.NewJsonErrorFromError(
			http.StatusGone, InvalidObjectId("vendor"))
		rqhttp.JsonWrite(w, jerr.Status, jerr)
		return
	}
	if !readObjectId(r, "pid", &pid) {
		jerr := rqhttp.NewJsonErrorFromError(
			http.StatusGone, InvalidObjectId("product"))
		rqhttp.JsonWrite(w, jerr.Status, jerr)
		return
	}

	err := self.dbCollection.
		UpdateId(vid, bson.M{"$pull": bson.M{"products": bson.M{"_id": pid}}})
	if err != nil {
		writeObjectIdError(w, vid.Hex(), err)
		return
	}

	rqhttp.JsonWrite(w, http.StatusNoContent, nil)
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
		rqhttp.Route{
			"UpdateProduct",
			"PUT",
			"/vendor/{vid}/product/{pid}",
			false,
			self.UpdateProduct,
		},
	}
}
