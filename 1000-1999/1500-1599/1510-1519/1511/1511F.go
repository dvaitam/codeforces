package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD = 998244353

type pair struct{ a, b int }

type Automaton struct {
	trans []map[byte][]int
}

func buildAutomaton(words []string) Automaton {
	id := map[string]int{"": 0}
	suffixes := []string{""}
	for _, w := range words {
		for i := 0; i < len(w); i++ {
			s := w[i:]
			if _, ok := id[s]; !ok {
				id[s] = len(suffixes)
				suffixes = append(suffixes, s)
			}
		}
	}
	trans := make([]map[byte][]int, len(suffixes))
	for i := range trans {
		trans[i] = make(map[byte][]int)
	}
	for i, s := range suffixes {
		if s == "" {
			for _, w := range words {
				ch := w[0]
				dest := id[w[1:]]
				trans[i][ch] = append(trans[i][ch], dest)
			}
		} else {
			ch := s[0]
			dest := id[s[1:]]
			trans[i][ch] = append(trans[i][ch], dest)
		}
	}
	return Automaton{trans}
}

func count(words []string, m int64) int64 {
	auto := buildAutomaton(words)
	start := 0
	pairID := map[pair]int{{start, start}: 0}
	pairs := []pair{{start, start}}
	edges := []map[int]int64{make(map[int]int64)}
	queue := []int{0}
	for len(queue) > 0 {
		id := queue[0]
		queue = queue[1:]
		p := pairs[id]
		t1 := auto.trans[p.a]
		t2 := auto.trans[p.b]
		for ch, l1 := range t1 {
			l2, ok := t2[ch]
			if !ok {
				continue
			}
			for _, d1 := range l1 {
				for _, d2 := range l2 {
					np := pair{d1, d2}
					idx, ok := pairID[np]
					if !ok {
						idx = len(pairs)
						pairID[np] = idx
						pairs = append(pairs, np)
						edges = append(edges, make(map[int]int64))
						queue = append(queue, idx)
					}
					edges[id][idx] = (edges[id][idx] + 1) % MOD
				}
			}
		}
	}
	n := len(pairs)
	mat := make([][]int64, n)
	for i := 0; i < n; i++ {
		mat[i] = make([]int64, n)
		for j, v := range edges[i] {
			mat[i][j] = v % MOD
		}
	}
	res := make([][]int64, n)
	for i := range res {
		res[i] = make([]int64, n)
		res[i][i] = 1
	}
	base := mat
	e := m
	mul := func(a, b [][]int64) [][]int64 {
		n := len(a)
		out := make([][]int64, n)
		for i := range out {
			out[i] = make([]int64, n)
		}
		for i := 0; i < n; i++ {
			for k := 0; k < n; k++ {
				if a[i][k] == 0 {
					continue
				}
				ak := a[i][k]
				for j := 0; j < n; j++ {
					if b[k][j] == 0 {
						continue
					}
					out[i][j] = (out[i][j] + ak*b[k][j]) % MOD
				}
			}
		}
		return out
	}
	for e > 0 {
		if e&1 == 1 {
			res = mul(res, base)
		}
		base = mul(base, base)
		e >>= 1
	}
	return res[0][0]
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var n int
	var m int64
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	words := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &words[i])
	}
	fmt.Fprintln(writer, count(words, m))
}
