package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testcase struct {
	a      int
	b      int
	expect int
}

const testcasesRaw = `-55 -10 -10
14 -66 14
-11 -31 -11
-77 -56 -56
21 -34 21
-11 88 88
23 -38 23
-65 2 2
-66 -17 -17
27 65 65
58 27 58
-7 74 74
90 -96 90
-59 -33 -33
-57 25 25
8 -30 8
98 -31 98
-16 -91 -16
-63 61 61
98 37 98
9 90 90
49 -78 49
88 69 88
-4 -21 -4
8 -58 8
88 48 88
66 51 66
26 -96 26
-50 -29 -29
-4 20 20
-95 100 100
-83 7 7
11 6 11
-30 -6 -6
-49 -53 -49
73 -80 73
64 -10 64
-92 -25 -25
8 -60 8
-14 -14 -14
41 -40 41
55 -68 55
-73 -12 -12
87 81 87
80 17 80
8 -32 8
-4 -7 -4
-5 98 98
-18 47 47
48 -87 48
-30 53 53
90 78 90
78 -67 78
-5 -82 -5
86 77 86
82 5 82
39 -5 39
-16 83 83
55 45 55
-74 50 50
-81 -32 -32
11 57 57
48 -74 48
11 -33 11
-83 19 19
2 43 43
38 -55 38
1 25 25
-70 3 3
-44 -15 -15
-60 -67 -60
-54 -85 -54
-84 86 86
64 98 98
96 -35 96
-30 52 52
90 38 90
94 29 94
72 76 76
-19 60 60
-9 67 67
32 -1 32
81 -10 81
-96 -80 -80
31 -26 31
20 -100 20
-81 -37 -37
-98 91 91
52 -50 52
22 49 49
35 -7 35
-37 71 71
32 -39 32
25 -48 25
-59 -72 -59
22 89 89
-90 16 16
-34 -87 -34
81 -49 81
-33 -67 -33`

var testcases = mustParseTestcases(testcasesRaw)

func mustParseTestcases(raw string) []testcase {
	scanner := bufio.NewScanner(strings.NewReader(strings.TrimSpace(raw)))
	scanner.Split(bufio.ScanWords)

	readInt := func() (int, bool) {
		if !scanner.Scan() {
			return 0, false
		}
		v, err := strconv.Atoi(scanner.Text())
		if err != nil {
			panic(fmt.Sprintf("invalid integer %q: %v", scanner.Text(), err))
		}
		return v, true
	}

	var res []testcase
	for {
		a, ok := readInt()
		if !ok {
			break
		}
		b, ok := readInt()
		if !ok {
			panic("dangling value: expected b")
		}
		expect, ok := readInt()
		if !ok {
			panic("dangling value: expected expected result")
		}
		res = append(res, testcase{a: a, b: b, expect: expect})
	}

	if err := scanner.Err(); err != nil {
		panic(fmt.Sprintf("scanner error: %v", err))
	}
	if len(res) == 0 {
		panic("no testcases parsed")
	}
	return res
}

// solve replicates the logic from 1089C.go: print the larger of two integers.
func solve(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return out.String(), nil
}

func parseCandidateOutput(out string) (int, error) {
	scanner := bufio.NewScanner(strings.NewReader(out))
	scanner.Split(bufio.ScanWords)
	if !scanner.Scan() {
		return 0, fmt.Errorf("no output")
	}
	v, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return 0, fmt.Errorf("failed to parse output: %v", err)
	}
	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("scanner error: %v", err)
	}
	return v, nil
}

func checkCase(bin string, idx int, tc testcase) error {
	input := fmt.Sprintf("%d %d\n", tc.a, tc.b)
	expected := solve(tc.a, tc.b)
	out, err := runCandidate(bin, input)
	if err != nil {
		return err
	}
	got, err := parseCandidateOutput(out)
	if err != nil {
		return err
	}
	if got != expected {
		return fmt.Errorf("case %d: expected %d got %d", idx+1, expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for i, tc := range testcases {
		if err := checkCase(bin, i, tc); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
