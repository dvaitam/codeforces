package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// fast integer reader
func readInt(r *bufio.Reader) int {
	b, _ := r.ReadByte()
	for b < '0' || b > '9' {
		b, _ = r.ReadByte()
	}
	x := 0
	for b >= '0' && b <= '9' {
		x = x*10 + int(b-'0')
		b, _ = r.ReadByte()
	}
	return x
}

func main() {
	r := bufio.NewReader(os.Stdin)
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	n := readInt(r)
	c := make([]int, n)
	for i := 0; i < n; i++ {
		c[i] = readInt(r)
	}
	// indices for sorting
	q := make([]int, n)
	for i := range q {
		q[i] = i
	}
	sort.Slice(q, func(i, j int) bool { return c[q[i]] < c[q[j]] })
	// compress values
	a := make([]int, n)
	now := 1
	for i, idx := range q {
		if i > 0 && c[idx] != c[q[i-1]] {
			now++
		}
		a[idx] = now
	}
	// initialize arrays
	fly := make([]int, n)
	for i := 0; i < n; i++ {
		fly[i] = i - 1
	}
	co := make([]int, now+2)
	mark := make([]int, now+2)
	for i := range mark {
		mark[i] = -1
	}
	st := make([]int, 0, n)
	pop := func() {
		for _, v := range st {
			co[v] = 0
			mark[v] = -1
		}
		st = st[:0]
	}
	var ans [][2]int
	// main logic
	for i := 0; i < n; i++ {
		ai := a[i]
		st = append(st, ai)
		co[ai]++
		if mark[ai] != -1 {
			ans = append(ans, [2]int{mark[ai], i})
			pop()
		} else if co[ai] == 4 {
			ans = append(ans, [2]int{i, i})
			pop()
		} else if co[ai] >= 2 {
			j := i - 1
			for j >= 0 {
				if a[j] == ai {
					fly[i] = j
					break
				}
				mark[a[j]] = i
				j = fly[j]
			}
		}
	}
	// output result
	total := len(ans) * 4
	fmt.Fprintf(w, "%d\n", total)
	for i, p := range ans {
		x := c[p[0]]
		y := c[p[1]]
		fmt.Fprintf(w, "%d %d %d %d", x, y, x, y)
		if i+1 < len(ans) {
			w.WriteByte(' ')
		} else {
			w.WriteByte('\n')
		}
	}
}
