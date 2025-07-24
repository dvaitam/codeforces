package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int = 1000000007

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	type Pair struct{ X, Y int }
	pairs := make([]Pair, n)
	seatIndex := make(map[int]int)
	idx := 0
	for i := 0; i < n; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		pairs[i] = Pair{x, y}
		if _, ok := seatIndex[x]; !ok {
			seatIndex[x] = idx
			idx++
		}
		if _, ok := seatIndex[y]; !ok {
			seatIndex[y] = idx
			idx++
		}
	}
	m := idx
	parent := make([]int, n+m)
	size := make([]int, n+m)
	for i := range parent {
		parent[i] = i
		size[i] = 1
	}
	var find func(int) int
	find = func(x int) int {
		for parent[x] != x {
			parent[x] = parent[parent[x]]
			x = parent[x]
		}
		return x
	}
	union := func(a, b int) {
		a = find(a)
		b = find(b)
		if a == b {
			return
		}
		if size[a] < size[b] {
			a, b = b, a
		}
		parent[b] = a
		size[a] += size[b]
	}
	for i, p := range pairs {
		union(i, n+seatIndex[p.X])
		union(i, n+seatIndex[p.Y])
	}

	type comp struct {
		t, s int
		self bool
	}
	comps := make(map[int]*comp)
	for i, p := range pairs {
		r := find(i)
		if comps[r] == nil {
			comps[r] = &comp{}
		}
		c := comps[r]
		c.t++
		if p.X == p.Y {
			c.self = true
		}
	}
	for _, idx := range seatIndex {
		r := find(n + idx)
		if comps[r] == nil {
			comps[r] = &comp{}
		}
		comps[r].s++
	}

	res := 1
	for _, c := range comps {
		if c.s == c.t+1 {
			res = res * c.s % MOD
		} else if c.s == c.t {
			if c.self {
				res = res * 1 % MOD
			} else {
				res = res * 2 % MOD
			}
		}
	}
	fmt.Println(res % MOD)
}
