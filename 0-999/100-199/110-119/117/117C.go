package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	a       [][]bool
	l       []int
	t       []int
	n       int
	A, B, C int
)

// find searches for a 3-cycle in the tournament graph between indices [p, q)
func find(p, q int) bool {
	if q-p <= 2 {
		return false
	}
	v := l[p]
	t[p] = v
	ef1, ef2 := p+1, q-1
	// partition by edge from v
	for i := p + 1; i < q; i++ {
		u := l[i]
		if a[v][u] {
			t[ef2] = u
			ef2--
		} else {
			t[ef1] = u
			ef1++
		}
	}
	// check for cycle v->B->C->v
	for i := p + 1; i < ef1; i++ {
		c := t[i]
		for j := ef2 + 1; j < q; j++ {
			b := t[j]
			if a[b][c] {
				A = v
				B = b
				C = c
				return true
			}
		}
	}
	// reorder l for recursion
	for i := p; i < q; i++ {
		l[i] = t[i]
	}
	if find(p+1, ef1) {
		return true
	}
	if find(ef2+1, q) {
		return true
	}
	return false
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	a = make([][]bool, n)
	for i := range a {
		a[i] = make([]bool, n)
	}
	l = make([]int, n)
	t = make([]int, n)
	// read adjacency matrix
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			var ch byte
			for {
				b, err := reader.ReadByte()
				if err != nil {
					return
				}
				if b == '0' || b == '1' {
					ch = b
					break
				}
			}
			a[i][j] = (ch == '1')
		}
		l[i] = i
	}
	if find(0, n) {
		fmt.Println(A+1, B+1, C+1)
	} else {
		fmt.Println(-1)
	}
}
