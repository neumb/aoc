package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func hash(s string) int {
	var res int

	for _, b := range s {
		res += int(b)
		res *= 17
		res %= 256
	}

	return res
}

func part1(input string) any {
	steps := strings.Split(strings.TrimSpace(input), ",")
	var res int

	for _, s := range steps {
		res += hash(s)
	}

	return res
}

func part2(input string) any {
	type lens struct {
		label  string
		number int
	}

	var boxes [256][]*lens
	steps := strings.Split(strings.TrimSpace(input), ",")
	var res int

	for _, s := range steps {
		if strings.Contains(s, "=") {
			split := strings.SplitN(s, "=", 2)
			label, num := split[0], split[1]

			numInt, err := strconv.Atoi(num)
			if err != nil {
				log.Fatal(err)
			}

			boxIdx := hash(label)
			box := &boxes[boxIdx]

			var found bool

			for _, v := range *box {
				if v.label == label {
					v.number = numInt
					found = true
					break
				}
			}

			if !found {
				*box = append(*box, &lens{label, numInt})
			}

		} else {
			label := strings.TrimRight(s, "-")
			boxIdx := hash(label)
			box := &boxes[boxIdx]

			removeIdx := -1

			for idx, v := range *box {
				if v.label == label {
					removeIdx = idx
					break
				}
			}

			if removeIdx >= 0 {
				*box = append((*box)[:removeIdx], (*box)[removeIdx+1:]...)
			}
		}
	}

	for bIdx, b := range boxes {
		if len(b) > 0 {
			for eIdx, e := range b {
				res += (1 + bIdx) * (eIdx + 1) * e.number
			}
		}
	}

	return res
}

func main() {
	content, err := os.ReadFile("input.txt")

	if err != nil {
		log.Fatal(err)
	}

	input := string(content)

	fmt.Println("part1:", part1(input))
	fmt.Println("part2:", part2(input))
}
