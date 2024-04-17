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
	Workflows         map[string]*Workflow
	Parts             []Part
	AccRanges         AccRangeMap
	AccRangesComputed bool
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
		nil,
		false,
	}
}

func (a *AoC) CheckPart(p Part) bool {
outer:
	for _, v := range a.computeRanges() {
		for k, r := range v {
			if p[k] < r.Start || p[k] > r.End {
				continue outer
			}
		}

		return true
	}

	return false
}

//func (a *AoC) CheckPart(p Part) bool {
//workflow := a.Workflows["in"]

//outer:
//for {
//for _, r := range workflow.Rules {
//if r.Clause != nil {
//lhsVal := p[r.Clause.Lhs]
//var clauseSatisfied bool
//switch r.Clause.Op {
//case ">":
//clauseSatisfied = lhsVal > r.Clause.Rhs
//case "<":
//clauseSatisfied = lhsVal < r.Clause.Rhs
//default:
//panic("invalid clause operator")
//}

//if clauseSatisfied {
//if len(r.Then) == 1 {
//return r.Then == "A"
//} else {
//workflow = a.Workflows[r.Then]
//continue outer
//}
//}

//} else if len(r.Then) == 1 {
//return r.Then == "A"
//} else {
//workflow = a.Workflows[r.Then]
//continue outer
//}
//}
//}
//}

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

type AccRangeMap map[string]map[string]*Range

type Range struct {
	Start int
	End   int
}

func NewRanges() map[string]*Range {
	return map[string]*Range{
		"x": {1, 4000},
		"m": {1, 4000},
		"a": {1, 4000},
		"s": {1, 4000},
	}
}

func (r *Range) String() string {
	return fmt.Sprintf("(%d..=%d)", r.Start, r.End)
}

func (a *AoC) unravelEnter(w *Workflow, rIdx int, rng map[string]*Range) {
	r := w.Rules[rIdx]

	if r.Clause != nil {
		if r.Clause.Op == ">" {
			// m > 1548 --> m <= 2006
			rng[r.Clause.Lhs].End = min(rng[r.Clause.Lhs].End, r.Clause.Rhs)
		} else if r.Clause.Op == "<" {
			// a < 2006 --> a >= 2006
			rng[r.Clause.Lhs].Start = max(rng[r.Clause.Lhs].Start, r.Clause.Rhs)
		} else {
			panic("unexpected rule clause operator")
		}
	}

	if rIdx > 0 {
		a.unravelEnter(w, rIdx-1, rng)
		return
	}

	if w.Name == "in" {
		// reached the enter workflow
		return
	}

	for _, pw := range a.Workflows {
		for pIdx, pr := range pw.Rules {
			if pr.Then == w.Name {
				a.unravel(pw, pIdx, rng)
				return
			}
		}
	}

	panic(fmt.Sprintf("could not find previous workflow for %s", w.Name))
}

func (a *AoC) unravel(w *Workflow, rIdx int, rng map[string]*Range) {
	r := w.Rules[rIdx]

	if r.Clause != nil {
		if r.Clause.Op == ">" {
			rng[r.Clause.Lhs].Start = max(rng[r.Clause.Lhs].Start, r.Clause.Rhs+1)
		} else if r.Clause.Op == "<" {
			rng[r.Clause.Lhs].End = min(rng[r.Clause.Lhs].End, r.Clause.Rhs-1)
		} else {
			panic("unexpected rule clause operator")
		}
	}

	if rIdx > 0 {
		a.unravelEnter(w, rIdx-1, rng)
	} else if w.Name != "in" {
		for _, pw := range a.Workflows {
			for pIdx, pr := range pw.Rules {
				if pr.Then == w.Name {
					a.unravel(pw, pIdx, rng)
					return
				}
			}
		}

		panic("could not fidenter workflow")
	}
}

func (a *AoC) computeRanges() AccRangeMap {
	if a.AccRangesComputed {
		return a.AccRanges
	}

	a.AccRanges = make(AccRangeMap)

	for _, w := range a.Workflows {
		for rIdx, r := range w.Rules {
			if r.Then != "A" {
				continue
			}
			rRng := NewRanges()
			a.unravel(w, rIdx, rRng)
			a.AccRanges[fmt.Sprintf("%s:%d", w.Name, rIdx)] = rRng
		}
	}

	a.AccRangesComputed = true
	return a.AccRanges
}

func (a *AoC) part2() any {
	var sum int
	for _, rr := range a.computeRanges() {
		prod := 1
		for _, r := range rr {
			prod *= (r.End - r.Start + 1)
		}
		sum += prod
	}

	return sum
}

func main() {
	content, err := os.ReadFile("input.txt")
	// content, err := os.ReadFile("input.txt")

	if err != nil {
		log.Fatal(err)
	}

	aoc := parse(string(content))

	fmt.Println("part1:", aoc.part1())
	fmt.Println("part2:", aoc.part2())
}
