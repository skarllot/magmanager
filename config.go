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
	"io"
	"io/ioutil"
	"time"
)

type (
	Config struct {
		HttpServer HttpServer `json:"httpServer"`
		Database   Database   `json:"database"`
	}

	HttpServer struct {
		Address string `json:"address"`
		Port    uint16 `json:"port"`
	}

	Database struct {
		Addrs    []string      `json:"addrs"`
		Timeout  time.Duration `json:timeout""`
		Database string        `json:"database"`
		Username string        `json:"username"`
		Password string        `json:"password"`
	}
)

func ParseConfig(r io.Reader) (*Config, error) {
	content, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	if err := json.Unmarshal(content, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
