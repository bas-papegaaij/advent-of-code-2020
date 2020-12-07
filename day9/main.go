package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

const inputFile = "input.txt"
const preambleSize = 25

func main() {
	file, err := os.Open(inputFile)
	if err != nil {
		panic(err)
	}

	numbers := []int{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		t := scanner.Text()
		val, err := strconv.Atoi(t)
		if err != nil {
			panic("found invalid number")
		}
		numbers = append(numbers, val)
	}

	var firstBadNumber int
	for i, val := range numbers[preambleSize:] {
		if !isValidNumber(val, numbers[i:i+preambleSize]) {
			firstBadNumber = val
			break
		}
	}
	fmt.Println("First bad number is:", firstBadNumber)

	vuln := findContiguousSum(firstBadNumber, numbers)
	fmt.Println("The vulnerability is:", vuln)
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

func findContiguousSum(sumValue int, values []int) int {
	startIndex, endIndex := 0, 1
	sum := 0
	for sum != sumValue {
		sum = sumValues(startIndex, endIndex, values)
		if sum < sumValue {
			// We haven't hit the wanted value yet. Grow the window by 1 and check the next sum
			endIndex++
		} else if sum > sumValue {
			// reset the window to its minimum size, starting from the next number
			startIndex++
			endIndex = startIndex + 1
		}
	}

	smallest, largest := findExtremes(values[startIndex : endIndex+1])

	return smallest + largest
}

const maxInt = int(^uint(0) >> 1)

func findExtremes(values []int) (int, int) {
	smallest, largest := maxInt, 0
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

func sumValues(start int, end int, sourceValues []int) int {
	sum := 0
	for _, val := range sourceValues[start : end+1] {
		sum += val
	}
	return sum
}
