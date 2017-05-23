// Copyright © 2017 Andreas Gajdosik <andreas.gajdosik@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var cfgFile string

//FBlamentPath = where config and SQLite lies
var FBlamentPath string

//YamlPath = where config file is
var YamlPath string

//SQLPath = where database is
var SQLPath string

//OutputFolderPath = where resulting files will be created
var OutputFolderPath string

//Configuration = stores data from config.yaml
var Configuration map[interface{}]interface{}

//ConfigExist = does config exist?
var ConfigExist bool

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "fblament",
	Short: "CLI tool which searches for angry comments on Facebook",
	Long: `CLI tool which searches for angry comments on Facebook.
Ideal choice when you need to support your criminal charge with some fresh data.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	//RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kidnei.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	userHome, err := homedir.Dir()
	if err != nil {
		os.Exit(1)
	}
	FBlamentPath = filepath.Join(userHome, ".fblament")
	os.MkdirAll(FBlamentPath, os.ModePerm)

	YamlPath = filepath.Join(FBlamentPath, "config.yaml")
	SQLPath = filepath.Join(FBlamentPath, "main.db")
	OutputFolderPath = filepath.Join(FBlamentPath, "results")
	os.MkdirAll(OutputFolderPath, os.ModePerm)

	data, err := ioutil.ReadFile(YamlPath)
	if err != nil {
		fmt.Println("config not found")
	}

	Configuration = make(map[interface{}]interface{})
	err = yaml.Unmarshal([]byte(data), &Configuration)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
