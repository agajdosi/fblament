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
	"strings"

	fb "github.com/huandu/facebook"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Gets comments from posts of defined pages and saves them into a database",
	Long:  `Gets comments from posts of defined pages and saves them into a database`,
	Run: func(cmd *cobra.Command, args []string) {
		refreshToken()
		commentCount := 0
		pages := Configuration["pages"].([]interface{})
		for _, page := range pages {
			pageID := ""
			switch v := page.(type) {
			case int:
				pageID = fmt.Sprintf("%d", page.(int))
				fmt.Printf("\nGetting comments for page: %v\n", page)
				getObjects(pageID, "posts", "comments", &commentCount)
			case string:
				pageID = page.(string)
				fmt.Println("\nGetting comments for page: %v\n", page)
				getObjects(pageID, "posts", "comments", &commentCount)
				fmt.Println()
			default:
				fmt.Printf("skipping page id=%v (id must be int or string)\n", v)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(getCmd)
}

func getObjects(ownerID, objectType, nextType string, count *int) {
	clientID := ""
	clientSecret := Configuration["clientSecret"].(string)
	switch t := Configuration["clientID"].(type) {
	case int:
		clientID = fmt.Sprintf("%d", Configuration["clientID"].(int))
	case string:
		clientID = Configuration["clientID"].(string)
	default:
		fmt.Println("Error - clientID must be int or string, is: %v", t)
	}
	var app = fb.New(clientID, clientSecret)
	session := app.Session(Configuration["accessToken"].(string))
	err := session.Validate()
	if err != nil {
		fmt.Println(err)
	}
	call := "/" + ownerID + "/" + objectType
	result, _ := session.Get(call, nil)
	if err != nil {
		fmt.Println("Error results:", err, result)
		return
	}
	paging, err := result.Paging(session)
	if err != nil {
		fmt.Println("Error paging:", err)
		return
	}

	noMore := false
	for !noMore {
		objects := paging.Data()
		for _, object := range objects {
			switch nextType {
			case "":
				*count++
				fmt.Printf("\rtotal number of comments saved: %d", *count)
				saveComment(object["id"].(string), object["from"].(map[string]interface{})["id"].(string), object["message"].(string))
			default:
				getObjects(object["id"].(string), nextType, "", count)
			}

		}
		noMore = !paging.HasNext()
		if !noMore {
			noMore, err = paging.Next()
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	return
}

func saveComment(commentID, userID, comment string) {
	db, err := sql.Open("sqlite3", SQLPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	comment = strings.Replace(comment, "'", "''", -1)
	sqlStmt := fmt.Sprintf("INSERT INTO comments(comment_id, user_id, comment) VALUES('%s', '%s', '%s');", commentID, userID, comment)
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
	}
}

func refreshToken() {
	res, err := fb.Get("/oauth/access_token", fb.Params{
		"grant_type":        "fb_exchange_token",
		"client_id":         Configuration["clientID"],
		"client_secret":     Configuration["clientSecret"],
		"fb_exchange_token": Configuration["accessToken"],
	})
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	Configuration["accessToken"] = res["access_token"]
	d, err := yaml.Marshal(&Configuration)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}
	ioutil.WriteFile(YamlPath, d, os.ModePerm)
	return
}
