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

const testcasesRaw = `100
2 4 3 3 2 2 16 9 1 17 19 19
2 2 1 17 17 14 17 10
1 2 4 19 14 3 4 14 3 4
4 2 1 15 14 14 1 16 11 9
1 3 1 4 12 1 12 12
2 1 2 12 3 20 5 7
1 2 1 1 10 12 1
5 2 2 6 15 4 16 12 9 5 1 7
3 3 4 10 10 18 11 6 19 3 4 18 19
3 2 4 5 5 8 11 17 8 8 6 10
3 4 1 5 20 1 13 3 3 5 14
3 5 4 5 19 14 10 12 3 8 15 12 17 2 13
4 1 4 11 15 7 12 10 16 3 6 4
3 1 5 20 5 15 13 6 14 14 6 8
4 3 5 5 12 15 3 16 7 10 1 15 20 15 1
2 3 1 10 18 20 5 14 16
1 4 2 18 13 9 1 4 9 2
1 3 4 17 19 13 15 4 9 12 10
2 5 1 2 3 9 10 18 11 4 17
2 2 1 14 10 10 17 5
5 5 2 18 4 14 18 13 9 10 15 12 19 5 6
1 1 4 13 19 15 5 18 10
3 4 4 7 16 16 17 11 16 2 15 10 5 16
1 5 2 1 12 16 13 1 17 3 3
4 1 3 2 4 20 1 9 10 8 5
5 3 2 4 14 15 11 13 6 11 14 14 5
4 2 5 11 5 7 6 15 12 13 14 16 13 8
2 4 2 19 2 13 2 8 3 6 12
1 2 2 20 10 20 3 17
3 3 4 15 2 17 18 14 19 15 16 9 16
2 3 3 2 2 2 6 12 1 10 1
2 1 4 8 20 13 18 8 15 7
3 5 1 20 3 11 11 18 15 11 9 1
5 1 2 12 3 7 17 12 7 7 9
3 3 5 13 9 16 12 8 2 10 18 3 1 15
4 4 1 14 16 15 15 4 3 3 8 4
2 4 2 15 20 3 14 18 13 2 6
2 4 2 5 9 12 11 14 4 18 10
5 5 2 10 15 17 20 15 18 9 9 8 1 4 20
1 2 4 8 7 10 1 18 17 14
1 1 4 9 4 19 12 8 18
3 2 2 3 17 10 11 8 12 16
3 5 2 5 1 18 17 11 12 19 1 5 13
2 2 5 3 5 7 16 19 7 8 5 8
4 3 5 19 5 16 4 20 1 17 20 12 16 15 10
1 2 5 6 16 16 18 11 3 9 5
5 4 2 11 10 13 2 7 2 11 8 11 15 8
3 3 2 10 1 12 19 18 2 5 12
1 4 1 1 8 2 1 8 11
1 1 3 14 5 7 15 14
2 3 3 6 11 14 13 1 14 9 18
5 4 1 19 4 14 13 6 1 17 5 20 17
2 1 3 8 6 8 1 6 18
2 1 4 20 4 20 15 5 20 20
1 3 3 13 1 2 16 3 12 10
2 4 2 17 12 6 13 11 9 16 13
1 3 5 10 18 16 2 18 19 18 9 2
4 4 1 13 12 16 2 1 9 2 9 19
3 2 5 17 11 13 9 7 4 19 11 8 19
5 3 2 5 11 1 19 2 19 5 12 12 10
3 3 4 13 20 14 6 1 5 19 2 15 5
3 1 4 9 20 7 3 18 14 9 6
5 2 1 6 19 4 17 18 20 13 14
3 3 3 1 14 9 9 18 17 18 11 11
2 4 2 1 17 5 19 13 12 15 2
5 4 5 8 1 12 17 6 7 12 16 1 8 19 8 9 6
4 1 5 15 8 15 17 4 7 6 15 3 14
4 3 3 14 12 20 11 3 10 1 16 1 9
2 4 4 14 13 2 19 15 12 19 5 19 9
3 1 4 16 17 5 2 3 19 12 12
1 1 2 4 18 16 2
3 1 3 13 5 9 14 5 20 5
4 3 5 2 6 5 5 16 2 17 2 18 13 6 12
5 1 1 18 6 9 7 9 11 9
3 5 4 5 15 18 5 2 19 6 17 2 11 3 7
4 5 2 15 17 6 11 5 16 18 2 18 3 17
3 1 1 4 14 20 12 19
4 3 4 17 12 4 5 11 1 6 5 1 11 20
2 1 4 2 10 13 2 20 6 12
1 4 1 15 12 20 20 9 10
5 4 4 6 1 15 9 7 13 3 12 4 4 1 12 1
2 4 5 1 11 15 18 16 16 3 2 18 13 9
1 5 1 3 11 12 4 16 2 5
5 3 1 1 13 11 6 18 5 6 6 6
2 5 3 1 16 13 2 8 8 10 11 6 8
3 2 2 14 15 12 19 5 13 19
1 2 5 1 13 6 5 1 1 11 17
1 1 1 4 19 20
2 2 4 1 14 14 19 11 8 5 12
5 2 5 13 3 5 14 19 12 4 14 14 8 16 13
2 4 2 16 13 19 3 9 9 17 12
5 1 5 20 16 8 9 2 20 11 13 4 18 2
2 4 1 14 13 14 4 15 20 15
2 2 3 16 14 6 19 10 17 4
3 3 2 12 16 17 2 7 9 6 19
3 3 4 2 10 18 14 2 14 9 13 7 12
2 2 1 20 12 6 1 14
5 4 4 3 3 14 18 18 5 6 5 7 6 8 1 17
2 4 3 20 10 11 4 14 9 6 11 17
3 5 2 13 18 10 8 13 12 13 16 17 10
`

