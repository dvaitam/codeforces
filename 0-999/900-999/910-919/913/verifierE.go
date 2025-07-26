package main

import (
	"container/heap"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Node struct {
	expr string
	mask int
	typ  int // 0=F,1=T,2=E
}

type PQ []Node

func (pq PQ) Len() int { return len(pq) }
func (pq PQ) Less(i, j int) bool {
	if len(pq[i].expr) != len(pq[j].expr) {
		return len(pq[i].expr) < len(pq[j].expr)
	}
	return pq[i].expr < pq[j].expr
}
func (pq PQ) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PQ) Push(x interface{}) { *pq = append(*pq, x.(Node)) }
func (pq *PQ) Pop() interface{} {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[:n-1]
	return x
}

var bestF [256]string
var bestT [256]string
var bestE [256]string
var pq PQ

func better(old, new string) bool {
	if old == "" {
		return true
	}
	if len(new) != len(old) {
		return len(new) < len(old)
	}
	return new < old
}

func tryUpdate(typ, mask int, expr string) {
	switch typ {
	case 0:
		if !better(bestF[mask], expr) {
			return
		}
		bestF[mask] = expr
	case 1:
		if !better(bestT[mask], expr) {
			return
		}
		bestT[mask] = expr
	case 2:
		if !better(bestE[mask], expr) {
			return
		}
		bestE[mask] = expr
	}
	heap.Push(&pq, Node{expr: expr, mask: mask, typ: typ})
}

func precompute() {
	maskX, maskY, maskZ := 0, 0, 0
	for i := 0; i < 8; i++ {
		if (i>>2)&1 == 1 {
			maskX |= 1 << i
		}
		if (i>>1)&1 == 1 {
			maskY |= 1 << i
		}
		if i&1 == 1 {
			maskZ |= 1 << i
		}
	}
	heap.Init(&pq)
	tryUpdate(0, maskX, "x")
	tryUpdate(0, maskY, "y")
	tryUpdate(0, maskZ, "z")
	for pq.Len() > 0 {
		node := heap.Pop(&pq).(Node)
		expr := node.expr
		mask := node.mask
		typ := node.typ
		switch typ {
		case 0:
			if bestF[mask] != expr {
				continue
			}
			tryUpdate(1, mask, expr)
			tryUpdate(0, (^mask)&255, "!"+expr)
			for m2 := 0; m2 < 256; m2++ {
				if bestT[m2] != "" {
					tryUpdate(1, m2&mask, bestT[m2]+"&"+expr)
				}
			}
		case 1:
			if bestT[mask] != expr {
				continue
			}
			tryUpdate(2, mask, expr)
			for m2 := 0; m2 < 256; m2++ {
				if bestF[m2] != "" {
					tryUpdate(1, mask&m2, expr+"&"+bestF[m2])
				}
			}
			for m2 := 0; m2 < 256; m2++ {
				if bestE[m2] != "" {
					tryUpdate(2, m2|mask, bestE[m2]+"|"+expr)
				}
			}
		case 2:
			if bestE[mask] != expr {
				continue
			}
			tryUpdate(0, mask, "("+expr+")")
			for m2 := 0; m2 < 256; m2++ {
				if bestT[m2] != "" {
					tryUpdate(2, mask|m2, expr+"|"+bestT[m2])
				}
			}
		}
	}
}

func randomInput() []string {
	n := 1 + rand.Intn(5)
	res := make([]string, n)
	for i := 0; i < n; i++ {
		var sb strings.Builder
		for j := 0; j < 8; j++ {
			if rand.Intn(2) == 1 {
				sb.WriteByte('1')
			} else {
				sb.WriteByte('0')
			}
		}
		res[i] = sb.String()
	}
	return res
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(1)
	precompute()
	const cases = 100
	for i := 0; i < cases; i++ {
		fs := randomInput()
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", len(fs)))
		for _, s := range fs {
			sb.WriteString(s)
			sb.WriteByte('\n')
		}
		input := sb.String()
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n", i+1, err)
			fmt.Printf("program output:\n%s\n", string(out))
			return
		}
		gotLines := strings.Split(strings.TrimSpace(string(out)), "\n")
		if len(gotLines) != len(fs) {
			fmt.Printf("case %d failed: wrong number of lines\ninput:\n%soutput:\n%s", i+1, input, string(out))
			return
		}
		for j, s := range fs {
			mask := 0
			for k := 0; k < 8; k++ {
				if s[k] == '1' {
					mask |= 1 << k
				}
			}
			want := bestE[mask]
			if strings.TrimSpace(gotLines[j]) != want {
				fmt.Printf("case %d failed: line %d expected %s got %s\n", i+1, j+1, want, gotLines[j])
				fmt.Printf("input:\n%s", input)
				return
			}
		}
	}
	fmt.Printf("OK %d cases\n", cases)
}
