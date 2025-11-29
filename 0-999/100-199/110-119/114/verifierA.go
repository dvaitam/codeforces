package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesData = `
14 432
3 729
17 83521
17 4913
18 5832
5 258
19 831
6 36
12 144
13 4826809
8 490
16 268
3 563
2 32
2 849
12 20736
4 4096
9 729
19 361
12 957
17 24137569
11 562
12 555
8 619
19 457
4 397
12 299
7 343
3 674
10 100
6 947
3 922
19 859
18 5832
8 605
15 463
17 719
13 4826809
5 625
8 4096
5 382
7 49
5 877
9 531441
19 77
2 8
20 400
13 120
3 202
7 492
8 65
2 637
5 73
9 6561
13 169
18 34012224
5 402
8 32768
20 8000
3 164
7 2401
5 454
7 16807
20 941
11 399
10 100
16 345
3 140
9 496
13 691
13 915
6 399
15 84
2 717
12 1728
16 16777216
15 11390625
15 728
3 9
10 459
18 931
19 2
3 81
16 1048576
8 988
4 135
2 32
12 144
2 693
18 197
5 206
11 1331
5 3125
4 1024
5 138
18 662
13 2197
10 45
3 81
19 130321
`

type testCase struct {
	k int64
	l int64
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(testcasesData, "\n")
	var cases []testCase
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid testcase line: %q", line)
		}
		k, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return nil, err
		}
		l, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return nil, err
		}
		cases = append(cases, testCase{k: k, l: l})
	}
	return cases, nil
}

// solve mirrors 114A.go to avoid external oracle.
func solve(tc testCase) string {
	k, l := tc.k, tc.l
	w := int64(0)
	for k != 0 && l%k == 0 {
		l /= k
		w++
	}
	if w == 0 || l > 1 {
		return "NO"
	}
	return fmt.Sprintf("YES\n%d", w-1)
}

func run(bin string, input string) (string, error) {
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	testcases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range testcases {
		input := fmt.Sprintf("%d\n%d\n", tc.k, tc.l)
		want := solve(tc)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
