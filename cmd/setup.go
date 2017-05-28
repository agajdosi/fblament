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
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Helps to setup FBlament database and config file",
	Long:  `Helps to setup FBlament database and config file in .fblament folder inside of the home folder.`,
	Run: func(cmd *cobra.Command, args []string) {
		deleteDatabase()
		createDatabase()
		deleteConfig()
		createConfig()
	},
}

func init() {
	RootCmd.AddCommand(setupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func deleteDatabase() {
	os.Remove(SQLPath)
	return
}

func createDatabase() {
	db, err := sql.Open("sqlite3", SQLPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := "CREATE TABLE comments (comment_id TEXT PRIMARY KEY, user_id TEXT, match_count TEXT, comment TEXT);"
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	return
}

func deleteConfig() {
	os.Remove(YamlPath)
	return
}

func createConfig() {
	configText := `# This is a configuration file for fblament.
# More information about setting this configuration file can be found here: https://github.com/agajdosi/fblamentt#getting-the-tokens.

# Warning: all comments in this file will be lost after running the "fblament get" command. 

# ID of your fake application
clientID: ""

# Secret for application under which fblament will authenticate to Facebook
clientSecret: ""

# Your personal token which gives fblament ability to read information from Facebook under your account
accessToken: ""

# IDs of Facebook pages which you want to crawl
pages:
- 1531739643707341
- 777739285668748
- 828801157196779

# Minimum number of matched comments which are needed to save user into "results" folder
minimumLimit: 1

# Regular expressions or common strings for which it is searched in comments
regexps:
- "example of string to match"
- "example of (regular expression|regexp) to match"
`
	err := ioutil.WriteFile(YamlPath, []byte(configText), os.ModePerm)
	if err != nil {
		fmt.Println("Error creating the config file: ", err)
		return
	}
	fmt.Printf(`Configuration file with helpful comments created at %v. Please fill in required values to run "fblament get" command successfully.
More info about getting the token, clientID and clientSecret from Facebook can be found here: https://github.com/agajdosi/fblamentt#getting-the-tokens. 
`, YamlPath)
	return
}
