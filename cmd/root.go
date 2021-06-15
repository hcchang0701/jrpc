/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strings"

	"github.com/hcchang0701/jrpc/model"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var (
	cfgFile string
	client  *http.Client
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:          "jrpc",
	Short:        "A CLI tool for JSON-RPC 2.0",
	RunE:         handler,
	SilenceUsage: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.jrpc.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolVarP(&isBatch, "batch", "b", false, "Batch mode")
	client = http.DefaultClient
}

func handler(cmd *cobra.Command, args []string) error {

	obj := model.Request{}
	err := readYaml(args[0], &obj)
	if err != nil {
		return errors.Wrap(err, "Request file")
	}

	b, err := json.Marshal(obj.Body)
	if err != nil {
		return errors.Wrap(err, "Marshal request")
	}

	req, err := http.NewRequest(http.MethodPost, obj.Url, bytes.NewReader(b))
	if err != nil {
		return errors.Wrap(err, "Init http request")
	}

	if obj.Header != nil {
		for k, vals := range obj.Header {
			for _, val := range vals {
				req.Header.Add(k, val)
			}
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "Http")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "Read http body")
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(string(body))
	}

	fmt.Println(string(body))
	return nil
}

func readYaml(path string, obj interface{}) error {

	if !strings.HasSuffix(path, "yml") && !strings.HasSuffix(path, "yaml") {
		return errors.New("Expect a yaml file")
	}

	if reflect.TypeOf(obj).Kind() != reflect.Ptr {
		return errors.New("Expect a pointer for file content")
	}

	file, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(file, obj)
	if err != nil {
		return err
	}

	return nil
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".jrpc" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".jrpc")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
