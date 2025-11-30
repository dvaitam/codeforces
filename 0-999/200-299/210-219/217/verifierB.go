package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Embedded copy of testcasesB.txt.
const testcasesBData = `
3 19
2 9
2 16
8 16
7 26
4 4
8 1
7 14
10 25
1 23
8 9
4 19
2 29
6 1
1 1
9 1
7 22
4 14
1 17
4 25
8 16
9 8
6 8
4 25
8 10
1 14
9 30
2 6
5 4
6 29
9 30
7 17
4 10
5 19
8 28
9 13
10 28
1 16
4 24
7 14
3 12
9 29
6 3
8 22
9 4
3 17
7 12
8 24
1 16
1 10
10 19
10 13
3 6
9 8
1 25
4 18
9 8
7 17
6 28
10 12
8 30
5 22
9 20
1 13
9 26
3 17
9 7
7 2
8 28
6 19
9 7
9 14
8 27
6 14
6 1
9 18
10 26
10 11
8 20
1 26
4 21
3 18
10 6
2 26
9 26
5 2
2 3
1 15
1 25
5 8
5 4
10 6
6 10
2 6
3 9
9 6
5 21
5 15
6 16
8 4
`

type testCase struct {
	n int
	r int
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(testcasesBData, "\n")
	var cases []testCase
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 2 {
			return nil, fmt.Errorf("line %d: expected 2 fields, got %d", idx+1, len(fields))
		}
		var n, r int
		if _, err := fmt.Sscan(fields[0], &n); err != nil {
			return nil, fmt.Errorf("line %d: parse n: %w", idx+1, err)
		}
		if _, err := fmt.Sscan(fields[1], &r); err != nil {
			return nil, fmt.Errorf("line %d: parse r: %w", idx+1, err)
		}
		cases = append(cases, testCase{n: n, r: r})
	}
	return cases, nil
}

func toggle(b byte) byte {
	if b == 'T' {
		return 'B'
	}
	return 'T'
}

// solve mirrors the logic from 217B.go, returning (errCount, sequence).
func solve(n, r int) (int, string) {
	if n == 1 {
		if r == 1 {
			return 0, "T"
		}
		return -1, ""
	}
	if n == 2 {
		if r == 2 {
			return 0, "TB"
		}
		return -1, ""
	}
	bi := -1
	const INF = 1 << 60
	bval := INF
	for i := 1; i+i <= r; i++ {
		x := i
		y := r - i
		ans := 0
		errc := 0
		for x > 0 && y > 0 {
			if x < y {
				x, y = y, x
			}
			dd := x / y
			ans += dd
			errc += dd - 1
			x %= y
		}
		x += y
		errc--
		if x == 1 && ans+1 == n {
			if errc < bval {
				bval = errc
				bi = i
			}
		}
	}
	if bi == -1 {
		return -1, ""
	}
	u := bi
	d := r - bi
	moves := make([]byte, 0, n+2)
	for u > 1 || d > 1 {
		if u > d {
			moves = append(moves, 'T')
			u -= d
		} else {
			moves = append(moves, 'B')
			d -= u
		}
	}
	for i, j := 0, len(moves)-1; i < j; i, j = i+1, j-1 {
		moves[i], moves[j] = moves[j], moves[i]
	}
	if len(moves) == 0 {
		return -1, ""
	}
	first := toggle(moves[0])
	last := toggle(moves[len(moves)-1])
	seq := make([]byte, 0, len(moves)+2)
	seq = append(seq, first)
	seq = append(seq, moves...)
	seq = append(seq, last)
	if seq[0] == 'B' {
		for i := range seq {
			seq[i] = toggle(seq[i])
		}
	}
	return bval, string(seq)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse embedded testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range tests {
		expVal, expSeq := solve(tc.n, tc.r)
		var input strings.Builder
		fmt.Fprintf(&input, "%d %d\n", tc.n, tc.r)

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input.String())
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx+1, err, stderr.String())
			os.Exit(1)
		}
		parts := strings.Fields(strings.TrimSpace(out.String()))
		if expVal == -1 {
			if len(parts) != 1 || parts[0] != "IMPOSSIBLE" {
				fmt.Printf("test %d failed: expected IMPOSSIBLE got %s\n", idx+1, strings.Join(parts, " "))
				os.Exit(1)
			}
			continue
		}
		if len(parts) != 2 {
			fmt.Printf("test %d: expected 2 outputs got %d\n", idx+1, len(parts))
			os.Exit(1)
		}
		if parts[0] != fmt.Sprint(expVal) || parts[1] != expSeq {
			fmt.Printf("test %d failed\nexpected: %d %s\n     got: %s %s\n", idx+1, expVal, expSeq, parts[0], parts[1])
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
