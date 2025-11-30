package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type query struct {
	p int
	x int
}

type testCase struct {
	n int
	a []int
	q []query
}

func goodCount(a []int) int64 {
	w := 1
	var ans int64
	for i := 0; i < len(a); i++ {
		val := i + 1 - a[i] + 1
		if val > w {
			w = val
		}
		if w < 1 {
			w = 1
		}
		ans += int64(i + 1 - w + 1)
	}
	return ans
}

// solve embeds the logic from 1736C2.go for a single testcase.
func solve(tc testCase) []string {
	a := append([]int(nil), tc.a...)
	out := make([]string, 0, len(tc.q))
	for _, qu := range tc.q {
		old := a[qu.p-1]
		a[qu.p-1] = qu.x
		out = append(out, strconv.FormatInt(goodCount(a), 10))
		a[qu.p-1] = old
	}
	return out
}

// Embedded copy of testcasesC2.txt.
const testcaseData = `
4 2 3 4 1 5 1 4 3 2 2 4 4 4 2 2
3 3 2 3 1 3 1
3 3 1 2 1 2 2
7 6 7 4 4 6 7 5 4 2 3 1 1 2 4 2 3
7 7 6 7 3 4 5 7 4 5 3 5 5 4 5 2 3
1 1 5 1 1 1 1 1 1 1 1 1 1
2 2 1 4 1 1 2 2 2 1 1 1
7 6 5 3 5 3 5 2 1 3 1
2 1 1 2 2 2 2 1
1 1 3 1 1 1 1 1 1
2 2 2 2 2 2 2 2
6 1 4 5 3 1 4 5 5 6 2 1 6 6 3 4 3 6
6 5 6 3 6 4 1 5 1 6 1 3 3 6 4 3 5 5
6 2 3 2 3 3 5 3 3 4 1 1 5 6
3 2 3 1 3 1 2 1 3 2 3
2 1 2 3 1 2 1 1 2 1
8 5 4 2 1 4 6 3 5 3 2 6 3 7 5 5
8 6 7 5 7 7 1 7 3 2 1 8 7 4
1 1 5 1 1 1 1 1 1 1 1 1 1
7 5 1 1 4 6 1 2 5 3 2 6 1 5 5 4 1 5 1
6 2 3 5 4 1 3 2 2 1 5 1
3 1 2 1 1 2 3
7 1 7 3 2 3 5 5 5 4 1 4 3 7 7 1 7 1 7
3 1 1 1 1 2 1
2 2 2 2 2 1 2 2
7 5 3 3 3 2 3 4 1 2 5
1 1 1 1 1
6 4 5 6 5 4 6 1 5 4
1 1 4 1 1 1 1 1 1 1 1
2 2 1 4 1 1 2 2 2 1 2 1
7 6 1 6 3 5 5 1 5 6 1 4 2 2 7 4 1 5 1
2 2 1 1 2 1
1 1 4 1 1 1 1 1 1 1 1
6 5 2 1 2 2 6 2 4 5 6 4
5 3 5 4 3 5 4 1 4 5 2 4 2 4 5
8 3 7 3 3 2 8 8 8 5 3 3 5 4 3 6 4 5 7 5
4 3 1 3 4 4 2 2 3 2 3 4 2 4
8 4 8 1 8 2 7 1 8 2 4 2 4 5
4 2 3 2 2 5 1 3 2 1 3 2 4 1 1 1
2 2 2 1 2 2
6 1 1 3 3 4 4 4 1 2 6 5 6 4 4 2
6 1 3 1 6 4 1 4 5 3 1 5 6 3 6 3
8 5 5 2 6 2 8 6 1 3 3 3 3 6 8 2
2 1 2 5 2 2 1 2 2 2 1 1 1 1
7 3 1 3 3 4 1 7 2 1 4 5 3
4 3 3 4 4 5 1 3 4 1 2 3 1 1 1 2
4 4 4 2 2 4 2 2 2 2 3 3 3 1
8 7 7 6 5 8 5 8 1 5 4 1 2 4 8 3 8 4 4 4
1 1 1 1 1
3 2 1 3 1 1 2
4 1 4 1 3 3 3 3 2 1 1 1
6 2 1 2 4 6 2 4 3 1 1 4 1 2 6 2
7 4 4 6 1 5 7 4 2 6 4 3 1
8 8 7 8 3 8 1 5 6 3 8 6 7 4 1 4
5 3 2 4 5 2 2 2 1 2 5
7 5 2 6 1 2 1 5 2 4 4 2 1
1 1 4 1 1 1 1 1 1 1 1
8 1 4 4 2 7 8 4 3 3 2 6 2 1 5 5
8 5 8 4 5 1 6 6 6 1 1 7
2 1 1 1 1 1
3 3 1 2 1 1 3
6 2 4 3 4 3 6 1 4 3
7 1 3 2 7 4 1 5 4 5 3 5 6 7 4 2 7
7 7 5 3 2 3 7 4 4 2 4 7 6 4 3 3 7
3 3 3 3 4 2 1 1 1 3 1 2 1
2 2 1 5 1 1 2 2 2 2 1 1 1 2
2 2 2 1 2 1
4 2 1 2 2 5 1 1 2 3 2 2 4 3 2 4
3 3 1 1 1 3 1
2 1 2 1 1 2
7 2 7 7 6 5 1 6 4 7 6 7 6 3 4 5 6
8 2 5 8 7 4 2 1 8 3 2 6 6 4 8 2
6 4 6 1 5 5 2 2 3 1 3 2
7 4 1 3 2 3 2 6 2 6 2 7 6
8 7 6 4 7 8 7 8 1 3 1 3 2 1 3 5
1 1 1 1 1
5 4 4 1 3 4 2 3 2 1 4
5 4 4 1 2 2 4 3 4 3 2 4 3 4 4
4 3 3 3 1 4 3 3 3 2 1 1 2 2
1 1 1 1 1
8 6 6 1 8 3 4 1 5 4 6 1 2 8 5 5 4 8
7 6 3 3 1 1 4 1 2 7 7 5 2
3 3 2 3 3 3 1 2 1 1 1
1 1 1 1 1
7 3 5 6 2 3 2 1 2 7 5 1 4
6 4 4 3 3 4 3 5 2 5 3 2 5 1 4 4 2 4
5 1 3 5 5 5 1 4 2
5 3 1 4 1 3 1 5 4
8 4 5 6 3 7 7 6 1 3 4 1 6 6 7 5
1 1 4 1 1 1 1 1 1 1 1
2 1 2 1 2 2
8 8 5 5 1 3 4 8 3 2 3 3 8 7
1 1 4 1 1 1 1 1 1 1 1
3 3 1 2 1 2 2
3 1 2 1 2 1 1 2 1
7 3 5 1 2 1 3 2 4 2 1 3 5 4 1 6 7
6 4 1 6 3 5 6 4 5 4 4 5 5 3 4 2
8 7 4 5 3 5 3 6 5 2 3 1 1 7
3 1 3 1 2 1 2 3 3
6 5 6 3 1 4 5 5 4 3 2 5 5 2 3 4 5 5
`

