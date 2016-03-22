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

package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/skarllot/magmanager/models"
	"github.com/skarllot/raiqub/test"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/raiqub/web.v0"
)

var (
	testSampleProducts = []models.Product{
		models.Product{
			bson.ObjectId(""),
			models.TECH_LTO6,
			"LTO-6",
		},
		models.Product{
			bson.ObjectId(""),
			models.TECH_FILE,
			"Raw file",
		},
	}
)

func TestProductsBasic(t *testing.T) {
	var buf bytes.Buffer
	logger = log.New(&buf, "", log.LstdFlags)

	mongo := test.NewMongoDBEnvironment(t)
	if !mongo.Applicability() {
		t.Skip("This test connot be run because Docker is not acessible")
	}

	if !mongo.Run() {
		t.Fatal("Could not start MongoDB server")
	}
	defer mongo.Stop()

	net, err := mongo.Network()
	if err != nil {
		t.Fatalf("Error getting MongoDB IP address: %s\n", err)
	}

	mgourl := fmt.Sprintf(MONGODB_URL_TPL, net[0].IpAddress, net[0].Port)
	session, err := openSession(mgourl)
	if err != nil {
		t.Fatalf("Error opening a MongoDB session: %s\n", err)
	}
	defer session.Close()

	err = session.FillCollectionsIfEmpty()
	if err != nil {
		t.Fatalf("Error inserting data on empty collections: %s\n", err)
	}

	reference := vendorsCollection[:]
	vendorId := reference[0].Id

	router := NewRouter(session)
	ts := httptest.NewServer(router)
	defer ts.Close()

	client := NewHttpClient(t)

	// =========================================================================
	// Test GetProductList
	// =========================================================================
	t.Log("Testing GetProductList")
	for i, _ := range reference {
		url := fmt.Sprintf("%s/vendor/%s/product", ts.URL, reference[i].Id.Hex())
		var dbProducts []models.Product
		client.Get(url, http.StatusOK, &dbProducts)

		if len(reference[i].Products) != len(dbProducts) {
			t.Fatal("Length of product list do not match")
		}

		for j, p := range dbProducts {
			if err := compareProduct(
				&reference[i].Products[j], &p); err != nil {
				t.Errorf("product[%d] do not match: %s\n", j, err)
			}
		}
	}
	// =========================================================================

	// =========================================================================
	// Test GetProduct
	// =========================================================================
	t.Log("Testing GetProduct")
	for i, _ := range reference {
		for j, _ := range reference[i].Products {
			url := fmt.Sprintf("%s/vendor/%s/product/%s",
				ts.URL, reference[i].Id.Hex(),
				reference[i].Products[j].Id.Hex())
			var dbProduct models.Product
			client.Get(url, http.StatusOK, &dbProduct)

			if err := compareProduct(
				&reference[i].Products[j], &dbProduct); err != nil {
				t.Errorf("Product do not match: %s\n", err)
			}
		}
	}

	url := fmt.Sprintf("%s/vendor/%s/product/%s",
		ts.URL, reference[0].Id.Hex(), bson.NewObjectId().Hex())
	var jerr web.JSONError
	client.Get(url, http.StatusNotFound, &jerr)
	if jerr.Status != http.StatusNotFound {
		t.Fatalf("Invalid JsonError object: %#v", jerr)
	}
	// =========================================================================

	// =========================================================================
	// Test CreateProduct
	// =========================================================================
	t.Log("Testing CreateProduct")
	for i, _ := range testSampleProducts {
		url := fmt.Sprintf("%s/vendor/%s/product", ts.URL, vendorId.Hex())
		var dbProduct models.Product
		client.Post(url, http.StatusCreated, &testSampleProducts[i], &dbProduct)

		if err := compareProduct(&testSampleProducts[i], &dbProduct); err != nil {
			t.Errorf("Product do not match: %v\n", err)
		}
		testSampleProducts[i].Id = dbProduct.Id
	}
	// =========================================================================

	// =========================================================================
	// Test UpdateProduct
	// =========================================================================
	t.Log("Testing UpdateProduct")
	for _, p := range testSampleProducts {
		url := fmt.Sprintf("%s/vendor/%s/product/%s",
			ts.URL, vendorId.Hex(), p.Id.Hex())
		p.Name = p.Id.Hex() + p.Name
		p.Technology = models.Technology(p.Id.Hex())
		client.Put(url, http.StatusNoContent, &p)

		var dbProduct models.Product
		client.Get(url, http.StatusOK, &dbProduct)
		if err := compareProduct(&p, &dbProduct); err != nil {
			t.Errorf("Product was not updated: %v\n", err)
		}
	}
	// =========================================================================

	// =========================================================================
	// Test RemoveProduct
	// =========================================================================
	t.Log("Testing RemoveProduct")
	for _, p := range testSampleProducts {
		url := fmt.Sprintf("%s/vendor/%s/product/%s",
			ts.URL, vendorId.Hex(), p.Id.Hex())
		client.Delete(url, http.StatusNoContent)

		client.Get(url, http.StatusNotFound, nil)
	}
	// =========================================================================
}

func compareProduct(p1 *models.Product, p2 *models.Product) error {
	if string(p1.Id) != "" && p1.Id != p2.Id {
		return fmt.Errorf(
			"Product ID do not match. Found '%s' instead of '%s'",
			p2.Id.Hex(), p1.Id.Hex())
	}
	if p1.Name != p2.Name {
		return fmt.Errorf(
			"Product name do not match. Found '%s' instead of '%s'",
			p2.Name, p1.Name)
	}
	if p1.Technology != p2.Technology {
		return fmt.Errorf(
			"Product technology do not match. Found %d items instead of %d",
			p2.Technology, p1.Technology)
	}

	return nil
}
