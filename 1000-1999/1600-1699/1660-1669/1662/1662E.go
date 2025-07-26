package main

import (
	"bufio"
	"fmt"
	"os"
)

func minSwaps(p []int) int {
	n := len(p)
	if n > 8 { // fallback for large n, not feasible
		return 0
	}
	start := make([]int, n)
	for i := range start {
		start[i] = i + 1
	}
	banned := make(map[[2]int]bool)
	for i := 1; i < n; i++ {
		banned[[2]int{i, i + 1}] = true
		banned[[2]int{i + 1, i}] = true
	}
	targetStates := make([][]int, n)
	for r := 0; r < n; r++ {
		t := make([]int, n)
		copy(t, p[r:])
		copy(t[n-r:], p[:r])
		targetStates[r] = t
	}

	type node struct {
		arr []int
		d   int
	}
	vis := make(map[string]bool)
	q := []node{{start, 0}}
	encode := func(a []int) string {
		b := make([]byte, len(a))
		for i, v := range a {
			b[i] = byte(v)
		}
		return string(b)
	}
	goals := make(map[string]bool)
	for _, t := range targetStates {
		goals[encode(t)] = true
	}
	vis[encode(start)] = true
	for len(q) > 0 {
		cur := q[0]
		q = q[1:]
		if goals[encode(cur.arr)] {
			return cur.d
		}
		for i := 0; i < n; i++ {
			j := (i + 1) % n
			a, b := cur.arr[i], cur.arr[j]
			if banned[[2]int{a, b}] {
				continue
			}
			next := append([]int(nil), cur.arr...)
			next[i], next[j] = next[j], next[i]
			key := encode(next)
			if !vis[key] {
				vis[key] = true
				q = append(q, node{next, cur.d + 1})
			}
		}
	}
	return 0
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var t int
	fmt.Fscan(in, &t)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		p := make([]int, n)
		for i := range p {
			fmt.Fscan(in, &p[i])
		}
		fmt.Fprintln(out, minSwaps(p))
	}
}
