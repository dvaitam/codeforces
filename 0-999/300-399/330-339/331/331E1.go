package main

import (
	"bufio"
	"fmt"
	"os"
)

type Edge struct {
	to     int
	vision []int
}

var (
	n, m      int
	graph     [][]Edge
	limit     int
	actual    []int
	visionSeq []int
	matchLen  int
	answer    []int
	found     bool
)

func updateMatch() bool {
	for matchLen < len(actual) && matchLen < len(visionSeq) {
		if actual[matchLen] != visionSeq[matchLen] {
			return false
		}
		matchLen++
	}
	return true
}

func dfs(u int) {
	if found {
		return
	}
	if len(actual) == len(visionSeq) && len(actual) > 1 && matchLen == len(actual) {
		answer = append([]int(nil), actual...)
		found = true
		return
	}
	if len(actual) >= limit {
		return
	}
	for _, e := range graph[u] {
		if len(visionSeq)+len(e.vision) > limit {
			continue
		}
		prevActualLen := len(actual)
		prevVisionLen := len(visionSeq)
		prevMatch := matchLen

		visionSeq = append(visionSeq, e.vision...)
		if !updateMatch() {
			visionSeq = visionSeq[:prevVisionLen]
			matchLen = prevMatch
			continue
		}

		actual = append(actual, e.to)
		if len(actual) > limit {
			actual = actual[:prevActualLen]
			visionSeq = visionSeq[:prevVisionLen]
			matchLen = prevMatch
			continue
		}
		if !updateMatch() {
			actual = actual[:prevActualLen]
			visionSeq = visionSeq[:prevVisionLen]
			matchLen = prevMatch
			continue
		}

		dfs(e.to)
		if found {
			return
		}

		actual = actual[:prevActualLen]
		visionSeq = visionSeq[:prevVisionLen]
		matchLen = prevMatch
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}

	graph = make([][]Edge, n+1)
	for i := 0; i < m; i++ {
		var x, y, k int
		fmt.Fscan(in, &x, &y, &k)
		seq := make([]int, k)
		for j := 0; j < k; j++ {
			fmt.Fscan(in, &seq[j])
		}
		graph[x] = append(graph[x], Edge{to: y, vision: seq})
	}

	limit = 2 * n
	found = false
	for start := 1; start <= n && !found; start++ {
		actual = actual[:0]
		visionSeq = visionSeq[:0]
		matchLen = 0
		actual = append(actual, start)
		dfs(start)
	}

	if !found {
		fmt.Fprintln(out, 0)
		return
	}
	fmt.Fprintln(out, len(answer))
	for i, v := range answer {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, v)
	}
	fmt.Fprintln(out)
}
