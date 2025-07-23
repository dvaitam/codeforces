package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Hash struct {
	a uint64
	b uint64
}

type Pair struct {
	ch   byte
	hash Hash
	idx  int
}

const (
	base1 uint64 = 911382323
	base2 uint64 = 972663749
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	c := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &c[i])
	}

	var s string
	fmt.Fscan(in, &s)
	letters := []byte(s)
	if len(letters) != n {
		// just in case
		letters = append(letters, make([]byte, n-len(letters))...)
	}

	g := make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}

	parent := make([]int, n+1)
	order := make([]int, 0, n)
	stack := []int{1}
	parent[1] = 0
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		order = append(order, v)
		for _, to := range g[v] {
			if to == parent[v] {
				continue
			}
			parent[to] = v
			stack = append(stack, to)
		}
	}
	// process in reverse order
	hashVal := make([]Hash, n+1)
	dif := make([]int64, n+1)
	for i := len(order) - 1; i >= 0; i-- {
		v := order[i]
		pairs := make([]Pair, 0, len(g[v]))
		for _, to := range g[v] {
			if to == parent[v] {
				continue
			}
			pairs = append(pairs, Pair{ch: letters[to-1], hash: hashVal[to], idx: to})
		}
		sort.Slice(pairs, func(i, j int) bool {
			if pairs[i].ch != pairs[j].ch {
				return pairs[i].ch < pairs[j].ch
			}
			if pairs[i].hash.a != pairs[j].hash.a {
				return pairs[i].hash.a < pairs[j].hash.a
			}
			return pairs[i].hash.b < pairs[j].hash.b
		})
		uniq := make([]Pair, 0, len(pairs))
		for j := 0; j < len(pairs); j++ {
			if j == 0 || pairs[j].ch != pairs[j-1].ch || pairs[j].hash.a != pairs[j-1].hash.a || pairs[j].hash.b != pairs[j-1].hash.b {
				uniq = append(uniq, pairs[j])
			}
		}
		dv := int64(1)
		for _, p := range uniq {
			dv += dif[p.idx]
		}
		dif[v] = dv
		h := Hash{a: uint64(letters[v-1]) + 1, b: uint64(letters[v-1]) + 1}
		for _, p := range uniq {
			h.a = h.a*base1 + uint64(p.ch) + 1
			h.a = h.a*base1 + p.hash.a
			h.b = h.b*base2 + uint64(p.ch) + 1
			h.b = h.b*base2 + p.hash.b
		}
		h.a = h.a*base1 + 7
		h.b = h.b*base2 + 7
		hashVal[v] = h
	}

	var maxVal int64 = -1
	var count int64
	for i := 1; i <= n; i++ {
		val := dif[i] + c[i]
		if val > maxVal {
			maxVal = val
			count = 1
		} else if val == maxVal {
			count++
		}
	}
	out := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(out, maxVal)
	fmt.Fprintln(out, count)
	out.Flush()
}
