package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var N, M int
	fmt.Fscan(reader, &N, &M)
	P := make([]int, N+1)
	for i := 1; i <= N; i++ {
		fmt.Fscan(reader, &P[i])
	}
	C := make([]int, N+1)
	for i := 1; i <= N; i++ {
		fmt.Fscan(reader, &C[i])
	}
	var D int
	fmt.Fscan(reader, &D)
	K := make([]int, N+1)
	found := make([]bool, N+1)
	for i := 1; i <= D; i++ {
		fmt.Fscan(reader, &K[i])
		found[K[i]] = true
	}
	newD := D
	for i := 1; i <= N; i++ {
		if !found[i] {
			newD++
			K[newD] = i
		}
	}
	mm := M
	for i := 1; i <= N; i++ {
		if P[i] > mm {
			mm = P[i]
		}
		if C[i] > mm {
			mm = C[i]
		}
	}
	G := make([][]int, mm+1)
	cnt := make([]bool, (mm+1)*(mm+1))
	used := make([]bool, mm+1)
	Le := make([]int, mm+1)
	Ri := make([]int, mm+1)
	for i := range Le {
		Le[i] = -1
	}
	var pairUp func(node int) bool
	pairUp = func(node int) bool {
		if used[node] {
			return false
		}
		used[node] = true
		for _, to := range G[node] {
			if Le[to] == -1 {
				Le[to] = node
				Ri[node] = to
				return true
			}
		}
		for _, to := range G[node] {
			if pairUp(Le[to]) {
				Le[to] = node
				Ri[node] = to
				return true
			}
		}
		return false
	}
	ans := make([]int, D+1)
	L := -1
	for i := newD; i >= 1; i-- {
		if i <= D {
			ans[i] = L + 1
		}
		p := P[K[i]]
		c := C[K[i]]
		idx := c*(mm+1) + p
		if !cnt[idx] {
			G[p] = append(G[p], c)
		}
		for ok := true; ok; {
			ok = false
			for j := range used {
				used[j] = false
			}
			if pairUp(L + 1) {
				L++
				ok = true
			}
		}
		cnt[idx] = true
	}
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	for i := 1; i <= D; i++ {
		fmt.Fprintln(writer, ans[i])
	}
}
