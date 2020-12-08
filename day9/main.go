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
	startProgram := time.Now()
	defer func() {
		fmt.Println("Total Program time:", time.Since(startProgram))
	}()

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

	start := time.Now()
	var firstBadNumber int
	for i, val := range numbers[preambleSize:] {
		if !isValidNumber(val, numbers[i:i+preambleSize]) {
			firstBadNumber = val
			break
		}
	}
	fmt.Println("Finding bad number took:", time.Since(start))
	fmt.Println("First bad number is:", firstBadNumber)

	start = time.Now()
	vuln := findContiguousSum(firstBadNumber, numbers)
	fmt.Println("Finding vuln took:", time.Since(start))
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
	sum := values[startIndex] + values[endIndex]
	for sum != sumValue {
		if sum < sumValue {
			// extend the window, add the new number to our running sum
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
