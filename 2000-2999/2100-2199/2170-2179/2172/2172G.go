package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353

const (
	genA = 0
	genB = 1
	genb = 2
)

var inv = [3]int{0, 2, 1}

type groupInfo struct {
	nxtA []int
	nxtB []int
	adj  [][]int64
}

type enumerator struct {
	next   [][]int
	parent []int
	size   []int
	queue  []int
}

func newEnumerator() *enumerator {
	e := &enumerator{}
	e.makeCoset()
	return e
}

func (e *enumerator) makeCoset() int {
	id := len(e.next)
	e.next = append(e.next, []int{-1, -1, -1})
	e.parent = append(e.parent, id)
	e.size = append(e.size, 1)
	e.queue = append(e.queue, id)
	return id
}

func (e *enumerator) find(x int) int {
	if e.parent[x] != x {
		e.parent[x] = e.find(e.parent[x])
	}
	return e.parent[x]
}

func (e *enumerator) link(a, g, b int) {
	ra := e.find(a)
	rb := e.find(b)
	if ta := e.next[ra][g]; ta != -1 {
		e.union(ta, rb)
		return
	}
	e.next[ra][g] = rb
	ig := inv[g]
	if tb := e.next[rb][ig]; tb != -1 {
		e.union(tb, ra)
	} else {
		e.next[rb][ig] = ra
	}
}

func (e *enumerator) union(a, b int) int {
	ra := e.find(a)
	rb := e.find(b)
	if ra == rb {
		return ra
	}
	if e.size[ra] < e.size[rb] {
		ra, rb = rb, ra
	}
	e.parent[rb] = ra
	e.size[ra] += e.size[rb]
	for g := 0; g < 3; g++ {
		if ta := e.next[ra][g]; ta != -1 {
			e.next[ra][g] = e.find(ta)
		}
		if tb := e.next[rb][g]; tb != -1 {
			tb = e.find(tb)
			if e.next[ra][g] == -1 {
				e.link(ra, g, tb)
			} else {
				e.union(e.next[ra][g], tb)
			}
		}
	}
	e.queue = append(e.queue, ra)
	return ra
}

func (e *enumerator) step(a, g int) int {
	ra := e.find(a)
	if tgt := e.next[ra][g]; tgt != -1 {
		rt := e.find(tgt)
		e.next[ra][g] = rt
		return rt
	}
	d := e.makeCoset()
	e.link(ra, g, d)
	return e.find(e.next[ra][g])
}

func (e *enumerator) impose(start int, rel []int) {
	cur := e.find(start)
	for _, g := range rel {
		cur = e.step(cur, g)
	}
	e.union(cur, start)
}

func buildGroup(s string) *groupInfo {
	e := newEnumerator()
	relS := make([]int, len(s))
	for i := range s {
		if s[i] == 'A' {
			relS[i] = genA
		} else {
			relS[i] = genB
		}
	}
	relations := [][]int{
		{genA, genA},
		{genB, genB, genB},
		{genB, genb},
		{genb, genB},
		relS,
	}
	for head := 0; head < len(e.queue); head++ {
		root := e.find(e.queue[head])
		for _, rel := range relations {
			e.impose(root, rel)
		}
		e.step(root, genA)
		e.step(root, genB)
	}
	rootMap := make(map[int]int)
	identity := e.find(0)
	rootMap[identity] = 0
	roots := []int{identity}
	for i := range e.next {
		r := e.find(i)
		if _, ok := rootMap[r]; !ok {
			rootMap[r] = len(roots)
			roots = append(roots, r)
		}
	}
	m := len(roots)
	nxtA := make([]int, m)
	nxtB := make([]int, m)
	for idx, r := range roots {
		a := e.next[r][genA]
		b := e.next[r][genB]
		nxtA[idx] = rootMap[e.find(a)]
		nxtB[idx] = rootMap[e.find(b)]
	}
	adj := make([][]int64, m)
	for i := range adj {
		adj[i] = make([]int64, m)
	}
	for i := 0; i < m; i++ {
		adj[nxtA[i]][i] = (adj[nxtA[i]][i] + 1) % mod
		adj[nxtB[i]][i] = (adj[nxtB[i]][i] + 1) % mod
	}
	return &groupInfo{
		nxtA: nxtA,
		nxtB: nxtB,
		adj:  adj,
	}
}

func matMulVec(mat [][]int64, vec []int64) []int64 {
	n := len(mat)
	res := make([]int64, n)
	for i := 0; i < n; i++ {
		var sum int64
		row := mat[i]
		for j := 0; j < n; j++ {
			if row[j] == 0 || vec[j] == 0 {
				continue
			}
			sum += row[j] * vec[j] % mod
			if sum >= (1 << 62) {
				sum %= mod
			}
		}
		res[i] = sum % mod
	}
	return res
}

func matMul(a, b [][]int64) [][]int64 {
	n := len(a)
	res := make([][]int64, n)
	for i := 0; i < n; i++ {
		res[i] = make([]int64, n)
		for k := 0; k < n; k++ {
			if a[i][k] == 0 {
				continue
			}
			val := a[i][k]
			row := res[i]
			col := b[k]
			for j := 0; j < n; j++ {
				if col[j] == 0 {
					continue
				}
				row[j] = (row[j] + val*col[j]) % mod
			}
		}
	}
	return res
}

func powerVec(mat [][]int64, power int64, start int) []int64 {
	n := len(mat)
	resVec := make([]int64, n)
	resVec[start] = 1
	if power == 0 {
		return resVec
	}
	trans := mat
	for power > 0 {
		if power&1 == 1 {
			resVec = matMulVec(trans, resVec)
		}
		power >>= 1
		if power > 0 {
			trans = matMul(trans, trans)
		}
	}
	return resVec
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var q int
	if _, err := fmt.Fscan(in, &q); err != nil {
		return
	}
	cache := make(map[string]*groupInfo)
	for ; q > 0; q-- {
		var s, t string
		var n int64
		fmt.Fscan(in, &s)
		fmt.Fscan(in, &t)
		fmt.Fscan(in, &n)
		grp, ok := cache[s]
		if !ok {
			grp = buildGroup(s)
			cache[s] = grp
		}
		state := 0
		for i := 0; i < len(t); i++ {
			if t[i] == 'A' {
				state = grp.nxtA[state]
			} else {
				state = grp.nxtB[state]
			}
		}
		vec := powerVec(grp.adj, n, 0)
		fmt.Fprintln(out, vec[state]%mod)
	}
}
