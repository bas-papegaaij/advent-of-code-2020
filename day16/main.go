package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

const inputFile = "input.txt"

type rule []validRange
type validRange struct {
	start int
	end   int
}

type ticketRules map[string]rule
type ticket []int

const (
	parseRules int = iota
	parseMyTicket
	parseOtherTickets
)

func main() {
	file, err := os.Open(inputFile)
	if err != nil {
		panic(err)
	}

	timings := make(map[string]time.Duration)
	programStart := time.Now()

	rules := make(ticketRules)
	parseState := parseRules
	var myTicket ticket
	var otherTickets []ticket
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		t := scanner.Text()
		if t == "" {
			continue
		}
		if t == "your ticket:" {
			parseState = parseMyTicket
			continue
		}
		if t == "nearby tickets:" {
			parseState = parseOtherTickets
			continue
		}

		switch parseState {
		case parseRules:
			name, ranges := parseRule(t)
			rules[name] = ranges
		case parseMyTicket:
			myTicket = parseTicket(t)
		case parseOtherTickets:
			otherTickets = append(otherTickets, parseTicket(t))
		default:
			break
		}
	}
	timings["parse"] = time.Since(programStart)

	// part 1
	part1Start := time.Now()
	rate := calculateErrorRate(rules, otherTickets)
	timings["part1"] = time.Since(part1Start)

	// part 2
	part2Start := time.Now()
	validTickets := removeInvalidTickets(rules, otherTickets)
	fieldMap := mapFields(rules, validTickets)
	value := -1
	for name, idx := range fieldMap {
		if strings.HasPrefix(name, "departure") {
			if value == -1 {
				value = myTicket[idx]
			} else {
				value *= myTicket[idx]
			}
		}
	}
	timings["part2"] = time.Since(part2Start)

	fmt.Println("part1", rate)
	fmt.Println("part2", value)
	fmt.Println("timings:", timings)
}

func parseRule(ruleString string) (string, rule) {
	parts := strings.Split(ruleString, ":")
	if len(parts) != 2 {
		panic("bad rule")
	}
	name := parts[0]
	ranges := strings.Split(parts[1], "or")
	validRanges := make(rule, len(ranges))
	for i, r := range ranges {
		nums := strings.Split(strings.TrimSpace(r), "-")
		start, err := strconv.Atoi(nums[0])
		if err != nil {
			panic(err)
		}
		end, err := strconv.Atoi(nums[1])
		if err != nil {
			panic(err)
		}
		validRanges[i] = validRange{start, end}
	}
	return name, validRanges
}

func parseTicket(t string) ticket {
	nums := strings.Split(t, ",")
	res := make(ticket, len(nums))
	for i, numStr := range nums {
		num, err := strconv.Atoi(numStr)
		if err != nil {
			panic(err)
		}
		res[i] = num
	}

	return res
}

func calculateErrorRate(rules ticketRules, tickets []ticket) int {
	rate := 0
	for _, t := range tickets {
		rate += t.calculateErrorRate(rules)
	}
	return rate
}

func (t ticket) calculateErrorRate(rules ticketRules) int {
	rate := 0
	for _, value := range t {
		valid := false
		// checkValue:
		for _, rule := range rules {
			if rule.isValid(value) {
				valid = true
				break
			}
		}
		if !valid {
			rate += value
		}
	}
	return rate
}

func removeInvalidTickets(rules ticketRules, tickets []ticket) []ticket {
	results := make([]ticket, 0)
	for _, t := range tickets {
		if t.calculateErrorRate(rules) == 0 {
			results = append(results, t)
		}
	}
	return results
}

func mapFields(rules ticketRules, tickets []ticket) map[string]int {
	fields := make(map[string]int)

	// first we match each index on the ticket to every possible field name
	// it could map to
	type possibleFields struct {
		ticketIndex int
		fieldNames  []string
	}
	possibilities := make([]possibleFields, len(tickets[0]))

	for i := 0; i < len(tickets[0]); i++ {
		possibilities[i] = possibleFields{
			fieldNames:  findPossibleFieldNames(rules, i, tickets),
			ticketIndex: i,
		}
	}

	// some testing confirmed that the possible fields for each value index
	// allows us to be ascertain which index maps to which field one at a time
	// i.e. there is one index which can only belong to 1 rule, one index which can belong to 2
	// etc.
	allNames := make(map[string]struct{})
	for name := range rules {
		allNames[name] = struct{}{}
	}

	sort.Slice(possibilities, func(i, j int) bool { return len(possibilities[i].fieldNames) < len(possibilities[j].fieldNames) })
	for _, entry := range possibilities {
		for _, name := range entry.fieldNames {
			_, ok := allNames[name]
			if !ok {
				continue
			}
			// since we're iterating the entries in ascending order of possibilities
			// the first valid name we find should always be the only valid name
			// since we've removed previously used entries from the map
			fields[name] = entry.ticketIndex
			delete(allNames, name)
		}
	}

	return fields
}

func (r rule) isValid(value int) bool {
	for _, vr := range r {
		if value >= vr.start && value <= vr.end {
			return true
		}
	}
	return false
}

func findPossibleFieldNames(rules ticketRules, fieldIndex int, tickets []ticket) []string {
	possibleNames := make([]string, 0)
	for name := range rules {
		possibleNames = append(possibleNames, name)
	}

	// so many loops
	for _, t := range tickets {
		val := t[fieldIndex]
		for i, name := range possibleNames {
			if !rules[name].isValid(val) {
				possibleNames = append(possibleNames[:i], possibleNames[i+1:]...)
			}
		}
	}

	return possibleNames
}
