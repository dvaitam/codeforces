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
	n      int
	expect int
}

const testcasesRaw = `1 14
2 12
3 13
4 8
5 9
6 10
7 11
8 0
9 1
10 2
11 3
12 4
13 5
14 6
15 7
16 15
17 15
18 15
19 15
20 15
21 15
22 15
23 15
24 15
25 15
26 15
27 15
28 15
29 15
30 15
31 15
32 15
33 15
34 15
35 15
36 15
37 15
38 15
39 15
40 15
41 15
42 15
43 15
44 15
45 15
46 15
47 15
48 15
49 15
50 15
51 15
52 15
53 15
54 15
55 15
56 15
57 15
58 15
59 15
60 15
61 15
62 15
63 15
64 15
65 15
66 15
67 15
68 15
69 15
70 15
71 15
72 15
73 15
74 15
75 15
76 15
77 15
78 15
79 15
80 15
81 15
82 15
83 15
84 15
85 15
86 15
87 15
88 15
89 15
90 15
91 15
92 15
93 15
94 15
95 15
96 15
97 15
98 15
99 15
100 15`

var testcases = mustParseTestcases(testcasesRaw)

func mustParseTestcases(raw string) []testcase {
	scanner := bufio.NewScanner(strings.NewReader(strings.TrimSpace(raw)))
	scanner.Split(bufio.ScanWords)

	cases := []testcase{}
	for {
		if !scanner.Scan() {
			break
		}
		nVal, err := strconv.Atoi(scanner.Text())
		if err != nil {
			panic(fmt.Sprintf("invalid n %q: %v", scanner.Text(), err))
		}
		if !scanner.Scan() {
			panic("dangling expected value")
		}
		expVal, err := strconv.Atoi(scanner.Text())
		if err != nil {
			panic(fmt.Sprintf("invalid expected %q: %v", scanner.Text(), err))
		}
		cases = append(cases, testcase{n: nVal, expect: expVal})
	}
	if err := scanner.Err(); err != nil {
		panic(fmt.Sprintf("scanner error: %v", err))
	}
	if len(cases) == 0 {
		panic("no testcases parsed")
	}
	return cases
}

// solve replicates 1145C.go behavior: map n to predefined array else 15.
func solve(n int) int {
	arr := []int{14, 12, 13, 8, 9, 10, 11, 0, 1, 2, 3, 4, 5, 6, 7}
	if n >= 1 && n <= len(arr) {
		return arr[n-1]
	}
	return 15
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
	input := fmt.Sprintf("%d\n", tc.n)
	expected := solve(tc.n)
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
