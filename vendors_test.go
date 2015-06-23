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
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/skarllot/magmanager/models"
	rqhttp "github.com/skarllot/raiqub/http"
	"github.com/skarllot/raiqub/test"
	"gopkg.in/mgo.v2/bson"
)

const (
	CONFIG_SAMPLE_FILE = "config.sample.json"
)

func TestVendorsBasic(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", log.LstdFlags)
	file, err := os.Open(CONFIG_FILE_NAME)
	if err != nil {
		t.Fatalf("Error opening configuration file: %s\n", err)
	}
	defer file.Close()

	cfg, err := ParseConfig(file)
	if err != nil {
		t.Fatalf("Error parsing configuration file: %s\n", err)
	}

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

	cfg.Database.Addrs = []string{net[0].FormatDialAddress()}
	cfg.Database.Username = ""
	cfg.Database.Password = ""
	session, err := getSession(cfg.Database, logger)
	if err != nil {
		t.Fatalf("Error opening a MongoDB session: %s\n", err)
	}
	defer session.Close()

	reference := vendorsCollection[:]

	router := createMux(cfg, session)
	ts := httptest.NewServer(router)
	defer ts.Close()

	// =========================================================================
	// Test GetVendorList
	// =========================================================================
	url := fmt.Sprintf("%s/vendor", ts.URL)
	var dbVendors []models.Vendor
	httpGet(url, http.StatusOK, &dbVendors, t)

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
	for i, _ := range reference {
		url := fmt.Sprintf("%s/vendor/%s", ts.URL, reference[i].Id.Hex())
		var dbVendor models.Vendor
		httpGet(url, http.StatusOK, &dbVendor, t)

		if err := compareVendor(&reference[i], &dbVendor); err != nil {
			t.Errorf("Vendor do not match: %s\n", err)
		}
	}

	url = fmt.Sprintf("%s/vendor/%s", ts.URL, bson.NewObjectId().Hex())
	var jerr rqhttp.JsonError
	httpGet(url, http.StatusNotFound, &jerr, t)
	if jerr.Status != http.StatusNotFound {
		t.Fatalf("Invalid JsonError object: %#v", jerr)
	}
	// =========================================================================
}

func httpGet(url string, code int, v interface{}, t *testing.T) {
	res, err := http.Get(url)
	if err != nil {
		t.Fatal(err.Error())
	}

	if res.StatusCode != code {
		t.Fatalf("Unexpected HTTP status. Expected '%d' got '%d'.",
			code, res.StatusCode)
	}

	if err := json.NewDecoder(res.Body).Decode(v); err != nil {
		t.Fatalf("Could not parse server response: %s\n", err)
	}
}

func compareVendor(v1 *models.Vendor, v2 *models.Vendor) error {
	if v1.Id != v2.Id {
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
