/*
 * Copyright 2022 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package adaptor

import (
	"bytes"
	"net/http"

	"github.com/cloudwego/hertz/pkg/protocol"
)

// GetCompatRequest only support basic function of Request, not for all.
func GetCompatRequest(req *protocol.Request) (*http.Request, error) {
	r, err := http.NewRequest(string(req.Method()), req.URI().String(), bytes.NewReader(req.Body()))
	if err != nil {
		return r, err
	}

	h := make(map[string][]string)
	req.Header.VisitAll(func(k, v []byte) {
		h[string(k)] = append(h[string(k)], string(v))
	})

	r.Header = h
	return r, nil
}

func SwapToHertzRequest(req *http.Request, hreq *protocol.Request) error {
	hreq.Header.SetRequestURI(req.RequestURI)
	hreq.Header.SetHost(req.Host)
	hreq.Header.SetProtocol(req.Proto)
	for k, v := range req.Header {
		for _, vv := range v {
			hreq.Header.Add(k, vv)
		}
	}
	req.Header = nil

	hreq.SetBodyStream(req.Body, hreq.Header.ContentLength())
	req.Body = nil
	return nil
}
