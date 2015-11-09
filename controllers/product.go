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

	rqhttp "github.com/skarllot/magmanager/Godeps/_workspace/src/github.com/raiqub/http"
	"github.com/skarllot/magmanager/Godeps/_workspace/src/github.com/raiqub/rest"
	"github.com/skarllot/magmanager/Godeps/_workspace/src/gopkg.in/mgo.v2"
	"github.com/skarllot/magmanager/Godeps/_workspace/src/gopkg.in/mgo.v2/bson"
	"github.com/skarllot/magmanager/models"
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
	vars := rest.Vars(r)

	vid, ok := vars.GetObjectId("vid")
	if !ok {
		jerr := rqhttp.NewJsonErrorFromError(
			http.StatusGone, InvalidObjectId("vendor"))
		rqhttp.JsonWrite(w, jerr.Status, jerr)
		return
	}
	pid, ok := vars.GetObjectId("pid")
	if !ok {
		jerr := rqhttp.NewJsonErrorFromError(
			http.StatusGone, InvalidObjectId("product"))
		rqhttp.JsonWrite(w, jerr.Status, jerr)
		return
	}

	v := models.Vendor{}
	// db.vendors.findOne({ _id: vid },
	//	{ products: { $elemMatch: { _id: pid } } })
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

func (self *ProductController) GetProductList(
	w http.ResponseWriter,
	r *http.Request,
) {
	vid, ok := rest.Vars(r).GetObjectId("vid")
	if !ok {
		jerr := rqhttp.NewJsonErrorFromError(
			http.StatusGone, InvalidObjectId("vendor"))
		rqhttp.JsonWrite(w, jerr.Status, jerr)
		return
	}

	v := models.Vendor{}
	// db.vendors.findOne({ _id: vid })
	err := self.dbCollection.FindId(vid).One(&v)
	if err != nil {
		writeObjectIdError(w, vid.Hex(), err)
		return
	}

	rqhttp.JsonWrite(w, http.StatusOK, v.Products)
}

func (self *ProductController) CreateProduct(
	w http.ResponseWriter,
	r *http.Request,
) {
	vid, ok := rest.Vars(r).GetObjectId("vid")
	if !ok {
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

	// db.vendors.update({ _id: vid }, { $push: { products: p } })
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
	vars := rest.Vars(r)

	vid, ok := vars.GetObjectId("vid")
	if !ok {
		jerr := rqhttp.NewJsonErrorFromError(
			http.StatusGone, InvalidObjectId("vendor"))
		rqhttp.JsonWrite(w, jerr.Status, jerr)
		return
	}
	pid, ok := vars.GetObjectId("pid")
	if !ok {
		jerr := rqhttp.NewJsonErrorFromError(
			http.StatusGone, InvalidObjectId("product"))
		rqhttp.JsonWrite(w, jerr.Status, jerr)
		return
	}

	p := models.Product{}
	if !rqhttp.JsonRead(r.Body, &p, w) {
		return
	}

	p.Id = pid
	// db.vendors.update(
	//  { _id: vid, products: { $elemMatch: { _id: pid } } },
	//	{ $set: { "products.$": p } })
	err := self.dbCollection.Update(
		bson.M{"_id": vid, "products": bson.M{"$elemMatch": bson.M{"_id": pid}}},
		bson.M{"$set": bson.M{"products.$": p}})
	if err != nil {
		writeObjectIdError(w, pid.Hex(), err)
		return
	}

	rqhttp.JsonWrite(w, http.StatusNoContent, nil)
}

func (self *ProductController) RemoveProduct(
	w http.ResponseWriter,
	r *http.Request,
) {
	vars := rest.Vars(r)

	vid, ok := vars.GetObjectId("vid")
	if !ok {
		jerr := rqhttp.NewJsonErrorFromError(
			http.StatusGone, InvalidObjectId("vendor"))
		rqhttp.JsonWrite(w, jerr.Status, jerr)
		return
	}
	pid, ok := vars.GetObjectId("pid")
	if !ok {
		jerr := rqhttp.NewJsonErrorFromError(
			http.StatusGone, InvalidObjectId("product"))
		rqhttp.JsonWrite(w, jerr.Status, jerr)
		return
	}

	// db.vendors.update({ _id: vid }, { $pull: { products: { _id: pid } } })
	err := self.dbCollection.
		UpdateId(vid, bson.M{"$pull": bson.M{"products": bson.M{"_id": pid}}})
	if err != nil {
		writeObjectIdError(w, vid.Hex(), err)
		return
	}

	rqhttp.JsonWrite(w, http.StatusNoContent, nil)
}

func (self *ProductController) Routes() rest.Routes {
	return rest.Routes{
		rest.Route{
			"GetProductList",
			"GET",
			"/vendor/{vid}/product",
			false,
			self.GetProductList,
		},
		rest.Route{
			"GetProduct",
			"GET",
			"/vendor/{vid}/product/{pid}",
			false,
			self.GetProduct,
		},
		rest.Route{
			"CreateProduct",
			"POST",
			"/vendor/{vid}/product",
			false,
			self.CreateProduct,
		},
		rest.Route{
			"RemoveProduct",
			"DELETE",
			"/vendor/{vid}/product/{pid}",
			false,
			self.RemoveProduct,
		},
		rest.Route{
			"UpdateProduct",
			"PUT",
			"/vendor/{vid}/product/{pid}",
			false,
			self.UpdateProduct,
		},
	}
}
