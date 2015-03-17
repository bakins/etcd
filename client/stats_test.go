// Copyright 2015 CoreOS, Inc.
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

package client

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"
)

func TestStatsAPIActionLeader(t *testing.T) {
	ep := url.URL{Scheme: "http", Host: "example.com"}
	act := &statsAPIActionLeader{}

	wantURL := &url.URL{
		Scheme: "http",
		Host:   "example.com",
		Path:   "/v2/stats/leader",
	}

	got := *act.HTTPRequest(ep)
	err := assertRequest(got, "GET", wantURL, http.Header{}, nil)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestV2StatsURL(t *testing.T) {
	got := v2StatsURL(url.URL{
		Scheme: "http",
		Host:   "foo.example.com:4002",
		Path:   "/pants",
	}, "leader")
	want := &url.URL{
		Scheme: "http",
		Host:   "foo.example.com:4002",
		Path:   "/pants/v2/stats/leader",
	}

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("v2StatsURL got %#v, want %#v", got, want)
	}
}

func TestHTTPStatsAPILeaderSuccess(t *testing.T) {
	wantAction := &statsAPIActionLeader{}
	sAPI := &httpStatsAPI{
		client: &actionAssertingHTTPClient{
			t:   t,
			act: wantAction,
			resp: http.Response{
				StatusCode: http.StatusOK,
			},
			body: []byte(`{"leader":"94088180e21eb87b", "followers":{ }}`),
		},
	}

	wantResponseStats := LeaderStats{
		Leader:    "94088180e21eb87b",
		Followers: make(map[string]*FollowerStats),
	}

	l, err := sAPI.Leader(context.Background())
	if err != nil {
		t.Errorf("got non-nil err: %#v", err)
	}
	if !reflect.DeepEqual(wantResponseStats, *l) {
		t.Errorf("incorrect Stats: want=%#v got=%#v", wantResponseStats, l)
	}
}
