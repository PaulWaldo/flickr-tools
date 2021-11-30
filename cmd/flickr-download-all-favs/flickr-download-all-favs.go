package main

import (
	"flag"
	"fmt"
	"log"

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

	client, err := flickr.NewClient()
	if err != nil {
		log.Fatalf("Unable to create client: %s", err)
	}

	user, err := client.FindUser(userName)
	if err != nil {
		log.Fatalf("Error finding user %s: %s", userName, err)
	} else {
		utils.SLog(fmt.Sprintf("ID is %s", user.Id))
	}

	paginatedClient, err := flickr.NewPhotosClient()
	if err != nil {
		log.Fatalf("Unable to create photos client: %s", err)
	}

	favs, err := paginatedClient.Favs(user.Id)
	if err != nil {
		log.Fatalf("Error getting Favs: %s", err)
	}
	utils.SLog(fmt.Sprintf("Page %d/%d: %d items", favs.Page, favs.Pages, len(favs.Photos)))

	// var favs []flickr.Fav
	for err == nil {
		favs, err = paginatedClient.NextPage()
		if err == flickr.ErrPaginatorExhausted {
			break
		}
		if err != nil {
			log.Fatalf("Error getting Favs: %s", err)
		}
		utils.SLog(fmt.Sprintf("Page %d/%d: %d items", favs.Page, favs.Pages, len(favs.Photos)))
		// fmt.Printf("%v\n", favs)
		for _, fav:=range favs.Photos {
			utils.SLog(fmt.Sprintf("Would have downloaded %s", fav.Title))
			// utils.DownloadFile(fav.)
		}
	}
	if err != flickr.ErrPaginatorExhausted {
		log.Fatalf("Error getting favs: %s", err)
	}
}
