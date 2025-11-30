package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type testCase struct {
	n   int
	arr []int64
}

// Embedded testcases from testcasesD.txt.
const testcaseData = `
5 97 44 13 150 43 2 48
5 70 58 237 1 43 21 65
5 52 66 14 227 54 47 41
6 51 322 71 97 74 37 32 34
2 71 93 37 108
1 93 36 36
6 70 283 41 66 39 10 27 71
6 4 41 72 54 221 96 21 5
2 17 76 91 15
4 59 26 43 8 47 140
1 5 63 63
5 30 98 11 147 6 2 68
3 63 32 48 160 65
2 77 24 25 101
5 61 47 3 63 3 14 177
3 10 83 54 147 25
5 88 414 78 58 74 95 99
5 33 179 34 2 40 70 78
5 90 90 57 43 61 27 278
3 132 70 33 97 2
3 76 160 7 44 40
6 53 100 89 17 336 83 9 38
4 192 26 60 54 68 10
6 51 32 49 29 280 80 27 92
4 64 22 192 4 19 87
3 168 78 15 83 75
3 63 191 82 12 97
4 7 14 48 30 77 99
1 9 28 9
4 9 198 92 3 8 94
5 72 61 2 35 189 26 28
2 59 8 51 31
3 73 3 72 53 16
2 62 13 39 23
1 57 20 57
6 235 8 26 12 37 45 57 76
6 37 283 12 94 15 96 92 29
2 42 45 87 46
6 44 256 63 55 32 50 23 12
3 97 198 48 59 53
4 17 66 134 39 48 3
2 81 142 59 83
1 33 79 79
6 10 22 210 41 61 55 40 36
1 24 88 88
1 93 93 25
6 46 81 390 65 65 47 48 86
5 42 198 92 8 27 32 39
3 1 164 26 73 90
5 35 21 19 30 115 40 10
2 112 46 66 26
2 19 128 30 98
2 110 24 47 63
5 12 28 31 90 87 57 248
2 103 62 71 41
4 72 45 63 9 18 198
5 10 62 54 34 44 182 12
2 41 33 34 75
1 38 38 21
1 56 15 15
3 46 99 84 229 93
5 46 19 204 87 23 42 33
2 90 32 152 62
1 10 95 95
5 22 42 90 62 73 67 283
2 48 76 152 76
3 177 31 57 89 83
3 57 21 64 18 25
2 63 62 136 73
5 236 66 57 84 8 27 21
2 24 17 43 60
4 66 59 40 71 22 236
3 195 1 20 87 88
6 57 23 405 80 84 66 95 58
1 50 69 69
4 76 67 58 61 210 24
4 97 100 50 45 233 38
5 86 8 30 204 71 9 58
4 199 96 53 7 60 36
5 62 55 65 356 20 79 95
2 84 87 28 59
4 134 59 7 48 72 20
4 92 92 58 95 284 39
2 76 63 13 24
1 69 69 16
3 64 87 223 82 72
3 203 51 63 69 89
5 80 261 7 98 42 100 32
3 239 53 77 85 77
4 50 23 326 86 100 90
2 88 79 9 41
1 1 56 56
4 20 23 198 94 15 69
5 36 67 82 302 27 25 90
4 67 75 56 18 215 17
1 46 100 46
2 64 153 89 76
2 130 73 69 61
6 50 95 61 253 36 45 9 18
6 67 10 25 6 35 90 262 39
`

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	res := make([]testCase, 0, len(lines))
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("case %d bad n: %v", i+1, err)
		}
		expect := n + 3
		if len(fields) != expect {
			return nil, fmt.Errorf("case %d expected %d values, got %d", i+1, expect, len(fields))
		}
		arr := make([]int64, n+2)
		for j := 0; j < n+2; j++ {
			v, err := strconv.ParseInt(fields[1+j], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("case %d bad value %d: %v", i+1, j+1, err)
			}
			arr[j] = v
		}
		res = append(res, testCase{n: n, arr: arr})
	}
	if len(res) == 0 {
		return nil, fmt.Errorf("no test data")
	}
	return res, nil
}

func solve(n int, arr []int64) string {
	total := n + 2
	if len(arr) != total {
		return "-1"
	}
	var sum int64
	for _, v := range arr {
		sum += v
	}
	sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
	y := arr[total-1]
	x := sum - 2*y
	if idx := indexOf(arr, 0, total-1, x); idx >= 0 {
		var res []int64
		for i, v := range arr {
			if i == idx || i == total-1 {
				continue
			}
			res = append(res, v)
		}
		return joinInts(res)
	}
	y = arr[total-2]
	x = sum - 2*y
	if x == arr[total-1] {
		var res []int64
		for i, v := range arr {
			if i == total-2 || i == total-1 {
				continue
			}
			res = append(res, v)
		}
		return joinInts(res)
	}
	return "-1"
}

func joinInts(a []int64) string {
	if len(a) == 0 {
		return ""
	}
	var sb strings.Builder
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	return sb.String()
}

func indexOf(a []int64, lo, hi int, x int64) int {
	idx := sort.Search(hi-lo, func(i int) bool { return a[lo+i] >= x })
	if idx < hi-lo && a[lo+idx] == x {
		return lo + idx
	}
	return -1
}

func run(bin, input string) (string, error) {
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseTestcases()
	if err != nil {
		fmt.Println("failed to parse testcases:", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for i, v := range tc.arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
		expected := solve(tc.n, tc.arr)
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
