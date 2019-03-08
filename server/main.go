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
	go StartPinger()

	var PORT = os.Getenv("PORT")
	var IP string
	if IP = os.Getenv("IP"); IP == "" {
		IP = GetLocalIP()
	}
	flag.String("p", PORT, "port to serve on")
	flag.Parse()

	http.HandleFunc("/upload", upload)
	dir := http.Dir("/go-file-server/files")
	http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(dir)))

	fmt.Println("PORT: ", PORT)

	log.Printf("Serving %s on HTTP port: %s\n", dir, IP+":9090")

	// log.Fatal(http.ListenAndServe(IP+":9090", nil))
	log.Fatal(http.ListenAndServe(":"+PORT, nil))
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
