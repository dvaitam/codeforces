package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `42 41 1722
27 40 1080
13 18 234
39 44 1716
4 38 76
6 11 66
1 49 49
8 2 8
44 10 220
11 23 253
14 3 42
29 44 1276
26 4 52
19 14 266
47 5 235
45 26 1170
30 38 570
22 19 418
50 16 400
5 46 230
28 25 700
12 50 300
50 44 1100
12 2 12
32 31 992
11 14 154
42 8 168
2 49 98
45 42 630
36 22 396
17 41 697
33 48 528
25 23 575
16 9 144
18 47 846
27 1 27
31 28 868
44 11 44
28 18 252
28 2 28
23 15 345
49 10 490
11 26 286
12 4 12
4 28 28
15 24 120
35 8 280
25 26 650
16 5 80
6 49 294
18 9 18
26 49 1274
49 16 784
6 24 24
6 16 48
3 6 6
45 24 360
40 9 360
24 14 168
11 31 341
21 25 525
31 42 1302
21 27 189
23 26 598
11 4 44
34 17 34
44 44 44
12 11 132
28 20 140
17 42 714
44 49 2156
32 16 32
23 50 1150
27 25 675
44 24 264
48 46 1104
18 31 558
27 48 432
5 38 190
19 2 38
24 16 48
8 15 120
35 25 175
47 26 1222
5 36 180
50 28 700
13 4 52
22 22 22
42 41 1722
20 27 540
10 23 230
25 43 1075
14 30 210
30 3 30
28 35 140
36 16 144
12 6 12
34 20 340
49 20 980
4 12 12`

type testCase struct {
	a int64
	b int64
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

// solve embeds the lcm logic from 1089I.go.
func solve(a, b int64) int64 {
	g := gcd(a, b)
	return a / g * b
}

func runCase(bin string, tc testCase) error {
	input := fmt.Sprintf("%d %d\n", tc.a, tc.b)
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	want := strconv.FormatInt(solve(tc.a, tc.b), 10)
	if got != want {
		return fmt.Errorf("expected %s got %s", want, got)
	}
	return nil
}

func parseTestcases(raw string) ([]testCase, error) {
	fields := strings.Fields(raw)
	if len(fields)%3 != 0 {
		return nil, fmt.Errorf("invalid testcase data: expected triples")
	}
	tests := make([]testCase, 0, len(fields)/3)
	for i := 0; i < len(fields); i += 3 {
		a, err := strconv.ParseInt(fields[i], 10, 64)
		if err != nil {
			return nil, err
		}
		b, err := strconv.ParseInt(fields[i+1], 10, 64)
		if err != nil {
			return nil, err
		}
		exp, err := strconv.ParseInt(fields[i+2], 10, 64)
		if err != nil {
			return nil, err
		}
		if solve(a, b) != exp {
			return nil, fmt.Errorf("embedded expected mismatch for %d %d", a, b)
		}
		tests = append(tests, testCase{a: a, b: b})
	}
	return tests, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run verifierI.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests, err := parseTestcases(testcasesRaw)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for i, tc := range tests {
		if err := runCase(bin, tc); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d cases passed\n", len(tests))
}
