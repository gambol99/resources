/*
Copyright 2018 All rights reserved - Appvia.io

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package utils

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type httpLogging struct {
	rt http.RoundTripper
}

// NewHTTPLogger returns a new HTTP logger
func NewHTTPLogger(t http.RoundTripper) http.RoundTripper {
	return &httpLogging{rt: t}
}

// RoundTrip implements the RoundTripper interface
func (h *httpLogging) RoundTrip(req *http.Request) (*http.Response, error) {
	content, reader, err := readBody(req.Body)
	if err != nil {
		return nil, err
	}
	req.Body = reader

	log.WithFields(log.Fields{
		"body":   content,
		"method": req.Method,
		"param":  req.URL.Query().Encode(),
		"uri":    req.RequestURI,
	}).Debugf("controller request to api")

	resp, err := h.rt.RoundTrip(req)
	if err != nil {
		return resp, err
	}
	content, reader, err = readBody(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body = reader

	log.WithFields(log.Fields{
		"body":     content,
		"response": resp.StatusCode,
	}).Debugf("controller response")

	return resp, err
}

// readBody duplicates the content
func readBody(source io.ReadCloser) ([]byte, io.ReadCloser, error) {
	if source == nil {
		return []byte{}, nil, nil
	}

	buf, err := ioutil.ReadAll(source)
	if err != nil {
		return []byte{}, nil, err
	}

	return buf, ioutil.NopCloser(bytes.NewBuffer(buf)), nil
}
