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
	"html/template"
	"net/http"

	"github.com/skarllot/magmanager/Godeps/_workspace/src/github.com/raiqub/rest"
)

const (
	rootHtml = `
  <!DOCTYPE html>
    <html>
      <head>
        <meta charset="utf-8">
        <title>MagManager API</title>
        <link rel="stylesheet" href="//cdnjs.cloudflare.com/ajax/libs/pure/0.6.0/pure-min.css">
      </head>
      <body style="margin: 20px;">
        <h2>Endpoints</h2>
        {{.}}
      </body>
    </html>
`
	endpointLine = `<p>[%s] %s</p>`
)

type ApiRoutes rest.Routes

func (s ApiRoutes) RootHandler(w http.ResponseWriter, r *http.Request) {
	content := template.Must(template.New("RootPage").Parse(rootHtml))
	routesHtml := ""
	for _, v := range s {
		if v.Method == rest.DEFAULT_CORS_PREFLIGHT_METHOD {
			continue
		}

		routesHtml += fmt.Sprintf(endpointLine, v.Method, v.Path)
	}

	content.Execute(w, template.HTML(routesHtml))
}

func (s ApiRoutes) Routes() rest.Routes {
	return rest.Routes{
		rest.Route{
			"RootPage",
			"GET",
			"/",
			false,
			s.RootHandler,
		},
	}
}
