package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
)

const inputFile = "input.txt"
const maxAdaptorDifference = 3

func main() {
	file, err := os.Open(inputFile)
	if err != nil {
		panic(err)
	}

	ratings := []int{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		l := scanner.Text()
		val, err := strconv.Atoi(l)
		if err != nil {
			panic(err)
		}
		ratings = append(ratings, val)
	}
	sort.Ints(ratings)
	_, oneDiffs, threeDiffs := getDifferences(ratings)
	fmt.Println("part1:", oneDiffs*(threeDiffs+1))
}

// gets the difference between each consecutive rating
// also counts the number of difference exactly equal to 1 or 3
func getDifferences(ratings []int) ([]int, int, int) {
	oneDiffs, threeDiffs := 0, 0
	res := make([]int, len(ratings))
	// this is the difference between the first adaptor and 0
	res[0] = ratings[0]
	if res[0] == 1 {
		oneDiffs++
	}
	if res[0] == 3 {
		threeDiffs++
	}

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
