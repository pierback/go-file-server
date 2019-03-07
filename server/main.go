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
	"strconv"
	"time"
)

func main() {
	var PORT = os.Getenv("PORT")
	var IP string
	if IP = os.Getenv("IP"); IP == "" {
		ip := GetLocalIP()
		ipStr := fmt.Sprintf("%s", ip)
		IP = ipStr
	}
	flag.String("p", PORT, "port to serve on")
	directory := flag.String("d", ".", "the directory of static file to host")
	flag.Parse()

	http.Handle("/", http.FileServer(http.Dir(*directory)))
	// http.Handle("/go-file-server/files", http.FileServer(http.Dir("./go-file-server/files")))
	http.HandleFunc("/upload", upload)

	fs := http.FileServer(http.Dir("/go-file-server/files"))

	http.Handle("/files/", http.StripPrefix("/files/", fs))

	fmt.Println("PORT: ", PORT)

	log.Printf("Serving %s on HTTP port: %s\n", *directory, PORT)

	log.Fatal(http.ListenAndServe(":"+PORT, nil))
}

func GetLocalIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
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
		if _, err := os.Stat("/go-file-server/files/"); os.IsNotExist(err) {
			fmt.Println("// path/to/whatever does not exist")
		}
		f, err := os.Create("/go-file-server/files/" + handler.Filename)
		// f, err := os.OpenFile(handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
	}
}
