package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcaseData = `
675 665 2
655 191 1
869 571 1
487 981 1
82 730 4
917 265 0
53 387 0
401 474 3
680 144 1
60 943 0
155 647 0
855 413 4
914 96 2
84 590 4
103 222 2
866 845 3
924 176 4
464 719 2
783 345 1
663 827 3
994 333 2
70 549 0
833 382 3
957 809 2
922 893 0
556 827 4
289 727 0
954 282 2
310 609 4
115 114 3
805 902 2
112 334 2
300 955 0
16 39 4
314 407 2
977 638 3
35 310 3
827 538 0
620 555 1
100 917 0
727 507 4
745 268 4
287 17 1
481 132 0
323 842 2
972 966 0
577 149 2
594 106 4
723 226 0
209 452 1
17 267 0
604 779 2
278 401 1
728 55 1
787 923 3
797 341 2
111 445 2
98 174 4
262 461 4
554 824 2
481 373 4
306 878 4
436 559 2
535 987 4
488 6 4
977 866 2
931 84 1
131 246 3
598 269 4
965 360 2
161 926 2
952 717 1
450 600 4
79 241 2
208 27 4
950 776 1
594 756 0
481 488 3
917 942 1
734 68 3
40 330 0
715 603 4
928 126 3
903 135 2
746 33 3
784 953 1
369 642 2
845 495 4
944 590 2
235 795 3
913 37 2
174 646 4
145 160 1
41 17 1
520 165 4
481 59 1
49 434 3
218 9 0
957 652 4
22 919 0
422 420 2
65 192 3
94 257 1
`

func solve(n, m, k int64) int64 {
	var ans int64
	if k == 1 {
		ans = 1
	} else if k == 2 {
		part1 := m
		if n < m {
			part1 = n
		}
		var part2 int64
		if m >= n {
			part2 = m/n - 1
		}
		if part2 < 0 {
			part2 = 0
		}
		ans = part1 + part2
	} else if k == 3 {
		part1 := m
		if n < m {
			part1 = n
		}
		var part2 int64
		if m >= n {
			part2 = m/n - 1
		}
		if part2 < 0 {
			part2 = 0
		}
		ans = m - (part1 + part2)
	} else {
		ans = 0
	}
	return ans
}

type testCase struct {
	input    string
	expected string
}

func loadCases() ([]testCase, error) {
	lines := strings.Split(testcaseData, "\n")
	cases := []testCase{}
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 3 {
			return nil, fmt.Errorf("line %d: expected 3 numbers, got %d", idx+1, len(parts))
		}
		n, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("line %d: bad n: %w", idx+1, err)
		}
		m, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("line %d: bad m: %w", idx+1, err)
		}
		k, err := strconv.ParseInt(parts[2], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("line %d: bad k: %w", idx+1, err)
		}
		var sb strings.Builder
		sb.WriteString("1\n")
		fmt.Fprintf(&sb, "%d %d %d\n", n, m, k)
		cases = append(cases, testCase{
			input:    sb.String(),
			expected: strconv.FormatInt(solve(n, m, k), 10),
		})
	}
	return cases, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierC /path/to/binary")
		os.Exit(1)
	}
	cases, err := loadCases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load cases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		got, err := run(os.Args[1], tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", idx+1, tc.expected, got, tc.input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
