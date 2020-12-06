package main

import (
	"bufio"
	"fmt"
	"os"
)

const inputFile = "input.txt"

type answerGroup struct {
	members int
	answers answers
}
type answers map[rune]int

func newAnswerGroup() *answerGroup {
	return &answerGroup{
		members: 0,
		answers: make(answers),
	}
}

func main() {
	f, err := os.Open(inputFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	groups := make([]*answerGroup, 0)
	curGroup := newAnswerGroup()
	groups = append(groups, curGroup)

	for scanner.Scan() {
		t := scanner.Text()
		if t == "" {
			curGroup = newAnswerGroup()
			groups = append(groups, curGroup)
			continue
		}

		curGroup.parse(t)
	}

	fmt.Println("Sum of answers for any member in the group", calculateAnyAnswerTotal(groups))
	fmt.Println("Sum of answers for all members in the group", claculateAllAnswerTotal(groups))
}

func (g *answerGroup) parse(data string) {
	g.members++
	for _, r := range data {
		g.answers[r]++
	}
}

func (g answerGroup) getAnyAnswerCount() int {
	return len(g.answers)
}

func (g answerGroup) getAllAnswerCount() int {
	sum := 0
	for _, a := range g.answers {
		if a == g.members {
			sum++
		}
	}
	return sum
}

func calculateAnyAnswerTotal(groups []*answerGroup) int {
	sum := 0
	for _, a := range groups {
		sum += a.getAnyAnswerCount()
	}
	return sum
}

func claculateAllAnswerTotal(groups []*answerGroup) int {
	sum := 0
	for _, a := range groups {
		sum += a.getAllAnswerCount()
	}
	return sum
}
