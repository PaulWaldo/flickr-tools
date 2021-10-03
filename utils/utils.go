package utils

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
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

	client := http.Client{}
	resp, err := client.Get(fullURLFile)
	if err != nil {
		return noFile, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return noFile, fmt.Errorf("unable to download %s : got response %s", fullURLFile, resp.Status)
	}

	fullPath := filepath.Join(dir, fileName)
	file, err := os.Create(fullPath)
	if err != nil {
		return noFile, err
	}

	size, err := io.Copy(file, resp.Body)
	if err != nil {
		return noFile, err
	}
	defer file.Close()

	SLog(fmt.Sprintf("Downloaded file %s with size %d bytes", fileName, size))
	return fullPath, nil
}

// parseDir parses a user-entered directory and properly formats it
func ParseDir(path string) (string, error) {
	// Try tilde exapnsion
	var newPath string
	if strings.Contains(path, "~") {
		usr, err := user.Current()
		if err != nil {
			return "", err
		}
		homeDir := usr.HomeDir
		if path == "~" {
			newPath = homeDir
		} else if strings.HasPrefix(path, "~/") {
			newPath = filepath.Join(homeDir, path[2:])
		}
	} else {
		newPath = filepath.Clean(path)
	}

	_, err := os.Stat(newPath)
	if err != nil {
		return "", err
	}

	return newPath, nil
}

func DivMod(numerator, denominator int) (q, r int) {
	q = numerator / denominator
	r = numerator % denominator
	return
}

