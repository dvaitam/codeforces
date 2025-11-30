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

func runBinary(bin string, input []byte) ([]byte, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	err := cmd.Run()
	if err != nil {
		return append(out.Bytes(), errb.Bytes()...), fmt.Errorf("%v", err)
	}
	return out.Bytes(), nil
}

func minCreepiness(a, b int) int {
	if a == 0 || b == 0 {
		if a > b {
			return a
		}
		return b
	}
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	if diff == 0 {
		return 1
	}
	return diff
}

func creepiness(s string) (int, bool) {
	zeros, ones, maxDiff := 0, 0, 0
	for i := 0; i < len(s); i++ {
		if s[i] == '0' {
			zeros++
		} else if s[i] == '1' {
			ones++
		} else {
			return 0, false
		}
		diff := zeros - ones
		if diff < 0 {
			diff = -diff
		}
		if diff > maxDiff {
			maxDiff = diff
		}
	}
	return maxDiff, true
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := loadTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load embedded testcases: %v\n", err)
		os.Exit(1)
	}

	// Build input for all cases at once.
	var inBuf bytes.Buffer
	fmt.Fprintln(&inBuf, len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&inBuf, "%d %d\n", tc.a, tc.b)
	}

	outBytes, err := runBinary(bin, inBuf.Bytes())
	if err != nil {
		fmt.Fprintf(os.Stderr, "binary execution failed: %v\n%s", err, outBytes)
		os.Exit(1)
	}
	outBuf := bytes.NewBuffer(outBytes)

	t := len(tests)
	outScanner := bufio.NewScanner(outBuf)
	outScanner.Split(bufio.ScanWords)
	for caseIdx := 1; caseIdx <= t; caseIdx++ {
		a := tests[caseIdx-1].a
		b := tests[caseIdx-1].b

		if !outScanner.Scan() {
			fmt.Fprintf(os.Stderr, "not enough output for case %d\n", caseIdx)
			os.Exit(1)
		}
		ans := outScanner.Text()
		if len(ans) != a+b {
			fmt.Fprintf(os.Stderr, "case %d wrong length\n", caseIdx)
			os.Exit(1)
		}
		zeros := strings.Count(ans, "0")
		ones := len(ans) - zeros
		if zeros != a || ones != b {
			fmt.Fprintf(os.Stderr, "case %d wrong number of zeros/ones\n", caseIdx)
			os.Exit(1)
		}
		c, ok := creepiness(ans)
		if !ok {
			fmt.Fprintf(os.Stderr, "case %d output has invalid characters\n", caseIdx)
			os.Exit(1)
		}
		exp := minCreepiness(a, b)
		if c != exp {
			fmt.Fprintf(os.Stderr, "case %d wrong creepiness got %d expected %d\n", caseIdx, c, exp)
			os.Exit(1)
		}
	}
	if outScanner.Scan() {
		fmt.Fprintln(os.Stderr, "extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}

type testCase struct {
	a int
	b int
}

// Embedded copy of testcasesA.txt (t followed by a b pairs).
const testcaseData = `
100
50 98
54 6
34 66
63 52
39 62
46 75
28 65
18 37
18 97
13 80
33 69
91 78
19 40
13 94
10 88
43 61
72 13
46 56
41 79
82 27
71 62
57 67
34 8
71 2
12 93
52 91
86 81
1 79
64 43
32 94
42 91
9 25
73 29
31 19
70 58
12 11
41 66
63 14
39 71
38 91
16 71
43 70
27 78
71 76
37 57
12 77
50 41
74 31
38 24
25 24
17 24
49 45
80 76
13 94
70 12
33 76
28 90
32 30
83 90
56 43
64 74
39 47
88 81
21 50
87 58
5 68
6 45
71 53
44 29
32 25
45 73
53 58
44 52
64 26
95 65
65 20
53 96
71 21
93 33
49 60
20 7
58 22
87 37
72 35
90 68
84 55
68 68
47 27
81 55
70 21
98 24
57 50
77 42
30 95
5 27
60 74
12 55
62 16
33 14
36 35
16 69
18 17
63 92
61 64
70 57
2 10
95 13
3 72
`

func loadTestcases() ([]testCase, error) {
	fields := strings.Fields(testcaseData)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no testcases")
	}
	pos := 0
	nextInt := func() (int, error) {
		if pos >= len(fields) {
			return 0, fmt.Errorf("unexpected end of data")
		}
		v, err := strconv.Atoi(fields[pos])
		pos++
		return v, err
	}
	t, err := nextInt()
	if err != nil {
		return nil, fmt.Errorf("bad count: %v", err)
	}
	tests := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		a, err := nextInt()
		if err != nil {
			return nil, fmt.Errorf("case %d missing a: %v", i+1, err)
		}
		b, err := nextInt()
		if err != nil {
			return nil, fmt.Errorf("case %d missing b: %v", i+1, err)
		}
		tests = append(tests, testCase{a: a, b: b})
	}
	return tests, nil
}
