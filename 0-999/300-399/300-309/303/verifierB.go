package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded testcases from testcasesB.txt.
const embeddedTestcasesB = `100
5 79 5 33 4 9
12 87 12 16 3 5
11 90 8 87 7 68
36 67 15 27 27 36
58 64 42 45 6 42
79 15 62 10 25 4
3 94 2 14 3 29
48 22 21 13 4 4
19 90 7 5 19 82
69 78 9 3 16 25
78 74 15 50 12 48
15 5 0 1 3 1
62 27 46 25 4 22
3 70 3 12 2 9
29 10 20 4 12 7
24 8 16 7 2 2
90 51 25 16 46 47
61 73 10 26 50 8
87 21 20 10 68 9
16 77 14 22 1 61
88 53 72 32 40 42
46 50 42 16 10 36
89 2 58 2 11 2
95 6 69 2 18 2
98 62 45 39 37 44
46 76 40 16 46 40
50 96 26 83 6 1
77 25 42 5 31 8
82 58 48 45 73 56
54 5 25 5 37 4
99 85 90 5 22 58
9 34 2 28 9 32
72 78 0 4 64 42
40 60 3 51 27 13
71 82 10 16 2 52
87 54 40 0 28 1
92 97 0 86 68 79
13 25 1 19 11 7
39 36 11 6 31 26
81 11 2 4 58 2
33 18 33 11 8 5
36 3 2 0 14 3
34 72 20 46 3 64
92 83 58 81 56 48
69 23 26 12 38 1
18 20 8 10 11 12
92 12 43 12 80 1
6 35 1 9 5 19
47 51 35 8 19 8
62 94 15 6 20 23
67 94 9 38 52 43
39 54 6 6 36 31
61 44 53 21 8 31
15 90 15 54 1 39
43 95 43 19 11 81
73 49 11 4 11 13
96 29 7 12 2 4
51 72 33 37 29 63
75 92 27 54 11 48
29 34 18 10 14 13
46 15 4 0 34 8
97 87 25 15 64 51
33 27 2 25 14 20
19 14 6 7 13 6
70 20 13 19 63 5
73 52 54 33 64 44
42 64 31 25 35 29
2 44 2 20 2 3
68 19 32 19 20 13
75 38 60 4 11 34
6 9 1 2 1 5
2 98 1 42 1 20
84 59 47 32 49 58
68 65 4 11 67 10
96 55 96 13 38 35
77 54 61 54 50 39
76 30 2 21 1 24
24 39 16 36 9 22
9 64 4 38 7 50
50 8 10 2 16 5
94 43 7 2 62 27
19 63 19 10 5 52
46 53 2 39 30 25
59 7 6 7 50 2
3 5 1 5 2 1
90 71 83 44 25 50
100 100 62 14 8 79
90 60 78 40 44 42
16 88 9 16 13 38
96 88 15 66 25 5
51 57 23 48 13 30
46 81 4 5 3 63
33 4 33 4 14 2
12 100 10 99 9 90
68 54 64 19 15 10
55 73 27 10 7 54
9 13 6 12 3 12
4 58 3 43 4 2
64 42 32 5 46 5
16 46 0 22 12 12`

// gcd returns the greatest common divisor of a and b.
func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func floorDiv(a, b int64) int64 {
	if b < 0 {
		a = -a
		b = -b
	}
	if a >= 0 {
		return a / b
	}
	return -((-a + b - 1) / b)
}

func solve303B(n, m, x, y, a, b int64) (int64, int64, int64, int64) {
	g := gcd(a, b)
	a /= g
	b /= g
	k := n / a
	if km := m / b; km < k {
		k = km
	}
	wk := a * k
	hk := b * k
	lx := x - wk
	if lx < 0 {
		lx = 0
	}
	rx := x
	if n-wk < rx {
		rx = n - wk
	}
	ly := y - hk
	if ly < 0 {
		ly = 0
	}
	ry := y
	if m-hk < ry {
		ry = m - hk
	}
	D := 2*x - wk
	t0 := floorDiv(D, 2)
	candX := []int64{t0, t0 + 1}
	x1 := lx
	bestDx := int64(-1)
	for _, c := range candX {
		xx := c
		if xx < lx {
			xx = lx
		}
		if xx > rx {
			xx = rx
		}
		diff := 2*xx + wk - 2*x
		if diff < 0 {
			diff = -diff
		}
		if bestDx == -1 || diff < bestDx || (diff == bestDx && xx < x1) {
			bestDx = diff
			x1 = xx
		}
	}
	D = 2*y - hk
	t0 = floorDiv(D, 2)
	candY := []int64{t0, t0 + 1}
	y1 := ly
	bestDy := int64(-1)
	for _, c := range candY {
		yy := c
		if yy < ly {
			yy = ly
		}
		if yy > ry {
			yy = ry
		}
		diff := 2*yy + hk - 2*y
		if diff < 0 {
			diff = -diff
		}
		if bestDy == -1 || diff < bestDy || (diff == bestDy && yy < y1) {
			bestDy = diff
			y1 = yy
		}
	}
	x2 := x1 + wk
	y2 := y1 + hk
	return x1, y1, x2, y2
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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	lines := strings.Split(strings.TrimSpace(embeddedTestcasesB), "\n")
	if len(lines) < 1 {
		fmt.Fprintln(os.Stderr, "no testcases found")
		os.Exit(1)
	}
	// skip first line count
	tLine := strings.TrimSpace(lines[0])
	var t int
	fmt.Sscan(tLine, &t)
	if t != len(lines)-1 {
		// proceed but note mismatch
	}
	for idx, line := range lines[1:] {
		parts := strings.Fields(line)
		if len(parts) != 6 {
			fmt.Fprintf(os.Stderr, "case %d: expected 6 values got %d\n", idx+1, len(parts))
			os.Exit(1)
		}
		vals := make([]int64, 6)
		for i := 0; i < 6; i++ {
			v, err := strconv.ParseInt(parts[i], 10, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "case %d: parse error: %v\n", idx+1, err)
				os.Exit(1)
			}
			vals[i] = v
		}
		n, m, x, y, a, b := vals[0], vals[1], vals[2], vals[3], vals[4], vals[5]
		input := fmt.Sprintf("%d %d %d %d %d %d\n", n, m, x, y, a, b)
		x1, y1, x2, y2 := solve303B(n, m, x, y, a, b)
		want := fmt.Sprintf("%d %d %d %d", x1, y1, x2, y2)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(lines)-1)
}
