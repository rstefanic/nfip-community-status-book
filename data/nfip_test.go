package data

import (
	"testing"
	"time"
)

func TestParseBoolFromYesNo(t *testing.T) {
	// Basic parse
	testString := "no"
	b, err := parseBoolFromYesNo(testString)
	if err != nil || b != false {
		t.Errorf("expected %s to be parsed to %t", testString, b)
	}

	testString = "yes"
	b, err = parseBoolFromYesNo(testString)
	if err != nil || b != true {
		t.Errorf("expected %s to be parsed to %t", testString, b)
	}

	// The case of the word should be ignored
	testString = "yES"
	b, err = parseBoolFromYesNo(testString)
	if err != nil || b != true {
		t.Errorf("expected %s to be parsed to %t", testString, b)
	}

	// Return an empty string error when the string is empty
	testString = ""
	_, err = parseBoolFromYesNo(testString)
	if err != ErrEmptyString {
		t.Errorf("expected an empty string error")
	}

	// Return a parse failure if the string is neither "yes" or "no"
	testString = "fail me"
	_, err = parseBoolFromYesNo(testString)
	if err == nil {
		t.Errorf("%s should not be able to be parsed", testString)
	}
}

func TestParseDate(t *testing.T) {
	// Should be able to do simple strings
	testString := "08/08/99"
	tm, err := parseDate(testString)
	if err != nil {
		t.Errorf("error occurred parsing \"%s\"", testString)
	}

	year, month, day := tm.Date()
	if year != 1999 || month != time.August || day != 8 {
		t.Errorf("\"%s\" was parsed incorrectly", testString)
	}

	// If the two digit year is between 00 and 22, then it's in the 21st century
	testString = "08/08/01"
	tm, err = parseDate(testString)
	if err != nil {
		t.Errorf("error occurred parsing \"%s\"", testString)
	}

	year, month, day = tm.Date()
	if year != 2001 || month != time.August || day != 8 {
		t.Errorf("\"%s\" was parsed incorrectly", testString)
	}

	// Letters in the string should be ignored
	testString = "09/11/09(M)"
	tm, err = parseDate(testString)
	if err != nil {
		t.Errorf("error occurred parsing \"%s\"", testString)
	}

	year, month, day = tm.Date()
	if year != 2009 || month != time.September || day != 11 {
		t.Errorf("\"%s\" was parsed incorrectly", testString)
	}

	// Strings that are not dates should be ignored
	// Example: many curr eff dates have "(NSFHA)" entered in
	testString = "(NSFHA)"
	tm, err = parseDate(testString)
	if err != ErrInvalidDateString {
		t.Errorf("expected an error parsing \"%s\" has a date", testString)
	}
}
