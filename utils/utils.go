package utils

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var DownloadDir string
var Verbose bool

// SLog is a simple logger based on the verbose flag
func SLog(msg string) {
	if Verbose {
		log.Println(msg)
	}
}

func DownloadFile(fullURLFile string, dir string) (string, error) {
	const noFile = ""
	fileURL, err := url.Parse(fullURLFile)
	if err != nil {
		return noFile, err
	}
	path := fileURL.Path
	segments := strings.Split(path, "/")
	fileName := segments[len(segments)-1]

	// Create blank file
	file, err := os.Create(fmt.Sprintf("%s/%s", dir, fileName))
	if err != nil {
		return noFile, err
	}

	client := http.Client{}
	resp, err := client.Get(fullURLFile)
	if err != nil {
		return noFile, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return noFile, fmt.Errorf("unable to download %s : got response %s", fullURLFile, resp.Status)
	}

	size, err := io.Copy(file, resp.Body)
	if err != nil {
		return noFile, err
	}
	defer file.Close()

	SLog(fmt.Sprintf("Downloaded file %s with size %d bytes", fileName, size))
	return fileName, nil
}
