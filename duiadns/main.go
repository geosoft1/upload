// duia (Dynamic Updates for Internet Addressing) client
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
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"time"
)

const UserAgent = "duia-golang-1.0.0.3"

//change debug to false in real life
//in debug mode the update request is not send to server but shown
const debug = false

func getIpFromSite(version int) (s string, err error) {
	//get my ip from server (ipv4/ipv6 compatible)
	req, err := http.Get("http://" + "ipv" + strconv.Itoa(version) + ".duia.ro")
	if err != nil {
		println("no ipv" + strconv.Itoa(version) + " connection")
		return "", err
	}
	ip, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return "", err
	}
	//println("get ip" + strconv.Itoa(version) + " " + string(ip))
	return string(ip), nil
}

func updateDNS(version int, host, password, ip string) (err error) {
	if debug {
		//line commented in real life. uncomment to simulate DNS update
		println("http://ipv" + strconv.Itoa(version) +
			".duia.ro/dynamic.duia?host=" + host +
			"&password=" + password +
			"&ip" + strconv.Itoa(version) + "=" + ip)

	} else {
		//connect to duia server and update ip
		//http://stackoverflow.com/questions/13263492/set-useragent-in-http-request
		client := &http.Client{}
		req, err := http.NewRequest(
			"GET",
			"http://ipv"+strconv.Itoa(version)+
				".duia.ro/dynamic.duia?host="+host+
				"&password="+password+
				"&ip"+strconv.Itoa(version)+"="+ip, nil)
		if err != nil {
			return err
		}
		//tested with http://httpbin.org
		req.Header.Set("User-Agent", UserAgent)
		resp, err := client.Do(req)
		if err != nil {
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
		}
		//TODO: parse response
		println(string(body))

		// duia server update response
		// TODO: parse for a nice success message
		// BUG: inconsistent response message from duia.ro
		//<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01//EN">
		//<html>
		//<head><title>DUIA Module</title></head>
		//<ul><li>Hostname geosoft1.duia.ro</li>
		//<ul><li>Ipv4 79.113.171.69</li></ul>
		//</ul></body></html>
	}
	return nil
}

func readCache() (ip4, ip6 string, err error) {
	path, _ := os.Getwd()
	file, err := os.Open(path + string(filepath.Separator) + "duia.cache")
	if err != nil {
		return "", "", err
	}
	fmt.Fscanf(file, "%s\t%s", &ip4, &ip6)
	return ip4, ip6, nil
}

func updateCache(ip4, ip6 string) (err error) {
	path, _ := os.Getwd()
	file, err := os.Create(path + string(filepath.Separator) + "duia.cache")
	if err != nil {
		return err
	}
	// write creditentials in unix clasic style format
	fmt.Fprintf(file, "%s\t%s", ip4, ip6)
	return nil
}

func getSecrets() (host, password string, err error) {
	path, _ := os.Getwd()
	// get creditentials from file
	file, err := os.Open(path + string(filepath.Separator) + "duia.cfg")
	if err != nil {
		return "", "", err
	}
	fmt.Fscanf(file, "%s\t%s", &host, &password)
	return host, password, nil
}

func updateSecrets(host, password string) {
	// md5 password encoding
	h := md5.New()
	io.WriteString(h, password)
	md5 := fmt.Sprintf("%x", h.Sum(nil))
	// create creditentials file
	path, _ := os.Getwd()
	file, _ := os.Create(path + string(filepath.Separator) + "duia.cfg")
	// write creditentials in unix clasic style format
	fmt.Fprintf(file, "%s\t%s", host, md5)
}

func main() {
	// get current path
	file, _ := exec.LookPath(os.Args[0])
	dir, _ := path.Split(file)
	// important!
	// change to curent path to avoid problems when the program
	// is launch from other location  without directory changed
	os.Chdir(dir)
	fmt.Println(dir)

	// if first run, get creditentials from stdin
	host, password, err := getSecrets()
	if err != nil {
		fmt.Print("Host: ")
		fmt.Scan(&host)
		fmt.Print("Password: ")
		fmt.Scan(&password)
		updateSecrets(host, password)
		//programs must be run again after that
		return
	}

	// ticker make easy  to use  on various platforms
	// without need to use cron or similar schedulers
	// please don't use a value < 60 for ticker or your DuiaDNS account will be automatically disabled
	ticker := time.NewTicker(time.Second * 60)
	for t := range ticker.C {
		needUpdate := false

		ip4, ip6, _ := readCache()
		fmt.Printf("cache checked: %02d:%02d.%02d\n", t.Hour(), t.Minute(), t.Second())

		//ipv4 support
		ip4FromSite, err := getIpFromSite(4)
		if err == nil {
			// check was changed ip4 and update to the site
			if ip4FromSite != ip4 {
				updateDNS(4, host, password, ip4FromSite)
				ip4 = ip4FromSite
				needUpdate = true
			}
		}

		//ipv6 support
		ip6FromSite, err := getIpFromSite(6)
		if err == nil {
			// check was changed ip6 and update to the site
			if ip6FromSite != ip6 {
				updateDNS(6, host, password, ip6FromSite)
				ip6 = ip6FromSite
				needUpdate = true
			}
		}

		//finally update the cache
		if needUpdate {
			updateCache(ip4, ip6)
			fmt.Printf("cache updated: %02d:%02d.%02d\n", t.Hour(), t.Minute(), t.Second())
		}
	}
}
