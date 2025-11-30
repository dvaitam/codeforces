package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// solution logic from 1927G.go
func solveCase(a []int) string {
	n := len(a)
	type interval struct{ l, r int }
	intervals := make([]interval, 0, 2*n)
	for i := 1; i <= n; i++ {
		l1 := i - a[i-1] + 1
		if l1 < 1 {
			l1 = 1
		}
		intervals = append(intervals, interval{l1, i})
		r2 := i + a[i-1] - 1
		if r2 > n {
			r2 = n
		}
		intervals = append(intervals, interval{i, r2})
	}
	const INF = int(1e9)
	dp := make([]int, n+2)
	for i := 0; i <= n; i++ {
		dp[i] = INF
	}
	dp[n+1] = 0
	for i := n; i >= 1; i-- {
		best := INF
		for _, seg := range intervals {
			if seg.l <= i && seg.r >= i {
				if val := dp[seg.r+1] + 1; val < best {
					best = val
				}
			}
		}
		dp[i] = best
	}
	return fmt.Sprintf("%d", dp[1])
}

const testcasesData = `
13 10 2 8 13 5 1 1 3 11 10 8 13 12
6 3 1 3 4 2 6
14 7 9 9 11 2 4 10 9 12 13 12 5 11 13
10 2 7 6 2 6 7 5 8 2 4
12 11 5 2 1 10 4 11 6 8 4 9 10
11 9 1 11 6 4 10 7 5 6 10 2
2 1 1
10 5 5 4 7 8 4 3 10 4 9
1 1
13 3 1 11 6 9 11 10 10 5 6 7 9 7
5 2 4 1 2 4
10 7 2 8 4 2 10 8 8 7 2
9 7 8 5 7 2 4 5 8 8
12 3 1 1 9 2 5 10 6 4 5 9 8
6 5 3 4 4 5 4
5 5 4 4 4 2
12 7 8 5 11 8 6 6 11 3 10 7 12
10 5 6 7 8 3 5 9 1 10 8
1 1
1 1
11 6 6 8 10 1 4 3 5 10 1 7
13 13 13 9 8 2 8 4 2 6 6 3 11 11
4 3 2 1 1
12 2 1 11 8 8 12 2 11 1 12 9 3
10 2 3 10 5 8 8 1 3 4 8
8 8 8 5 8 2 5 6 5
6 1 1 5 3 3 3
8 7 8 6 1 2 4 4 1
12 3 6 6 1 1 12 4 12 7 6 2 3
12 11 2 7 9 3 8 12 8 11 2 3 2
6 5 1 2 5 6 2
3 2 3 1
9 7 4 4 3 7 3 6 4 9
10 1 10 7 4 6 7 9 8 1 7
14 9 9 14 11 3 12 4 11 9 13 12 6 5 11
3 1 1 3
11 3 7 2 1 1 1 8 1 2 11 9
2 2 1
1 1
13 7 4 13 10 12 5 2 8 2 12 10 1 3
15 11 8 2 2 12 3 1 11 15 1 13 14 11 12 7
11 5 5 8 4 4 11 9 11 4 5 10
1 1
9 6 9 2 1 6 6 2 4 1
6 6 3 4 6 1 3
3 1 1 2
2 2 1
8 5 2 5 4 3 3 7 3
4 4 4 4 3
5 2 3 1 2 4
2 1 1
15 13 1 9 3 2 3 10 7 8 8 14 14 9 13 12
9 8 4 7 3 8 1 6 7 5
6 2 3 3 5 6 3
4 1 3 2 1
7 5 5 7 2 6 1 2
7 4 2 7 7 2 5 1
9 9 8 6 6 7 5 2 5 3
1 1
10 10 10 3 9 9 2 2 2 2 4
6 3 2 6 2 2 2
4 2 4 4 1
12 10 3 8 4 10 6 12 3 6 11 5 2
5 1 1 1 4 4
6 1 1 3 6 1 4
15 10 6 2 12 14 14 6 3 14 13 3 6 11 13 8
12 7 1 9 3 2 8 6 1 3 11 9 8
5 1 4 1 4 2
13 9 3 8 6 5 13 9 13 6 5 9 12 5
8 1 4 1 8 5 3 1 6
1 1
10 5 4 3 5 1 5 5 6 4 3
10 1 2 6 10 3 6 3 8 4 8
7 5 5 3 1 6 6 1
7 1 5 5 2 5 4 3
4 1 4 3 3
11 6 5 7 6 2 9 7 11 3 5 1
11 9 6 11 9 1 2 11 2 1 8 7
9 2 3 5 6 1 2 5 5 6
7 4 5 6 7 7 6 1
12 2 5 8 8 12 5 5 4 1 1 9 9
14 7 3 4 13 5 11 11 14 12 1 2 5 13 3
9 2 1 7 5 8 4 2 5 3
15 3 3 2 14 1 13 4 10 6 15 7 15 12 8 10
10 8 8 6 10 8 4 3 9 2 4
11 3 7 4 10 6 8 10 11 5 10 10
14 4 10 4 12 1 7 4 7 10 12 2 11 1 6
6 4 6 3 4 4 3
6 5 6 5 3 5 1
7 6 5 2 3 3 5 6
8 4 2 8 3 6 7 8 6
2 2 2
8 1 2 2 8 7 1 1 6
6 2 2 4 5 1 4
10 10 2 1 6 4 10 9 5 10 9
15 13 15 5 14 13 7 6 12 1 10 6 7 5 3 8
14 13 6 1 1 9 13 4 1 13 1 14 1 8 11
2 2 2
6 4 2 2 3 2 1
4 1 2 3 4
`

func runProg(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	scanner := bufio.NewScanner(strings.NewReader(testcasesData))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		if len(fields) < 1 {
			fmt.Printf("test %d invalid\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(fields[0])
		if len(fields) != 1+n {
			fmt.Printf("test %d wrong count\n", idx)
			os.Exit(1)
		}
		nums := make([]int, n)
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fields[1+i])
			nums[i], _ = strconv.Atoi(fields[1+i])
		}
		sb.WriteByte('\n')
		input := sb.String()
		expect := solveCase(nums)
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Printf("test %d: %v\n", idx, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("test %d failed\nexpected: %s\ngot: %s\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
