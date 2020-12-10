package main

import (
	"bufio"
	"fmt"
	"os"
)

const inputFile = "input.txt"
const floor = '.'
const occupied = '#'
const empty = 'L'

type seatMap struct {
	width  int
	height int
	data   []rune
}

type direction struct {
	x int
	y int
}

var cardinalDirections = []direction{
	{-1, 0},  // left
	{1, 0},   // right
	{0, -1},  // up
	{0, 1},   // down
	{-1, -1}, // up-left
	{1, -1},  // up-right
	{-1, 1},  // down-left
	{1, 1},   // down-right
}

func main() {
	file, err := os.Open(inputFile)
	if err != nil {
		panic(err)
	}

	lines := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	width := len(lines[0])
	m := newSeatMap(width, len(lines))
	m.parseSeats(lines)
	m2 := m.copySeatMap()

	// Part 1
	hasChanged := true
	simCount := 0
	for hasChanged {
		simCount++
		hasChanged = m.simulateRules(1)
	}

	fmt.Println("Simulation took", simCount, "rounds to stabilise")
	fmt.Println(m.getOccupiedSeatCount(), "seats are occupied")

	// Part 2
	hasChanged = true
	simCount = 0
	for hasChanged {
		simCount++
		hasChanged = m2.simulateRules(2)
	}

	fmt.Println("Simulation took", simCount, "rounds to stabilise")
	fmt.Println(m2.getOccupiedSeatCount(), "seats are occupied")
}

func newSeatMap(width int, height int) *seatMap {
	return &seatMap{
		width:  width,
		height: height,
		data:   make([]rune, width*height),
	}
}

func (sm *seatMap) copySeatMap() *seatMap {
	res := seatMap{
		width:  sm.width,
		height: sm.height,
		data:   make([]rune, len(sm.data)),
	}
	copy(res.data, sm.data)
	return &res
}

func (sm *seatMap) parseSeats(data []string) {
	for y, line := range data {
		for x, r := range line {
			sm.setSeatState(x, y, r)
		}
	}
}

func (sm *seatMap) setSeatState(x int, y int, state rune) {
	sm.data[sm.width*y+x] = state
}

func (sm *seatMap) getSeatState(x int, y int) rune {
	return sm.data[sm.width*y+x]
}

// returns whether any state has changed
func (sm *seatMap) simulateRules(ruleset int) bool {
	changed := false
	oldSate := sm.copySeatMap()
	for i := range oldSate.data {
		x, y := i%sm.width, i/sm.width
		if ruleset == 1 {
			sm.data[i] = oldSate.getNewStateV1(x, y)
		} else {
			sm.data[i] = oldSate.getNewStateV2(x, y)
		}
		if sm.data[i] != oldSate.data[i] {
			changed = true
		}
	}

	return changed
}

func (sm *seatMap) getNewStateV1(x int, y int) rune {
	curState := sm.getSeatState(x, y)
	if curState == floor {
		return floor
	}

	occupiedCount := 0
	for x2 := x - 1; x2 <= x+1; x2++ {
		for y2 := y - 1; y2 <= y+1; y2++ {
			// cut out anything out of bounds and the seat itself
			if x2 < 0 || x2 >= sm.width ||
				y2 < 0 || y2 >= sm.height ||
				(x2 == x && y2 == y) {
				continue
			}
			if sm.getSeatState(x2, y2) == occupied {
				occupiedCount++
			}
		}
	}

	if curState == empty && occupiedCount == 0 {
		return occupied
	}
	if curState == occupied && occupiedCount >= 4 {
		return empty
	}
	return curState
}

func (sm *seatMap) getNewStateV2(x int, y int) rune {
	curState := sm.getSeatState(x, y)
	if curState == floor {
		return floor
	}

	occupiedCount := 0
	for _, dir := range cardinalDirections {
		state := sm.getSeatStateInDirection(x, y, dir)
		if state == occupied {
			occupiedCount++
		}
	}

	if curState == empty && occupiedCount == 0 {
		return occupied
	}
	if curState == occupied && occupiedCount >= 5 {
		return empty
	}
	return curState
}

func (sm *seatMap) getOccupiedSeatCount() int {
	count := 0
	for _, seat := range sm.data {
		if seat == occupied {
			count++
		}
	}

	return count
}

// returns the state of the first (non-floor) seat in the given direction
func (sm *seatMap) getSeatStateInDirection(x int, y int, dir direction) rune {
	for x2, y2 := x+dir.x, y+dir.y; x2 >= 0 && x2 < sm.width && y2 >= 0 && y2 < sm.height; x2, y2 = x2+dir.x, y2+dir.y {
		state := sm.getSeatState(x2, y2)
		if state == floor {
			continue
		}
		return state
	}
	// consider anything out of bounds floor
	return floor
}
