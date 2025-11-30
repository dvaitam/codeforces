package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesB.txt so the verifier is self-contained.
const testcasesRaw = `20
3
72 97 8
5
15 63 97 57 60
7
100 26 12 62 3 49 55
10
97 98 0 89 57 34 92 29 75 13
6
3 2 3 83 69 1
7
87 27 54 92 3 67 28
8
63 70 29 44 29 86 28 97
8
37 2 53 71 82 12 23 80
5
15 95 42 92 91
9
54 64 85 24 38 36 75 63 64
7
75 4 61 31 95 51 53
3
46 70 89
6
11 56 84 65 13 99
3
66 50 47
8
93 3 60 5 39 90 78 75
10
50 82 21 21 64 29 1 98 25 69
9
29 51 65 44 73 45 58 34 84
9
77 93 0 49 100 94 65 16 66
9
26 54 7 61 46 72 70 25 64
7
62 45 53 44 0 68 69`

type testCase struct {
	n    int
	nums []int
}

// solveCase mirrors 200B.go.
func solveCase(n int, nums []int) float64 {
	sum := 0
	for i := 0; i < n; i++ {
		sum += nums[i]
	}
	return float64(sum) / float64(n)
}

func loadTestcases() ([]testCase, error) {
	scan := bufio.NewScanner(strings.NewReader(testcasesRaw))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		return nil, fmt.Errorf("no testcases found")
	}
	t, err := strconv.Atoi(scan.Text())
	if err != nil {
		return nil, fmt.Errorf("parse t: %w", err)
	}
	cases := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			return nil, fmt.Errorf("test %d: missing n", i+1)
		}
		n, err := strconv.Atoi(scan.Text())
		if err != nil {
			return nil, fmt.Errorf("test %d: parse n: %w", i+1, err)
		}
		nums := make([]int, n)
		for j := 0; j < n; j++ {
			if !scan.Scan() {
				return nil, fmt.Errorf("test %d: missing value %d", i+1, j+1)
			}
			nums[j], err = strconv.Atoi(scan.Text())
			if err != nil {
				return nil, fmt.Errorf("test %d: parse value %d: %w", i+1, j+1, err)
			}
		}
		cases = append(cases, testCase{n: n, nums: nums})
	}
	if scan.Scan() {
		return nil, fmt.Errorf("unexpected trailing data")
	}
	return cases, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	testcases, err := loadTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}
	bin := os.Args[1]
	for idx, tc := range testcases {
		expect := solveCase(tc.n, tc.nums)

		var input strings.Builder
		fmt.Fprintf(&input, "%d\n", tc.n)
		for i, v := range tc.nums {
			if i > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprintf(&input, "%d", v)
		}
		input.WriteByte('\n')

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input.String())
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx+1, err, stderr.String())
			os.Exit(1)
		}

		fields := strings.Fields(stdout.String())
		if len(fields) == 0 {
			fmt.Printf("test %d: no output\n", idx+1)
			os.Exit(1)
		}
		got, err := strconv.ParseFloat(fields[0], 64)
		if err != nil {
			fmt.Printf("test %d: failed to parse output %q\n", idx+1, fields[0])
			os.Exit(1)
		}
		if math.Abs(got-expect) > 1e-4 {
			fmt.Printf("test %d failed: expected %.6f got %.6f\n", idx+1, expect, got)
			os.Exit(1)
		}
		if len(fields) > 1 {
			fmt.Printf("test %d: extra output detected\n", idx+1)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
