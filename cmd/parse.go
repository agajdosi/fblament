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
	"log"
	"os"
	"path/filepath"
	"reflect"
	"regexp"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

type match struct {
	commentID string
	comment   string
}

// parseCmd represents the parse command
var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "Parses the information in database and creates a lament",
	Long: `Parses the information in database and creates a lament
Lament files for individual users will be stored in .fblament/report`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Matching words of users' comments...")
		expressions := Configuration["regexps"].([]interface{})
		fmt.Println(expressions, reflect.TypeOf(expressions))
		var regExps []*regexp.Regexp
		for _, expression := range expressions {
			fmt.Println(expression, reflect.TypeOf(expression))
			compiled := regexp.MustCompile(expression.(string))
			regExps = append(regExps, compiled)
		}

		getUsers()

	},
}

func init() {
	RootCmd.AddCommand(parseCmd)
}

func getUsers() {
	usersParsed := 0
	db, err := sql.Open("sqlite3", SQLPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT user_id FROM comments GROUP BY user_id")
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var userID string
		err := rows.Scan(&userID)
		if err != nil {
			log.Fatal("error reading user_id:", err)
		}
		getUserComments(userID)
		usersParsed++
		fmt.Printf("\rusers parsed: %d", usersParsed)
	}
	fmt.Println()

}

func getUserComments(userID string) {
	var matches []match
	db, err := sql.Open("sqlite3", SQLPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	search := fmt.Sprintf("SELECT comment_id, comment FROM comments WHERE user_id IS '%s';", userID)
	rows, err := db.Query(search)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var commentID, comment string
		err := rows.Scan(&commentID, &comment)
		if err != nil {
			log.Fatal("v pici :(", err)
		}
		//HERE WE ARE GOING TO SEARCH IF THE COMMENT MATCHES WORDS/REGEXPS
		matches = append(matches, match{commentID: commentID, comment: comment})
	}
	//SAVE RESULTS INTO A FILE
	if len(matches) >= Configuration["minimumLimit"].(int) {
		saveResults(userID, matches)
	}
}

func saveResults(userID string, matches []match) {
	outputFile := filepath.Join(OutputFolderPath, userID+".txt")
	file, err := os.Create(outputFile)
	if err != nil {
		fmt.Println("error creating file:", err)
	}
	defer file.Close()

	//HEADER
	str := "Results for user with ID: " + userID + "\n\n\n"
	_, err = file.WriteString(str)
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}

	//COMMENTS
	for _, match := range matches {
		str = match.commentID + ":\n" + match.comment + "\n\n"
		_, err = file.WriteString(str)
		if err != nil {
			fmt.Println("Error writing to file:", err)
		}
	}
}
