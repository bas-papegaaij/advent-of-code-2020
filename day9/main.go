package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"
)

const inputFile = "input.txt"
const preambleSize = 25

func main() {
	timings := make(map[string]time.Duration)

	file, err := os.Open(inputFile)
	if err != nil {
		panic(err)
	}

	numbers := []int{}
	scanner := bufio.NewScanner(file)
	startProgram := time.Now()
	start := time.Now()
	for scanner.Scan() {
		t := scanner.Text()
		val, err := strconv.Atoi(t)
		if err != nil {
			panic("found invalid number")
		}
		numbers = append(numbers, val)
	}
	timings["parse"] = time.Since(start)

	start = time.Now()
	var firstBadNumber int
	for i, val := range numbers[preambleSize:] {
		if !isValidNumber(val, numbers[i:i+preambleSize]) {
			firstBadNumber = val
			break
		}
	}
	timings["part1"] = time.Since(start)

	start = time.Now()
	vuln := findvulnerability(firstBadNumber, numbers)
	timings["part2"] = time.Since(start)
	timings["total"] = time.Since(startProgram)

	fmt.Println("First bad number is:", firstBadNumber)
	fmt.Println("The vulnerability is:", vuln)

	fmt.Println("Timings:", timings)
}

// checks if num is the sum of any 2 numbers from previousValues
func isValidNumber(num int, previousValues []int) bool {
	for i, val := range previousValues {
		for _, val2 := range previousValues[i:] {
			if val+val2 == num {
				return true
			}
		}
	}

	return false
}

func findvulnerability(sumValue int, values []int) int {
	startIndex, endIndex := 0, 1
	sum := values[startIndex] + values[endIndex]
	for sum != sumValue {
		if sum < sumValue {
			endIndex++
			sum += values[endIndex]
		}
		if sum > sumValue {
			sum -= values[startIndex]
			startIndex++
		}
	}

	smallest, largest := minMax(values[startIndex : endIndex+1])

	return smallest + largest
}

func minMax(values []int) (int, int) {
	smallest, largest := values[0], values[0]
	for _, val := range values {
		if val < smallest {
			smallest = val
		}
		if val > largest {
			largest = val
		}
	}
	return smallest, largest
}
