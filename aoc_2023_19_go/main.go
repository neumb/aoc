package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type RuleClause struct {
	Lhs string
	Op  string
	Rhs int
}

type WorkflowRule struct {
	Clause *RuleClause
	Then   string
}

func WorkflowRuleFromString(s string) *WorkflowRule {
	if !strings.Contains(s, ":") {
		return &WorkflowRule{nil, s}
	}

	// contains clause
	if strings.Contains(s, "<") {
		lhs, rhs, _ := strings.Cut(s, "<")
		rhs, then, _ := strings.Cut(rhs, ":")
		rhsInt, err := strconv.Atoi(rhs)
		if err != nil {
			log.Fatal(err)
		}

		clause := &RuleClause{lhs, "<", rhsInt}
		return &WorkflowRule{clause, then}
	}

	lhs, rhs, _ := strings.Cut(s, ">")
	rhs, then, _ := strings.Cut(rhs, ":")
	rhsInt, err := strconv.Atoi(rhs)
	if err != nil {
		log.Fatal(err)
	}

	clause := &RuleClause{lhs, ">", rhsInt}
	return &WorkflowRule{clause, then}
}

type Workflow struct {
	Name  string
	Rules []*WorkflowRule
}

func WorkflowFromString(s string) *Workflow {
	split := strings.SplitAfter(s, "{")
	workflowName, rhs := strings.TrimRight(split[0], "{"), strings.TrimRight(split[1], "}")

	ruleStrings := strings.Split(rhs, ",")

	rules := make([]*WorkflowRule, 0, 2)

	for _, rs := range ruleStrings {
		rules = append(rules, WorkflowRuleFromString(rs))
	}

	return &Workflow{workflowName, rules}
}

type Part map[string]int

func PartFromString(s string) Part {
	p := make(Part)
	props := strings.Split(strings.Trim(s, "{}"), ",")
	for _, prop := range props {
		name, val, _ := strings.Cut(prop, "=")
		valInt, err := strconv.Atoi(val)
		if err != nil {
			log.Fatal(err)
		}
		p[name] = valInt
	}
	return p
}

func (p Part) Sum() int {
	var sum int
	for _, v := range p {
		sum += v
	}
	return sum
}

type AoC struct {
	Workflows map[string]*Workflow
	Parts     []Part
}

func parse(input string) AoC {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	wrkflws := make(map[string]*Workflow)

	var cur int

	for cur < len(lines) {
		l := lines[cur]
		if len(l) == 0 {
			cur++
			break
		}
		wrkflw := WorkflowFromString(l)
		wrkflws[wrkflw.Name] = wrkflw
		cur++
	}
	parts := make([]Part, 0, len(lines)-len(wrkflws)-1)

	for cur < len(lines) {
		l := lines[cur]
		parts = append(parts, PartFromString(l))
		cur++
	}

	return AoC{
		wrkflws,
		parts,
	}
}

func (a *AoC) CheckPart(p Part) bool {
	workflow := a.Workflows["in"]

outer:
	for {
		for _, r := range workflow.Rules {
			if r.Clause != nil {
				lhsVal := p[r.Clause.Lhs]
				var clauseSatisfied bool
				switch r.Clause.Op {
				case ">":
					clauseSatisfied = lhsVal > r.Clause.Rhs
				case "<":
					clauseSatisfied = lhsVal < r.Clause.Rhs
				default:
					panic("invalid clause operator")
				}

				if clauseSatisfied {
					if len(r.Then) == 1 {
						return r.Then == "A"
					} else {
						workflow = a.Workflows[r.Then]
						continue outer
					}
				}

			} else if len(r.Then) == 1 {
				return r.Then == "A"
			} else {
				workflow = a.Workflows[r.Then]
				continue outer
			}
		}
	}
}

func (a *AoC) part1() any {
	var res int
	for _, p := range a.Parts {
		if a.CheckPart(p) {
			res += p.Sum()
		}
	}
	return res
}

func (a *AoC) CanBeAc(wrk *Workflow, rCur int) bool {
	r := wrk.Rules[rCur]

	if r.Then == "A" {
		return true
	}

	// if implies go to the specific workflow, 
	// then change the workflow

	if len(r.Then) > 1 {
		if a.CanBeAc(a.Workflows[r.Then], 0) {
			return true
		}
	}

	if rCur+1 < len(wrk.Rules) {
		return a.CanBeAc(wrk, rCur+1)
	}

	return false
}

func (a *AoC) printAc(wrk *Workflow) {
	if a.CanBeAc(wrk, 0) {
		fmt.Println(wrk.Rules)
	}
}

func (a *AoC) part2() any {
	cur := a.Workflows["in"]

	a.printAc(cur)
	return 0
}

func main() {
	content, err := os.ReadFile("test-input.txt")
	// content, err := os.ReadFile("input.txt")

	if err != nil {
		log.Fatal(err)
	}

	aoc := parse(string(content))

	fmt.Println("part1:", aoc.part1())
	fmt.Println("part2:", aoc.part2())
}