// Expected outputs for each query in every testcase (precomputed).
var expectedOutputs = [][]string{
	{"7", "6", "7", "10", "7"},
	{"4"},
	{"5"},
	{"26", "26", "26", "26"},
	{"21", "24", "25", "24"},
	{"1", "1", "1", "1", "1"},
	{"2", "3", "2", "2"},
	{"15"},
	{"3", "2"},
	{"1", "1", "1"},
	{"3", "3"},
	{"17", "10", "12", "12", "12"},
	{"15", "15", "15", "14", "16"},
	{"16", "15", "17"},
	{"4", "4", "4"},
	{"3", "3", "2"},
	{"18", "19", "18"},
	{"21", "21"},
	{"1", "1", "1", "1", "1"},
	{"14", "11", "11", "9", "9"},
	{"10", "13"},
	{"4"},
	{"16", "21", "20", "20", "20"},
	{"3"},
	{"2", "3"},
	{"18"},
	{"1"},
	{"19"},
	{"1", "1", "1", "1"},
	{"2", "3", "2", "2"},
	{"13", "14", "19", "11", "11"},
	{"2"},
	{"1", "1", "1", "1"},
	{"11", "11"},
	{"13", "11", "11", "15"},
	{"23", "31", "23", "24", "23"},
	{"10", "7", "7", "10"},
	{"14", "14"},
	{"7", "6", "7", "6", "7"},
	{"3"},
	{"15", "16", "15", "13"},
	{"10", "10", "12", "12"},
	{"20", "20", "19"},
	{"3", "3", "3", "3", "3"},
	{"14", "13"},
	{"10", "7", "10", "10", "10"},
	{"7", "7", "8", "6"},
	{"17", "27", "29", "30", "27"},
	{"1"},
	{"4"},
	{"9", "5", "6"},
	{"10", "13", "13", "13"},
	{"16", "14"},
	{"19", "19", "19"},
	{"9", "12"},
	{"15", "10"},
	{"1", "1", "1", "1"},
	{"22", "20", "22"},
	{"20"},
	{"2"},
	{"4"},
	{"16"},
	{"14", "15", "15", "15"},
	{"19", "20", "20", "19"},
	{"4", "6", "4", "4"},
	{"2", "3", "3", "2", "2"},
	{"2"},
	{"6", "7", "7", "7", "7"},
	{"3"},
	{"3"},
	{"18", "18", "18", "18"},
	{"19", "21", "19"},
	{"11", "14"},
	{"14", "14"},
	{"29", "23", "29"},
	{"1"},
	{"12", "9"},
	{"10", "9", "8", "8"},
	{"7", "6", "7", "7"},
	{"1"},
	{"13", "16", "16", "16"},
	{"13", "13"},
	{"4", "4", "6"},
	{"1"},
	{"16", "14"},
	{"16", "15", "12", "17", "16"},
	{"11"},
	{"7"},
	{"17", "25", "24"},
	{"1", "1", "1", "1"},
	{"3"},
	{"19", "21"},
	{"1", "1", "1", "1"},
	{"5"},
	{"4", "3"},
	{"10", "13", "10", "11"},
	{"16", "16", "14", "13"},
	{"21", "25"},
	{"4", "6"},
	{"18", "12", "12", "12", "12"},
}

