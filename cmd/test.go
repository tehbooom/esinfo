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
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/tehbooom/esinfo/internal/config"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		esClient := config.SetClient(endpoint, username, password, cacert, unsafe)
		getInfo(esClient)
		color.Green("Connection successful!")
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}

func getInfo(esClient *elasticsearch.Client) {
	var (
		out = new(bytes.Buffer)
		b1  = bytes.NewBuffer([]byte{})
		b2  = bytes.NewBuffer([]byte{})
		tr  io.Reader
	)
	res, err := esClient.Info()
	if err != nil {
		color.Red("Error executing the request: %s", err)
		os.Exit(1)
	}
	tr = io.TeeReader(res.Body, b1)
	defer res.Body.Close()
	io.Copy(b2, tr)
	defer func() { res.Body = ioutil.NopCloser(b1) }()
	out.ReadFrom(b2)
	arr := out.String()
	fmt.Println(arr)
}
