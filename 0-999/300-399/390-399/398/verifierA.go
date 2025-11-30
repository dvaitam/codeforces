package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded testcases (one "a b" per line).
const embeddedTestcases = `6 6
0 4
8 7
6 4
7 5
9 3
8 2
4 2
1 9
4 8
9 2
4 1
1 10
5 7
8 1
5 6
5 9
10 3
8 7
7 8
4 0
8 0
1 6
10 10
0 9
7 5
3 5
1 3
9 3
3 2
8 7
1 1
5 8
7 1
4 8
4 1
8 5
8 3
9 8
9 4
7 1
9 6
5 9
3 4
2 3
2 0
9 10
4 7
1 1
10 2
2 0
1 8
10 6
8 4
8 3
3 10
9 6
9 4
7 7
10 10
5 1
5 9
1 7
9 10
5 3
3 0
4 1
3 5
2 5
6 0
1 2
3 0
9 10
8 9
10 1
0 1
10 3
9 9
1 6
1 5
1 0
9 0
3 2
1 7
3 0
10 0
8 6
9 1
4 1
3 1
10 4
5 6
2 0
8 7
0 9
1 6
3 4
5 7
9 2
10 3`

func bestScore(a, b int) (int64, string) {
	if a == 0 {
		score := -int64(b) * int64(b)
		return score, strings.Repeat("x", b)
	}
	if b == 0 {
		score := int64(a) * int64(a)
		return score, strings.Repeat("o", a)
	}
	bestS := int64(-1 << 63)
	bestK := 1
	for k := 1; k <= a; k++ {
		big := int64(a - k + 1)
		sumO := big*big + int64(k-1)
		m := k + 1
		if b < m {
			m = b
		}
		var sumX int64
		if m > 0 {
			base := b / m
			r := b % m
			sumX += int64(r) * int64(base+1) * int64(base+1)
			sumX += int64(m-r) * int64(base) * int64(base)
		}
		s := sumO - sumX
		if s > bestS {
			bestS = s
			bestK = k
		}
	}
	k := bestK
	oLens := make([]int, k)
	oLens[0] = a - k + 1
	for i := 1; i < k; i++ {
		oLens[i] = 1
	}
	slots := make([]int, k+1)
	m := k + 1
	if b < m {
		m = b
	}
	if m > 0 {
		base := b / m
		r := b % m
		for i := 0; i < m; i++ {
			if i < r {
				slots[i] = base + 1
			} else {
				slots[i] = base
			}
		}
	}
	var sb strings.Builder
	for i := 0; i <= k; i++ {
		if slots[i] > 0 {
			sb.WriteString(strings.Repeat("x", slots[i]))
		}
		if i < k {
			sb.WriteString(strings.Repeat("o", oLens[i]))
		}
	}
	return bestS, sb.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	lines := strings.Split(strings.TrimSpace(embeddedTestcases), "\n")
	idx := 0
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		idx++
		parts := strings.Fields(line)
		if len(parts) != 2 {
			fmt.Printf("Test %d invalid input line\n", idx)
			os.Exit(1)
		}
		a, _ := strconv.Atoi(parts[0])
		b, _ := strconv.Atoi(parts[1])
		expScore, expStr := bestScore(a, b)
		input := fmt.Sprintf("%d %d\n", a, b)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Test %d runtime error: %v\nstderr: %s\n", idx, err, string(out))
			os.Exit(1)
		}
		outs := strings.Split(strings.TrimSpace(string(out)), "\n")
		if len(outs) < 2 {
			fmt.Printf("Test %d invalid output: %s\n", idx, string(out))
			os.Exit(1)
		}
		gotScore, err := strconv.ParseInt(strings.TrimSpace(outs[0]), 10, 64)
		if err != nil {
			fmt.Printf("Test %d: first line not integer\n", idx)
			os.Exit(1)
		}
		gotStr := strings.TrimSpace(outs[1])
		if gotScore != expScore || gotStr != expStr {
			fmt.Printf("Test %d failed\nexpected:\n%d\n%s\n\ngot:\n%s", idx, expScore, expStr, string(out))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
