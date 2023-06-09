package config

import "github.com/elastic/go-elasticsearch/v8"

type EsClient struct {
	client *elasticsearch.Client
}
