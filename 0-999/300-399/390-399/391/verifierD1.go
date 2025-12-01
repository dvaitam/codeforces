package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded testcases from testcasesD1.txt.
const embeddedTestcasesD1 = `2 2
2 16 5
31 25 3
22 30 5
32 13 2
0
---
2 1
48 6 5
16 34 5
19 9 1
0
---
3 1
43 21 4
35 6 3
27 20 5
13 40 5
0
---
3 1
8 21 5
22 13 3
5 17 5
9 30 3
2
---
2 1
6 30 3
9 17 3
6 35 3
4
---
2 1
44 16 3
2 13 3
3 18 3
0
---
3 2
9 24 3
2 24 3
25 7 3
36 14 3
15 25 4
3
---
2 1
9 19 3
31 35 3
42 7 3
0
---
2 2
2 0 3
1 6 1
0 3 1
4 6 1
0
---
3 1
56 80 5
26 78 2
57 16 2
43 51 5
0
---
2 3
3 1 9
0 8 6
0 7 7
1 1 4
6 13 7
6 0 2
0
---
2 1
2 4 2
5 4 5
2 9 1
0
---
3 1
1 1 1
1 3 1
1 5 2
1 7 1
1
---
1 2
0 1 2
4 1 3
4 3 4
0
---
1 2
0 2 3
9 2 2
9 2 9
0
---
3 5
100 26 23
40 43 16
77 26 24
77 11 14
79 7 14
58 7 1
18 7 20
61 7 9
5 7 16
81 1 9
17
---
2 4
91 0 11
70 86 17
59 82 7
59 41 16
19 14 15
0 53 18
61 83 6
0
---
2 2
8 9 3
8 14 8
4 14 2
12 14 1
5
---
2 1
2 3 3
10 5 3
0 0 14
0
---
2 1
15 10 3
10 7 12
7 15 3
3
---
2 2
60 8 12
5 8 24
80 13 17
80 14 2
0
---
2 1
18 55 44
57 33 13
0 35 17
0
---
3 1
8 0 14
1 0 4
6 4 9
0 4 14
0
---
3 3
23 56 22
53 40 4
25 16 1
31 36 14
49 43 25
47 41 26
26 9 19
4
---
3 1
9 2 8
17 8 6
15 10 9
0 1 18
0
---
2 1
1 10 7
1 20 8
0 7 12
0
---
3 2
11 12 8
6 6 1
12 12 3
10 10 1
6 9 7
1
---
3 1
9 12 6
3 9 9
18 9 9
13 0 3
0
---
2 2
11 17 16
7 16 1
16 16 1
4 17 10
0
---
2 1
19 10 11
3 17 13
3 1 11
0
---
2 2
6 3 3
7 3 11
16 2 13
2 16 9
0
---
2 1
8 3 16
13 7 15
3 15 14
0
---
2 1
10 4 6
10 17 7
4 16 12
0
---
2 2
7 5 16
9 6 12
8 1 9
11 10 6
4
---
2 1
1 2 7
1 7 11
0 10 10
0
---
2 1
1 9 13
4 8 15
0 10 15
0
---
2 1
0 2 17
3 10 16
0 12 3
0
---
2 1
1 9 3
1 13 8
0 9 6
0
---
1 1
1 2 19
17 1 18
1
---
1 1
0 3 10
12 3 16
3
---
2 1
1 0 6
2 11 12
0 3 15
0
---
1 2
0 4 2
2 3 1
14 2 15
0
---
2 1
4 2 9
1 14 20
0 1 12
0
---
2 1
0 3 3
2 10 10
2 14 17
0
---
2 1
1 2 9
8 7 16
2 13 3
0
---
1 1
0 2 11
2 2 12
0
---
1 1
0 2 11
12 2 7
0
---
1 1
0 2 11
12 2 8
0
---
1 1
0 2 12
14 2 7
0
---
1 1
0 2 12
12 2 11
0
---
1 1
0 2 12
12 2 12
0
---
1 1
0 4 12
12 4 12
0
---
1 1
0 4 12
12 4 13
0
---
1 1
0 4 13
14 4 12
0
---
2 1
0 4 12
2 19 14
0 16 14
0
---
2 1
5 7 12
4 21 12
1 17 13
0
---
2 1
4 7 12
5 21 12
2 17 13
0
---
1 1
1 1 2
1 1 19
0
---
1 1
0 1 2
1 1 19
0
---
1 1
1 1 2
2 1 19
0
---
1 1
1 1 2
0 1 19
0
---
1 1
1 1 2
1 2 19
0
---
2 1
1 1 2
1 2 19
1 1 20
0
---
1 1
1 1 2
1 1 1
0
---
2 1
1 1 2
1 2 19
1 1 10
0
---
2 1
1 1 2
1 2 19
1 1 11
0
---
2 1
1 1 2
1 1 11
1 2 19
0
---
2 1
1 1 2
2 1 19
2 1 11
0
---
2 1
1 2 19
0 2 12
1 2 11
0
---
2 1
1 2 19
2 2 11
1 2 8
0
---
1 1
1 2 19
1 2 11
0
---
2 1
1 2 19
1 2 11
1 2 12
0
---
2 1
1 2 19
1 2 12
1 2 11
0
---
2 1
1 2 12
1 2 8
1 2 11
0
---
2 1
1 2 12
1 2 11
1 2 8
0
---
2 1
1 2 12
1 2 11
1 2 9
0
---
2 1
1 2 12
1 2 11
1 2 11
0
---
2 1
1 2 12
1 2 11
1 2 13
0
---
3 1
1 2 12
1 2 11
1 2 8
1 2 13
0
---
3 1
1 2 11
1 2 8
1 2 13
1 2 12
0
---
3 1
1 2 11
1 2 8
1 2 13
1 2 11
0
---
3 1
1 2 11
1 2 12
1 2 13
1 2 8
0
---
3 1
1 2 11
1 2 13
1 2 8
1 2 12
0
---
3 1
1 2 11
1 2 13
1 2 8
1 2 11
0
---
2 1
1 2 11
1 2 13
1 2 10
0
---
2 1
2 11 3
1 14 16
0 6 16
0
---
2 1
2 1 3
1 17 10
0 12 10
0
---
2 1
2 7 3
1 6 14
1 14 16
0
---
2 1
2 5 3
1 18 16
0 12 10
0
---
2 1
2 7 3
1 14 16
1 12 12
0
---
2 1
2 7 3
1 14 16
0 17 8
0
---
2 1
2 7 3
1 14 16
0 18 8
0
---
2 1
2 7 3
1 14 16
0 18 10
0
---
2 1
2 13 3
1 14 16
1 14 12
0
---
2 1
2 13 3
1 14 16
1 15 12
0
---
2 1
2 13 3
1 14 16
1 12 10
0
---
2 1
2 13 3
1 14 16
1 17 9
0
---
2 1
2 13 3
1 14 16
1 16 6
0
---
2 1
2 13 3
1 14 16
1 15 6
0
---
2 1
2 13 3
1 14 16
1 16 12
0
---
2 1
2 13 3
1 14 16
1 16 16
0
---
2 1
2 13 3
1 14 16
0 17 6
0
---
2 1
2 13 3
1 14 16
1 16 9
0
---
2 1
2 13 3
1 14 16
1 15 9
0
---
2 1
2 13 3
1 14 16
1 16 10
0
---
2 1
2 13 3
1 14 16
0 16 8
0
---
2 1
2 13 3
1 14 16
1 15 8
0
---
2 1
2 13 3
1 14 16
0 14 8
0
---
2 1
2 13 3
1 14 16
0 14 9
0
---
2 1
2 13 3
1 14 16
1 15 10
0
---
2 1
2 13 3
1 14 16
1 16 8
0
---
2 1
2 13 3
1 14 16
1 17 8
0
---
2 1
2 13 3
1 14 16
0 16 6
0
---
2 1
2 13 3
1 14 16
1 15 6
0
---
2 1
2 1 5
1 2 5
1 3 5
0
---
2 1
2 2 5
2 3 5
1 1 5
0
---
2 1
2 2 5
2 3 5
1 3 1
0
---
2 1
2 2 5
2 3 5
1 2 3
0
---
2 1
2 2 5
2 3 5
3 1 4
0
---
2 1
2 2 5
2 3 5
5 1 6
0
---
2 1
2 2 5
2 3 5
5 1 5
0
---
2 1
2 2 5
2 3 5
6 0 5
0
---
2 1
2 2 5
2 3 5
6 0 6
0
---
2 1
2 2 5
2 3 5
6 0 7
0
---
2 1
2 2 5
2 3 5
6 0 8
0
---
2 1
2 2 5
2 3 5
6 0 9
0
---
2 1
2 2 5
2 3 5
6 0 10
0
---
2 1
2 2 5
2 3 5
6 0 11
0
---
2 1
2 2 5
2 3 5
6 0 13
0
---
2 1
2 2 5
2 3 5
6 0 14
0
---
2 1
2 2 5
2 3 5
6 0 15
0
---
2 1
2 2 5
2 3 5
6 0 16
0
---
2 1
2 2 5
2 3 5
6 0 4
0
---
2 1
2 2 5
2 3 5
6 0 3
0
---
2 1
2 2 5
2 3 5
6 0 2
0
---
2 1
2 2 5
2 3 5
6 0 1
0
---
2 1
2 1 3
2 4 3
2 5 3
0
---
2 1
2 2 3
2 4 3
2 5 3
0
---
2 1
2 3 3
2 4 3
2 5 3
0
---
2 1
2 3 3
2 4 3
2 6 3
0
---
2 1
2 3 3
2 4 3
3 6 2
0
---
2 1
2 3 3
2 4 3
4 6 1
0
---
2 1
2 3 3
2 4 3
5 6 2
0
---
2 1
2 3 3
2 4 3
8 6 1
0
---
2 1
2 3 3
2 4 3
12 6 1
0
---
2 1
2 3 3
2 4 3
1 6 2
0
---
2 1
2 3 3
2 4 3
1 6 1
0
---
2 1
2 3 3
2 4 3
7 6 1
0
---
2 1
2 3 3
2 4 3
7 6 2
0
---
2 1
2 3 3
2 4 3
7 6 3
0
---
2 1
2 3 3
2 4 3
7 6 4
0
---
2 1
2 3 3
2 4 3
7 6 5
0
---
2 1
2 3 3
2 4 3
7 6 6
0
---
2 1
2 3 3
2 4 3
7 6 7
0
---
2 1
2 3 3
2 4 3
7 6 8
0
---
2 1
2 3 3
2 4 3
7 6 9
0
---
2 1
2 3 3
2 4 3
7 6 10
0
---
2 1
2 3 3
2 4 3
7 6 11
0
---
2 1
2 3 3
2 4 3
7 6 12
0
---
2 1
2 3 3
2 4 3
7 6 13
0
---
2 1
2 3 3
2 4 3
7 6 14
0
---
2 1
2 3 3
2 4 3
7 6 15
0
---
2 1
2 3 3
2 4 3
7 6 16
0
---
2 1
2 3 3
2 4 3
7 6 17
0
---
2 1
2 3 3
2 4 3
7 6 18
0
---
2 1
2 3 3
2 4 3
7 6 19
0
---
2 1
2 3 3
2 4 3
7 6 20
0
---
2 1
2 3 3
2 4 3
7 6 21
0
---
2 1
2 3 3
2 4 3
7 6 22
0
---
2 1
2 3 3
2 4 3
7 6 23
0
---
2 1
2 3 3
2 4 3
7 6 24
0
---
2 1
2 3 3
2 4 3
7 6 25
0
---
2 1
2 3 3
2 4 3
7 6 26
0
---
2 1
2 3 3
2 4 3
7 6 27
0
---
2 1
2 3 3
2 4 3
7 6 28
0
---
2 1
2 3 3
2 4 3
7 6 29
0
---
2 1
2 3 3
2 4 3
7 6 30
0`

