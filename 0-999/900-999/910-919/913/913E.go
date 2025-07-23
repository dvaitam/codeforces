package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
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

func main() {
	// Precompute masks for variables x,y,z
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
			// F -> T
			tryUpdate(1, mask, expr)
			// !F
			tryUpdate(0, (^mask)&255, "!"+expr)
			// combine with existing T on the left
			for m2 := 0; m2 < 256; m2++ {
				if bestT[m2] != "" {
					tryUpdate(1, m2&mask, bestT[m2]+"&"+expr)
				}
			}
		case 1:
			if bestT[mask] != expr {
				continue
			}
			// T -> E
			tryUpdate(2, mask, expr)
			// combine with existing F on the right
			for m2 := 0; m2 < 256; m2++ {
				if bestF[m2] != "" {
					tryUpdate(1, mask&m2, expr+"&"+bestF[m2])
				}
			}
			// combine with existing E on the left for OR
			for m2 := 0; m2 < 256; m2++ {
				if bestE[m2] != "" {
					tryUpdate(2, m2|mask, bestE[m2]+"|"+expr)
				}
			}
		case 2:
			if bestE[mask] != expr {
				continue
			}
			// (E) -> F
			tryUpdate(0, mask, "("+expr+")")
			// combine with existing T on the right using OR
			for m2 := 0; m2 < 256; m2++ {
				if bestT[m2] != "" {
					tryUpdate(2, mask|m2, expr+"|"+bestT[m2])
				}
			}
		}
	}

	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	fmt.Fscan(reader, &n)
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(reader, &s)
		mask := 0
		for j := 0; j < 8; j++ {
			if s[j] == '1' {
				mask |= 1 << j
			}
		}
		fmt.Fprintln(writer, bestE[mask])
	}
}
