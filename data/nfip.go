package data

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var NFIPCommunityBookFileName = "nation.csv"

const (
	POSCID = iota
	POSCommunityName
	POSCounty
	POSFHBMIdentified
	POSFIRMIdentified
	POSCurrEffMapDate
	POSRegEmerDate
	POSTribal
	POSCRSEntryDate
	POSCurrEffDate
	POSCurClass
	POSPercentDiscSFHA
	POSPercentNonSFHA
	POSProgram
	POSParticipatingCommunity
)

type NFIPCommunities []NFIPCommunity

type NFIPCommunity struct {
	CID                    int       `json:"cid"`
	CommunityName          string    `json:"community_name"`
	County                 string    `json:"county"`
	FHBMIdentified         time.Time `json:"fhbm_identified"`
	FIRMIdentified         time.Time `json:"firm_identified"`
	CurrEffMapDate         time.Time `json:"curr_eff_map_date"`
	RegEmerDate            time.Time `json:"reg_emer_date"`
	Tribal                 bool      `json:"tribal"`
	CRSEntryDate           string    `json:"crs_entry_date"`
	CurrEffDate            string    `json:"curr_eff_date"`
	CurClass               string    `json:"cur_class"`
	PercentDiscSFHA        string    `json:"percent_disc_sfha"`
	PercentNonSFHA         string    `json:"percent_non_sfha"`
	Program                string    `json:"program"`
	ParticipatingCommunity bool      `json:"participating_community"`
}

var ErrEmptyString = fmt.Errorf("string is empty")
var ErrInvalidDateString = fmt.Errorf("invalid date string")

