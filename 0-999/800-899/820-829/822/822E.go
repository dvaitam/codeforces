package main

import (
	"bufio"
	"fmt"
	"os"
)

type state struct {
	t int
	s int
	c int
}

var (
	n, m, x    int
	sStr, tStr string
	nextPos    [][]int
	pw         []uint64
	hs, ht     []uint64
	memo       map[state]bool
)

const base uint64 = 911382323

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func buildHashes(str string) []uint64 {
	h := make([]uint64, len(str)+1)
	for i := 0; i < len(str); i++ {
		h[i+1] = h[i]*base + uint64(str[i])
	}
	return h
}

func getHash(h []uint64, l, r int) uint64 {
	return h[r] - h[l]*pw[r-l]
}

func lcp(iS, iT int) int {
	hi := n - iS
	if m-iT < hi {
		hi = m - iT
	}
	lo := 0
	for lo < hi {
		mid := (lo + hi + 1) / 2
		if getHash(hs, iS, iS+mid) == getHash(ht, iT, iT+mid) {
			lo = mid
		} else {
			hi = mid - 1
		}
	}
	return lo
}

func solve(posT, posS, cnt int) bool {
	if cnt > x {
		return false
	}
	if posT == m {
		return true
	}
	st := state{posT, posS, cnt}
	if v, ok := memo[st]; ok {
		return v
	}
	c := int(tStr[posT] - 'a')
	for start := nextPos[posS][c]; start < n; start = nextPos[start+1][c] {
		l := lcp(start, posT)
		if l == 0 {
			continue
		}
		if solve(posT+l, start+l, cnt+1) {
			memo[st] = true
			return true
		}
	}
	memo[st] = false
	return false
}

func buildNext() {
	nextPos = make([][]int, n+1)
	for i := range nextPos {
		nextPos[i] = make([]int, 26)
	}
	for c := 0; c < 26; c++ {
		nextPos[n][c] = n
	}
	for i := n - 1; i >= 0; i-- {
		for c := 0; c < 26; c++ {
			nextPos[i][c] = nextPos[i+1][c]
		}
		nextPos[i][int(sStr[i]-'a')] = i
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Fscan(reader, &n)
	fmt.Fscan(reader, &sStr)
	fmt.Fscan(reader, &m)
	fmt.Fscan(reader, &tStr)
	fmt.Fscan(reader, &x)

	limit := max(n, m)
	pw = make([]uint64, limit+1)
	pw[0] = 1
	for i := 1; i <= limit; i++ {
		pw[i] = pw[i-1] * base
	}
	hs = buildHashes(sStr)
	ht = buildHashes(tStr)
	buildNext()
	memo = make(map[state]bool)
	if solve(0, 0, 0) {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
}
