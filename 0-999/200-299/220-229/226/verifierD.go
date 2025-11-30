package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func solve(input string) (string, error) {
	reader := bufio.NewReader(strings.NewReader(strings.TrimSpace(input)))
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return "", err
	}
	a := make([][]int, n)
	for i := 0; i < n; i++ {
		a[i] = make([]int, m)
		for j := 0; j < m; j++ {
			if _, err := fmt.Fscan(reader, &a[i][j]); err != nil {
				return "", err
			}
		}
	}
	r := make([]int, n)
	c := make([]int, m)
	rowSum := make([]int, n)
	colSum := make([]int, m)
	for i := 0; i < n; i++ {
		sum := 0
		for j := 0; j < m; j++ {
			sum += a[i][j]
		}
		rowSum[i] = sum
	}
	for j := 0; j < m; j++ {
		sum := 0
		for i := 0; i < n; i++ {
			sum += a[i][j]
		}
		colSum[j] = sum
	}
	changed := true
	for changed {
		changed = false
		for i := 0; i < n; i++ {
			if rowSum[i] < 0 {
				changed = true
				r[i] ^= 1
				rowSum[i] = -rowSum[i]
				for j := 0; j < m; j++ {
					oldParity := (r[i]^1 + c[j]) & 1
					var oldVal int
					if oldParity == 1 {
						oldVal = -a[i][j]
					} else {
						oldVal = a[i][j]
					}
					colSum[j] += -2 * oldVal
				}
			}
		}
		for j := 0; j < m; j++ {
			if colSum[j] < 0 {
				changed = true
				c[j] ^= 1
				colSum[j] = -colSum[j]
				for i := 0; i < n; i++ {
					oldParity := (r[i] + (c[j]^1)) & 1
					var oldVal int
					if oldParity == 1 {
						oldVal = -a[i][j]
					} else {
						oldVal = a[i][j]
					}
					rowSum[i] += -2 * oldVal
				}
			}
		}
	}
	var rows []int
	for i := 0; i < n; i++ {
		if r[i] == 1 {
			rows = append(rows, i+1)
		}
	}
	var cols []int
	for j := 0; j < m; j++ {
		if c[j] == 1 {
			cols = append(cols, j+1)
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprint(len(rows)))
	for _, v := range rows {
		sb.WriteString(" ")
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	sb.WriteString(fmt.Sprint(len(cols)))
	for _, v := range cols {
		sb.WriteString(" ")
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	return sb.String(), nil
}

var testcases = []string{
	`4 4
-5 -1 3 2
1 -1 2 0
4 -2 3 -3
-1 -3 -4 4`,
	`3 5
4 -3 -1 -4 -4
5 0 2 3 -4
0 1 0 4 5`,
	`2 5
2 2 3 -1 -5
3 -5 -4 1 5`,
	`1 5
2 0 -2 0 -4`,
	`2 5
-2 -2 -3 3 2
-4 -4 0 3 2`,
	`1 3
3 -1 -4`,
	`5 3
3 -2 4
3 4 -1
2 -4 4
1 0 4
-2 -1 -3`,
	`2 2
-5 4
5 -1`,
	`4 1
-4
5
-3
-3`,
	`1 1
3`,
	`4 5
-1 3 -2 -2 5
4 1 4 -1 2
2 5 5 0 -4
0 4 -4 2 4`,
	`3 2
-2 -5
-1 -4
-2 0`,
	`2 3
1 -5 -4
-3 -2 -5`,
	`5 5
4 5 -4 -5 -4
5 -2 4 4 -4
1 -4 0 -4 -5
4 -5 -2 -3 -4
2 -2 -5 5 -5`,
	`5 4
4 -4 -1 -4
-2 -4 5 -1
0 1 -3 -5
3 2 -5 4
-4 1 -2 -1`,
	`3 4
4 -3 5 -2
-5 5 -3 -3
0 3 -1 -4`,
	`5 4
5 -3 -5 2
5 1 4 3
-1 5 0 1
5 -1 -3 3
-5 2 -4 0`,
	`1 5
-1 -3 -2 2 0`,
	`5 3
5 0 4
5 4 -3
-1 1 1
5 -4 -5
4 -2 0`,
	`2 2
-2 5
2 1`,
	`5 4
-5 1 4 1
5 -5 -3 2
-4 -1 -3 2
3 2 3 4
-5 -5 2 0`,
	`3 4
-5 1 -2 3
5 -4 -3 -5
1 5 1 0`,
	`1 2
-5 -5`,
	`5 5
-4 -2 -4 4 5
-2 -1 -1 -3 -4
2 1 5 -4 -5
-1 2 -4 -1 -3
5 3 5 5 0`,
	`1 2
-1 -5`,
	`1 1
-2`,
	`3 5
0 0 4 -5 4
5 2 5 2 5
1 0 3 -3 -2`,
	`4 5
-1 -5 -3 -3 -1
0 0 0 -4 0
4 -5 -5 -1 -3
-3 4 -1 0 1`,
	`5 2
-1 -4
2 -2
-5 -1
-3 3
-4 -1`,
	`4 3
-1 1 -4
-4 3 2
2 0 0
-4 2 -4`,
	`4 4
-5 -1 0 5
-3 -3 5 4
1 5 -4 -4
-4 -2 -2 -5`,
	`4 1
-4
1
3
3`,
	`3 4
2 4 5 -2
1 -4 0 -2
-1 4 -3 1`,
	`2 3
-4 -4 -5
3 2 5`,
	`2 1
2
1`,
	`3 2
5 -5
-2 4
-3 -4`,
	`2 4
1 0 3 -3
-4 4 2 -3`,
	`5 4
5 5 1 3
2 5 0 2
2 5 5 -2
3 4 -2 -5
0 0 0 -5`,
	`5 2
-1 4
-3 1
4 -1
2 -4
-4 3`,
	`1 1
-2`,
	`2 1
-1
-5`,
	`4 3
-3 -3 5
2 0 3
1 3 3
-5 4 -4`,
	`5 5
-4 1 -2 -1 3
4 1 2 1 4
4 -2 -5 5 -5
-3 -1 3 4 -1
0 -4 2 -1 -1`,
	`4 4
1 -5 -3 5
-3 -2 -1 0
-5 -5 2 1
-3 2 4 -4`,
	`2 3
1 -5 4
2 1 2`,
	`1 1
2`,
	`2 1
-5
4`,
	`5 2
5 0
-4 3
5 0
-2 1
2 -4`,
	`1 5
2 4 5 0 5`,
	`1 5
-1 -3 1 -1 5`,
	`1 5
-2 -5 1 2 0`,
	`2 4
0 5 -4 -5
-5 2 -1 -5`,
	`5 5
4 -2 -2 -4 5
3 3 1 3 -1
-4 -3 1 4 1
-4 -4 1 -4 -4
1 -3 -5 2 1`,
	`4 1
2
0
-1
-4`,
	`3 1
-4
0
-5`,
	`3 3
-3 -5 -2
0 -4 4
-3 -2 -5`,
	`2 1
-5
-1`,
	`3 1
4
-2
-3`,
	`2 4
-4 2 0 -1
-3 -5 -2 0`,
	`3 4
-1 -1 3 5
0 -3 4 -4
-4 3 4 -1`,
	`2 4
-3 -3 -2 0
3 -2 -2 -3`,
	`3 3
1 5 -5
-3 4 -5
1 -4 -4`,
	`2 4
-1 3 1 -3
4 1 -1 5`,
	`3 1
-2
2
5`,
	`3 5
-5 1 1 -5 1
0 2 -2 0 -1
2 -4 -3 -4 -1`,
	`1 5
4 -3 2 1 -3`,
	`4 4
-3 -2 2 0
3 -3 0 2
5 5 -4 2
-2 -1 -5 2`,
	`5 4
-5 -2 -1 -4
5 -1 3 4
-3 1 2 -4
5 2 -2 3
1 -1 5 -5`,
	`1 3
5 -5 -5`,
	`3 4
3 4 1 2
-4 -1 0 -1
5 -2 4 -4`,
	`1 1
-1`,
	`3 5
0 -4 3 -2 -3
-4 1 -1 -1 3
-3 4 3 5 -2`,
	`5 1
1
5
3
1
-1`,
	`3 4
0 4 5 -3
-3 -4 -4 1
1 4 2 -3`,
	`5 3
0 5 2
1 -2 2
2 3 0
2 5 -5
2 -1 -3`,
	`4 1
4
-2
-5
0`,
	`4 4
-5 3 -4 5
-4 5 5 1
-5 0 -5 -4
4 -5 -1 5`,
	`3 2
-3 4
-1 -2
-4 1`,
	`4 3
1 -3 0
1 5 5
1 -3 2
-3 3 0`,
	`2 2
-3 2
0 1`,
	`4 4
1 -2 -2 2
-2 4 -5 1
-5 -2 5 -4
-3 0 -5 5`,
	`2 2
4 -1
4 -4`,
	`5 3
0 1 2
-5 5 3
5 5 3
1 4 2
2 -1 2`,
	`2 3
-1 -5 -5
-5 -3 0`,
	`1 3
5 -5 -3`,
	`1 4
5 -2 4 1`,
	`5 2
2 -2
0 4
-4 4
-4 0
0 3`,
	`4 3
-1 -5 3
-5 -2 0
-4 -2 3
0 -2 -2`,
	`3 3
-1 3 1
-1 2 0
-2 -5 -1`,
	`5 1
-5
2
2
2
-5`,
	`4 4
2 2 -4 -4
-4 -2 -4 -3
1 -2 2 4
-4 1 3 1`,
	`1 2
-2 2`,
	`2 2
-1 0
0 1`,
	`1 5
-1 4 3 -2 -1`,
	`4 5
4 2 3 5 -1
-1 -2 -5 -4 4
-4 -3 1 -2 -2
-1 5 -5 3 3`,
	`4 1
-4
1
5
-1`,
	`1 5
0 -2 5 3 5`,
	`3 2
-2 -4
3 -1
5 0`,
	`2 3
5 2 -1
4 -3 -3`,
	`1 5
3 0 0 4 5`,
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	for idx, tc := range testcases {
		input := strings.TrimSpace(tc) + "\n"

		expected, err := solve(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: oracle error: %v\n", idx+1, err)
			os.Exit(1)
		}

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", idx+1, err)
			os.Exit(1)
		}

		got := strings.TrimSpace(string(out))
		if got != strings.TrimSpace(expected) {
			fmt.Fprintf(os.Stderr, "test %d failed\nexpected: %s\n got: %s\n", idx+1, strings.TrimSpace(expected), got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(testcases))
}
