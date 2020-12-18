package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

const inputFile = "input.txt"
const simCount = 6

type cube struct {
	dimensions []int
	data       []bool
}

type rulesFunc func(curState bool, activeNeighbours int) bool

func main() {
	data, err := ioutil.ReadFile(inputFile)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(data), "\n")
	state := parseData(lines, 3)

	for i := 0; i < simCount; i++ {
		state.runSim(rulesV1)
	}
	fmt.Println("part 1:", state.countActive())

	state = parseData(lines, 4)
	for i := 0; i < simCount; i++ {
		state.runSim(rulesV1)
	}
	fmt.Println("part 2", state.countActive())
}

func parseData(data []string, numDims int) *cube {
	dims := make([]int, numDims)
	dims[0] = len(data[0])
	dims[1] = len(data)
	for i := 2; i < numDims; i++ {
		dims[i] = 1
	}
	c := cube{
		dimensions: dims,
	}

	size := dims[0]
	for _, dim := range dims[1:] {
		size *= dim
	}

	c.data = make([]bool, size)
	for y, line := range data {
		for x, r := range line {
			coords := make([]int, numDims)
			coords[0] = x
			coords[1] = y
			switch r {
			case '#':
				c.setState(coords, true)
			case '.':
				c.setState(coords, false)
			default:
				panic("bad state")
			}
		}
	}
	return &c
}

func (c *cube) copyAndGrowMultiDim() *cube {
	cp := cube{
		dimensions: make([]int, len(c.dimensions)),
	}
	for i, dim := range c.dimensions {
		cp.dimensions[i] = dim + 2
	}
	size := cp.dimensions[0]
	for i := 1; i < len(cp.dimensions); i++ {
		size *= cp.dimensions[i]
	}
	cp.data = make([]bool, size)

	for idx, state := range c.data {
		coords := c.getCoords(idx)
		offset := make([]int, len(coords))
		for i, coord := range coords {
			offset[i] = coord + 1
		}
		cp.setState(offset, state)
	}

	// swap that data
	c.dimensions, cp.dimensions = cp.dimensions, c.dimensions
	c.data, cp.data = cp.data, c.data

	return &cp
}

func (c *cube) runSim(rules rulesFunc) {
	tmp := c.copyAndGrowMultiDim()
	for idx := range c.data {
		coords := c.getCoords(idx)
		offset := make([]int, len(coords))
		for i, coord := range coords {
			offset[i] = coord - 1
		}

		newState := tmp.getNewState(offset, rules)
		c.setState(coords, newState)
	}
}
func (c *cube) countActive() int {
	count := 0
	for _, state := range c.data {
		if state == true {
			count++
		}
	}

	return count
}

func (c *cube) getIndex(coords []int) int {
	idx := 0
	for i, coord := range coords {
		curIdx := coord
		for j := 0; j < i; j++ {
			curIdx *= c.dimensions[j]
		}
		idx += curIdx
	}

	return idx
}

func (c *cube) setState(coords []int, state bool) {
	c.data[c.getIndex(coords)] = state
}

func (c *cube) getState(coords []int) bool {
	// out of bounds is bad
	for i, coord := range coords {
		if coord < 0 || coord >= c.dimensions[i] {
			return false
		}
	}

	return c.data[c.getIndex(coords)]
}

func (c *cube) getCoords(index int) []int {
	coords := make([]int, len(c.dimensions))
	for i := 0; i < len(c.dimensions); i++ {
		val := index
		for j := 0; j < i; j++ {
			val /= c.dimensions[j]
		}
		if i != len(c.dimensions)-1 {
			val %= c.dimensions[i]
		}
		coords[i] = val
	}
	return coords
}

func (c *cube) getNewState(coords []int, rules rulesFunc) bool {
	curState := c.getState(coords)
	count := c.getNeighbourCount(coords, make([]int, 0))
	return rules(curState, count)
}

func rulesV1(curState bool, activeNeighbours int) bool {
	if curState == true {
		return activeNeighbours == 2 || activeNeighbours == 3
	}

	return activeNeighbours == 3
}

func (c *cube) getNeighbourCount(reference []int, check []int) int {
	dimIdx := len(check)
	val := reference[dimIdx]
	count := 0
	for i := val - 1; i <= val+1; i++ {
		next := make([]int, len(check))
		copy(next, check)
		next = append(next, i)
		if len(next) == len(reference) {
			if arrayEquals(next, reference) {
				continue
			}
			if c.getState(next) == true {
				count++
			}
		} else {
			count += c.getNeighbourCount(reference, next)
		}
	}

	return count
}

func arrayEquals(a []int, b []int) bool {
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
