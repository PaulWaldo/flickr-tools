package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/PaulWaldo/flickr-tools/utils"
	flickr "github.com/azer/go-flickr"
)

var (
	userName         string
	apiKey           string
	envFile          string
	size             int
	minSize          int
	failNotAvailable bool
)

func setupFlags() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "%s accepts a Flickr user name and downloads all images from the favorites chosen by that user.\nUsage:\n", os.Args[0])
		flag.PrintDefaults()
	}

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
	flag.IntVar(&minSize, "minsize", 2000, "Minimum acceptable long edge size if desired size is not available")
	flag.BoolVar(&failNotAvailable, "fail", false, "Fails if min desired size not available.  Default is false")
	flag.Parse()
}

func main() {
	setupFlags()
	ch := make(chan flickr.PhotoListItem, 100)

	var client *flickr.Client

	if len(apiKey) > 0 {
		client = flickr.NewClientApiKey(apiKey)
	} else {
		var err error
		client, err = flickr.NewClient()
		if err != nil {
			log.Fatalf("Unable to create client: %s", err)
		}
	}

	user, err := client.FindUser(userName)
	if err != nil {
		log.Fatalf("Error finding user %s: %s", userName, err)
	}
	utils.SLog(fmt.Sprintf("ID is %s", user.Id))

	var paginatedClient *flickr.PhotosClient
	if len(apiKey) > 0 {
		paginatedClient = flickr.NewPhotosClientApiKey(apiKey)
	} else {
		paginatedClient, err = flickr.NewPhotosClient()
		if err != nil {
			log.Fatalf("Unable to create photos client: %s", err)
		}
	}

	favs, err := paginatedClient.Favs(user.Id)
	if err != nil {
		log.Fatalf("Error getting Favs: %s", err)
	}

	go func() {
		defer close(ch)
		for err == nil {
			utils.SLog(fmt.Sprintf("Page %d/%d: %d items", favs.Page, favs.Pages, len(favs.Photos)))
			for _, fav := range favs.Photos {
				ch <- fav
			}
			favs, err = paginatedClient.NextPage()
			if err == flickr.ErrPaginatorExhausted {
				break
			}
			if err != nil {
				log.Fatalf("Error getting Favs: %s", err)
			}
		}
	}()

	var (
		downloaded          = 0
		minSizeNotAvailable = 0
		errors              = 0
	)
	for f := range ch {
		go func(p flickr.PhotoListItem) {
			utils.SLog(fmt.Sprintf("Getting title '%s'", p.Title))
			filename, err := utils.DownloadPhoto(*client, p, size, minSize, true)
			if err == flickr.ErrMinSizeNotAvailable {
				utils.SLog(err.Error())
				minSizeNotAvailable++
				return
			}
			if err != nil {
				utils.SLog(fmt.Sprintf("error downloading photo: %s", err))
				errors++
				return
			}
			utils.SLog(fmt.Sprintf("Downloaded file %s", filename))
			downloaded++
		}(f)
	}
	fmt.Printf("Downloaded %d\n", downloaded)
	fmt.Printf("Minimum size not available %d\n", minSizeNotAvailable)
	fmt.Printf("Error downloading %d\n", errors)
	fmt.Println("Downloader done!")
}
