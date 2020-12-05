package main

import (
	"bufio"
	"fmt"
	"os"
)

const inputFile = "input.txt"
const rowCount = 128
const columnCount = 8

// SeatMap Maps which seats are occupied by id
type SeatMap []bool

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

	seatMap := makeSeatMap(lines)
	highest := findHighestSeatID(seatMap)
	yourSeat := findMissingSeatID(seatMap)

	fmt.Println("Highest seat ID:", highest)
	fmt.Println("your seat:", yourSeat)
}

func findMissingSeatID(seatMap SeatMap) int {

	// no need to check the first or last rows
	for i, occupied := range seatMap[columnCount : len(seatMap)-columnCount] {
		i := i + columnCount
		if occupied {
			continue
		}

		if seatMap[i-1] && seatMap[i+1] {
			return i
		}
	}

	panic("couldn't find your seat")
}

func findHighestSeatID(seatMap SeatMap) int {
	for i := len(seatMap) - 1; i >= 0; i-- {
		if seatMap[i] {
			return i
		}
	}

	panic("no seats occupied")
}

func makeSeatMap(seatData []string) []bool {
	seatMap := make([]bool, rowCount*columnCount)
	for _, line := range seatData {
		id := getSeatID(line)
		seatMap[id] = true
	}
	return seatMap
}

func getSeatID(data string) int {
	// quick sanity check
	if len(data) != 10 {
		panic(fmt.Sprintf("invalid seat specifier: %s", data))
	}

	row := getRow(data[:7])
	column := getColumn(data[7:])

	return row*columnCount + column
}

func getRow(rowChars string) int {
	rowStart, rowEnd := 0, rowCount-1
	for _, c := range rowChars {
		diff := (rowEnd - rowStart + 1) / 2
		if c == 'F' {
			rowEnd -= diff
		} else if c == 'B' {
			rowStart += diff
		} else {
			panic("bad row character")
		}
	}

	if rowStart != rowEnd {
		panic("bad row data")
	}

	return rowStart
}

func getColumn(columnChars string) int {
	colStart, colEnd := 0, columnCount-1
	for _, c := range columnChars {
		diff := (colEnd - colStart + 1) / 2
		if c == 'R' {
			colStart += diff
		} else if c == 'L' {
			colEnd -= diff
		} else {
			panic("bad column character")
		}
	}

	if colStart != colEnd {
		panic("bad column data")
	}

	return colStart
}
