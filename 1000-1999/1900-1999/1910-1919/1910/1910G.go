package main

import (
	"bufio"
	"fmt"
	"os"
)

func isValid(t []int64) bool {
	n := len(t)
	q := t[0]
	p := int64(-1)
	for _, v := range t {
		if v%q != 0 {
			p = v
			break
		}
	}
	if p == -1 {
		for i := 0; i < n; i++ {
			if t[i] != int64(i+1)*q {
				return false
			}
		}
		return true
	}
	if p <= q {
		return false
	}
	idx := 0
	pq := int64(1)
	pp := int64(1)
	for idx < n {
		nextQ := pq * q
		nextP := pp * p
		next := nextQ
		if nextP < next {
			next = nextP
		}
		if t[idx] != next {
			return false
		}
		if next == nextQ {
			pq++
		}
		if next == nextP {
			pp++
		}
		idx++
	}
	nextQ := pq * q
	nextP := pp * p
	if nextP < nextQ {
		nextQ = nextP
	}
	return nextQ > t[n-1]
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var c int
	fmt.Fscan(in, &c)
	for ; c > 0; c-- {
		var n int
		fmt.Fscan(in, &n)
		t := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &t[i])
		}
		if isValid(t) {
			fmt.Fprintln(out, "VALID")
		} else {
			fmt.Fprintln(out, "INVALID")
		}
	}
}
