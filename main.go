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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/skarllot/magmanager/controllers"
	rqhttp "github.com/skarllot/raiqub/http"
)

const (
	CONFIG_FILE_NAME = "config.json"
)

func main() {
	content, err := ioutil.ReadFile(CONFIG_FILE_NAME)
	if err != nil {
		log.Fatalf("ReadConfigFile: %s\n", err)
		os.Exit(1)
	}

	cfg := Config{}
	if err := json.Unmarshal(content, &cfg); err != nil {
		log.Fatalf("UnmarshalConfigFile: %s\n", err)
		os.Exit(1)
	}

	session, err := getSession(cfg.Database)
	if err != nil {
		log.Fatalf("CreateDbSession: %s\n", err)
		os.Exit(1)
	}
	defer session.Close()

	router := mux.NewRouter()
	routes := rqhttp.MergeRoutes(
		controllers.NewVendorController(cfg.Database.Database, session),
	)
	for _, r := range routes {
		router.
			Methods(r.Method).
			Path(r.Path).
			Name(r.Name).
			Handler(r.ActionFunc)
	}

	fmt.Println("HTTP server listening on port", cfg.HttpServer.Port)
	http.ListenAndServe(
		fmt.Sprintf("%s:%d", cfg.HttpServer.Address, cfg.HttpServer.Port),
		router)
}
