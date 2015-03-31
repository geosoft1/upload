// upload service
// Copyright (C) 2014  geosoft1@gmail.com
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Page struct {
	Status template.HTML
}

var t *template.Template

func checkSecrets(username, password string) (err error) {
	var usernameFromFile, passwordFromFile string
	var n int

	file, err := os.Open("secrets")
	if err != nil {
		return err
	}
	for {
		n, _ = fmt.Fscanf(file, "%s\t%s", &usernameFromFile, &passwordFromFile)
		if username == usernameFromFile && password == passwordFromFile {
			return nil
		}
		if n == 0 {
			break
		}
	}
	err = errors.New("secrets: wrong username or password " + username)
	return err
}

func index(w http.ResponseWriter, r *http.Request) {
	//create|open log file
	f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
	}

	p := Page{}
	defer func() {
		t.ExecuteTemplate(w, "index.html", p)
	}()

	if r.Method == "POST" {
		err := checkSecrets(r.FormValue("username"), r.FormValue("password"))

		defer func() {
			p.Status = template.HTML(err.Error())
			//write to log file
			fmt.Fprintln(f, time.Now().Format("Mon, 02 Jan 2006 15:04 ")+err.Error())
		}()

		if err != nil {
			return
		}

		file, handler, err := r.FormFile("file")
		if err != nil {
			return
		}
		out, err := os.Create("files/" + handler.Filename)
		if err != nil {
			return
		}
		bytes, err := io.Copy(out, file)
		if err != nil {
			return
		}

		//seems to be a good ideea to generate download link into an error message
		//to treat this and previous possible errors in one defer code
		err = errors.New(fmt.Sprintf("<a href=/files/%s>%s</a> %db transferred by user %s",
			handler.Filename, handler.Filename, bytes, r.FormValue("username")))
	}
}

func main() {
	os.Chdir(filepath.Dir(os.Args[0]))
	t, _ = template.ParseFiles("index.html")
	//web server
	http.HandleFunc("/", index)
	//file server
	http.Handle("/files/", http.FileServer(http.Dir("./")))
	http.ListenAndServe(":8080", nil)
}