func GetNFIPCommunityBook(l *log.Logger) (NFIPCommunities, error) {
	if _, err := os.Stat(NFIPCommunityBookFileName); os.IsNotExist(err) {
		l.Println("NFIP Community book does not exist. Downloading...")
		resp, err := http.Get("https://www.fema.gov/cis/nation.csv")

		if err != nil {
			l.Println("** Err -", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		f, err := os.Create(NFIPCommunityBookFileName)

		if err != nil {
			l.Println("** Err -", err)
			os.Exit(1)
		}

		io.Copy(f, resp.Body)
		defer f.Close()
	}

	f, err := os.Open(NFIPCommunityBookFileName)

	if err != nil {
		l.Println("Could not open NFIP Community book")
		os.Exit(1)
	}

	defer f.Close()

	csvReader := csv.NewReader(f)
	csvReader.LazyQuotes = true
	communities, err := unmarshal(csvReader)

	if err != nil {
		return nil, fmt.Errorf("could not parse NFIP Community book CSV File. Reason: %s", err.Error())
	}

	return communities, nil
}

func (c NFIPCommunities) Search(term string) *NFIPCommunities {
	var matchingCommunities NFIPCommunities
	term = strings.ToLower(term)

	for _, community := range c {
		if strings.Contains(strings.ToLower(community.CommunityName), term) ||
			strings.Contains(strings.ToLower(community.County), term) ||
			strings.Contains(strconv.Itoa(community.CID), term) {
			matchingCommunities = append(matchingCommunities, community)
		}
	}

	return &matchingCommunities
}

func (c *NFIPCommunities) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(c)
}

func (c *NFIPCommunities) addCommunity(comm *NFIPCommunity) {
	*c = append(*c, *comm)
}

func unmarshal(reader *csv.Reader) (NFIPCommunities, error) {
	var communities NFIPCommunities
	var lineNumber int = 1
	var firstPass bool = true

	for {
		record, err := reader.Read()

		// Skip the header
		if firstPass {
			firstPass = false
			lineNumber++
			continue
		}

		if err != nil {
			if err == io.EOF {
				break
			}

			// if we get an error other than EOF, then return it
			return nil, fmt.Errorf("** ERR: %s on line %d", err.Error(), lineNumber)
		}

		var date time.Time
		var boolVal bool

		// Clean all of the data by trimming off '=' and '"' characters
		for i := 0; i < len(record); i++ {
			record[i] = strings.Trim(record[i], "\"=")
		}

		nc := NFIPCommunity{}

		// Trim the leading "=" before each CID number
		cidString := record[POSCID]

		if len(cidString) > 0 {
			cid, err := strconv.Atoi(cidString)

			if err != nil {
				return nil, fmt.Errorf("** ERR: %s on line %d", err.Error(), lineNumber)
			}

			nc.CID = cid
		}

		nc.CommunityName = record[POSCommunityName]
		nc.County = record[POSCounty]

		date, err = parseDate(record[POSFHBMIdentified])
		if err == nil {
			nc.FHBMIdentified = date
		} else if err != nil && err != ErrEmptyString && err != ErrInvalidDateString {
			return nil, fmt.Errorf("** ERR: %s on line %d", err.Error(), lineNumber)
		}

		date, err = parseDate(record[POSFIRMIdentified])
		if err == nil {
			nc.FIRMIdentified = date
		} else if err != nil && err != ErrEmptyString && err != ErrInvalidDateString {
			return nil, fmt.Errorf("** ERR: %s on line %d", err.Error(), lineNumber)
		}

		date, err = parseDate(record[POSCurrEffMapDate])
		if err == nil {
			nc.CurrEffMapDate = date
		} else if err != nil && err != ErrEmptyString && err != ErrInvalidDateString {
			return nil, fmt.Errorf("** ERR: %s on line %d", err.Error(), lineNumber)
		}

		date, err = parseDate(record[POSRegEmerDate])
		if err == nil {
			nc.RegEmerDate = date
		} else if err != nil && err != ErrEmptyString && err != ErrInvalidDateString {
			return nil, fmt.Errorf("** ERR: %s on line %d", err.Error(), lineNumber)
		}

		boolVal, err = parseBoolFromYesNo(record[POSTribal])
		if err == nil {
			nc.Tribal = boolVal
		} else if err != nil && err != ErrEmptyString {
			return nil, fmt.Errorf("** ERR: %s on line %d", err.Error(), lineNumber)
		}

		nc.CRSEntryDate = record[POSCRSEntryDate]
		nc.CurrEffDate = record[POSCurrEffDate]
		nc.CurClass = record[POSCurClass]
		nc.PercentDiscSFHA = record[POSPercentDiscSFHA]
		nc.PercentNonSFHA = record[POSPercentNonSFHA]
		nc.Program = record[POSProgram]

		boolVal, err = parseBoolFromYesNo(record[POSParticipatingCommunity])
		if err == nil {
			nc.ParticipatingCommunity = boolVal
		} else if err != nil && err != ErrEmptyString {
			return nil, fmt.Errorf("** ERR: %s on line %d", err.Error(), lineNumber)
		}

		lineNumber++
		communities.addCommunity(&nc)
	}

	return communities, nil
}

func parseDate(s string) (time.Time, error) {
	if len(s) == 0 {
		return time.Time{}, ErrEmptyString
	}

	reg := regexp.MustCompile("([0-9]+)")
	matches := reg.FindAllString(s, 3)

	// If we don't have 3 sets of numbers,
	// then it isn't a valid date string
	// and we can stop trying to parse it
	if len(matches) < 3 {
		return time.Time{}, ErrInvalidDateString
	}

	m, err := strconv.Atoi(matches[0])
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse month to integer")
	}
	month, err := iToMonth(m)
	if err != nil {
		return time.Time{}, err
	}

	day, err := strconv.Atoi(matches[1])
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse day to integer")
	}

	year, err := strconv.Atoi(matches[2])
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse year to integer")
	}

	// The date is only stored in 2 digit format. So I'm taking a guess
	// on whether it represents a year from the 20th or the 21st century.
	if year <= 22 {
		year = 2000 + year
	} else {
		year = 1900 + year
	}

	return time.Date(year, month, day, 0, 0, 0, 0, time.Local), nil
}

func iToMonth(i int) (time.Month, error) {
	switch i {
	case 1:
		return time.January, nil
	case 2:
		return time.February, nil
	case 3:
		return time.March, nil
	case 4:
		return time.April, nil
	case 5:
		return time.May, nil
	case 6:
		return time.June, nil
	case 7:
		return time.July, nil
	case 8:
		return time.August, nil
	case 9:
		return time.September, nil
	case 10:
		return time.October, nil
	case 11:
		return time.November, nil
	case 12:
		return time.December, nil
	default:
		return -1, fmt.Errorf("invalid integer given for conversion to month")
	}
}

func parseBoolFromYesNo(s string) (bool, error) {
	if len(s) == 0 {
		return false, ErrEmptyString
	}

	// Normalize the string
	s = strings.ToLower(s)

	if s == "yes" {
		return true, nil
	} else if s == "no" {
		return false, nil
	}

	return false, fmt.Errorf("failed to parse bool from string %s", s)
}
