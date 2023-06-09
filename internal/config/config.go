/*
Copyright © 2023 Alec Carpenter @tehbooom

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
package config

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/elastic/go-elasticsearch/v8"
)

func SetClient(endpoint, username, password, cacert string, unsafe bool) *elasticsearch.Client {
	var err error
	var esClient *elasticsearch.Client
	retryBackoff := backoff.NewExponentialBackOff()

	if unsafe {
		esClient, err = elasticsearch.NewClient(
			elasticsearch.Config{
				RetryOnStatus: []int{502, 503, 504, 429},
				Addresses:     []string{endpoint},
				Username:      username,
				Password:      password,
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: true,
					},
				},
				RetryBackoff: func(i int) time.Duration {
					if i == 1 {
						retryBackoff.Reset()
					}
					return retryBackoff.NextBackOff()
				},
			})
		if err != nil {
			log.Fatalf("ERROR: Unable to create client: %s", err)
		}
	} else if cacert != "" {
		cert, err := ioutil.ReadFile(cacert)
		if err != nil {
			log.Fatal()
		}
		esClient, err = elasticsearch.NewClient(
			elasticsearch.Config{
				RetryOnStatus: []int{502, 503, 504, 429},
				Addresses:     []string{endpoint},
				Username:      username,
				Password:      password,
				CACert:        cert,
				RetryBackoff: func(i int) time.Duration {
					if i == 1 {
						retryBackoff.Reset()
					}
					return retryBackoff.NextBackOff()
				},
			})
		if err != nil {
			log.Fatalf("ERROR: Unable to create client: %s", err)
		}
	} else {
		esClient, err = elasticsearch.NewClient(
			elasticsearch.Config{
				RetryOnStatus: []int{502, 503, 504, 429},
				Addresses:     []string{endpoint},
				Username:      username,
				Password:      password,
				RetryBackoff: func(i int) time.Duration {
					if i == 1 {
						retryBackoff.Reset()
					}
					return retryBackoff.NextBackOff()
				},
			})
		if err != nil {
			log.Fatalf("ERROR: Unable to create client: %s", err)
		}

	}

	return esClient
}