func loadTestcases() ([]testCase, error) {
	fields := strings.Fields(testcaseData)
	nums := make([]int, len(fields))
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return nil, fmt.Errorf("bad number %q: %v", f, err)
		}
		nums[i] = v
	}
	idx := 0
	tests := make([]testCase, 0, 100)
	for idx < len(nums) {
		n := nums[idx]
		idx++
		if idx+n > len(nums) {
			return nil, fmt.Errorf("incomplete array for n at position %d", len(tests)+1)
		}
		a := make([]int, n)
		copy(a, nums[idx:idx+n])
		idx += n
		if idx >= len(nums) {
			return nil, fmt.Errorf("missing q for case %d", len(tests)+1)
		}
		q := nums[idx]
		idx++
		if idx+2*q > len(nums) {
			return nil, fmt.Errorf("missing queries for case %d", len(tests)+1)
		}
		qs := make([]query, q)
		for i := 0; i < q; i++ {
			qs[i] = query{p: nums[idx+2*i], x: nums[idx+2*i+1]}
		}
		idx += 2 * q
		tests = append(tests, testCase{n: n, a: a, q: qs})
	}
	return tests, nil
}

func runCase(bin string, tc testCase, expected []string) error {
	var input strings.Builder
	input.WriteString(strconv.Itoa(tc.n))
	input.WriteByte('\n')
	for i, v := range tc.a {
		if i > 0 {
			input.WriteByte(' ')
		}
		input.WriteString(strconv.Itoa(v))
	}
	input.WriteByte('\n')
	input.WriteString(strconv.Itoa(len(tc.q)))
	input.WriteByte('\n')
	for _, qu := range tc.q {
		input.WriteString(fmt.Sprintf("%d %d\n", qu.p, qu.x))
	}

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	gotLines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(gotLines) == 1 && gotLines[0] == "" {
		gotLines = []string{}
	}
	if len(gotLines) != len(expected) {
		return fmt.Errorf("expected %d lines got %d", len(expected), len(gotLines))
	}
	for i := range expected {
		if strings.TrimSpace(gotLines[i]) != expected[i] {
			return fmt.Errorf("line %d expected %s got %s", i+1, expected[i], gotLines[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC2.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := loadTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to load testcases:", err)
		os.Exit(1)
	}
	if len(tests) != len(expectedOutputs) {
		fmt.Fprintf(os.Stderr, "testcase/expected mismatch: %d vs %d\n", len(tests), len(expectedOutputs))
		os.Exit(1)
	}

	for i, tc := range tests {
		exp := expectedOutputs[i]
		if err := runCase(bin, tc, exp); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
