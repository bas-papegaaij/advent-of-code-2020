package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const inputFile = "input.txt"

type memoryMap map[uint64]uint64

func main() {
	file, err := os.Open(inputFile)
	if err != nil {
		panic(err)
	}

	// we'll just use int64 so we can store 36 bits easily
	part1Memory := make(memoryMap)
	part2Memory := make(memoryMap)
	scanner := bufio.NewScanner(file)
	var curMask string

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "mask") {
			curMask = parseMask(line)
		} else {
			addr, value := parseMemStore(line)
			part1Memory[addr] = maskValue(curMask, value)
			part2Memory.writeMemory(curMask, addr, value)
		}
	}

	fmt.Println("Part 1:", part1Memory.sum())
	fmt.Println("Part 2", part2Memory.sum())
}

func (m memoryMap) sum() uint64 {
	total := uint64(0)
	for _, val := range m {
		total += val
	}
	return total
}

func parseMask(maskString string) string {
	maskPart := strings.Split(maskString, "=")[1]
	return strings.TrimSpace(maskPart)
}

func parseMemStore(memString string) (uint64, uint64) {
	memRegex := regexp.MustCompile(`^mem\[([0-9]+)] = ([0-9]+)$`)
	matches := memRegex.FindStringSubmatch(memString)
	if len(matches) != 3 {
		panic("got some bad memory instructions :(")
	}

	addr, err := strconv.ParseUint(matches[1], 10, 64)
	if err != nil {
		panic(err)
	}
	val, err := strconv.ParseUint(matches[2], 10, 64)
	if err != nil {
		panic(err)
	}

	return addr, val
}

const maxuint = ^uint64(0)

func maskValue(mask string, value uint64) uint64 {
	for i, m := range mask {
		switch m {
		case 'X':
			continue
		case '1':
			value |= 1 << (35 - i)
		case '0':
			value &= maxuint ^ 1<<(35-i)
		}
	}

	return value
}

func (m memoryMap) writeMemory(mask string, baseAddr uint64, value uint64) {
	floatingIndex := strings.Index(mask, "X")
	if floatingIndex == -1 {
		maskVal, err := strconv.ParseUint(mask, 2, 64)
		if err != nil {
			panic(err)
		}
		addr := baseAddr | maskVal

		m[addr] = value
		return
	}

	cp := []byte(mask)
	cp[floatingIndex] = '0'
	m.writeMemory(string(cp), baseAddr&(maxuint^1<<(35-floatingIndex)), value)
	m.writeMemory(string(cp), baseAddr|(1<<(35-floatingIndex)), value)
	return
}
