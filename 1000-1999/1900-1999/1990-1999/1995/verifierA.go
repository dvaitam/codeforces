package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesA = `50 1722
6 16
66 3980
52 1242
62 1466
75 1789
65 1140
37 286
97 1553
80 2052
69 1203
40 202
94 1208
88 2704
61 2292
13 90
56 1295
79 5246
27 565
62 1813
98 2596
18 192
42 330
80 1941
58 620
71 1025
87 3037
55 1314
20 150
67 1547
70 2790
59 1416
7 24
71 1026
61 1760
8 18
90 3511
84 948
66 2100
57 1348
76 4537
65 1089
62 1738
67 1435
97 1516
90 3688
60 756
71 1250
34 240
12 25
80 2583
81 2960
90 3421
62 2002
28 414
98 5536
42 395
8 18
89 2736
13 52
56 1189
63 1121
47 385
82 1962
12 53
33 264
56 837
57 1106
61 1958
61 927
11 23
90 2372
96 3323
87 3563
73 1772
98 2191
100 4126
56 1522
53 685
68 2005
81 2239
97 4607
41 492
86 2032
42 270
49 519
100 3874
5 14
19 200
97 1333
25 185
51 912
8 31
5 4
60 1402
63 1628
85 1151
30 200
92 1658
56 1295
49 726
78 2389`

func expected(n, k int) int {
	if k == 0 {
		return 0
	}
	ans := 0
	if k >= n {
		k -= n
		ans++
	}
	for i := n - 1; i >= 1; i-- {
		if k >= i {
			k -= i
			ans++
		}
		if k >= i {
			k -= i
			ans++
		}
	}
	return ans
}

type testCase struct {
	n int
	k int
}

func parseTests() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesA), "\n")
	tests := make([]testCase, len(lines))
	for i, line := range lines {
		fields := strings.Fields(line)
		if len(fields) != 2 {
			return nil, fmt.Errorf("bad test line %d", i+1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, err
		}
		k, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, err
		}
		tests[i] = testCase{n: n, k: k}
	}
	return tests, nil
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.k)
	}
	return sb.String()
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierA /path/to/binary")
		os.Exit(1)
	}
	tests, err := parseTests()
	if err != nil {
		fmt.Println("parse error:", err)
		os.Exit(1)
	}
	input := buildInput(tests)
	output, err := runCandidate(os.Args[1], input)
	if err != nil {
		fmt.Println("runtime error:", err)
		os.Exit(1)
	}
	outFields := strings.Fields(output)
	if len(outFields) != len(tests) {
		fmt.Printf("expected %d outputs, got %d\n", len(tests), len(outFields))
		os.Exit(1)
	}
	for i, tc := range tests {
		want := strconv.Itoa(expected(tc.n, tc.k))
		if outFields[i] != want {
			fmt.Printf("case %d failed\nexpected: %s\ngot: %s\n", i+1, want, outFields[i])
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
