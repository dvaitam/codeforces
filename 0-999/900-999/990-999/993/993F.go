package main

import (
	"bufio"
	"fmt"
	"os"
)

// gate types
const (
	AND = iota
	OR
	NAND
	NOR
)

func gateEval(t int, a, b bool) bool {
	switch t {
	case AND:
		return a && b
	case OR:
		return a || b
	case NAND:
		return !(a && b)
	case NOR:
		return !(a || b)
	}
	return false
}

func parseGateType(s string) int {
	switch s {
	case "and":
		return AND
	case "or":
		return OR
	case "nand":
		return NAND
	case "nor":
		return NOR
	default:
		return AND
	}
}

type FirstGate struct {
	typ  int
	a, b int
}

type SecondGate struct {
	typ  int
	a, b int
}

var n, m, k int
var first []FirstGate
var second []SecondGate

func compute(features []bool, keep []bool, love bool) bool {
	fl := make([]bool, m)
	for i, g := range first {
		x := features[g.a]
		y := features[g.b]
		val := gateEval(g.typ, x, y)
		if love {
			val = !val
		}
		fl[i] = val
	}
	out := false
	for i, g := range second {
		if !keep[i] {
			continue
		}
		x := fl[g.a]
		y := fl[g.b]
		val := gateEval(g.typ, x, y)
		if love {
			val = !val
		}
		if val {
			out = true
		}
	}
	if love {
		return !out
	}
	return out
}

func checkSubset(keep []bool) bool {
	// gather used first layer gates and features
	usedFirst := make(map[int]struct{})
	for i, g := range second {
		if keep[i] {
			usedFirst[g.a] = struct{}{}
			usedFirst[g.b] = struct{}{}
		}
	}
	if len(usedFirst) == 0 {
		return false
	}
	featureSet := make(map[int]struct{})
	for idx := range usedFirst {
		fg := first[idx]
		featureSet[fg.a] = struct{}{}
		featureSet[fg.b] = struct{}{}
	}
	features := make([]int, 0, len(featureSet))
	for f := range featureSet {
		features = append(features, f)
	}
	if len(features) > 20 {
		return false
	}
	total := 1 << len(features)
	featVals := make([]bool, n)
	for mask := 0; mask < total; mask++ {
		for i, f := range features {
			featVals[f] = (mask>>i)&1 == 1
		}
		normal := compute(featVals, keep, false)
		love := compute(featVals, keep, true)
		if normal != love {
			return false
		}
	}
	return true
}

func nextComb(c []int, n int) bool {
	k := len(c)
	i := k - 1
	for i >= 0 && c[i] == n-k+i {
		i--
	}
	if i < 0 {
		return false
	}
	c[i]++
	for j := i + 1; j < k; j++ {
		c[j] = c[j-1] + 1
	}
	return true
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Fscan(reader, &n, &m, &k)
	first = make([]FirstGate, m)
	for i := 0; i < m; i++ {
		var typ, conn string
		fmt.Fscan(reader, &typ, &conn)
		p1, p2 := -1, -1
		for j, ch := range conn {
			if ch == 'x' {
				if p1 == -1 {
					p1 = j
				} else {
					p2 = j
				}
			}
		}
		first[i] = FirstGate{parseGateType(typ), p1, p2}
	}
	second = make([]SecondGate, k)
	for i := 0; i < k; i++ {
		var typ, conn string
		fmt.Fscan(reader, &typ, &conn)
		p1, p2 := -1, -1
		for j, ch := range conn {
			if ch == 'x' {
				if p1 == -1 {
					p1 = j
				} else {
					p2 = j
				}
			}
		}
		second[i] = SecondGate{parseGateType(typ), p1, p2}
	}

	keepAll := make([]bool, k)
	for i := range keepAll {
		keepAll[i] = true
	}
	if checkSubset(keepAll) {
		fmt.Println(0)
		return
	}
	maxRemove := 5
	if k < maxRemove {
		maxRemove = k
	}
	for d := 1; d <= maxRemove; d++ {
		comb := make([]int, d)
		for i := 0; i < d; i++ {
			comb[i] = i
		}
		for {
			keep := make([]bool, k)
			for i := 0; i < k; i++ {
				keep[i] = true
			}
			for _, rm := range comb {
				keep[rm] = false
			}
			if checkSubset(keep) {
				fmt.Println(d)
				return
			}
			if !nextComb(comb, k) {
				break
			}
		}
	}
	fmt.Println(-1)
}
