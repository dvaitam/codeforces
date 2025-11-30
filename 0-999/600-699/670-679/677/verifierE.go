package main

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testcase struct {
	n    int
	rows []string
}

// Embedded testcases from testcasesE.txt.
const testcasesRaw = `2 20 33
2 00 03
3 012 210 210
3 211 222 023
2 11 32
1 2
1 2
3 133 233 112
3 000 323 211
1 3
2 32 12
4 2120 0121 0212 3002
1 2
3 022 213 021
4 2123 1202 0312 2203
2 31 00
1 0
2 10 31
3 002 313 113
4 3013 3131 3100 2211
2 32 12
1 2
1 3
1 3
4 0311 2232 3122 3302
2 03 12
1 1
4 3331 0110 2033 1011
3 330 000 321
1 2
1 0
3 200 331 231
4 3000 2333 0123 3012
2 11 23
3 201 020 033
1 1
4 2211 0023 2333 1222
2 30 11
3 212 003 122
3 120 002 011
4 0013 1231 3122 1233
4 2233 0213 2330 1032
3 023 130 110
3 101 320 110
2 02 00
3 021 131 020
2 31 12
4 0311 2323 2130 3020
3 101 001 103
2 33 21
2 03 22
4 0122 3201 2002 3201
1 3
2 01 02
2 20 00
2 20 00
4 1121 2233 0133 0221
4 2110 3322 3321 2200
4 2212 0101 2130 2031
2 01 00
3 210 101 120
3 221 110 010
3 301 133 322
4 3001 1212 1212 3201
2 21 10
1 0
2 22 30
3 311 332 203
2 30 31
1 2
3 003 101 120
1 3
4 1101 1213 2212 2220
2 21 33
1 0
2 11 10
3 001 120 232
2 32 01
3 300 220 130
4 0211 1100 0023 0021
2 00 01
4 1022 2100 0102 3221
3 313 130 123
3 121 222 132
3 212 012 113
1 2
1 0
1 0
2 11 23
4 1332 0200 1030 3130
2 22 21
1 2
3 031 323 031
1 3
4 1002 0113 3232 3230
3 303 321 132
1 2
1 1
2 31 23
4 3030 3232 1122 1231`

func parseTestcases() ([]testcase, error) {
	lines := strings.Split(testcasesRaw, "\n")
	var res []testcase
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 1 {
			return nil, fmt.Errorf("line %d malformed", idx+1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid n", idx+1)
		}
		if len(fields) != n+1 {
			return nil, fmt.Errorf("line %d: expected %d rows got %d", idx+1, n, len(fields)-1)
		}
		rows := make([]string, n)
		copy(rows, fields[1:])
		res = append(res, testcase{n: n, rows: rows})
	}
	return res, nil
}

const mod int64 = 1000000007

// Embedded solver logic from 677E.go.
func powmod(base int64, exp int) int64 {
	result := int64(1)
	b := base % mod
	e := exp
	for e > 0 {
		if e&1 == 1 {
			result = result * b % mod
		}
		b = b * b % mod
		e >>= 1
	}
	return result
}

func min4(a, b, c, d int) int {
	m := a
	if b < m {
		m = b
	}
	if c < m {
		m = c
	}
	if d < m {
		m = d
	}
	return m
}

