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
	"github.com/raiqub/rest"
	"github.com/skarllot/magmanager/controllers"
)

type Router struct {
	*rest.Rest
}

func NewRouter(session *Session) *Router {
	router := rest.NewRest()

	// Middlewares
	router.AddMiddlewarePublic(LogMiddleware)
	router.AddMiddlewarePublic(rest.RecoverHandlerJson)
	router.EnableCORS()

	// Resources
	router.AddResource(controllers.NewVendorController(session.DB("")))
	router.AddResource(controllers.NewProductController(session.DB("")))
	router.AddResource(controllers.NewTechnologyController(session.DB("")))
	// Shows available endpoints on root page
	router.AddResource(ApiRoutes(router.ResourcesRoutes()))

	return &Router{router}
}
