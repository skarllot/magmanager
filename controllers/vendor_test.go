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
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/skarllot/magmanager/models"
	"github.com/skarllot/raiqub/docker"
	"gopkg.in/mgo.v2"
)

const (
	DB_NAME = "magmanager"
)

var (
	container          *docker.Container
	environmentRunning bool
	vendors            []models.Vendor
)

func TestCreateEnvironment(t *testing.T) {
	mgoImg := docker.NewImageMongoDB()
	if err := mgoImg.Setup(); err != nil {
		t.Error("Error setting up Docker:", err)
		return
	}

	var err error
	container, err = mgoImg.Run("mongodbgolangtest")
	if err != nil {
		t.Error("Error starting MongoDB instance:", err)
		return
	}

	if err := container.WaitStartup(5 * time.Minute); err != nil {
		t.Error(err.Error())
		return
	}
	environmentRunning = true
}

func TestMongoDbContainer(t *testing.T) {
	if !environmentRunning {
		t.Fatal("Needs a running environment to continue")
	}

	ip, err := container.IP()
	if err != nil {
		t.Error(err.Error())
		return
	}

	session, err := mgo.Dial(ip)
	if err != nil {
		t.Errorf("MongoDB connection failed, with address '%s'.", ip)
		return
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	if _, err := session.DatabaseNames(); err != nil {
		t.Error("Could not get database list from MongoDB")
		return
	}

	vendors = models.PreInitVendors()
	for _, v := range vendors {
		err := session.
			DB(DB_NAME).
			C(models.C_VENDORS_NAME).
			Insert(v)
		if err != nil {
			t.Error("Could not initialize vendors collection")
			return
		}
	}
}

func TestCRUD(t *testing.T) {
	if !environmentRunning {
		t.Fatal("Needs a running environment to continue")
	}

	ip, _ := container.IP()
	session, _ := mgo.Dial(ip)
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	ctrl := NewVendorController(session.DB(DB_NAME))
	req, err := http.NewRequest("GET", "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	ctrl.GetVendorList(w, req)
	var dbVendors []models.Vendor
	if err := json.NewDecoder(w.Body).Decode(&dbVendors); err != nil {
		t.Fatal("Error parsing HTTP response to GetVendorList:", err)
	}

	if len(dbVendors) != len(vendors) {
		t.Errorf("Vendor list do not match. Found %d instead of %d",
			len(dbVendors), len(vendors))
	}

	for i, _ := range dbVendors {
		var db, lc models.Vendor = dbVendors[i], vendors[i]
		if db.Name != lc.Name {
			t.Errorf(
				"Vendor name do not match on '%d' position. Found '%s' instead of '%s'",
				i, db.Name, lc.Name)
		}
		if len(db.Products) != len(lc.Products) {
			t.Errorf(
				"Product list do not match on '%s' vendor. Found %d instead of %d",
				db.Name, len(db.Products), len(lc.Products))
		}
	}
}

func TestDestroyEnvironment(t *testing.T) {
	environmentRunning = false
	container.Kill()
	container.Remove()
}
