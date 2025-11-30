package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded testcases (one per line) in the same flattened format the solver expects.
const testcasesB2 = `6 582 10 4 4 2 9 9 7 9 5 2 3 7
5 432 2 2 7 2 2 7 3 1 8 7
6 427 1 8 6 5 2 6 2 2 6 1 6 6
2 10 4 6 2 10
2 212 1 4 2 1
3 377 1 10 4 3 3 8
1 488 6 5
2 28 4 6 6 8
3 303 9 6 3 10 2 2
5 594 5 3 7 3 3 4 6 9 4 4
2 298 6 7 1 3
5 21 7 2 2 3 7 5 9 7 3 10
4 305 6 2 4 8 6 9 1 7
4 8 7 6 8 4 6 5 8 2
2 815 2 5 2 9
5 704 3 8 7 3 7 7 3 4 8 6
5 145 6 8 2 8 4 5 1 8 10 8
1 223 5 2
6 308 9 10 3 7 8 2 8 4 9 7 5 1
1 276 1 1
3 407 9 10 7 8 2 5
3 290 4 10 2 1 2 5
3 546 6 2 9 4 3 2
4 881 5 5 9 3 10 9 4 9
1 420 9 7
6 797 5 5 8 6 10 3 3 2 2 7 7 10
4 142 9 5 6 8 7 4 8 8
6 513 6 8 1 8 5 3 8 1 10 4 1 6
4 400 1 9 2 2 7 1 6 1
1 635 1 5
6 716 5 4 3 10 5 4 2 7 8 6 7 3
3 431 7 3 8 3 9 6
2 213 3 8 6 7
4 826 8 7 4 4 8 4 10 1
4 34 4 2 3 6 1 3 4 10
3 627 2 9 5 6 7 8
1 647 9 9
6 956 7 10 8 8 5 8 4 6 5 1 1 1
2 358 1 5 1 3
1 805 7 4
5 406 9 4 8 4 6 10 2 10 2 6
3 548 8 6 5 1 9 1
2 377 2 4 9 6
2 857 4 5 5 5
5 880 7 5 8 6 4 1 5 9 2 1
4 507 8 1 7 8 8 8 2 2
1 246 2 3
4 905 4 8 10 2 7 9 7 1
2 255 8 4 3 5
3 327 7 2 9 5 10 9
2 728 5 8 9 10
4 549 5 5 4 1 2 10 2 3
6 424 4 4 5 1 9 9 7 1 2 7 5 2
6 578 6 4 9 5 4 4 2 9 5 6 4 6
6 491 5 10 3 3 1 9 9 6 6 10 1 3
4 158 3 9 2 3 4 8 10 4
2 750 3 4 7 6
5 605 3 8 2 10 1 9 10 6 8 8
3 12 4 9 3 8 8 9
3 721 2 5 3 10 7 4
3 804 5 7 1 4 1 6
6 765 4 6 8 4 5 6 3 5 1 6 10 9
1 749 3 6
1 502 1 1
2 46 1 4 6 2
1 353 7 3
2 459 7 3 6 5
2 664 6 7 7 1
4 912 5 9 9 8 1 10 2 7
4 175 1 9 3 10 9 3 2 6
2 860 3 4 1 3
6 806 9 3 2 7 10 2 10 8 3 10 10 1
3 348 7 1 1 8 2 6
3 687 3 8 4 9 6 3
6 771 7 6 5 8 7 1 5 9 5 9 8 1
5 584 9 5 1 8 7 2 7 6 8 1
1 279 1 5
6 696 10 5 4 9 9 6 7 5 4 2 10 6
2 600 9 6 3 3
3 866 1 10 1 10 3 6
3 298 5 6 8 7 10 7
2 1 3 10 1 8
2 349 1 8 5 10
2 71 9 7 5 3
5 173 2 3 10 2 9 9 10 7 7 5
3 291 1 7 5 5 9 9
5 325 6 4 7 3 1 9 3 10 7 6
4 972 1 9 7 10 4 1 6 9
2 696 4 6 8 1
6 752 4 10 4 5 3 7 2 10 8 4 8 9
6 102 4 3 8 2 7 7 5 5 7 6 10 6
1 314 1 8
1 780 5 4
4 395 7 7 1 10 8 6 10 3
5 727 5 6 1 7 8 9 3 1 2 10
3 369 1 2 4 2 9 8
1 322 1 6
4 926 3 5 7 3 10 3 7 5
5 61 3 3 3 8 1 9 1 9 7 3
3 819 10 2 2 9 3 5`

