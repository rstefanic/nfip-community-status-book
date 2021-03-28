# NFIP Community Status Book Service

A small service to search through [FEMA's NFIP Community Status Book](https://www.fema.gov/flood-insurance/work-with-nfip/community-status-book). The Community Status Book is downloaded from FEMA's site on start up if it doesn't exist locally.

Once the service is ran, make a GET request to `/search?term=<search_term>` to search by CID, Community Name, or County. Results are returned in JSON.

## Installation

Docker:
```shell
docker container run --rm -p 9001:9001 rstefanic/nfip-community-status-book:1.0
```

Build and serve:
```shell
go run main.go
```