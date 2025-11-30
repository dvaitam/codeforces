package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesD.txt so the verifier is self-contained.
const testcasesD = `3 3
110
011
001

5 5
00010
00010
11111
10100
01011

4 5
11101
01010
01100
11010

5 3
011
100
011
100
100

2 2
01
11

3 2
11
10
11

5 5
01101
11110
11010
01111
11000

4 4
1110
1011
1110
0010

4 3
101
001
101
001

3 5
10000
10110
10111

5 4
1111
0100
0110
0111
0010

3 2
00
10
01

2 3
100
100

4 3
110
100
000
010

2 5
10101
10110

2 3
000
010

2 5
10101
11111

3 4
1000
1000
1110

5 2
11
11
11
00
01

2 5
01001
11010

5 2
10
01
00
10
00

5 3
010
001
110
000
001

3 2
00
01
11

4 5
11010
10110
10001
11001

3 3
101
110
101

5 5
00101
10110
10101
01000
01001

3 3
010
010
100

2 2
11
01

5 4
0011
1110
0110
1010
1011

2 4
1111
0101

4 2
10
00
11
00

3 4
1101
1100
0011

2 4
1100
0110

4 4
1101
1001
0000
0110

3 5
00001
11011
11111

5 2
00
00
10
10
00

2 5
01001
00010

2 5
01111
00001

3 2
01
01
10

2 5
00011
10101

4 2
11
11
01
01

4 4
1110
0011
0100
0001

3 2
00
01
10

2 2
11
11

2 2
01
01

5 2
00
01
10
01
01

2 2
11
11

5 3
101
110
010
000
000

2 5
00101
11101

4 5
01010
11101
10111
01111

3 5
10000
10010
00111

4 2
00
10
11
01

2 3
000
000

3 4
0011
0100
0000

3 4
0011
0011
1101

5 5
00011
01101
01100
01010
11010

4 3
001
110
111
101

2 3
000
101

2 3
011
101

5 3
100
111
100
011
110

5 4
1101
1101
1110
0000
0000

3 5
11010
00111
00111

3 5
11100
10000
11010

3 2
01
01
11

4 3
100
001
111
111

4 3
100
110
110
101

3 4
1010
1011
0110

5 5
10100
11110
01110
10000
00101

4 5
11100
01000
01100
10001

3 3
000
010
110

2 4
0000
1011

2 3
010
100

4 5
01101
11111
01101
01011

3 3
001
000
001

4 4
0111
0011
1111
1111

3 2
01
01
11

5 3
110
100
111
001
100

4 3
010
111
011
111

4 2
10
10
01
00

3 2
01
11
00

2 4
1101
1010

5 3
110
000
101
101
011

3 5
01111
00011
00111

2 4
1010
1000

3 2
00
11
11

5 3
011
011
101
011
010

3 2
10
10
10

2 4
1010
0111

3 3
111
101
000

4 5
00010
01001
10111
00110

3 4
1101
1111
0011

3 4
0101
1011
1010

2 4
1111
0111

5 4
0010
1010
1101
1011
0010

3 3
000
111
101

4 4
0100
0111
0110
1000

4 4
1000
0011
0100
0110

4 5
11111
11000
01000
01101

3 2
00
01
11

3 2
10
10
00`

// Embedded solver from 435D.go.
const MOD = 1000000007

func add(a, b int64) int64 {
	a += b
	if a >= MOD {
		a -= MOD
	}
	return a
}

