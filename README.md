# Flickr Tools

A set of tools for interacting with the [Flickr](https://www.flickr.com/) photo sharing service.

## Tools

### Random Favorite

`flickr-random-fav` accepts a Flickr user name and downloads a random image from the favorites chosen by that user.

Usage:

- `-dir` string
  Directory of the downloaded image (default ".")
- `-env` string
  Name of an environment file that contains the FLICKR_API_KEY value. If omitted and no "key" argument specified, uses "./.env" for enviroment file.
- `-key` string
  [Flickr API Key](https://www.flickr.com/services/apps/create/). May be specified in an evironment file (see 'env' option)
- `-minsize` int
  Minimum acceptable long edge size if desired size is not available (default 2000)
- `-size` int
  Desired size of the long edge of the image. Resultant image may be larger if size does not exist. (default 2048)
- `-user` string
  Name of user for which to load favorites
- `-v` Verbose

### Download All Favorites

`flickr-download-all-favs` accepts a Flickr user name and downloads all images from the favorites chosen by that user.

Usage:

- `-dir` string
  Directory of the downloaded image (default ".")
- `-env` string
  Name of an environment file that contains the `FLICKR_API_KEY` value. If omitted and no "key" argument specified, uses "./.env" for enviroment file.
- `-fail`
  Fails if min desired size not available. Default is false
- `-key` string
  [Flickr API Key](https://www.flickr.com/services/apps/create/). May be specified in an evironment file (see 'env' option)
- `-minsize` int
  Minimum acceptable long edge size if desired size is not available (default 2000)
- `-size` int
  Desired size of the long edge of the image. Resultant image may be larger if size does not exist. (default 2048)
- `-user` string
  Name of user for which to load favorites
- `-v` Verbose

## Flickr API Key (Required)

In order to access the Flickr API, you must obtain a [Flickr API Key](https://www.flickr.com/services/apps/create/).  This key has to be know by all of the tools in this package.

Setting the API key can be done in a number of ways and is checked in this order:

1. Specifying with the `-key` command line parameter.
1. Specifying a file with the `-env` command line parameter.  This file should contain the key specification in the form `FLICKR_API_KEY=XXXXXXX`.
1. A file in the current directory named `.env` that contains the key as above.
1. An existing environment variable that contains the key as above.  Note that this overrides any env file specified above.
