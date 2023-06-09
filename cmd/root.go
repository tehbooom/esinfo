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
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var (
	cfgFile  string
	endpoint string
	username string
	password string
	cacert   string
	format   string
	unsafe   bool
	rootCmd  = &cobra.Command{
		Use:   "esinfo",
		Short: "Grabs the types of indexes or integrations for a cluster",
		Long: `When running large elasticsearch clusters it can be difficult to know what indexes you have in the cluster without manually 
	searching through index management or using dev tools and scrolling. Esinfo queries elasticsearch for all indexes in the cluster and 
	outputs them in a nice format(csv, json, yaml). `,
	}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is esinfo.yaml)")
	rootCmd.PersistentFlags().StringVarP(&endpoint, "endpoint", "e", "localhost:9200", "Address to elasticsearch")
	rootCmd.PersistentFlags().BoolVarP(&unsafe, "unsafe", "U", false, "Ignore certificate errors")
	rootCmd.PersistentFlags().StringVarP(&username, "username", "u", "elastic", "Username for elasticsearch")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "p", "changeme", "Password for elasticsearch")
	rootCmd.PersistentFlags().StringVarP(&format, "format", "f", "csv", "Output type for file")
	rootCmd.PersistentFlags().StringVar(&cacert, "cacert", "", "Certificate Authority for cluster")

	viper.BindPFlag("cacert", rootCmd.PersistentFlags().Lookup("cacert"))
	viper.BindPFlag("endpoint", rootCmd.PersistentFlags().Lookup("endpoint"))
	viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))
	viper.BindPFlag("username", rootCmd.PersistentFlags().Lookup("username"))
	viper.BindPFlag("format", rootCmd.PersistentFlags().Lookup("format"))
	viper.BindPFlag("unsafe", rootCmd.PersistentFlags().Lookup("unsafe"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName("esinfo")
	}

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error reading config file: %w", err))
	}

	password = viper.GetString("password")
	endpoint = viper.GetString("endpoint")
	username = viper.GetString("username")
	cacert = viper.GetString("cacert")
	unsafe = viper.GetBool("unsafe")
	format = viper.GetString("format")

	if format == "json" || format == "csv" || format == "yml" || format == "yaml" {

	} else {
		fmt.Printf("Format is neither csv or json\n")
		os.Exit(1)
	}
}
