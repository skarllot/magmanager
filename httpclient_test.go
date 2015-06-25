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
	"net/http"
	"testing"

	rqhttp "github.com/skarllot/raiqub/http"
)

type HttpClient struct {
	client *http.Client
	t      *testing.T
}

func NewHttpClient(t *testing.T) *HttpClient {
	return &HttpClient{
		&http.Client{},
		t,
	}
}

func (s *HttpClient) Delete(url string, code int) {
	req, _ := http.NewRequest("DELETE", url, nil)
	rqhttp.HttpHeader_ContentType_Json().SetWriter(req.Header)

	res, err := s.client.Do(req)
	if err != nil {
		s.t.Fatal(err.Error())
	}

	s.parseOutput(res, code, nil)
}

func (s *HttpClient) Get(url string, code int, output interface{}) {
	res, err := s.client.Get(url)
	if err != nil {
		s.t.Fatal(err.Error())
	}

	s.parseOutput(res, code, output)
}

func (s *HttpClient) Post(url string, code int, input, output interface{}) {
	buf := s.parseInput(input)

	res, err := s.client.Post(
		url, rqhttp.HttpHeader_ContentType_Json().Value, buf)
	if err != nil {
		s.t.Fatal(err.Error())
	}

	s.parseOutput(res, code, output)
}

func (s *HttpClient) Put(url string, code int, input interface{}) {
	req, _ := http.NewRequest("PUT", url, nil)
	rqhttp.HttpHeader_ContentType_Json().SetWriter(req.Header)
	req.Body = &NullCloser{s.parseInput(input)}

	res, err := s.client.Do(req)
	if err != nil {
		s.t.Fatal(err.Error())
	}

	s.parseOutput(res, code, nil)
}

func (s *HttpClient) parseInput(input interface{}) *bytes.Buffer {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(input); err != nil {
		s.t.Fatalf("Could not encode content: %v\n", err)
	}

	return &buf
}

func (s *HttpClient) parseOutput(
	res *http.Response, code int, output interface{},
) {
	defer res.Body.Close()
	
	if res.StatusCode != code {
		var jerr rqhttp.JsonError
		json.NewDecoder(res.Body).Decode(&jerr)
		s.t.Fatalf("Unexpected HTTP status. Expected '%d' got '%d'\n",
			code, res.StatusCode)
	}

	if output != nil {
		if err := json.NewDecoder(res.Body).Decode(output); err != nil {
			s.t.Fatalf("Could not parse server response: %s\n", err)
		}
	}
}
