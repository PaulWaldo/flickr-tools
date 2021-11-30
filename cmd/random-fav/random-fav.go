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

func randomFav(client *flickr.PhotosClient, userId string) (flickr.PhotoListItem, error) {
	const allRightReserved = "0"
	// Get the favs list metadata
	client.PerPage = 1
	favs, err := client.Favs(userId)
	if err != nil {
		return flickr.PhotoListItem{}, err
	}
	client.PerPage = 1

	// Loop through random Favs, looking for ones which are not restricted, i.e.
	// license value != 0 ("All Rights Reserved")}
	found := false
	// var offset int
	for !found {
		photoNum := rand.Intn(favs.Total)
		// var page int
		page, offset := utils.DivMod(photoNum, client.PerPage)
		page += 1 // Account for API pages starting at 1
		// Get specified page
		client.Page = page
		favs, err := client.Favs(userId)
		if err != nil {
			return flickr.PhotoListItem{}, err
		}

		if favs.Photos[offset].License != allRightReserved {
			return favs.Photos[offset], nil
		} else {
			utils.SLog(fmt.Sprintf("Photo %d is restricted, trying another", photoNum))
		}
	}
	return flickr.PhotoListItem{}, fmt.Errorf("unable to find a fav")
}

func main() {
	setupFlags()
	rand.Seed(time.Now().UnixNano())

	client,err := flickr.NewClient()
	if err != nil {
		log.Fatalf("Error creating client: %s", err)
	}

	user, err := client.FindUser(userName)
	if err != nil {
		log.Fatalf("Error finding user %s: %s", userName, err)
	} else {
		utils.SLog(fmt.Sprintf("ID is %s", user.Id))
	}

	paginatedClient,err := flickr.NewPhotosClient()
	if err != nil {
		log.Fatalf("Error creating photos client: %s", err)
	}

	// paginatedClient.Cache = true
	fav, err := randomFav(paginatedClient, user.Id)
	if err != nil {
		log.Fatalf("Unable to get random fav: %s", err)
	}
	utils.SLog(fmt.Sprintf("Got \"%s\"", fav.Title))

	favId, err := strconv.Atoi(fav.ID)
	if err != nil {
		log.Fatalf("Error converting Favorite Id '%s' to integer: %s", fav.ID, err)
	}

	sizes, err := client.GetPhotoSizes(favId)
	if err != nil {
		log.Fatalf("Error getting photo sizes: %s", err)
	}

	url, err := sizes.ClosestWidthUrl(size)
	if err != nil {
		log.Fatalf("Error getting width based URL: %s", err)
	}

	scrubbedDir, err := utils.ParseDir(utils.DownloadDir)
	if err != nil {
		log.Fatal(err)
	}
	
	file, err := utils.DownloadFile(url, scrubbedDir)
	if err != nil {
		log.Fatalf("Unable to download URL %s : %s", url, err)
	}
	fmt.Println(file)
}