type testCase struct {
	n, m    int64
	a, freq []int64
}

// Solver logic embedded from 1995B2.go.
func min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

// solve1995B2 mirrors the reference solution for a single testcase.
func solve1995B2(tc testCase) int64 {
	mp := make(map[int64]int64, tc.n)
	for i := int64(0); i < tc.n; i++ {
		mp[tc.a[i]] = tc.freq[i]
	}

	var mx int64
	for i := int64(0); i < tc.n; i++ {
		f := min64(tc.m/tc.a[i], mp[tc.a[i]])
		f1 := min64((tc.m-f*tc.a[i])/(tc.a[i]+1), mp[tc.a[i]+1])
		ans := f*tc.a[i] + f1*(tc.a[i]+1)

		f3 := min64(f, mp[tc.a[i]+1]-f1)
		ans += f3

		mx = max64(mx, min64(ans, tc.m))
		if mx == tc.m {
			break
		}
	}
	return mx
}

func parseLine(line string) (testCase, error) {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return testCase{}, fmt.Errorf("not enough fields")
	}
	n, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return testCase{}, fmt.Errorf("invalid n: %v", err)
	}
	m, err := strconv.ParseInt(fields[1], 10, 64)
	if err != nil {
		return testCase{}, fmt.Errorf("invalid m: %v", err)
	}
	expectLen := 2 + 2*int(n)
	if len(fields) != expectLen {
		return testCase{}, fmt.Errorf("expected %d fields, got %d", expectLen, len(fields))
	}
	a := make([]int64, n)
	freq := make([]int64, n)
	for i := int64(0); i < n; i++ {
		v, err := strconv.ParseInt(fields[2+i], 10, 64)
		if err != nil {
			return testCase{}, fmt.Errorf("invalid a[%d]: %v", i, err)
		}
		a[i] = v
	}
	for i := int64(0); i < n; i++ {
		v, err := strconv.ParseInt(fields[2+int(n)+int(i)], 10, 64)
		if err != nil {
			return testCase{}, fmt.Errorf("invalid freq[%d]: %v", i, err)
		}
		freq[i] = v
	}
	return testCase{n: n, m: m, a: a, freq: freq}, nil
}

func loadTestcases(data string) ([]testCase, error) {
	lines := strings.Split(data, "\n")
	cases := make([]testCase, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		tc, err := parseLine(line)
		if err != nil {
			return nil, fmt.Errorf("parse error: %w", err)
		}
		cases = append(cases, tc)
	}
	return cases, nil
}

func buildInput(tc testCase) string {
	var b strings.Builder
	b.WriteString("1\n")
	b.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
	for i := int64(0); i < tc.n; i++ {
		b.WriteString(fmt.Sprintf("%d ", tc.a[i]))
	}
	b.WriteString("\n")
	for i := int64(0); i < tc.n; i++ {
		b.WriteString(fmt.Sprintf("%d ", tc.freq[i]))
	}
	b.WriteString("\n")
	return b.String()
}

func runBinary(bin, input string) (string, string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB2.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := loadTestcases(testcasesB2)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load embedded testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		expect := solve1995B2(tc)

		out, stderr, err := runBinary(bin, buildInput(tc))
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx+1, err, stderr)
			os.Exit(1)
		}
		gotStr := strings.TrimSpace(out)
		got, err := strconv.ParseInt(gotStr, 10, 64)
		if err != nil {
			fmt.Printf("test %d: failed to parse output %q\n", idx+1, gotStr)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("test %d failed: expected %d got %d\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
