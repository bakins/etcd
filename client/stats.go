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
	"encoding/json"
	"net/http"
	"net/url"
	"path"

	"github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"
)

var (
	defaultV2StatsPrefix = "/v2/stats"
)

type (
	// LeaderStats represents the stats of the current leader.
	LeaderStats struct {
		// Leader is the id of the current leader.
		Leader    string                    `json:"leader"`
		Followers map[string]*FollowerStats `json:"followers"`
	}

	// FollowerStats represents the stats of a current follower.
	FollowerStats struct {
		Counts struct {
			Fail    uint64 `json:"fail"`
			Success uint64 `json:"success"`
		} `json:"counts"`

		Latency struct {
			Current           float64 `json:"current"`
			Average           float64 `json:"average"`
			StandardDeviation float64 `json:"standardDeviation"`
			Minimum           float64 `json:"minimum"`
			Maximum           float64 `json:"maximum"`
		} `json:"latency"`
	}

	// StatsAPI is used to fetch the stats.
	StatsAPI interface {
		// Leader returns the current leader stats.
		Leader(ctx context.Context) (*LeaderStats, error)
	}

	httpStatsAPI struct {
		client httpClient
	}
)

// UnmarshalJSON unmarshals json data into a valid LeaderStats struct.
func (l *LeaderStats) UnmarshalJSON(data []byte) error {
	d := struct {
		Leader    string
		Followers map[string]*FollowerStats
	}{}

	if err := json.Unmarshal(data, &d); err != nil {
		return err
	}

	if d.Followers == nil {
		d.Followers = make(map[string]*FollowerStats)
	}

	l.Leader = d.Leader
	l.Followers = d.Followers

	return nil
}

// NewStatsAPI constructs a new StatsAPI that uses HTTP to
// interact with etcd's statship API.
func NewStatsAPI(c Client) StatsAPI {
	return &httpStatsAPI{
		client: c,
	}
}

func (s *httpStatsAPI) Leader(ctx context.Context) (*LeaderStats, error) {
	req := &statsAPIActionLeader{}
	resp, body, err := s.client.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	if err := assertStatusCode(resp.StatusCode, http.StatusOK); err != nil {
		return nil, err
	}

	var l LeaderStats
	if err := json.Unmarshal(body, &l); err != nil {
		return nil, err
	}

	return &l, nil
}

type statsAPIActionLeader struct{}

func (l *statsAPIActionLeader) HTTPRequest(ep url.URL) *http.Request {
	u := v2StatsURL(ep, "leader")
	req, _ := http.NewRequest("GET", u.String(), nil)
	return req
}

func v2StatsURL(ep url.URL, key string) *url.URL {
	ep.Path = path.Join(ep.Path, defaultV2StatsPrefix, key)
	return &ep
}
