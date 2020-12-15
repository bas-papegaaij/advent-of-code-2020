package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

const inputFile = "input.txt"
const part1Turns = 2020
const part2Turns = 30000000

func main() {
	data, err := ioutil.ReadFile(inputFile)
	if err != nil {
		panic(err)
	}

	timings := make(map[string]time.Duration)
	nums := parseNums(data)

	t := time.Now()
	// part1 := findNthNum(part1Turns, nums)
	part1 := findNthNumMemoized(part1Turns, nums)
	timings["part1"] = time.Since(t)

	fmt.Println("--------")

	t = time.Now()
	// part2 := findNthNum(part2Turns, nums)
	part2 := findNthNumMemoized(part2Turns, nums)
	timings["part2"] = time.Since(t)

	fmt.Println("number", part1Turns, "was", part1)
	fmt.Println("number", part2Turns, "was", part2)

	fmt.Println(timings)
}

func parseNums(data []byte) []int {
	split := strings.Split(string(data), ",")
	nums := make([]int, len(split))
	var err error
	for i, v := range split {
		nums[i], err = strconv.Atoi(v)
		if err != nil {
			panic(err)
		}
	}
	return nums
}

func findNthNum(n int, seed []int) int {
	lastSpoken := make(map[int]int)
	for i, num := range seed[:len(seed)-1] {
		lastSpoken[num] = i
	}
	lastNum := seed[len(seed)-1]
	for i := len(seed); i < n; i++ {
		turn, ok := lastSpoken[lastNum]
		lastSpoken[lastNum] = i - 1
		if !ok {
			lastNum = 0
			continue
		}

		lastNum = (i - 1) - turn
	}
	return lastNum
}

func findNthNumMemoized(n int, seed []int) int {
	// since a new number is only ever the number of iterations since
	// the last time the previous number was seen. An array of length n
	// is guaranteed to represent every possible number
	memo := make([]int, n)
	for i, num := range seed[:len(seed)-1] {
		// we use 1-based index for the "turn" number so that 0 can flag
		// a number we've never seen before
		memo[num] = i + 1
	}
	lastNum := seed[len(seed)-1]
	for i := len(seed); i < n; i++ {
		seenIdx := memo[lastNum]
		memo[lastNum] = i
		if seenIdx == 0 {
			lastNum = 0
			continue
		}

		lastNum = i - seenIdx
	}
	return lastNum
}
