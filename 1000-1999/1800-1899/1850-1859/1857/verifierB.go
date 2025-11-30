package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcaseData = `100
419866524
41
747338484219
3308330165712082677
7777766334373894784
8
7
382285084801
32
99999963325
591221932218
51814623951881
3968
6439870550373361807
3
7
1
999224811
7535
1
1
48706464
14
15
35547
2844666415
11
548775019454
7371705688
2
27
15
94
7312713545
91
8729583556
57931795916180304
5
5
838136237659185
83144711
6
7129292742698600633
4
6
9
87872924627
5
1057283092531
9
60437649
936518840
3276486602
2058129672
13
44824835653
4257990141418744102
64159923
100000114
6
2
3
723744
88281
97869529
1
392760006099782403
1
8555151356764419
3703536493
5397811048878
6549
6
71
99
8
5
453078
264418
12
9240
22074
726271
9
15
89
18
1
19
43
1161
1
1
9761080083
8
805
2761
447379
45869344902
21049716
27446408523
1395158980133041491
11411111
555555565
203630
7
999990
2
1887870317`

type testCase struct {
	input    string
	expected string
}

func solve(num string) string {
	digits := make([]int, len(num)+1)
	for i := 0; i < len(num); i++ {
		digits[i+1] = int(num[i] - '0')
	}
	n := len(num)
	for i := n; i > 0; i-- {
		if digits[i] >= 5 {
			digits[i] = 0
			j := i - 1
			digits[j]++
			for j > 0 && digits[j] == 10 {
				digits[j] = 0
				j--
				digits[j]++
			}
			for k := i + 1; k <= n; k++ {
				digits[k] = 0
			}
		}
	}
	start := 0
	if digits[0] == 0 {
		start = 1
	}
	var b strings.Builder
	for i := start; i <= n; i++ {
		b.WriteByte(byte(digits[i]) + '0')
	}
	return b.String()
}

func loadCases() ([]testCase, error) {
	fields := strings.Fields(testcaseData)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no testcases")
	}
	pos := 0
	t, err := strconv.Atoi(fields[pos])
	if err != nil {
		return nil, fmt.Errorf("bad test count: %w", err)
	}
	pos++
	cases := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		if pos >= len(fields) {
			return nil, fmt.Errorf("case %d: missing number", i+1)
		}
		num := fields[pos]
		pos++
		input := "1\n" + num + "\n"
		cases = append(cases, testCase{
			input:    input,
			expected: solve(num),
		})
	}
	return cases, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierB /path/to/binary")
		os.Exit(1)
	}
	cases, err := loadCases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		got, err := run(os.Args[1], tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, tc.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
