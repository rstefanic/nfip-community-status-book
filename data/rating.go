package data

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/tealeg/xlsx/v3"
)

const NFIPCommunityRatingSystemFilename = "crs.xlsx"
const NFIPCommunityRatingSystemURL = "https://www.fema.gov/sites/default/files/2020-08/fema_crs_eligible-communities_oct-2020.xlsx"
const CRSSheetName = "Sheet1"

const (
	RatingState = iota
	RatingCommunityNumber
	RatingCommunityName
	RatingCRSEntryDate
	RatingCurrentEffectiveDate
	RatingCurrentClass
	RatingDiscountForSFHA
	RatingDiscountForNonSFHA
	RatingStatus
)

type NFIPCommunityRatings []NFIPCommunityRating

type NFIPCommunityRating struct {
	State                string `json:"state"`
	CommunityNumber      string `json:"community_number"`
	CommunityName        string `json:"community_name"`
	CRSEntryDate         string `json:"crs_entry_date"`
	CurrentEffectiveDate string `json:"current_effective_date"`
	CurrentClass         string `json:"current_class"`
	DiscountForSFHA      string `json:"discount_for_sfha"`
	DiscountForNonSFHA   string `json:"discount_for_non_sfha"`
	Status               string `json:"status"`
}

func GetNFIPCommunityRatingSystem(l *log.Logger) (NFIPCommunityRatings, error) {
	// First check if the file exists and download it if it does not exist
	if _, err := os.Stat(NFIPCommunityRatingSystemFilename); os.IsNotExist(err) {
		l.Println("NFIP CRS does not exist. Downloading...")
		resp, err := http.Get(NFIPCommunityRatingSystemURL)

		if err != nil {
			l.Println("** Err - ", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		f, err := os.Create(NFIPCommunityRatingSystemFilename)
		if err != nil {
			l.Println("** Err - ", err)
			os.Exit(1)
		}

		io.Copy(f, resp.Body)
		defer f.Close()
	}

	wb, err := xlsx.OpenFile(NFIPCommunityRatingSystemFilename)
	if err != nil {
		l.Println("** Err -", err)
		os.Exit(1)
	}

	// Get the worksheeet that contains the
	// data we're interested in from the xlsx.
	crsSheet, ok := wb.Sheet[CRSSheetName]
	if !ok {
		l.Println("** Err -", err)
		os.Exit(1)
	}

	var crs NFIPCommunityRatings

	// Once we have the sheet, we want to loop through
	// each row and make a NFIPCommunityRating out of it.
	err = crsSheet.ForEachRow(func(r *xlsx.Row) error {
		rowNumber := r.GetCoordinate()

		// We skip the first 5 rows because it's just the
		// header and extra and metadata about the NFIP CRS.
		if rowNumber < 5 {
			return nil
		}

		var cr NFIPCommunityRating
		cr.State = getFormattedCellValue(r, RatingState)
		cr.CommunityNumber = getFormattedCellValue(r, RatingCommunityNumber)
		cr.CommunityName = getFormattedCellValue(r, RatingCommunityName)
		cr.CRSEntryDate = getFormattedCellValue(r, StatusCRSEntryDate)
		cr.CurrentEffectiveDate = getFormattedCellValue(r, RatingCurrentEffectiveDate)
		cr.CurrentClass = getFormattedCellValue(r, RatingCurrentClass)
		cr.DiscountForSFHA = getFormattedCellValue(r, RatingDiscountForSFHA)
		cr.DiscountForNonSFHA = getFormattedCellValue(r, RatingDiscountForNonSFHA)
		cr.Status = getFormattedCellValue(r, RatingStatus)

		crs = append(crs, cr)
		return nil
	})

	return crs, err
}

func (crs NFIPCommunityRatings) Search(term string) *NFIPCommunityRatings {
	var matchingCommunities NFIPCommunityRatings
	term = strings.ToLower(term)

	for _, cr := range crs {
		if strings.Contains(strings.ToLower(cr.State), term) ||
			strings.Contains(strings.ToLower(cr.CommunityNumber), term) ||
			strings.Contains(strings.ToLower(cr.CommunityName), term) {
			matchingCommunities = append(matchingCommunities, cr)
		}
	}

	return &matchingCommunities
}

func (crs *NFIPCommunityRatings) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(crs)
}

func getFormattedCellValue(r *xlsx.Row, pos int) string {
	cell := r.GetCell(pos)
	fv, err := cell.FormattedValue()

	// If there is a problem formatting the value from
	// this cell, then we'll ignore it and move on.
	if err != nil {
		fv = ""
	}
	return fv
}
