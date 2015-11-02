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

// Environment variables:
// PORT		Defines listening port for HTTP server
// MONGODB	Defines MongoDB database address
//			(eg: mongodb://user:password@db.example.com:55699/magmanager)
//

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/raiqub/rest"
	"github.com/skarllot/magmanager/controllers"
)

var (
	logger *log.Logger
)

func init() {
	//logger = log.New(os.Stderr, "magmanager ", log.LstdFlags|log.Lshortfile)
	logger = log.New(os.Stderr, "magmanager ", log.LstdFlags)
}

func main() {
	session, err := getSession()
	if err != nil {
		logger.Fatalf("CreateDbSession: %s\n", err)
	}
	defer session.Close()

	restMux := rest.NewRest()

	// Middlewares
	restMux.AddMiddlewarePublic(LogMiddleware)
	restMux.AddMiddlewarePublic(rest.RecoverHandlerJson)
	restMux.EnableCORS()

	// Resources
	restMux.AddResource(controllers.NewVendorController(session.DB("")))
	restMux.AddResource(controllers.NewProductController(session.DB("")))
	restMux.AddResource(controllers.NewTechnologyController(session.DB("")))
	// Shows available endpoints on root page
	restMux.AddResourceDetached(ApiRoutes(restMux.ResourcesRoutes()))

	fmt.Println("HTTP server listening on", EnvPort())
	restMux.ListenAndServe(EnvPort())
}

func LogMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		next.ServeHTTP(w, r)
		t2 := time.Now()
		logger.Printf("[%s] %q %v\n", r.Method, r.URL.String(), t2.Sub(t1))
	}

	return http.HandlerFunc(fn)
}
