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
package cmd

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	"github.com/tehbooom/esinfo/internal/config"
	"github.com/tidwall/gjson"
	"golang.org/x/exp/slices"
)

var (
	indices     []string
	datastreams []string
)

var data = [][]string{
	{"Indices", "Data Streams"},
}

type esJSON struct {
	Indices     string `json:"index"`
	Datastreams string `json:"datastream"`
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Queries the elasticsearch cluster for all indices and datastreams and outputs them.",
	Long: `Queries the elasticsearch cluster for all indices and datastreams and outputs them to ./esinfo.{yaml,json,csv}.
					The defualt output type is csv for spreadsheet purposes.`,
	Run: func(cmd *cobra.Command, args []string) {
		esClient := config.SetClient(endpoint, username, password, cacert, unsafe)
		getDatastreams(esClient)
		getIndices(esClient)
		createMatrix()
		if format == "json" {
			createJSON()
		} else if format == "csv" {
			createCSV()
		} else if format == "yml" || format == "yaml" {
			createYAML()
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func createMatrix() {
	if len(indices) < len(datastreams) {
		diff := len(datastreams) - len(indices)
		for i := 0; i < diff; i++ {
			indices = append(indices, " ")
		}
	} else if len(datastreams) < len(indices) {
		diff := len(indices) - len(datastreams)
		for i := 0; i < diff; i++ {
			datastreams = append(datastreams, "")
		}
	}

	for i := 0; i < len(indices); i++ {
		index := indices[i]
		datastream := datastreams[i]
		data = append(data, []string{index, datastream})
	}
}

func createJSON() {
	var jsonArray []esJSON
	for i, line := range data {
		if i > 0 {
			var rec esJSON
			for j, field := range line {
				if j == 0 {
					rec.Indices = field
				} else if j == 1 {
					rec.Datastreams = field
				}
			}
			jsonArray = append(jsonArray, rec)
		}
	}
	jsonData, err := json.MarshalIndent(jsonArray, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	_ = ioutil.WriteFile("indices.json", jsonData, 0644)

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("JSON file located at %s/indices.json\n", path)
}

func createYAML() {
	var jsonArray []esJSON
	for i, line := range data {
		if i > 0 {
			var rec esJSON
			for j, field := range line {
				if j == 0 {
					rec.Indices = field
				} else if j == 1 {
					rec.Datastreams = field
				}
			}
			jsonArray = append(jsonArray, rec)
		}
	}
	jsonData, err := json.MarshalIndent(jsonArray, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	yamlData, err := yaml.JSONToYAML(jsonData)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	_ = ioutil.WriteFile("indices.yaml", yamlData, 0644)

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("YAML file located at %s/indices.yaml\n", path)
}

func createCSV() {
	file, e := os.Create("indices.csv")
	if e != nil {
		log.Fatal(e)
	}

	w := csv.NewWriter(file)

	w.WriteAll(data) // calls Flush internally

	if err := w.Error(); err != nil {
		log.Fatalln("error writing csv:", err)
	}
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("CSV file located at %s/indices.csv\n", path)
}

func getIndices(esClient *elasticsearch.Client) {

	var (
		out = new(bytes.Buffer)
		b1  = bytes.NewBuffer([]byte{})
		b2  = bytes.NewBuffer([]byte{})
		tr  io.Reader
	)
	req := esapi.CatIndicesRequest{
		Index:  []string{"_all"},
		Format: "json",
		Pri:    &[]bool{true}[0],
	}

	res, err := req.Do(context.Background(), esClient)
	if err != nil {
		fmt.Println("Error executing the request for indices:", err)
		return
	}

	tr = io.TeeReader(res.Body, b1)
	defer res.Body.Close()
	io.Copy(b2, tr)
	defer func() { res.Body = ioutil.NopCloser(b1) }()
	out.ReadFrom(b2)
	arr := out.String()

	result := gjson.Get(arr, `#.index`)

	for _, name := range result.Array() {
		if !strings.HasPrefix(name.String(), ".") {
			index := strings.Split(name.String(), "-")[0]
			if !slices.Contains(indices, index) {
				indices = append(indices, index)
			}
		}
	}
}

func getDatastreams(esClient *elasticsearch.Client) {
	var (
		out = new(bytes.Buffer)
		b1  = bytes.NewBuffer([]byte{})
		b2  = bytes.NewBuffer([]byte{})
		tr  io.Reader
	)
	req := esapi.IndicesGetDataStreamRequest{
		Name: []string{"*"},
	}

	res, err := req.Do(context.Background(), esClient)
	if err != nil {
		fmt.Println("Error executing the request:", err)
		return
	}

	tr = io.TeeReader(res.Body, b1)
	defer res.Body.Close()
	io.Copy(b2, tr)
	defer func() { res.Body = ioutil.NopCloser(b1) }()
	out.ReadFrom(b2)
	arr := out.String()

	result := gjson.Get(arr, `data_streams.#.name`)
	for _, name := range result.Array() {
		if !strings.HasPrefix(name.String(), ".") {
			index := strings.Split(name.String(), "-")[1]
			if !slices.Contains(datastreams, index) {
				datastreams = append(datastreams, index)
			}
		}
	}
}
