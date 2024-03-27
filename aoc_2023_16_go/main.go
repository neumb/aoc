package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	DIR_RIGHT = iota
	DIR_DOWN
	DIR_LEFT
	DIR_UP
)

type layout struct {
	width  int
	height int
	grid   [][]rune
}

type coords struct {
	y, x int
}

type pathStep struct {
	c coords
	//dir int
}

func NewPathStep(y, x, dir int) pathStep {
	_ = dir
	return pathStep{
		coords{y, x},
	}
}

type beamState struct {
	c       coords
	dir     int
	path    map[pathStep]struct{}
	blocked bool
}

func NewBeamState(y, x, dir int) *beamState {
	path := make(map[pathStep]struct{})
	path[NewPathStep(y, x, dir)] = struct{}{}
	return &beamState{
		coords{y, x},
		dir,
		path,
		false,
	}
}

func (b *beamState) AddStep() bool {
	_, ok := b.path[NewPathStep(b.c.y, b.c.x, b.dir)]
	b.path[NewPathStep(b.c.y, b.c.x, b.dir)] = struct{}{}
	return ok == false
}

func parse(input string) layout {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	width := len(lines[0])
	height := len(lines)

	grid := make([][]rune, height)
	var last int

	for _, l := range lines {
		if len(l) == 0 {
			continue
		}
		grid[last] = []rune(l)
		last++
	}

	return layout{
		width,
		height,
		grid,
	}
}

func part1(input layout) any {
	energized := make(map[coords]struct{})
	splitted := make(map[coords]struct{})
	beams := make([]*beamState, 0)
	beams = append(beams, NewBeamState(0, 0, DIR_RIGHT))

	shouldContinue := true
	var steps int

	for shouldContinue {
		steps++
		shouldContinue = steps < 2000
		for _, b := range beams {
			if b.blocked {
				continue
			}
			b.AddStep()
			shouldContinue = true
			energized[b.c] = struct{}{}

			tile := input.grid[b.c.y][b.c.x]
			switch tile {
			case '/', '\\': // mirror
				if tile == '/' {
					switch b.dir {
					case DIR_UP:
						b.dir = DIR_RIGHT
					case DIR_RIGHT:
						b.dir = DIR_UP
					case DIR_DOWN:
						b.dir = DIR_LEFT
					case DIR_LEFT:
						b.dir = DIR_DOWN
					}
				} else {
					switch b.dir {
					case DIR_UP:
						b.dir = DIR_LEFT
					case DIR_RIGHT:
						b.dir = DIR_DOWN
					case DIR_DOWN:
						b.dir = DIR_RIGHT
					case DIR_LEFT:
						b.dir = DIR_UP
					}
				}

			case '-': // splitter, check current dir
				if b.dir == DIR_RIGHT || b.dir == DIR_LEFT {
					// pass through
				} else {
					leftSideBeam := NewBeamState(b.c.y, b.c.x-1, DIR_LEFT)
					_, wasSplitted := splitted[b.c]
					//b.c.x = b.c.x + 1
					b.dir = DIR_RIGHT
					if !wasSplitted && leftSideBeam.c.x >= 0 {
						beams = append(beams, leftSideBeam)
						splitted[b.c] = struct{}{}
					} else {
					}
				}
			case '|': // splitter, check current dir
				if b.dir == DIR_UP || b.dir == DIR_DOWN {
					// pass through
				} else {
					upSideBeam := NewBeamState(b.c.y-1, b.c.x, DIR_UP)
					_, wasSplitted := splitted[b.c]
					//b.c.y = b.c.y + 1
					b.dir = DIR_DOWN

					if !wasSplitted && upSideBeam.c.y >= 0 {
						beams = append(beams, upSideBeam)
						splitted[b.c] = struct{}{}
					} else {
					}
				}
			case '.':
			default:
				panic("invalid tile")
			}

			energized[b.c] = struct{}{}

			switch b.dir {
			case DIR_UP:
				b.c = coords{
					max(b.c.y-1, 0),
					b.c.x,
				}
			case DIR_DOWN:
				b.c = coords{
					min(b.c.y+1, input.height-1),
					b.c.x,
				}
			case DIR_RIGHT:
				b.c = coords{
					b.c.y,
					min(b.c.x+1, input.width-1),
				}
			case DIR_LEFT:
				b.c = coords{
					b.c.y,
					max(b.c.x-1, 0),
				}
			default:
				panic(fmt.Sprintf("unexpected dir %d", b.dir))
			}
			energized[b.c] = struct{}{}

			if !b.AddStep() {
				b.blocked = true
			}
		}

		//shouldContinue = prevEnergized != len(energized)
	}

	//fmt.Println("steps:", steps)

	//for y := 0; y < input.height; y++ {
		//for x := 0; x < input.width; x++ {
			//_, ok := energized[coords{y, x}]
			//if ok {
				//fmt.Print("#")
			//} else {
				//fmt.Print(string(input.grid[y][x]))
			//}
		//}
		//fmt.Println()
	//}

	return len(energized)
}

func dirToString(dir int) string {
	switch dir {
	case DIR_RIGHT:
		return "RIGHT"
	case DIR_LEFT:
		return "LEFT"
	case DIR_UP:
		return "UP"
	case DIR_DOWN:
		return "DOWN"
	default:
		panic("invalid dir")
	}
}

func part2(input layout) any {
	return 0
}

func main() {
	content, err := os.ReadFile("input.txt")

	// 5706 is too low
	// 7246 is too high

	if err != nil {
		log.Fatal(err)
	}

	input := parse(string(content))

	fmt.Println("part1:", part1(input))
	fmt.Println("part2:", part2(input))
}