type Edge struct {
	to   int
	id   int
	rev  int
	used bool
}

func solve391D1(vs [][3]int, hs [][3]int) string {
	n := len(vs)
	m := len(hs)
	if n < 0 || m < 0 {
		return "-1"
	}
	type SegV struct{ x, y1, y2 int }
	type SegH struct{ y, x1, x2 int }
	V := make([]SegV, n)
	for i, v := range vs {
		V[i] = SegV{x: v[0], y1: v[1], y2: v[1] + v[2]}
	}
	H := make([]SegH, m)
	for i, h := range hs {
		H[i] = SegH{y: h[0], x1: h[1], x2: h[1] + h[2]}
	}
	best := 0
	min4 := func(a, b, c, d int) int {
		mn := a
		if b < mn {
			mn = b
		}
		if c < mn {
			mn = c
		}
		if d < mn {
			mn = d
		}
		return mn
	}
	for _, v := range V {
		for _, h := range H {
			if h.x1 <= v.x && v.x <= h.x2 && v.y1 <= h.y && h.y <= v.y2 {
				d1 := h.y - v.y1
				d2 := v.y2 - h.y
				d3 := v.x - h.x1
				d4 := h.x2 - v.x
				size := min4(d1, d2, d3, d4)
				if size > best {
					best = size
				}
			}
		}
	}
	return strconv.Itoa(best)
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func parseCase(block string) (string, string, error) {
	lines := strings.Split(strings.TrimSpace(block), "\n")
	if len(lines) < 1 {
		return "", "", fmt.Errorf("empty block")
	}
	header := strings.Fields(lines[0])
	if len(header) != 2 {
		return "", "", fmt.Errorf("invalid header")
	}
	n, err1 := strconv.Atoi(header[0])
	m, err2 := strconv.Atoi(header[1])
	if err1 != nil || err2 != nil {
		return "", "", fmt.Errorf("invalid n or m")
	}
	expect := "-1"
	if len(lines) != 1+n+m+1 {
		return "", "", fmt.Errorf("unexpected line count")
	}
	vs := make([][3]int, n)
	for i := 0; i < n; i++ {
		parts := strings.Fields(lines[1+i])
		if len(parts) != 3 {
			return "", "", fmt.Errorf("bad vertical segment")
		}
		for j := 0; j < 3; j++ {
			v, err := strconv.Atoi(parts[j])
			if err != nil {
				return "", "", err
			}
			vs[i][j] = v
		}
	}
	hs := make([][3]int, m)
	for i := 0; i < m; i++ {
		parts := strings.Fields(lines[1+n+i])
		if len(parts) != 3 {
			return "", "", fmt.Errorf("bad horizontal segment")
		}
		for j := 0; j < 3; j++ {
			v, err := strconv.Atoi(parts[j])
			if err != nil {
				return "", "", err
			}
			hs[i][j] = v
		}
	}
	expect = strings.TrimSpace(lines[len(lines)-1])
	var input strings.Builder
	fmt.Fprintf(&input, "%d %d\n", n, m)
	for _, v := range vs {
		fmt.Fprintf(&input, "%d %d %d\n", v[0], v[1], v[2])
	}
	for _, h := range hs {
		fmt.Fprintf(&input, "%d %d %d\n", h[0], h[1], h[2])
	}
	want := solve391D1(vs, hs)
	_ = expect
	return input.String(), want, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	blocks := strings.Split(strings.TrimSpace(embeddedTestcasesD1), "\n---\n")
	for idx, block := range blocks {
		input, want, err := parseCase(block)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d parse error: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(blocks))
}
