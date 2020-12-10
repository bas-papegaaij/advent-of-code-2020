package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

const inputFile = "input.txt"

func main() {
	timings := make(map[string]time.Duration)
	data, err := ioutil.ReadFile(inputFile)
	if err != nil {
		panic(err)
	}

	parseStart := time.Now()
	lines := strings.Split(string(data), "\n")
	if len(lines) != 2 {
		panic("wrong number of lines in input")
	}

	curTime := parseTime(lines[0])
	busIDs := parseBusData(lines[1])
	timings["parse"] = time.Since(parseStart)

	part1Start := time.Now()
	waitTime, ID := findSoonestBus(curTime, busIDs)
	timings["part1"] = time.Since(part1Start)

	part2Start := time.Now()
	t := findSpecialTimestamp(busIDs)
	timings["part2"] = time.Since(part2Start)

	fmt.Println("Soonest bus departing in", waitTime, "minutes with ID:", ID)
	fmt.Println("Part 1:", waitTime*ID)
	fmt.Println("Part 2", t)
	fmt.Println("timings:", timings)
}

// returns the number of minutes before the soonest bus arrives and its ID
func findSoonestBus(curTime int, buses []int) (int, int) {
	soonestDiff := -1
	soonestID := 0
	for _, bus := range buses {
		if bus == -1 {
			continue
		}

		mod := curTime % bus
		// bus leaves exactly at the current timestamp
		if mod == 0 {
			return curTime, bus
		}

		diff := bus - mod
		if soonestDiff == -1 || diff < soonestDiff {
			soonestDiff = diff
			soonestID = bus
		}
	}

	return soonestDiff, soonestID
}

// no idea what would be a sensible name for this function
// it finds the nearest timestamp so that bus[i] leaves i minutes after bus[0]
func findSpecialTimestamp(buses []int) int64 {
	time, step := int64(0), int64(buses[0])
	for i := int64(1); i < int64(len(buses)); i++ {
		if buses[i] == -1 {
			continue
		}

		curBus := int64(buses[i])
		for (time+i)%curBus != 0 {
			time += step
		}
		step *= curBus
	}
	return time
}

func parseTime(time string) int {
	curTime, err := strconv.Atoi(time)
	if err != nil {
		panic(err)
	}

	return curTime
}

func parseBusData(buses string) []int {
	split := strings.Split(buses, ",")
	busIDs := make([]int, len(split))
	for i, entry := range split {
		if entry == "x" {
			busIDs[i] = -1
			continue
		}
		val, err := strconv.Atoi(entry)
		if err != nil {
			panic(err)
		}
		busIDs[i] = val
	}

	return busIDs
}
