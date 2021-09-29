package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/PaulWaldo/flickr-tools/utils"
	flickr "github.com/azer/go-flickr"
)

var (
	userName string
	apiKey   string
	envFile  string
	size     int
)

func setupFlags() {
	flag.StringVar(&utils.DownloadDir, "dir", ".", "Directory of the downloaded image")
	flag.StringVar(&userName, "user", "", "Name of user for which to load favorites")
	flag.StringVar(&apiKey, "key", "", "Flickr API Key (see https://www.flickr.com/services/apps/create/).  "+
		"May be specified in an evironment file (see 'env' option)")
	flag.StringVar(&envFile, "env", "", fmt.Sprintf(
		"Name of an environment file that contains the %s value.  "+
			"If omitted and no \"key\" argument specified, uses \"./.env\" for enviroment file.",
		flickr.ApiKeyEnvVar))
	flag.BoolVar(&utils.Verbose, "v", false, "Verbose")
	flag.IntVar(&size, "size", 2048, "Desired size of the long edge of the image.  "+
		"Resultant image may be larger if size does not exist.")
	flag.Parse()
}

func main() {
	setupFlags()
	rand.Seed(time.Now().UnixNano())

	client := flickr.NewClient(apiKey, envFile)

	user, err := client.FindUser(userName)
	if err != nil {
		log.Fatalf("Error finding user %s: %s", userName, err)
	} else {
		utils.SLog(fmt.Sprintf("ID is %s", user.Id))
	}

	fav, err := client.RandomFav(user.Id)
	if err != nil {
		log.Fatalf("Unable to get random fav: %s", err)
	}
	utils.SLog(fmt.Sprintf("Got \"%s\"", fav.Title))

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
	file, err = utils.DownloadFile(url)
	if err != nil {
		log.Fatalf("Unable to download URL %s : %s", url, err)
	}
	fmt.Println(file)
}
