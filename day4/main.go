package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const inputFile = "input.txt"

type passportData map[string]string
type validateFunc func(string) bool

var requiredFields map[string]validateFunc = map[string]validateFunc{
	"byr": validateByr,
	"iyr": validateIyr,
	"eyr": validateEyr,
	"hgt": validateHgt,
	"hcl": validateHcl,
	"ecl": validateEcl,
	"pid": validatePid,
}

func main() {
	f, err := os.Open(inputFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	passport := make(passportData)
	var validCount int
	for scanner.Scan() {
		t := scanner.Text()
		if t == "" {
			if passport.isValidPassport() {
				validCount++
			}
			passport = make(passportData)
		} else {
			passport.parseData(t)
		}
	}
	// need to include the last entry
	if passport.isValidPassport() {
		validCount++
	}

	fmt.Println("found", validCount, "valid passports")
}

func (p passportData) parseData(data string) {
	entries := strings.Split(data, " ")
	for _, entry := range entries {
		vals := strings.Split(entry, ":")
		if len(vals) != 2 {
			panic(fmt.Sprintf("invalid passport entry %s", entry))
		}
		p[vals[0]] = vals[1]
	}
}

func (p passportData) isValidPassport() bool {
	for key, validate := range requiredFields {
		val, ok := p[key]
		if !ok || !validate(val) {
			return false
		}
	}

	return true
}

// can be used for the 1st part of the challenge where we only need to validate
// the existance of each field
func alwaysValidate(input string) bool {
	return true
}

func validateByr(input string) bool {
	return validateYear(input, 1920, 2002)
}

func validateIyr(input string) bool {
	return validateYear(input, 2010, 2020)
}

func validateEyr(input string) bool {
	return validateYear(input, 2020, 2030)
}

func validateHgt(input string) bool {
	if strings.HasSuffix(input, "cm") {
		return validateNumber(strings.TrimSuffix(input, "cm"), 150, 193)
	} else if strings.HasSuffix(input, "in") {
		return validateNumber(strings.TrimSuffix(input, "in"), 59, 76)
	}
	return false
}

func validateHcl(input string) bool {
	matched, err := regexp.Match("^#[0-9a-f]{6}$", []byte(input))
	return err == nil && matched
}

var validColors []string = []string{
	"amb",
	"blu",
	"brn",
	"gry",
	"grn",
	"hzl",
	"oth",
}

func validateEcl(input string) bool {
	for _, c := range validColors {
		if input == c {
			return true
		}
	}
	return false
}

func validatePid(input string) bool {
	matched, err := regexp.Match("^[0-9]{9}$", []byte(input))
	return err == nil && matched
}

func validateYear(input string, min int, max int) bool {
	if len(input) != 4 {
		return false
	}

	return validateNumber(input, min, max)
}

func validateNumber(input string, min int, max int) bool {
	val, err := strconv.Atoi(input)
	if err != nil {
		return false
	}

	return val >= min && val <= max
}
