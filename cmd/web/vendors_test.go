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
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	rqhttp "github.com/raiqub/http"
	"github.com/skarllot/magmanager/models"
	"github.com/skarllot/raiqub/test"
	"gopkg.in/mgo.v2/bson"
)

const (
	MONGODB_URL_TPL = "mongodb://%s:%d/magmanager"
)

var (
	testSampleVendors = []models.Vendor{
		models.Vendor{
			bson.ObjectId(""),
			"Vendor001",
			[]models.Product{},
		},
		models.Vendor{
			bson.ObjectId(""),
			"Vendor002",
			[]models.Product{},
		},
	}
)

func TestVendorsBasic(t *testing.T) {
	logger = log.New(NullWriter{}, "", log.LstdFlags)

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

	router := NewRouter(session)
	ts := httptest.NewServer(router)
	defer ts.Close()

	client := NewHttpClient(t)

	// =========================================================================
	// Test GetVendorList
	// =========================================================================
	t.Log("Testing GetVendorList")
	url := fmt.Sprintf("%s/vendor", ts.URL)
	var dbVendors []models.Vendor
	client.Get(url, http.StatusOK, &dbVendors)

	if len(reference) != len(dbVendors) {
		t.Fatal("Length of vendor list do not match")
	}

	for i, v := range dbVendors {
		if err := compareVendor(&reference[i], &v); err != nil {
			t.Errorf("vendor[%d] do not match: %s\n", i, err)
		}
	}
	// =========================================================================

	// =========================================================================
	// Test GetVendor
	// =========================================================================
	t.Log("Testing GetVendor")
	for i, _ := range reference {
		url := fmt.Sprintf("%s/vendor/%s", ts.URL, reference[i].Id.Hex())
		var dbVendor models.Vendor
		client.Get(url, http.StatusOK, &dbVendor)

		if err := compareVendor(&reference[i], &dbVendor); err != nil {
			t.Errorf("Vendor do not match: %s\n", err)
		}
	}

	url = fmt.Sprintf("%s/vendor/%s", ts.URL, bson.NewObjectId().Hex())
	var jerr rqhttp.JsonError
	client.Get(url, http.StatusNotFound, &jerr)
	if jerr.Status != http.StatusNotFound {
		t.Fatalf("Invalid JsonError object: %#v", jerr)
	}
	// =========================================================================

	// =========================================================================
	// Test CreateVendor
	// =========================================================================
	t.Log("Testing CreateVendor")
	for i, _ := range testSampleVendors {
		url := fmt.Sprintf("%s/vendor", ts.URL)
		var dbVendor models.Vendor
		client.Post(url, http.StatusCreated, &testSampleVendors[i], &dbVendor)

		if err := compareVendor(&testSampleVendors[i], &dbVendor); err != nil {
			t.Errorf("Vendor do not match: %v\n", err)
		}
		testSampleVendors[i].Id = dbVendor.Id
	}
	// =========================================================================

	// =========================================================================
	// Test UpdateVendor
	// =========================================================================
	t.Log("Testing UpdateVendor")
	for _, v := range testSampleVendors {
		url := fmt.Sprintf("%s/vendor/%s", ts.URL, v.Id.Hex())
		v.Name = v.Id.Hex() + v.Name
		v.Products = []models.Product{models.Product{}, models.Product{}}
		client.Put(url, http.StatusNoContent, &v)

		var dbVendor models.Vendor
		client.Get(url, http.StatusOK, &dbVendor)
		if err := compareVendor(&v, &dbVendor); err != nil {
			t.Errorf("Vendor was not updated: %v\n", err)
		}
	}
	// =========================================================================

	// =========================================================================
	// Test RemoveVendor
	// =========================================================================
	t.Log("Testing RemoveVendor")
	for _, v := range testSampleVendors {
		url := fmt.Sprintf("%s/vendor/%s", ts.URL, v.Id.Hex())
		client.Delete(url, http.StatusNoContent)

		client.Get(url, http.StatusNotFound, nil)
	}
	// =========================================================================
}

func compareVendor(v1 *models.Vendor, v2 *models.Vendor) error {
	if string(v1.Id) != "" && v1.Id != v2.Id {
		return fmt.Errorf(
			"Vendor ID do not match. Found '%s' instead of '%s'",
			v2.Id.Hex(), v1.Id.Hex())
	}
	if v1.Name != v2.Name {
		return fmt.Errorf(
			"Vendor name do not match. Found '%s' instead of '%s'",
			v2.Name, v1.Name)
	}
	if len(v1.Products) != len(v2.Products) {
		return fmt.Errorf(
			"Product list do not match. Found %d items instead of %d",
			len(v2.Products), len(v1.Products))
	}

	return nil
}