func solve(grid []string) int64 {
	n := len(grid)
	m := len(grid[0])

	black := make([][]int, n+2)
	for i := range black {
		black[i] = make([]int, m+2)
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			if grid[i-1][j-1] == '1' {
				black[i][j] = 1
			}
		}
	}

	rowPS := make([][]int, n+2)
	colPS := make([][]int, n+2)
	for i := range rowPS {
		rowPS[i] = make([]int, m+2)
		colPS[i] = make([]int, m+2)
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			rowPS[i][j] = rowPS[i][j-1] + black[i][j]
			colPS[i][j] = colPS[i-1][j] + black[i][j]
		}
	}

	diag1 := make([][]int, n+2)
	diag2 := make([][]int, n+2)
	for i := range diag1 {
		diag1[i] = make([]int, m+3)
		diag2[i] = make([]int, m+3)
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			diag1[i][j] = black[i][j] + diag1[i-1][j-1]
		}
		for j := m; j >= 1; j-- {
			diag2[i][j] = black[i][j] + diag2[i-1][j+1]
		}
	}

	hRun := make([][]int, n+2)
	hRunL := make([][]int, n+2)
	vRun := make([][]int, n+2)
	vRunU := make([][]int, n+2)
	for i := range hRun {
		hRun[i] = make([]int, m+2)
		hRunL[i] = make([]int, m+2)
		vRun[i] = make([]int, m+2)
		vRunU[i] = make([]int, m+2)
	}
	for i := n; i >= 1; i-- {
		for j := m; j >= 1; j-- {
			if grid[i-1][j-1] == '0' {
				hRun[i][j] = hRun[i][j+1] + 1
				vRun[i][j] = vRun[i+1][j] + 1
			}
		}
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			if grid[i-1][j-1] == '0' {
				hRunL[i][j] = hRunL[i][j-1] + 1
				vRunU[i][j] = vRunU[i-1][j] + 1
			}
		}
	}

	var ans int64

	// Type1
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			if grid[i-1][j-1] != '0' {
				continue
			}
			maxK := hRun[i][j]
			if vRun[i][j] < maxK {
				maxK = vRun[i][j]
			}
			for k := 1; k < maxK; k++ {
				if grid[i-1][j-1+k] != '0' || grid[i-1+k][j-1] != '0' {
					continue
				}
				if rowPS[i][j+k-1]-rowPS[i][j] != 0 {
					continue
				}
				if colPS[i+k-1][j]-colPS[i][j] != 0 {
					continue
				}
				if diag2[i+k-1][j+1]-diag2[i][j+k] != 0 {
					continue
				}
				ans++
			}

			maxK = hRun[i][j]
			if vRunU[i][j] < maxK {
				maxK = vRunU[i][j]
			}
			for k := 1; k < maxK; k++ {
				if grid[i-1][j-1+k] != '0' || grid[i-1-k][j-1] != '0' {
					continue
				}
				if rowPS[i][j+k-1]-rowPS[i][j] != 0 {
					continue
				}
				if colPS[i-1][j]-colPS[i-k][j] != 0 {
					continue
				}
				if diag1[i-1][j+k-1]-diag1[i-k][j] != 0 {
					continue
				}
				ans++
			}

			maxK = hRunL[i][j]
			if vRun[i][j] < maxK {
				maxK = vRun[i][j]
			}
			for k := 1; k < maxK; k++ {
				if grid[i-1][j-1-k] != '0' || grid[i-1+k][j-1] != '0' {
					continue
				}
				if rowPS[i][j-1]-rowPS[i][j-k] != 0 {
					continue
				}
				if colPS[i+k-1][j]-colPS[i][j] != 0 {
					continue
				}
				if diag1[i+k-1][j-1]-diag1[i][j-k] != 0 {
					continue
				}
				ans++
			}

			maxK = hRunL[i][j]
			if vRunU[i][j] < maxK {
				maxK = vRunU[i][j]
			}
			for k := 1; k < maxK; k++ {
				if grid[i-1][j-1-k] != '0' || grid[i-1-k][j-1] != '0' {
					continue
				}
				if rowPS[i][j-1]-rowPS[i][j-k] != 0 {
					continue
				}
				if colPS[i-1][j]-colPS[i-k][j] != 0 {
					continue
				}
				if diag2[i-1][j-k+1]-diag2[i-k][j] != 0 {
					continue
				}
				ans++
			}
		}
	}

	// Type2
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			if grid[i-1][j-1] != '0' {
				continue
			}
			maxL := hRun[i][j]
			for k := 2; k < maxL; k += 2 {
				mid := k / 2
				bj := j + k
				if bj > m || grid[i-1][bj-1] != '0' {
					continue
				}
				if rowPS[i][bj-1]-rowPS[i][j] != 0 {
					continue
				}
				ci := i - mid
				cj := j + mid
				if ci >= 1 && grid[ci-1][cj-1] == '0' {
					if diag2[i-1][j+1]-diag2[ci][cj] == 0 &&
						diag1[i-1][bj-1]-diag1[ci][cj] == 0 {
						ans++
					}
				}
				ci = i + mid
				if ci <= n && grid[ci-1][cj-1] == '0' {
					if diag1[ci-1][cj-1]-diag1[i][j] == 0 &&
						diag2[ci-1][cj+1]-diag2[i][bj] == 0 {
						ans++
					}
				}
			}
		}
	}

	// Type3
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			if grid[i-1][j-1] != '0' {
				continue
			}
			maxL := vRun[i][j]
			for k := 2; k < maxL; k += 2 {
				mid := k / 2
				bi := i + k
				if bi > n || grid[bi-1][j-1] != '0' {
					continue
				}
				if colPS[bi-1][j]-colPS[i][j] != 0 {
					continue
				}
				ci := i + mid
				cj := j - mid
				if cj >= 1 && grid[ci-1][cj-1] == '0' {
					if diag2[ci-1][cj+1]-diag2[i][j] == 0 &&
						diag1[bi-1][j-1]-diag1[ci][cj] == 0 {
						ans++
					}
				}
				cj = j + mid
				if cj <= m && grid[ci-1][cj-1] == '0' {
					if diag1[ci-1][cj-1]-diag1[i][j] == 0 &&
						diag2[bi-1][j+1]-diag2[ci][cj] == 0 {
						ans++
					}
				}
			}
		}
	}
	return ans % MOD
}

type testCase struct {
	n    int
	m    int
	grid []string
}

func parseCases() ([]testCase, error) {
	blocks := strings.Split(strings.TrimSpace(testcasesD), "\n\n")
	cases := make([]testCase, 0, len(blocks))
	for idx, block := range blocks {
		lines := strings.Split(strings.TrimSpace(block), "\n")
		if len(lines) < 1 {
			continue
		}
		var n, m int
		if _, err := fmt.Sscan(lines[0], &n, &m); err != nil {
			return nil, fmt.Errorf("case %d: bad n m", idx+1)
		}
		if len(lines)-1 != n {
			return nil, fmt.Errorf("case %d: expected %d rows got %d", idx+1, n, len(lines)-1)
		}
		grid := make([]string, n)
		for i := 0; i < n; i++ {
			grid[i] = strings.TrimSpace(lines[i+1])
		}
		cases = append(cases, testCase{n: n, m: m, grid: grid})
	}
	return cases, nil
}

func runCandidate(bin, input string) (string, error) {
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
		fmt.Println("usage: verifierD /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseCases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		expect := solve(tc.grid)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.m)
		for _, row := range tc.grid {
			sb.WriteString(row)
			sb.WriteByte('\n')
		}
		gotStr, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := strconv.ParseInt(strings.TrimSpace(gotStr), 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: bad output\n", idx+1)
			os.Exit(1)
		}
		if got%MOD != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