func solve(n int, rows []string) int64 {
	grid := make([][]int, n)
	for i := 0; i < n; i++ {
		row := make([]int, n)
		for j := 0; j < n; j++ {
			row[j] = int(rows[i][j] - '0')
		}
		grid[i] = row
	}
	row2 := make([][]int, n)
	row3 := make([][]int, n)
	for i := 0; i < n; i++ {
		row2[i] = make([]int, n+1)
		row3[i] = make([]int, n+1)
		for j := 0; j < n; j++ {
			val := grid[i][j]
			row2[i][j+1] = row2[i][j]
			row3[i][j+1] = row3[i][j]
			if val == 2 {
				row2[i][j+1]++
			} else if val == 3 {
				row3[i][j+1]++
			}
		}
	}
	col2 := make([][]int, n+1)
	col3 := make([][]int, n+1)
	for i := 0; i <= n; i++ {
		col2[i] = make([]int, n)
		col3[i] = make([]int, n)
	}
	for j := 0; j < n; j++ {
		for i := 0; i < n; i++ {
			val := grid[i][j]
			col2[i+1][j] = col2[i][j]
			col3[i+1][j] = col3[i][j]
			if val == 2 {
				col2[i+1][j]++
			} else if val == 3 {
				col3[i+1][j]++
			}
		}
	}
	diag12 := make([][]int, n+1)
	diag13 := make([][]int, n+1)
	diag22 := make([][]int, n+1)
	diag23 := make([][]int, n+1)
	for i := 0; i <= n; i++ {
		diag12[i] = make([]int, n+1)
		diag13[i] = make([]int, n+1)
		diag22[i] = make([]int, n+1)
		diag23[i] = make([]int, n+1)
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			val := grid[i][j]
			diag12[i+1][j+1] = diag12[i][j]
			diag13[i+1][j+1] = diag13[i][j]
			diag22[i+1][j] = diag22[i][j+1]
			diag23[i+1][j] = diag23[i][j+1]
			if val == 2 {
				diag12[i+1][j+1]++
				diag22[i+1][j]++
			} else if val == 3 {
				diag13[i+1][j+1]++
				diag23[i+1][j]++
			}
		}
	}
	left := make([][]int, n)
	right := make([][]int, n)
	up := make([][]int, n)
	down := make([][]int, n)
	ul := make([][]int, n)
	ur := make([][]int, n)
	dl := make([][]int, n)
	dr := make([][]int, n)
	for i := 0; i < n; i++ {
		left[i] = make([]int, n)
		right[i] = make([]int, n)
		up[i] = make([]int, n)
		down[i] = make([]int, n)
		ul[i] = make([]int, n)
		ur[i] = make([]int, n)
		dl[i] = make([]int, n)
		dr[i] = make([]int, n)
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if grid[i][j] != 0 {
				left[i][j] = 1
				up[i][j] = 1
				ul[i][j] = 1
				ur[i][j] = 1
				if j > 0 {
					left[i][j] += left[i][j-1]
				}
				if i > 0 {
					up[i][j] += up[i-1][j]
				}
				if i > 0 && j > 0 {
					ul[i][j] += ul[i-1][j-1]
				}
				if i > 0 && j+1 < n {
					ur[i][j] += ur[i-1][j+1]
				}
			}
		}
	}
	for i := n - 1; i >= 0; i-- {
		for j := n - 1; j >= 0; j-- {
			if grid[i][j] != 0 {
				right[i][j] = 1
				down[i][j] = 1
				dr[i][j] = 1
				dl[i][j] = 1
				if j+1 < n {
					right[i][j] += right[i][j+1]
				}
				if i+1 < n {
					down[i][j] += down[i+1][j]
				}
				if i+1 < n && j+1 < n {
					dr[i][j] += dr[i+1][j+1]
				}
				if i+1 < n && j > 0 {
					dl[i][j] += dl[i+1][j-1]
				}
			}
		}
	}

	ln2 := math.Log(2)
	ln3 := math.Log(3)
	bestLog := -1.0
	bestA, bestB := 0, 0
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if grid[i][j] == 0 {
				continue
			}
			center2 := 0
			center3 := 0
			if grid[i][j] == 2 {
				center2 = 1
			} else if grid[i][j] == 3 {
				center3 = 1
			}
			r := min4(left[i][j], right[i][j], up[i][j], down[i][j]) - 1
			if r >= 0 {
				a := row2[i][j+r+1] - row2[i][j-r]
				a += col2[i+r+1][j] - col2[i-r][j]
				a -= center2
				b := row3[i][j+r+1] - row3[i][j-r]
				b += col3[i+r+1][j] - col3[i-r][j]
				b -= center3
				logv := float64(a)*ln2 + float64(b)*ln3
				if logv > bestLog {
					bestLog = logv
					bestA = a
					bestB = b
				}
			}
			r = min4(ul[i][j], ur[i][j], dl[i][j], dr[i][j]) - 1
			if r >= 0 {
				a := diag12[i+r+1][j+r+1] - diag12[i-r][j-r]
				a += diag22[i+r+1][j-r] - diag22[i-r][j+r+1]
				a -= center2
				b := diag13[i+r+1][j+r+1] - diag13[i-r][j-r]
				b += diag23[i+r+1][j-r] - diag23[i-r][j+r+1]
				b -= center3
				logv := float64(a)*ln2 + float64(b)*ln3
				if logv > bestLog {
					bestLog = logv
					bestA = a
					bestB = b
				}
			}
		}
	}
	if bestLog < 0 {
		return 0
	}
	ans := powmod(2, bestA) * powmod(3, bestB) % mod
	return ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse embedded testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		expect := solve(tc.n, tc.rows)

		var sb strings.Builder
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte('\n')
		for _, row := range tc.rows {
			sb.WriteString(row)
			sb.WriteByte('\n')
		}

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(sb.String())
		var out bytes.Buffer
		var errb bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &errb
		if err := cmd.Run(); err != nil {
			fmt.Printf("case %d: runtime error: %v\nstderr: %s\n", idx+1, err, errb.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != strconv.FormatInt(expect, 10) {
			fmt.Printf("case %d failed\nexpected: %d\n got: %s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
