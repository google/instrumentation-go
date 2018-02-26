// Copyright 2018, OpenCensus Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ochttp_test

import (
	"log"
	"net/http"

	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/plugin/ochttp/propagation/google"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

func ExampleTransport() {

	// Subscribe to views

	err := view.Subscribe(
		ochttp.ClientRequestCountByMethod.Renamed("httpclient_requests_by_method"),
		ochttp.ClientResponseCountByStatusCode.Renamed("httpclient_responses_by_status_code"),
		ochttp.ClientLatencyView.Renamed("httpclient_latency_distribution"),
		view.New(
			"httpclient_latency_by_hostpath",
			"Client latency by URL path",
			[]tag.Key{ochttp.Host, ochttp.Path},
			ochttp.ClientLatency,
			ochttp.DefaultLatencyDistribution),
	)
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{
		Transport: &ochttp.Transport{},
	}
	_ = client // use client to perform requests
}

var usersHandler http.Handler

func ExampleHandler() {
	// Enables OpenCensus for the default serve mux.
	// By default, B3 propagation is used.
	http.Handle("/users", usersHandler)
	log.Fatal(http.ListenAndServe("localhost:8080", &ochttp.Handler{}))
}

func ExampleHandler_mux() {
	mux := http.NewServeMux()
	mux.Handle("/users", usersHandler)

	log.Fatal(http.ListenAndServe("localhost:8080", &ochttp.Handler{
		Handler:     mux,
		Propagation: &google.HTTPFormat{}, // Uses Google's propagation format.
	}))
}
