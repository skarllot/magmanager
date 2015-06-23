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
	"log"

	"github.com/skarllot/magmanager/models"
	"gopkg.in/mgo.v2"
)

var vendorsCollection []models.Vendor

func init() {
	models.PreInitVendors()
}

func getSession(cfg Database, logger *log.Logger) (*mgo.Session, error) {
	dialInfo := &mgo.DialInfo{
		Addrs:    cfg.Addrs,
		Timeout:  cfg.Timeout,
		Database: cfg.Database,
		Username: cfg.Username,
		Password: cfg.Password,
	}
	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		return nil, err
	}

	session.SetMode(mgo.Monotonic, true)

	cols, err := session.DB(cfg.Database).CollectionNames()
	if err != nil {
		return nil, err
	}
	if indexOfInStringSlice(cols, models.C_VENDORS_NAME) == -1 {
		logger.Println("The collection 'vendors' was not found")

		for _, v := range vendorsCollection {
			err = session.
				DB(cfg.Database).
				C(models.C_VENDORS_NAME).
				Insert(v)
			if err != nil {
				return nil, err
			}
		}
	}

	return session, err
}
