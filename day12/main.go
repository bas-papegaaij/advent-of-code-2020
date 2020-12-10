package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

const inputFile = "input.txt"

type instruction struct {
	Type  string
	Value int
}

type vector struct {
	x int
	y int
}

type ship struct {
	bearing  vector
	position vector
	waypoint vector
}

func main() {
	file, err := os.Open(inputFile)
	if err != nil {
		panic(err)
	}

	instructions := []instruction{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		instructions = append(instructions, parseInstruction(scanner.Text()))
	}

	// part1
	s := &ship{bearing: vector{1, 0}}
	for _, inst := range instructions {
		s.followInstruction(inst)
	}
	fmt.Printf("Ship is now in relative position %d, %d\n", s.position.x, s.position.y)
	fmt.Println("Manhattan distance is", abs(s.position.x)+abs(s.position.y))

	// part2
	s = &ship{waypoint: vector{10, 1}}
	for _, inst := range instructions {
		s.followWaypoint(inst)
	}
	fmt.Printf("Ship is now in relative position %d, %d\n", s.position.x, s.position.y)
	fmt.Println("Manhattan distance is", abs(s.position.x)+abs(s.position.y))
}

func parseInstruction(data string) instruction {
	t := data[0]
	val, err := strconv.Atoi(data[1:])
	if err != nil {
		panic(err)
	}
	return instruction{
		Type:  string(t),
		Value: val,
	}
}

func (s *ship) followInstruction(i instruction) {
	switch i.Type {
	case "L":
		s.bearing.rotate(-i.Value)
	case "R":
		s.bearing.rotate(i.Value)
	case "F":
		s.moveForward(i.Value)
	default:
		s.position.moveInDirection(i.Type, i.Value)
	}
}

func (s *ship) followWaypoint(i instruction) {
	switch i.Type {
	case "L":
		s.waypoint.rotate(-i.Value)
	case "R":
		s.waypoint.rotate(i.Value)
	case "F":
		s.moveToWaypoint(i.Value)
	default:
		s.waypoint.moveInDirection(i.Type, i.Value)
	}
}

func (vec *vector) rotate(rotation int) {
	// easiest to do all rotations in the same frame of reference
	if rotation < 0 {
		rotation += 360
	}

	switch rotation {
	case 90:
		vec.x, vec.y = vec.y, vec.x
		vec.y = -vec.y
	case 180:
		vec.x = -vec.x
		vec.y = -vec.y
	case 270:
		vec.x, vec.y = vec.y, vec.x
		vec.x = -vec.x
	default:
		panic("invalid rotation")
	}
}

func (s *ship) moveForward(val int) {
	s.position.x += s.bearing.x * val
	s.position.y += s.bearing.y * val
}

func (s *ship) moveToWaypoint(val int) {
	s.position.x += s.waypoint.x * val
	s.position.y += s.waypoint.y * val
}

func (vec *vector) moveInDirection(dir string, val int) {
	switch dir {
	case "N":
		vec.y += val
	case "S":
		vec.y -= val
	case "E":
		vec.x += val
	case "W":
		vec.x -= val
	default:
		panic("invalid direction")
	}
}

func abs(val int) int {
	if val < 0 {
		return -val
	}
	return val
}
