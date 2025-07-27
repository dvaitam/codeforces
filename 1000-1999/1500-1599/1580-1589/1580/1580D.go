package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type State struct {
	sum   int64
	score int64
}

func prune(states []State) []State {
	if len(states) == 0 {
		return states
	}
	sort.Slice(states, func(i, j int) bool {
		if states[i].sum == states[j].sum {
			return states[i].score > states[j].score
		}
		return states[i].sum < states[j].sum
	})
	res := make([]State, 0, len(states))
	maxScore := int64(-1 << 63)
	for _, st := range states {
		if st.score > maxScore {
			res = append(res, st)
			maxScore = st.score
		}
	}
	return res
}

var (
	n, m       int
	a          []int64
	leftChild  []int
	rightChild []int
)

func buildCartesian() int {
	st := []int{}
	root := 0
	for i := 1; i <= n; i++ {
		last := 0
		for len(st) > 0 && a[st[len(st)-1]] > a[i] {
			last = st[len(st)-1]
			st = st[:len(st)-1]
		}
		if len(st) > 0 {
			rightChild[st[len(st)-1]] = i
		} else {
			root = i
		}
		leftChild[i] = last
		st = append(st, i)
	}
	return root
}

func dfs(x int) [][]State {
	if x == 0 {
		return [][]State{{{0, 0}}}
	}
	L := dfs(leftChild[x])
	R := dfs(rightChild[x])
	nL := len(L) - 1
	nR := len(R) - 1
	maxSize := nL + nR + 1
	if maxSize > m {
		maxSize = m
	}
	res := make([][]State, maxSize+1)
	val := a[x]
	for i := 0; i <= nL && i <= m; i++ {
		for j := 0; j <= nR && i+j <= m; j++ {
			for _, sL := range L[i] {
				for _, sR := range R[j] {
					k := i + j
					if k <= m {
						sum := sL.sum + sR.sum
						score := sL.score + sR.score + int64(j)*sL.sum + int64(i)*sR.sum - 2*val*int64(i)*int64(j)
						res[k] = append(res[k], State{sum, score})
					}
					if k+1 <= m {
						sum2 := sL.sum + sR.sum + val
						score2 := sL.score + sR.score + int64(j+1)*sL.sum + int64(i+1)*sR.sum - 2*val*int64(i)*int64(j) - val*int64(i+j)
						res[k+1] = append(res[k+1], State{sum2, score2})
					}
				}
			}
		}
	}
	for k := range res {
		res[k] = prune(res[k])
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Fscan(reader, &n, &m)
	a = make([]int64, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	leftChild = make([]int, n+1)
	rightChild = make([]int, n+1)
	root := buildCartesian()
	states := dfs(root)
	best := int64(-1 << 63)
	for _, st := range states[m] {
		if st.score > best {
			best = st.score
		}
	}
	writer := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(writer, best)
	writer.Flush()
}
