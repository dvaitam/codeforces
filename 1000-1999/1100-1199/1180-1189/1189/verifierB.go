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

const testcasesRaw = `
5 19 28 26 25 3
7 4 16 25 15 16 21 13
6 4 16 1 29 27 13
9 20 25 25 1 23 15 9 24 26
6 19 4 29 11 1 1
3 21 18 1
9 22 7 14 24 1 17 8 25 15
10 18 8 12 8 22 8 25 15 10 30
3 14 27 30
4 6 21 24 28
7 4 24 11 29 24 23 17
9 17 27 30 22 7 10 10 19 29
10 28 17 13 19 28 2 16 8 24 26
9 14 22 6 12 18 29 23 25 22
8 3 15 22 17 4 25 6 17
9 12 16 24 1 16 2 10 23 28
9 21 6 6 17 8 1 25 7 18
6 13 17 12 28 19 12
10 30 9 22 18 20 24 1 13 26 28
5 17 25 18 7 14
3 16 28 12
6 17 14 16 27 12 14
8 1 18 18 20 26 20 11 15
3 26 8 21
5 18 19 6 28 3
7 2 27 22 3 3 28 1
10 1 25 25 9 8 9 4 26 20 6
8 10 3 6 6 9 17 6 22
7 21 23 10 15 23 11 16
10 4 1 10 13 11 14 26 7 9 4
7 29 24 17 7 20 14 27
3 8 1 13
5 2 24 6 15 23
9 18 27 8 21 26 23 17 15 8
3 13 22 19
8 22 21 14 2 24 10 5 7
3 10 3 28
4 10 30 10 24
5 14 19 9 5 1
3 19 27 7
10 6 27 28 28 25 23 20 17 2 13
6 12 4 7 19 22 29
9 19 7 16 4 22 13 10 17 16
3 11 20 28
9 29 10 1 6 7 28 11 26 19
5 11 14 7 9 22
4 27 13 30 18
8 30 29 27 22 18 16 25 18
6 3 24 2 3 5 6
5 30 18 7 9 25
8 20 17 27 9 12 11 11 4
7 8 28 20 25 23 29 16
5 19 18 25 4 11
3 14 3 13
5 27 5 11 4 20
9 3 19 18 8 19 3 9 12 29
7 19 18 30 4 15 29 9
4 26 2 27 10
3 20 22 1
4 14 4 27 29
3 7 8 26
9 6 4 15 6 22 8 6 24 28
4 14 30 13 26
7 18 9 23 16 11 4 7
8 2 1 1 26 30 10 24 20
8 15 13 11 13 3 3 30 11
10 4 9 7 26 20 25 29 18 28 23
10 22 12 9 6 18 7 10 7 8 12
4 27 9 3 25
10 3 21 19 21 11 8 13 10 2 11
5 11 26 28 19 29
7 8 11 4 18 20 19 26
4 8 8 1 26
6 13 3 9 18 28 3
4 1 21 1 10
8 16 16 28 28 5 4 17 25
8 3 17 22 6 6 25 5 5
8 10 4 23 17 27 30 20 10
5 29 7 5 18 30
3 25 11 27
6 6 10 14 18 6 2
6 9 25 3 22 15 26
9 18 9 18 15 28 18 15 1 13
8 6 9 16 1 26 21 30 14
3 2 23 12
5 19 5 5 9 27
7 13 19 13 6 20 3 8
10 1 6 17 11 17 29 21 30 15 30
6 8 11 16 22 16 8
9 11 18 20 30 24 30 21 9 21
6 2 30 3 25 17 21
8 6 17 25 26 29 7 10 10
7 28 18 12 6 23 23 24
10 20 3 28 4 29 20 17 19 13 6
5 9 14 7 19 24
3 16 22 13
8 13 17 28 6 18 24 2 17
4 26 9 21 4
7 24 30 3 5 25 20 27
4 15 28 30 8
`

type testCase struct {
	n int
	a []int
}

func possible(a []int) bool {
	if len(a) < 3 {
		return false
	}
	b := make([]int, len(a))
	copy(b, a)
	sort.Ints(b)
	return b[len(b)-1] < b[len(b)-2]+b[len(b)-3]
}

// solve embeds the arrangement logic from 1189B.go.
func solve(tc testCase) (string, []int) {
	n := tc.n
	a := append([]int(nil), tc.a...)
	sort.Ints(a)
	if n < 3 || a[n-3]+a[n-2] <= a[n-1] {
		return "NO", nil
	}
	b := make([]int, n)
	half := n / 2
	for k := 0; k < half; k++ {
		b[k] = a[2*k]
		b[n-1-k] = a[2*k+1]
	}
	if n%2 == 1 {
		b[half] = a[n-1]
	}
	return "YES", b
}

func runCase(bin string, tc testCase) error {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(tc.n))
	sb.WriteByte('\n')
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	input := sb.String()

	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	outFields := strings.Fields(out.String())
	if len(outFields) == 0 {
		return fmt.Errorf("no output")
	}
	wantAns, _ := solve(tc)
	gotAns := strings.ToUpper(outFields[0])
	if gotAns != wantAns {
		return fmt.Errorf("expected %s got %s", wantAns, gotAns)
	}
	if gotAns == "NO" {
		return nil
	}
	if len(outFields)-1 != tc.n {
		return fmt.Errorf("expected %d numbers got %d", tc.n, len(outFields)-1)
	}
	arr := make([]int, tc.n)
	for i := 0; i < tc.n; i++ {
		v, err := strconv.Atoi(outFields[i+1])
		if err != nil {
			return fmt.Errorf("invalid number %q", outFields[i+1])
		}
		arr[i] = v
	}
	cnt := make(map[int]int)
	for _, v := range tc.a {
		cnt[v]++
	}
	for _, v := range arr {
		cnt[v]--
	}
	for _, c := range cnt {
		if c != 0 {
			return fmt.Errorf("output elements mismatch input")
		}
	}
	for i := 0; i < tc.n; i++ {
		left := arr[(i-1+tc.n)%tc.n]
		right := arr[(i+1)%tc.n]
		if arr[i] >= left+right {
			return fmt.Errorf("triangle condition failed at index %d", i)
		}
	}
	if !possible(tc.a) {
		return fmt.Errorf("expected NO but got YES with valid arrangement")
	}
	return nil
}

func parseTestcases(raw string) ([]testCase, error) {
	lines := strings.Split(raw, "\n")
	tests := make([]testCase, 0)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		n, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, err
		}
		if len(parts) != n+1 {
			return nil, fmt.Errorf("expected %d numbers got %d", n+1, len(parts))
		}
		a := make([]int, n)
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(parts[i+1])
			if err != nil {
				return nil, err
			}
			a[i] = v
		}
		tests = append(tests, testCase{n: n, a: a})
	}
	if len(tests) == 0 {
		return nil, fmt.Errorf("no testcases")
	}
	return tests, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests, err := parseTestcases(testcasesRaw)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for i, tc := range tests {
		if err := runCase(bin, tc); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
