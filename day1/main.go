package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
)

const inputFile = "input.txt"

const numValues = 4
const findSum = 2017

func main() {
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

	numbers := convertToInt(lines)
	sort.Ints(numbers)
	complements := findComplements(0, findSum, numValues, numbers)
	if complements == nil || len(complements) == 0 {
		fmt.Println("Can't find", numValues, "numbers adding up to", findSum)
		os.Exit(0)
	}

	fmt.Println("numbers:", complements)
	product := complements[0]
	for _, val := range complements[1:] {
		product *= val
	}
	fmt.Println("final product", product)
}

func convertToInt(text []string) []int {
	output := make([]int, len(text))
	for i, t := range text {
		output[i], _ = strconv.Atoi(t)
	}
	return output
}

func findComplements(value int, wanted int, numComplements int, vals []int) []int {
	if len(vals) < numComplements || numComplements == 0 {
		return nil
	}

	for i, val := range vals[:len(vals)+1-numComplements] {
		sum := value + val
		if sum > wanted {
			return nil
		}

		if sum == wanted {
			if numComplements == 1 {
				return []int{val}
			}
			return nil
		}

		otherComplements := findComplements(sum, wanted, numComplements-1, vals[i+1:])
		if otherComplements != nil {
			return append(otherComplements, val)
		}
	}

	return nil
}
