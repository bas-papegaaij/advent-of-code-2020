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

type bagContents struct {
	colour string
	count  int
}
type contentList []bagContents
type bagRules map[string]contentList

var rulesRegexp = regexp.MustCompile(`^(?P<colour>[a-z]+ [a-z]+) bags contain (?P<contents>.*)\.$`)
var contentRegex = regexp.MustCompile(`^(?P<count>[0-9]+) (?P<colour>[a-z]+ [a-z]+)`)

const targetBagType = "shiny gold"

func main() {
	file, err := os.Open(inputFile)
	if err != nil {
		panic(err)
	}

	rules := make(bagRules)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		col, cont := parseBagRule(scanner.Text())
		rules[col] = cont
	}

	sum := 0
	for _, contents := range rules {
		if contents.canContainBagColour(rules, targetBagType) {
			sum++
		}
	}

	totalContainedBags := calculateContainedBags(rules, targetBagType)

	fmt.Println(sum, "bags can eventually contain a", targetBagType, "bag")
	fmt.Println(targetBagType, "must contain", totalContainedBags, "other bags")
}

func calculateContainedBags(rules bagRules, targetBagType string) int {
	content, ok := rules[targetBagType]
	if !ok {
		panic("unknown bag type")
	}
	return content.calculateTotalContainedBags(rules)
}

func (content contentList) calculateTotalContainedBags(rules bagRules) int {
	sum := 0
	for _, c := range content {
		count := c.count
		sum += count
		subContent, ok := rules[c.colour]
		if !ok {
			fmt.Println("unknown colour")
			continue
		}
		if subContent != nil {
			sum += count * subContent.calculateTotalContainedBags(rules)
		}
	}
	return sum
}

func (content contentList) canContainBagColour(rules bagRules, targetColour string) bool {
	for _, c := range content {
		if c.colour == targetColour {
			return true
		}
		subContent, ok := rules[c.colour]
		if !ok {
			fmt.Println("unknown bag colour:", c.colour)
		}

		if subContent != nil && subContent.canContainBagColour(rules, targetColour) {
			return true
		}
	}
	return false
}

func parseBagRule(rule string) (col string, contents contentList) {
	matches := getNamedRegexpMatches(rulesRegexp, rule)
	col, ok := matches["colour"]
	if !ok {
		panic("No colour in rule")
	}
	cont, ok := matches["contents"]
	if !ok {
		panic("No contents in rule")
	}

	if cont != "no other bags" {
		contents = parseContentsList(cont)
	}

	return col, contents
}

func parseContentsList(contentStr string) contentList {
	res := make(contentList, 0)
	splitContent := strings.Split(contentStr, ", ")

	for _, c := range splitContent {
		res = append(res, parseContents(c))
	}

	return res
}

func parseContents(content string) bagContents {
	matches := getNamedRegexpMatches(contentRegex, content)
	countStr, ok := matches["count"]
	if !ok {
		panic(fmt.Sprintf("No count in content: %s", content))
	}
	count, err := strconv.Atoi(countStr)
	if err != nil {
		panic(err)
	}

	col, ok := matches["colour"]
	if !ok {
		panic(fmt.Sprintf("No colour in content: %s", content))
	}

	return bagContents{
		colour: col,
		count:  count,
	}
}

func getNamedRegexpMatches(r *regexp.Regexp, str string) map[string]string {
	match := r.FindStringSubmatch(str)
	namedMatches := make(map[string]string)
	if match == nil {
		return namedMatches
	}
	for i, name := range r.SubexpNames() {
		if i != 0 {
			namedMatches[name] = match[i]
		}
	}
	return namedMatches
}
