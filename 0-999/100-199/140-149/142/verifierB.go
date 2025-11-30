package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// solve mirrors 142B.go logic.
func solve(n, m int) int {
	small, big := n, m
	if small > big {
		small, big = big, small
	}
	switch small {
	case 1:
		return big
	case 2:
		blocks := big / 4
		rem := big % 4
		if rem*2 > 4 {
			return blocks*4 + 4
		}
		return blocks*4 + rem*2
	default:
		return (n*m + 1) / 2
	}
}

// Embedded testcases from testcasesB.txt.
const testcaseData = `
1 1
2 1
3 1
4 1
5 1
6 1
7 1
8 1
9 1
10 1
1 2
2 2
3 2
4 2
5 2
6 2
7 2
8 2
9 2
10 2
342 250
748 334
721 892
65 196
940 582
228 245
823 991
146 823
557 459
94 83
328 897
521 956
502 112
309 565
299 724
128 561
341 835
945 554
209 987
819 618
561 602
295 456
94 611
818 395
325 590
248 298
189 194
842 192
34 628
673 267
488 71
92 696
776 134
898 154
946 40
863 83
920 717
946 850
554 700
401 858
723 538
283 535
832 242
870 221
917 696
604 846
973 430
594 282
462 505
677 657
718 939
813 366
85 333
628 119
499 602
646 344
866 195
249 17
750 278
120 723
226 381
814 175
341 437
836 64
104 802
150 876
715 225
47 837
588 650
932 959
548 617
697 76
28 128
651 194
621 851
590 123
401 94
380 854
119 38
621 23
`

type testCase struct {
	n, m int
}

func parseTestcases() ([]testCase, error) {
	fields := strings.Fields(testcaseData)
	if len(fields)%2 != 0 {
		return nil, fmt.Errorf("malformed test data")
	}
	res := make([]testCase, 0, len(fields)/2)
	for i := 0; i < len(fields); i += 2 {
		n, err := strconv.Atoi(fields[i])
		if err != nil {
			return nil, fmt.Errorf("bad n at pair %d: %v", i/2+1, err)
		}
		m, err := strconv.Atoi(fields[i+1])
		if err != nil {
			return nil, fmt.Errorf("bad m at pair %d: %v", i/2+1, err)
		}
		res = append(res, testCase{n: n, m: m})
	}
	return res, nil
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := args[0]

	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		input := fmt.Sprintf("%d %d\n", tc.n, tc.m)
		expected := strconv.Itoa(solve(tc.n, tc.m))
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Printf("test %d failed\ninput: %sexpected: %s\ngot: %s\n", idx+1, input, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
