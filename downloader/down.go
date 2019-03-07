package main

import (
	"flag"
	"io"
	"net/http"
	"os"
	"path"
)

func main() {
	var fFlag = flag.String("f", "cc.json", "get file")
	flag.Parse()
	fileURL := "http://192.168.178.34:9090/files/" + *fFlag
	if err := DownloadFile(fileURL); err != nil {
		panic(err)
	}
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(path.Base(url))
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
