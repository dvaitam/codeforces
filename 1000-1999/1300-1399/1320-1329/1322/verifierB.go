package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type testCase struct {
	n int
	a []int
}

const testcasesRaw = `6 12 16 11 20 3 4
6 0 5 4 0 10 19
4 1 13 20 1
5 12 1 19 5 11
2 13 1
8 11 19 19 8 9 18 14 13
3 0 14 8
4 12 2 11 3
2 0 11
1 5
7 19 20 0 10 14 17 15
8 2 1 17 12 8 0 20 16
2 2 10
6 3 15 1 4 16 20
5 1 0 12 10 5
3 5 5 5
4 19 10 0 15
7 1 7 7 20 9 10 5
4 11 7 5 13
8 11 18 4 12 18 0 5 18
1 12
3 4 0 0
6 16 0 1 1 3 18
3 4 12 0
7 13 18 10 7 4 11 16
4 17 12 2 4
7 18 11 3 13 13 7 15
7 7 12 7 20 15 12 18
2 8 8
6 17 0 19 19 15 7
5 1 19 10 12 20
2 17 1
3 12 0 13
7 13 3 14 19 14 5 5
6 15 13 5 18 9 16
2 11 11
3 20 11 15
1 6
5 5 18 10 9 12
1 9
7 1 13 8 12 6 11 4
3 3 19 11
3 0 13 18
7 14 2 20 2 13 17 17
3 5 4 6
3 7 0 16
3 15 11 19
5 10 3 13 8 5
6 20 16 10 17 4 12
5 7 12 11 12 15
5 13 13 3 4 4
1 18
2 20 6
2 8 19
3 12 2 1
1 3
6 15 10 3 14 11 18
5 15 7 5 17 17
2 16 5
1 20
3 16 13 19
4 14 12 8 0
3 12 5 14
1 12
2 20 18
7 10 7 16 14 1 15 19
2 8 16
8 17 20 12 15 8 5 7 17
6 5 9 19 4 14 2
2 15 12
7 17 2 8 15 7 3 9
3 11 3 4
1 4
4 0 1 12 17
8 3 15 17 11 10 3 0 7
4 15 9 8 7
1 15
6 16 10 2 2 9 18
7 7 11 12 4 7 9 6
8 11 9 12 19 4 3 12 11
8 7 20 11 20 11 13 8 11
7 9 3 15 9 3 14 4
6 7 5 10 15 7 3
7 12 14 16 14 18 19 7
7 16 9 15 7 10 16 0
2 15 10
7 7 13 1 18 1 13 2
5 6 10 5 3 5
6 0 7 1 0 12 17
1 4
2 19 19
4 2 14 6 0
7 2 17 5 7 7 13 12
8 0 13 6 12 1 19 8 0
6 11 10 14 20 4 19
2 8 3
2 8 0
3 19 4 12
4 18 10 6 13
2 17 3`

func parseTestcases(raw string) []testCase {
	lines := strings.Split(strings.TrimSpace(raw), "\n")
	tests := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 2 {
			panic(fmt.Sprintf("line %d malformed", idx+1))
		}
		n, err := strconv.Atoi(parts[0])
		if err != nil {
			panic(fmt.Sprintf("line %d bad n: %v", idx+1, err))
		}
		if len(parts) != n+1 {
			panic(fmt.Sprintf("line %d expected %d numbers got %d", idx+1, n+1, len(parts)))
		}
		a := make([]int, n)
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(parts[1+i])
			if err != nil {
				panic(fmt.Sprintf("line %d bad value: %v", idx+1, err))
			}
			a[i] = v
		}
		tests = append(tests, testCase{n: n, a: a})
	}
	if len(tests) == 0 {
		panic("no testcases parsed")
	}
	return tests
}

// Embedded solver logic from 1322B.go.
func expected(tc testCase) int {
	ans := 0
	n := tc.n
	a := tc.a
	for bit := 0; bit < 25; bit++ {
		mod := 1 << (bit + 1)
		b := make([]int, n)
		for i := 0; i < n; i++ {
			b[i] = a[i] % mod
		}
		sort.Ints(b)
		cnt := 0
		for i := 0; i < n; i++ {
			x := b[i]
			l1 := (1 << bit) - x
			r1 := (1 << (bit + 1)) - x
			l2 := (1 << (bit + 1)) + (1 << bit) - x
			r2 := (1 << (bit + 2)) - x

			idx1 := sort.SearchInts(b, l1)
			if idx1 < i+1 {
				idx1 = i + 1
			}
			idx2 := sort.SearchInts(b, r1)
			if idx2 < i+1 {
				idx2 = i + 1
			}
			idx3 := sort.SearchInts(b, l2)
			if idx3 < i+1 {
				idx3 = i + 1
			}
			idx4 := sort.SearchInts(b, r2)
			if idx4 < i+1 {
				idx4 = i + 1
			}
			cnt += idx2 - idx1
			cnt += idx4 - idx3
		}
		if cnt%2 == 1 {
			ans |= 1 << bit
		}
	}
	return ans
}

func run(bin string, tc testCase) (string, error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", tc.n))
	for i := 0; i < tc.n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", tc.a[i]))
	}
	sb.WriteByte('\n')
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := parseTestcases(testcasesRaw)
	for idx, tc := range tests {
		want := expected(tc)
		gotStr, err := run(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := strconv.Atoi(gotStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d bad output %q: %v\n", idx+1, gotStr, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
