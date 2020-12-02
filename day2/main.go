package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const inputFile = "input.txt"

type validateFunc func(passwordInfo) bool

var validateFuncs map[string]validateFunc = map[string]validateFunc{
	"index": isValidPasswordByIndex,
	"count": isValidPasswordByCount,
}
var validate validateFunc

func main() {
	parseFlags()

	f, err := os.Open(inputFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	lines := []string{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	passwords := make([]passwordInfo, len(lines))
	for i, line := range lines {
		passwords[i] = parsePassword(line)
	}

	validCount := 0
	for _, pw := range passwords {
		if validate(pw) {
			validCount++
		}
	}

	fmt.Println("found", validCount, "valid passwords")
}

func parseFlags() {
	v := flag.String("validate", "", "which validation function to use, must be one of 'index' or 'count'")
	flag.Parse()

	var ok bool
	validate, ok = validateFuncs[*v]
	if !ok {
		panic("bad validate function provided")
	}
}

func newPasswordPolicy(policy string) passwordPolicy {
	parts := strings.Split(policy, " ")
	if len(parts) != 2 {
		panic(fmt.Sprintf("invalid policy %s", policy))
	}

	counts := parts[0]
	splitCounts := strings.Split(counts, "-")
	if len(splitCounts) != 2 {
		panic(fmt.Sprintf("invalid policy %s", policy))
	}

	min, err := strconv.Atoi(splitCounts[0])
	if err != nil {
		panic(err)
	}

	max, err := strconv.Atoi(splitCounts[1])
	if err != nil {
		panic(err)
	}

	return passwordPolicy{
		targetLetter: parts[1],
		constraint1:  min,
		constraint2:  max,
	}
}

func parsePassword(pw string) passwordInfo {
	parts := strings.Split(pw, ":")
	if len(parts) != 2 {
		panic(fmt.Sprintf("bad password: %s", pw))
	}

	return passwordInfo{
		policy:   newPasswordPolicy(strings.Trim(parts[0], " ")),
		password: strings.Trim(parts[1], " "),
	}
}

type passwordPolicy struct {
	targetLetter string
	constraint1  int
	constraint2  int
}

type passwordInfo struct {
	password string
	policy   passwordPolicy
}

func isValidPasswordByCount(info passwordInfo) bool {
	count := 0
	for _, r := range info.password {
		// remaining := len(info.password) - i
		// if count+remaining < info.policy.minCount {
		// 	return false
		// }

		letter := string(r)
		if string(letter) != info.policy.targetLetter {
			continue
		}
		count++
		// if count > info.policy.maxCount {
		// 	return false
		// }
	}

	return count >= info.policy.constraint1 && count <= info.policy.constraint2
}

func isValidPasswordByIndex(info passwordInfo) bool {
	count := 0
	if string(info.password[info.policy.constraint1-1]) == info.policy.targetLetter {
		count++
	}
	if string(info.password[info.policy.constraint2-1]) == info.policy.targetLetter {
		count++
	}

	return count == 1
}
