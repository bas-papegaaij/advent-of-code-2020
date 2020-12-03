package main

import (
	"bufio"
	"fmt"
	"os"
)

const inputFile = "input.txt"

type treeData [][]bool

type offset struct {
	x int
	y int
}

var slopes []offset = []offset{
	{1, 1},
	{3, 1},
	{5, 1},
	{7, 1},
	{1, 2},
}

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

	data := make(treeData, len(lines))
	// Just going to assume this is the same for every line for now
	lineLength := len(lines[0])
	for i, line := range lines {
		data[i] = getTreeData(line)
	}

	allCounts := make([]int, len(slopes))
	product := 0
	for i, coords := range slopes {
		count := countTrees(coords, lineLength, data)
		allCounts[i] = count
		if i == 0 {
			product = count
		} else {
			product *= count
		}
	}
	fmt.Println("got counts", allCounts, "with total product", product)
}

func getTreeData(data string) []bool {
	res := make([]bool, len(data))
	for i, c := range data {
		if c == '.' {
			res[i] = false
		} else if c == '#' {
			res[i] = true
		} else {
			panic(fmt.Sprintf("got bad data %s", data))
		}
	}
	return res
}

// assuming rows are all the same length
func countTrees(offsets offset, rowSize int, data treeData) int {
	treeCount := 0
	for y, x := 0, 0; y < len(data); y, x = y+offsets.y, (x+offsets.x)%rowSize {
		point := data[y][x]
		if point {
			treeCount++
		}
	}

	return treeCount
}
