package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
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

func divMod(numerator, denominator int) (q, r int) {
	q = numerator / denominator
	r = numerator % denominator
	return
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

func randomFav(client *flickr.PaginatedClient, userId string) (flickr.Fav, error) {
	const allRightReserved = "0"
	// Get the favs list metadata
	client.NumPerPage = 1
	if _, err := client.Favs(userId); err != nil {
		return flickr.Fav{}, err
	}
	client.NumPerPage = 100

	// Loop through random Favs, looking for ones which are not restricted, i.e.
	// license value != 0 ("All Rights Reserved")}
	found := false
	// var offset int
	for !found {
		photoNum := rand.Intn(client.Total)
		// var page int
		page, offset := divMod(photoNum, client.RequestNumPerPage)
		page += 1 // Account for API pages starting at 1
		// Get specified page
		client.Page = client.NumPages // page
		if client.Page > client.NumPages {
			panic("Whoa!")
		}
		favs, err := client.Favs(userId)
		if err != nil {
			return flickr.Fav{}, err
		}

		if favs[offset].License != allRightReserved {
			return favs[offset], nil
		}
	}
	return flickr.Fav{}, fmt.Errorf("unable to find a fav")
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

	paginatedClient := flickr.NewDefaultPaginatedClient(apiKey, envFile)
	paginatedClient.Cache = true
	fav, err := randomFav(&paginatedClient, user.Id)
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
	file, err = utils.DownloadFile(url, utils.DownloadDir)
	if err != nil {
		log.Fatalf("Unable to download URL %s : %s", url, err)
	}
	fmt.Println(file)
}
