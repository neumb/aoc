package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Direction int

func (d Direction) String() string {
	switch d {
	case Dir_Up:
		return "U"
	case Dir_Down:
		return "D"
	case Dir_Left:
		return "L"
	case Dir_Right:
		return "R"
	default:
		panic("invalid direction")
	}
}

func (d Direction) Delta() [2]int {
	switch d {
	case Dir_Up:
		return [2]int{-1, 0}
	case Dir_Down:
		return [2]int{+1, 0}
	case Dir_Left:
		return [2]int{0, -1}
	case Dir_Right:
		return [2]int{0, +1}
	default:
		panic("invalid direction")
	}
}

const (
	Dir_Up Direction = iota
	Dir_Down
	Dir_Left
	Dir_Right
)

type Grid struct {
	grid          map[[2]int]bool
	height, width int
}

func (g *Grid) Print() {
	for y := 0; y < g.height; y++ {
		for x := 0; x < g.width; x++ {
			if g.grid[[2]int{y, x}] {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

type PlanLine struct {
	Dir   Direction
	Dist  int
	Color string
}

func (l PlanLine) String() string {
	return fmt.Sprintf("%s %d (%s)", l.Dir, l.Dist, l.Color)
}

func PlanLineFromString(line string) PlanLine {
	fields := strings.Fields(line)

	pl := PlanLine{}

	switch fields[0] {
	case "R":
		pl.Dir = Dir_Right
	case "L":
		pl.Dir = Dir_Left
	case "U":
		pl.Dir = Dir_Up
	case "D":
		pl.Dir = Dir_Down
	default:
		panic("invalid direction")
	}

	if dist, err := strconv.Atoi(fields[1]); err != nil {
		panic("invalid distance")
	} else {
		pl.Dist = dist
	}

	pl.Color = strings.Trim(fields[2], "()")

	return pl
}

type AoC struct {
	plan []PlanLine
}

func parse(input string) AoC {
	lines := strings.Split(strings.TrimSpace(input), "\n")

	plan := make([]PlanLine, 0, len(lines))

	for _, l := range lines {
		plan = append(plan, PlanLineFromString(l))
	}

	return AoC{
		plan,
	}
}

func (a *AoC) part1() any {
	grid := make(map[[2]int]bool)
	pos := [2]int{0, 0}
	grid[pos] = true

	var maxY, maxX int

	for _, l := range a.plan {
		posD := l.Dir.Delta()

		for d := 0; d < l.Dist; d++ {
			pos[0] += posD[0]
			pos[1] += posD[1]
			grid[pos] = true
		}
	}

	for p := range grid {
		maxY = max(maxY, p[0])
		maxX = max(maxX, p[1])
	}

	g := Grid{grid, maxY + 1, maxX + 1}
	//g.Print()
	//fmt.Println()
	//g.fillRay()

	//fmt.Println(pos)

	g.fill(1, 1)
	g.fill(1, 155)

	//g.Print()

	return len(grid)
}

func (g *Grid) fillRay() {
	for y := 0; y < g.height; y++ {
		for x := 0; x < g.width; x++ {
			var interTimes int
			var prev bool
			for nx := x; nx < g.width; nx++ {

				val := g.grid[[2]int{y, nx}]

				if !prev && val {
					interTimes++
				}

				prev = val
			}

			if interTimes%2 != 0 {
				g.grid[[2]int{y, x}] = true
			}
		}
	}
}

func (g *Grid) fill(y, x int) {
	if g.grid[[2]int{y, x}] {
		return
	}

	g.grid[[2]int{y, x}] = true

	if y > 0 {
		g.fill(y-1, x)
	}
	if y < g.height {
		g.fill(y+1, x)
	}
	if x > 0 {
		g.fill(y, x-1)
	}
	if x < g.width {
		g.fill(y, x+1)
	}
}

func (a *AoC) part2() any {
	return 0
}

func main() {
	//content, err := os.ReadFile("test-input.txt")
	content, err := os.ReadFile("input.txt")
	// 6136 is too low

	if err != nil {
		log.Fatal(err)
	}

	aoc := parse(string(content))

	fmt.Println("part1:", aoc.part1())
	fmt.Println("part2:", aoc.part2())
}
