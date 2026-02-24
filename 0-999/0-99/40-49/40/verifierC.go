package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sort"
)

// Embedded testcases (one per line: N x M y).
const embeddedTestcases = `4 -40 6 -4
11 35 48 -11
17 -23 39 27
3 37 38 -30
28 0 41 42
33 19 24 6
33 -46 18 -47
24 -10 30 -2
28 -29 34 21
12 -21 16 -47
12 -28 21 -33
33 -4 33 15
44 -27 36 7
27 17 48 47
24 -5 38 -4
29 46 11 1
46 9 48 33
34 12 16 -15
32 15 33 -5
43 9 30 -6
37 21 47 42
30 34 32 -22
21 -29 45 28
18 11 50 -11
20 14 46 21
34 33 33 28
38 -11 27 43
14 15 32 -4
44 -41 40 50
22 -49 47 -26
48 -43 7 23
42 -16 4 25
15 -37 44 46
34 -16 9 -19
14 4 4 41
49 -43 3 -4
24 -19 12 36
2 -36 6 -42
2 43 3 -48
24 -34 17 -30
48 16 12 38
1 25 25 -45
16 -46 10 -50
23 30 40 45
48 -14 8 -7
32 -11 2 7
36 27 50 44
3 46 17 1
40 -31 46 10
15 34 6 37
21 -47 7 7
9 24 34 49
26 15 32 -9
10 -17 22 -16
39 33 27 -48
45 -33 36 35
4 -46 17 -34
11 -38 11 8
41 15 15 40
3 -21 16 41
29 -18 5 -40
38 29 15 30
46 -18 24 37
28 17 18 46
1 -46 10 -1
27 -36 11 15
47 -20 6 -37
7 -27 2 46
15 -23 7 -47
34 9 43 8
20 32 35 -2
14 47 44 -24
47 4 28 15
2 25 38 -44
27 24 34 -27
7 11 43 -4
2 -35 34 28
24 38 19 -3
20 37 2 2
7 -11 7 -25
50 -48 44 7
4 31 27 12
30 25 14 28
5 -14 1 -47
24 42 20 -41
15 12 49 -26
8 -3 37 0
46 -33 30 46
23 -35 26 -18
8 -40 8 28
22 0 42 -23
45 -47 7 29
43 49 31 -45
47 13 46 -13
23 -32 30 -3
18 17 31 11
47 3 47 12
44 0 19 -21
11 26 32 -17
36 39 28 36`

type Interval struct {
	L, R int64
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func solve(N int, x int64, M int, y int64) int64 {
	D := abs(x - y)

	var intervals []Interval
	var isolatedU int64 = 0

	var I1, I2 int64 = 0, 0

	for u := int64(1); u <= int64(N); u++ {
		L := max(1, abs(u-D))
		R := min(int64(M), u+D)

		if L > R {
			isolatedU++
			continue
		}

		intervals = append(intervals, Interval{L, R})

		cL := max(1, abs(u-D)+1)
		cR := min(int64(M), u+D-1)

		var crossCount int64 = 0
		if cL <= cR {
			crossCount = cR - cL + 1
		}

		totalCount := R - L + 1
		touchCount := totalCount - crossCount

		I2 += crossCount
		I1 += touchCount
	}

	sort.Slice(intervals, func(i, j int) bool {
		if intervals[i].L == intervals[j].L {
			return intervals[i].R < intervals[j].R
		}
		return intervals[i].L < intervals[j].L
	})

	var mergedCount int64 = 0
	var coveredW int64 = 0
	var currentL int64 = -1
	var currentR int64 = -1

	for _, iv := range intervals {
		if currentR < iv.L {
			if currentL != -1 {
				mergedCount++
				coveredW += currentR - currentL + 1
			}
			currentL = iv.L
			currentR = iv.R
		} else {
			currentR = max(currentR, iv.R)
		}
	}

	if currentL != -1 {
		mergedCount++
		coveredW += currentR - currentL + 1
	}

	isolatedW := int64(M) - coveredW
	K_draw := mergedCount + isolatedU + isolatedW
	V_draw := I1 + 2*I2

	F := int64(1) + K_draw + V_draw

	return F
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	lines := strings.Split(strings.TrimSpace(embeddedTestcases), "\n")
	idx := 0
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		idx++
		var N, M int
		var x, y int64
		fmt.Sscan(line, &N, &x, &M, &y)
		exp := solve(N, x, M, y)
		input := fmt.Sprintf("%d %d %d %d\n", N, x, M, y)
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewBufferString(input)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\nstderr: %s\n", idx, err, string(out))
			os.Exit(1)
		}
		var got int64
		fmt.Sscan(strings.TrimSpace(string(out)), &got)
		if got != exp {
			fmt.Printf("Test %d failed: expected %d got %s\n", idx, exp, strings.TrimSpace(string(out)))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}