/*
Serve is a very simple static file server in go
Usage:
	-p="8100": port to serve on
	-d=".":    the directory of static files to host

Navigating to http://localhost:8100 will display the index.html or directory
listing file.
*/
package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func main() {
	// go StartPinger()

	name, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	fmt.Println("hostname:", name)

	var PORT = os.Getenv("PORT")
	var IP string
	if IP = os.Getenv("IP"); IP == "" {
		IP = GetLocalIP()
	}
	flag.String("p", PORT, "port to serve on")
	flag.Parse()

	http.HandleFunc("/upload", upload)

	var dir http.Dir
	if os.Getenv("ISDOCKER") == "true" {
		dir = http.Dir("/go-file-server/files")
	} else {
		dir = http.Dir("./go-file-server/files")
	}

	http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(dir)))

	log.Printf("Serving %s on HTTP port:", "oc-appsrv01.informatik.uni-augsburg.de:8081")

	// log.Fatal(http.ListenAndServe(IP+":9090", nil))
	//local listener
	log.Fatal(http.ListenAndServe("oc-appsrv01.informatik.uni-augsburg.de:8081", nil))

	//docker listener
	// log.Fatal(http.ListenAndServe(":"+PORT, nil))
}

func GetLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	ipStr := fmt.Sprintf("%s", localAddr)
	host, port, err := net.SplitHostPort(ipStr)
	fmt.Println(host, port)

	return host
}

// upload logic
func upload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("upload.gtpl")
		t.Execute(w, token)
	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Fprintf(w, "%v", handler.Header)
		fmt.Println("handler.Filename: ", handler.Filename)
		filpath := "/go-file-server/files/"
		filesDir := filepath.Join(".", filpath)
		fmt.Println("filesDir: ", filesDir)

		if _, err := os.Stat(filesDir); os.IsNotExist(err) {
			fmt.Println(filesDir)
			err = os.MkdirAll(filesDir, os.ModePerm)
			if err != nil {
				panic(err)
			}
		}

		var f *os.File
		if os.Getenv("ISDOCKER") == "true" {
			f, err = os.Create(filepath.Join(filpath, handler.Filename))
		} else {
			f, err = os.Create(filepath.Join(filesDir, handler.Filename))
		}

		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
	}
}
