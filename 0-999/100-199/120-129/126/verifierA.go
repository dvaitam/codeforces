package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `885441 935935 933489 441002 935114
42451 313944 509533 424605 310506
962839 982716 375442 611721 978456
934974 994488 529203 146040 949289
295529 442063 648407 838235 320388
262675 821108 325214 103561 416775
765285 784616 495078 587008 776104
105593 476570 331557 640562 333224
671533 778738 500182 464198 743953
907344 975678 65305 844133 941487
71 100 12 93 71
52 97 81 1 94
79 94 32 94 89
42 87 25 73 46
29 59 19 70 54
58 63 41 66 58
63 69 71 38 65
91 92 70 27 92
78 95 57 12 87
77 89 74 31 82
38 49 24 5 41
79 100 61 9 87
12 98 20 5 28
11 100 88 51 80
91 99 67 31 95
28 81 36 58 65
64 86 42 79 66
15 77 81 43 52
25 56 94 35 26
15 43 22 43 26
55 58 19 90 55
29 34 82 69 33
78 99 4 16 80
82 88 74 16 86
51 56 15 5 53
78 78 24 92 78
16 77 94 8 29
87 87 80 13 87
34 42 10 83 37
39 61 24 8 52
26 61 20 14 41
58 59 11 92 51
41 68 2 2 17
27 92 43 29 88
76 95 27 48 94
20 46 83 9 43
62 77 50 50 63
84 95 59 26 87
74 89 4 33 80
66 98 2 68 93
28 87 85 60 78
100 100 25 70 99
21 97 38 10 73
8 9 1 1 1
25 49 27 8 38
78 87 9 10 49
69 98 24 89 94
68 74 12 68 71
13 60 53 12 35
2 99 44 28 28
48 61 19 59 49
35 66 15 36 58
15 30 3 29 17
6 13 6 6 7
17 98 34 52 62
93 99 21 54 92
87 90 15 12 77
12 80 47 17 64
57 79 57 17 57
90 96 7 93 93
8 61 27 34 54
8 67 9 27 45
56 85 7 70 76
79 83 26 1 65
2 99 25 62 25
39 53 36 49 45
16 98 1 64 98
27 89 64 89 89
5 15 7 9 14
52 76 3 32 74
27 97 4 4 25
90 93 10 14 86
10 34 1 30 21
19 50 28 37 48
31 81 10 65 73
15 40 4 9 21
2 8 1 6 5
13 74 21 22 58
40 67 18 2 21
16 93 25 90 89
64 93 6 20 66
25 98 39 62 68
46 76 4 16 62
18 25 17 2 10
98 99 38 21 99
62 85 13 19 71
42 46 13 30 37
3 31 1 11 12
8 66 19 52 42
31 32 2 5 4`

// solve mirrors 126A.go logic.
func solve(t1, t2, x1, x2, t0 int64) (int64, int64) {
	if t0 == t1 && t0 == t2 {
		return x1, x2
	}
	if t0 == t2 {
		return 0, x2
	}
	bestY1, bestY2 := int64(0), x2
	bestN := (t2 - t0) * x2
	bestD := x2
	denomDiff := t2 - t0
	for y1 := int64(1); y1 <= x1; y1++ {
		need := y1 * (t0 - t1)
		var y2 int64
		if need > 0 {
			y2 = (need + denomDiff - 1) / denomDiff
		} else {
			y2 = 0
		}
		if y2 > x2 {
			continue
		}
		sum := y1 + y2
		N := (t1-t0)*y1 + (t2-t0)*y2
		if N < 0 {
			continue
		}
		D := sum
		if N*bestD < bestN*D || (N*bestD == bestN*D && sum > bestY1+bestY2) {
			bestY1, bestY2 = y1, y2
			bestN, bestD = N, D
		}
	}
	return bestY1, bestY2
}

type testCase struct {
	t1 int64
	t2 int64
	x1 int64
	x2 int64
	t0 int64
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(testcasesRaw, "\n")
	var cases []testCase
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 5 {
			return nil, fmt.Errorf("invalid test line")
		}
		vals := make([]int64, 5)
		for i, p := range parts {
			v, err := strconv.ParseInt(p, 10, 64)
			if err != nil {
				return nil, err
			}
			vals[i] = v
		}
		cases = append(cases, testCase{t1: vals[0], t2: vals[1], x1: vals[2], x2: vals[3], t0: vals[4]})
	}
	return cases, nil
}

func buildInput(tc testCase) string {
	return fmt.Sprintf("%d %d %d %d %d\n", tc.t1, tc.t2, tc.x1, tc.x2, tc.t0)
}

func runBinary(bin, input string) (string, error) {
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func runCase(bin string, tc testCase) error {
	input := buildInput(tc)
	expectedX1, expectedX2 := solve(tc.t1, tc.t2, tc.x1, tc.x2, tc.t0)
	expected := fmt.Sprintf("%d %d", expectedX1, expectedX2)
	got, err := runBinary(bin, input)
	if err != nil {
		return err
	}
	if strings.TrimSpace(got) != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests, err := parseTestcases()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for i, tc := range tests {
		if err := runCase(bin, tc); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
