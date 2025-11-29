package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var testcases = []struct{ n, d int }{
	{n: 5, d: 6},
	{n: 3, d: 3},
	{n: 9, d: 12},
	{n: 3, d: 2},
	{n: 2, d: 1},
	{n: 10, d: 27},
	{n: 2, d: 1},
	{n: 10, d: 43},
	{n: 7, d: 14},
	{n: 4, d: 3},
	{n: 6, d: 8},
	{n: 2, d: 1},
	{n: 6, d: 8},
	{n: 4, d: 5},
	{n: 6, d: 15},
	{n: 7, d: 8},
	{n: 7, d: 18},
	{n: 10, d: 24},
	{n: 4, d: 4},
	{n: 9, d: 16},
	{n: 3, d: 3},
	{n: 2, d: 1},
	{n: 6, d: 13},
	{n: 5, d: 7},
	{n: 8, d: 26},
	{n: 6, d: 11},
	{n: 9, d: 13},
	{n: 5, d: 6},
	{n: 6, d: 5},
	{n: 3, d: 2},
	{n: 9, d: 28},
	{n: 6, d: 13},
	{n: 10, d: 39},
	{n: 7, d: 10},
	{n: 5, d: 4},
	{n: 8, d: 13},
	{n: 9, d: 16},
	{n: 4, d: 5},
	{n: 8, d: 25},
	{n: 7, d: 12},
	{n: 7, d: 9},
	{n: 2, d: 1},
	{n: 6, d: 14},
	{n: 5, d: 4},
	{n: 7, d: 11},
	{n: 6, d: 12},
	{n: 2, d: 1},
	{n: 7, d: 8},
	{n: 6, d: 15},
	{n: 7, d: 6},
	{n: 7, d: 15},
	{n: 7, d: 10},
	{n: 8, d: 26},
	{n: 3, d: 3},
	{n: 5, d: 7},
	{n: 6, d: 7},
	{n: 6, d: 11},
	{n: 4, d: 5},
	{n: 2, d: 1},
	{n: 2, d: 1},
	{n: 4, d: 5},
	{n: 7, d: 15},
	{n: 3, d: 3},
	{n: 5, d: 7},
	{n: 5, d: 4},
	{n: 2, d: 1},
	{n: 2, d: 1},
	{n: 4, d: 3},
	{n: 10, d: 40},
	{n: 5, d: 6},
	{n: 2, d: 1},
	{n: 10, d: 27},
	{n: 8, d: 27},
	{n: 5, d: 7},
	{n: 5, d: 5},
	{n: 9, d: 21},
	{n: 9, d: 9},
	{n: 5, d: 7},
	{n: 9, d: 15},
	{n: 8, d: 13},
	{n: 9, d: 14},
	{n: 2, d: 1},
	{n: 6, d: 9},
	{n: 5, d: 8},
	{n: 5, d: 10},
	{n: 5, d: 7},
	{n: 6, d: 7},
	{n: 7, d: 7},
	{n: 7, d: 9},
	{n: 8, d: 27},
	{n: 2, d: 1},
	{n: 8, d: 9},
	{n: 8, d: 13},
	{n: 4, d: 5},
	{n: 6, d: 15},
	{n: 9, d: 33},
	{n: 7, d: 19},
	{n: 10, d: 22},
	{n: 6, d: 10},
	{n: 8, d: 22},
}

const testcasesCount = 100

func minAdditional(nodes int, cap int, depth int) int64 {
	var sum int64
	rem := nodes
	c := cap
	d := depth
	for rem > 0 {
		use := c
		if use > rem {
			use = rem
		}
		sum += int64(use * d)
		rem -= use
		c = use * 2
		d++
	}
	return sum
}

func maxAdditional(nodes int, depth int) int64 {
	n := int64(nodes)
	return n*int64(depth) + n*(n-1)/2
}

func constructTree(n, d int) ([]int, bool) {
	maxSum := int64(n*(n-1)) / 2
	minSum := minAdditional(n-1, 2, 1)
	if int64(d) < minSum || int64(d) > maxSum {
		return nil, false
	}

	levels := []int{1}
	curSum := int64(0)
	rem := n - 1
	avail := 1
	depth := 1
	for rem > 0 {
		maxNodes := avail * 2
		if maxNodes > rem {
			maxNodes = rem
		}
		chosen := 0
		for x := 1; x <= maxNodes; x++ {
			remaining := rem - x
			minAdd := minAdditional(remaining, x*2, depth+1)
			maxAdd := maxAdditional(remaining, depth+1)
			totalMin := curSum + int64(x*depth) + minAdd
			totalMax := curSum + int64(x*depth) + maxAdd
			if int64(d) >= totalMin && int64(d) <= totalMax {
				chosen = x
				break
			}
		}
		if chosen == 0 {
			return nil, false
		}
		levels = append(levels, chosen)
		curSum += int64(chosen * depth)
		rem -= chosen
		avail = chosen
		depth++
	}

	parent := make([]int, n+1)
	parent[1] = 0
	prev := []int{1}
	idx := 2
	for lvl := 1; lvl < len(levels); lvl++ {
		cnt := levels[lvl]
		next := make([]int, 0, cnt)
		pIdx := 0
		used := 0
		for i := 0; i < cnt; i++ {
			if used == 2 {
				pIdx++
				used = 0
			}
			if pIdx >= len(prev) || idx > n {
				return nil, false
			}
			parent[idx] = prev[pIdx]
			next = append(next, idx)
			idx++
			used++
		}
		prev = next
	}
	return parent, true
}

func expected(n, d int) string {
	parent, ok := constructTree(n, d)
	if !ok {
		return "NO"
	}
	var sb strings.Builder
	sb.WriteString("YES\n")
	for i := 2; i <= n; i++ {
		sb.WriteString(strconv.Itoa(parent[i]))
		if i < n {
			sb.WriteByte(' ')
		}
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	if len(testcases) != testcasesCount {
		fmt.Fprintf(os.Stderr, "unexpected testcase count: got %d want %d\n", len(testcases), testcasesCount)
		os.Exit(1)
	}

	for i, tc := range testcases {
		exp := expected(tc.n, tc.d)

		var input strings.Builder
		input.WriteString("1\n")
		input.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.d))

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input.String())
		var out bytes.Buffer
		var errBuf bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &errBuf
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", i+1, err, errBuf.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != exp {
			fmt.Printf("test %d failed\nexpected:\n%s\n\ngot:\n%s\n", i+1, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
