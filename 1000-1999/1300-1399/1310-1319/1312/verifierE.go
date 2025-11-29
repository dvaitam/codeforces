package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded testcases from testcasesE.txt to remove external dependency.
const testcasesRaw = `4 5 2 7 8
3 2 2 1
7 9 5 1 4 9 9 6
5 3 2 5 4 1
5 5 4 3 5 5
6 2 10 6 7 9 4
3 4 8 5
2 9 5
1 5
10 5 9 4 7 7 10 5 7 8 3
4 5 5 1 2
1 8
5 9 9 8 6 3
4 2 7 4 8
5 3 6 7 10 6
9 4 6 2 1 4 5 10 10 4
2 6 3
5 8 1 1 6 2
5 6 1 6 5 6
3 7 10 2
5 10 4 8 5 3
5 7 10 3 6 10
1 6
1 8
3 6 6 5
10 2 8 4 7 4 2 1 1 1 3
10 3 10 1 9 8 10 4 6 1 2
9 5 7 4 8 4 4 8 7 8
1 4
7 8 4 7 4 8 4 1
1 5
5 4 9 4 4 7
5 3 6 1 6 10
2 10 7
1 8
7 2 7 4 10 3 6 5
8 6 7 9 4 5 6 7 8
2 5 4
1 7
10 3 5 1 3 8 10 8 7 7 4
1 4
3 1 10 5
2 7 7
4 9 1 4 3
10 6 9 8 9 8 1 2 1 10 2
8 9 5 10 3 1 6 2 9
1 5
6 2 2 9 8 7 4
5 7 4 8 7 2
2 2 10
6 9 7 7 8 2 4
5 8 7 2 9 3
6 3 3 3 6 8 6
5 9 1 3 1 5
2 9 2
8 10 8 9 2 9 4 7 5
6 4 3 1 1 10 6
9 8 10 5 9 8 10 10 8 7
3 5 10 6
6 3 7 2 10 3 10
3 5 6 4
10 6 10 2 2 7 3 6 6 6 3
5 1 10 1 9 2
6 2 3 3 10 8 10
2 2 3
8 4 10 5 7 10 4 8 4
5 6 4 6 9 9
8 7 9 7 6 5 8 7 10
1 5
3 9 8 9
10 6 7 7 10 1 3 9 2 8 6
10 5 2 5 8 4 8 10 2 3 4
2 5 3
1 3
7 10 10 9 9 5 1 9
4 3 2 4 1
6 2 2 5 1 5 10
3 3 7 3
2 9 6
1 9
3 9 7 3
4 5 8 9 2
7 3 3 5 9 7 9 5
7 6 3 7 1 7 1 5
1 5
3 2 3 2
10 1 4 4 9 1 8 9 3 8 7
6 3 9 10 4 2 8
10 6 5 7 10 1 4 5 10 10 5
8 6 2 4 6 2 1 10 6
9 8 6 2 3 1 8 9 9 10
4 1 4 2 6
10 3 5 2 9 9 1 2 4 10 5
10 1 1 1 8 3 4 6 4 6 5
7 10 7 1 3 7 8 2
9 9 5 9 5 4 7 6 3 4
1 8
7 5 6 7 7 6 10 4
5 5 9 2 1 7
5 10 5 3 10 6`

type testCase struct {
	n   int
	arr []int
}

// referenceSolution embeds the DP logic from 1312E.go.
func referenceSolution(arr []int) int {
	n := len(arr)
	min := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}

	can := make([][]int, n)
	dp := make([][]int, n)
	for i := 0; i < n; i++ {
		can[i] = make([]int, n)
		dp[i] = make([]int, n)
		can[i][i] = arr[i]
		dp[i][i] = 1
	}
	const inf = int(1e9)
	for length := 2; length <= n; length++ {
		for l := 0; l+length-1 < n; l++ {
			r := l + length - 1
			dp[l][r] = inf
			for k := l; k < r; k++ {
				v1 := can[l][k]
				if v1 == 0 {
					continue
				}
				v2 := can[k+1][r]
				if v2 == v1 {
					can[l][r] = v1 + 1
					break
				}
			}
			if can[l][r] != 0 {
				dp[l][r] = 1
				continue
			}
			for k := l; k < r; k++ {
				dp[l][r] = min(dp[l][r], dp[l][k]+dp[k+1][r])
			}
		}
	}
	return dp[0][n-1]
}

func parseTestcases() []testCase {
	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	scanner.Buffer(make([]byte, 1024), 1<<20)
	tests := make([]testCase, 0)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			panic("invalid n in testcase")
		}
		if len(fields) != n+1 {
			panic("testcase length mismatch")
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			val, err := strconv.Atoi(fields[i+1])
			if err != nil {
				panic("invalid integer in testcase")
			}
			arr[i] = val
		}
		tests = append(tests, testCase{n: n, arr: arr})
	}
	return tests
}

func runBin(bin, input string) (string, error) {
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
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests := parseTestcases()
	for idx, tc := range tests {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for i, v := range tc.arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		input := sb.String()

		expected := referenceSolution(tc.arr)
		gotStr, err := runBin(bin, input)
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := strconv.Atoi(gotStr)
		if err != nil {
			fmt.Printf("test %d: cannot parse output %q\n", idx+1, gotStr)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("test %d failed: expected %d got %d\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
