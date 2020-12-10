package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
)

const inputFile = "input2.txt"
const maxAdaptorDifference = 3
const extraDeviceJolts = 3

func main() {
	file, err := os.Open(inputFile)
	if err != nil {
		panic(err)
	}

	// insert first element
	ratings := []int{0}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		l := scanner.Text()
		val, err := strconv.Atoi(l)
		if err != nil {
			panic(err)
		}
		ratings = append(ratings, val)
	}
	// insert device
	sort.Ints(ratings)
	ratings = append(ratings, ratings[len(ratings)-1]+extraDeviceJolts)
	diffs, oneDiffs, threeDiffs := getDifferences(ratings)
	fmt.Println("part1:", oneDiffs*threeDiffs)
	fmt.Println("part2", findAllowedPermutations(diffs, oneDiffs, maxAdaptorDifference))
}

// gets the difference between each consecutive rating
// also counts the number of difference exactly equal to 1 or 3
func getDifferences(ratings []int) ([]int, int, int) {
	oneDiffs, threeDiffs := 0, 0
	res := make([]int, len(ratings)-1)
	for i := 0; i < len(ratings)-1; i++ {
		res[i] = ratings[i+1] - ratings[i]
		if res[i] == 1 {
			oneDiffs++
		}
		if res[i] == 3 {
			threeDiffs++
		}
	}
	return res, oneDiffs, threeDiffs
}

func findAllowedPermutations(differences []int, oneCounts int, maxDifference int) int64 {
	var permutations int64 = 1 + int64(oneCounts)
	var lastPermutations int64 = 1
	// runningDifference := 0
	// lastPermutations := 0
	for i := len(differences) - 1; i >= 0; i-- {
		if differences[i] >= maxDifference {
			continue
		}
		permutations++
		permutations += lastPermutations
		lastPermutations = permutations
	}

	return permutations
}
