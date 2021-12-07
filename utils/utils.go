package utils

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/azer/go-flickr"
)

var DownloadDir string
var Verbose bool

// SLog is a simple logger based on the verbose flag
func SLog(msg string) {
	if Verbose {
		log.Println(msg)
	}
}

func fileExists(path string) (bool, error) {
	if _, err := os.Stat(path); err == nil {
		return true, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else {
		// Schrodinger: file may or may not exist. See err for details.
		return false, err
	}
}

func downloadFile(fullURLFile, dir string, cache bool) (string, error) {
	const noFile = ""
	fileURL, err := url.Parse(fullURLFile)
	if err != nil {
		return noFile, err
	}
	path := fileURL.Path
	segments := strings.Split(path, "/")
	fileName := segments[len(segments)-1]
	fullPath := filepath.Join(dir, fileName)

	exists, err := fileExists(fullPath)
	if err != nil {
		return noFile, err
	}
	if cache && exists {
		SLog(fmt.Sprintf("File %s already exists, not downloading", fullPath))
		return fullPath, nil
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
func parseDir(path string) (string, error) {
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

func DownloadPhoto(
	client flickr.Client, p flickr.PhotoListItem, size, minSize int, cache bool,
) (file string, err error) {
	favId, err := strconv.Atoi(p.ID)
	if err != nil {
		return "", fmt.Errorf("error converting Favorite Id '%s' to integer: %s", p.ID, err)
	}

	sizes, err := client.GetPhotoSizes(favId)
	if err != nil {
		return "", fmt.Errorf("error getting photo sizes: %s", err)
	}

	url, err := sizes.ClosestWidthUrl(size, minSize)
	if err == flickr.ErrMinSizeNotAvailable {
		return "", err
	}
	if err != nil {
		return "", fmt.Errorf("error getting width based URL: %s", err)
	}

	scrubbedDir, err := parseDir(DownloadDir)
	if err != nil {
		return "", fmt.Errorf("error parsing Download Dir: %s", err)
	}

	file, err = downloadFile(url, scrubbedDir, true)
	if err != nil {
		return "", fmt.Errorf("unable to download URL %s : %s", url, err)
	}

	return
}
