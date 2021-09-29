package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	flickr "github.com/azer/go-flickr"
)

func downloadFile(fullURLFile string) (string, error) {
	const noFile = ""
	fileURL, err := url.Parse(fullURLFile)
	if err != nil {
		return noFile, err
	}
	path := fileURL.Path
	segments := strings.Split(path, "/")
	fileName := segments[len(segments)-1]

	// Create blank file
	file, err := os.Create(fileName)
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

	sLog(fmt.Sprintf("Downloaded file %s with size %d bytes", fileName, size))
	return fileName, nil
}

var (
	downloadDir string
	userName    string
	apiKey      string
	envFile     string
	verbose     bool
	size        int
)

func setupFlags() {
	flag.StringVar(&downloadDir, "dir", ".", "Directory of the downloaded image")
	flag.StringVar(&userName, "user", "", "Name of user for which to load favorites")
	flag.StringVar(&apiKey, "key", "", "Flickr API Key (see https://www.flickr.com/services/apps/create/).  "+
		"May be specified in an evironment file (see 'env' option)")
	flag.StringVar(&envFile, "env", "", fmt.Sprintf(
		"Name of an environment file that contains the %s value.  "+
			"If omitted and no \"key\" argument specified, uses \"./.env\" for enviroment file.",
		flickr.ApiKeyEnvVar))
	flag.BoolVar(&verbose, "v", false, "Verbose")
	flag.IntVar(&size, "size", 2048, "Desired size of the long edge of the image.  "+
		"Resultant image may be larger if size does not exist.")
	flag.Parse()
}

// sLog is a simple logger based on the verbose flag
func sLog(msg string) {
	if verbose {
		log.Println(msg)
	}
}

func main() {
	setupFlags()
	rand.Seed(time.Now().UnixNano())

	client := flickr.NewClient(apiKey, envFile)

	user, err := client.FindUser(userName)
	if err != nil {
		log.Fatalf("Error finding user %s: %s", userName, err)
	} else {
		sLog(fmt.Sprintf("ID is %s", user.Id))
	}

	fav, err := client.RandomFav(user.Id)
	if err != nil {
		log.Fatalf("Unable to get random fav: %s", err)
	}
	sLog(fmt.Sprintf("Got \"%s\"", fav.Title))

	favId, err := strconv.Atoi(fav.Id)
	if err != nil {
		log.Fatalf("Error converting Favorite Id '%s' to integer: %s", fav.Id, err)
	}
	sizes, err := client.GetPhotoSizes(favId)
	if err != nil {
		log.Fatalf("Error getting photo sizes: %s", err)
	}

	url, err := sizes.ClosestWidthUrl(size)
	if err != nil {
		log.Fatalf("Error getting width based URL: %s", err)
	}

	var file string
	file, err = downloadFile(url)
	if err != nil {
		log.Fatalf("Unable to download URL %s : %s", url, err)
	}
	fmt.Println(file)
}