// gcd and min/calc/solve logic from 1336B.go (adapted for single test).
func calc(x, y, z int64) int64 {
	a := x - y
	b := y - z
	c := z - x
	return a*a + b*b + c*c
}

func best(a, b, c []int64) int64 {
	res := int64(1<<63 - 1)
	for _, y := range b {
		xi := sort.Search(len(a), func(i int) bool { return a[i] > y }) - 1
		if xi < 0 {
			continue
		}
		zi := sort.Search(len(c), func(i int) bool { return c[i] >= y })
		if zi == len(c) {
			continue
		}
		v := calc(a[xi], y, c[zi])
		if v < res {
			res = v
		}
	}
	return res
}

func solve(r, g, b []int64) int64 {
	sort.Slice(r, func(i, j int) bool { return r[i] < r[j] })
	sort.Slice(g, func(i, j int) bool { return g[i] < g[j] })
	sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })
	res := int64(1<<63 - 1)
	res = min(res, best(r, g, b))
	res = min(res, best(r, b, g))
	res = min(res, best(g, r, b))
	res = min(res, best(g, b, r))
	res = min(res, best(b, r, g))
	res = min(res, best(b, g, r))
	return res
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

type testCase struct {
	r []int64
	g []int64
	b []int64
}

func parseTestcases() ([]testCase, error) {
	fields := strings.Fields(testcasesRaw)
	if len(fields) == 0 {
		return nil, fmt.Errorf("invalid test data")
	}
	t, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, err
	}
	idx := 1
	var tests []testCase
	for i := 0; i < t; i++ {
		if idx+3 > len(fields) {
			return nil, fmt.Errorf("invalid test file")
		}
		nr, _ := strconv.Atoi(fields[idx])
		ng, _ := strconv.Atoi(fields[idx+1])
		nb, _ := strconv.Atoi(fields[idx+2])
		idx += 3
		total := nr + ng + nb
		if idx+total > len(fields) {
			return nil, fmt.Errorf("invalid test file")
		}
		arr := make([]int64, total)
		for j := 0; j < total; j++ {
			v, _ := strconv.Atoi(fields[idx+j])
			arr[j] = int64(v)
		}
		idx += total
		tests = append(tests, testCase{
			r: arr[:nr],
			g: arr[nr : nr+ng],
			b: arr[nr+ng:],
		})
	}
	return tests, nil
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString("1\n")
	fmt.Fprintf(&sb, "%d %d %d\n", len(tc.r), len(tc.g), len(tc.b))
	for i, v := range tc.r {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')
	for i, v := range tc.g {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')
	for i, v := range tc.b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runBinary(bin, input string) (string, error) {
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func runCase(bin string, tc testCase) error {
	input := buildInput(tc)
	expected := strconv.FormatInt(solve(tc.r, tc.g, tc.b), 10)
	got, err := runBinary(bin, input)
	if err != nil {
		return err
	}
	if strings.TrimSpace(got) != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests, err := parseTestcases()
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
