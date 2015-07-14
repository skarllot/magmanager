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
	"os"

	"github.com/gorilla/mux"
	"github.com/skarllot/magmanager/controllers"
	rqhttp "github.com/skarllot/raiqub/http"
	"gopkg.in/mgo.v2"
)

const (
	CONFIG_FILE_NAME = "config.json"
)

func main() {
	logger := log.New(os.Stderr, "magmanager", log.LstdFlags|log.Lshortfile)
	file, err := os.Open(CONFIG_FILE_NAME)
	if err != nil {
		logger.Fatalf("OpenConfigFile: %s\n", err)
		os.Exit(1)
	}

	cfg, err := ParseConfig(file)
	file.Close()
	if err != nil {
		logger.Fatalf("ParseConfig: %s\n", err)
		os.Exit(1)
	}

	session, err := getSession(cfg.Database, logger)
	if err != nil {
		logger.Fatalf("CreateDbSession: %s\n", err)
		os.Exit(1)
	}
	defer session.Close()

	router := createMux(cfg, session)
	fmt.Println("HTTP server listening on port", cfg.HttpServer.Port)
	http.ListenAndServe(
		fmt.Sprintf("%s:%d", cfg.HttpServer.Address, cfg.HttpServer.Port),
		router)
}

func createMux(cfg *Config, session *mgo.Session) http.Handler {
	router := mux.NewRouter()
	routes := rqhttp.MergeRoutes(
		controllers.NewVendorController(session.DB(cfg.Database.Database)),
		controllers.NewProductController(session.DB(cfg.Database.Database)),
	)
	for _, r := range routes {
		router.
			Methods(r.Method).
			Path(r.Path).
			Name(r.Name).
			Handler(r.ActionFunc)
	}

	return router
}
