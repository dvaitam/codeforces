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

const testcasesERaw = `100
5
-5 -3 5 -1 -5
-2 2 -2 -1 -1
4
4 -5 4 4
-3 5 1 2
7
-2 1 1 -3 3 4 -2
5 -3 0 2 4 1 -3
9
0 3 -5 -5 0 3 -3 5 -3
-5 -1 0 4 0 -2 -5 2 2
8
0 0 5 -4 -3 -2 4 -2
-1 3 -4 1 2 -4 -1 1
9
-5 0 5 0 4 3 0 2 3
-2 -4 3 -2 0 3 3 -2 1
5
4 0 1 -5 -2
-1 4 2 3 -3
8
-2 1 -2 2 1 -5 4 5
-2 5 3 -5 3 4 4 5
8
-3 -3 0 -5 3 -5 1 -1
4 -2 2 4 0 1 -2 -4
4
-5 -3 -2 -1
3 2 3 0
6
-5 3 0 4 0 -2
-5 5 -3 3 4 2
3
-5 -1 0
0 5 -5
3
-2 2 1
0 4 -2
4
3 5 3 2
-2 -4 1 2
4
3 4 0 -4
-2 -4 -5 -1
6
-3 3 1 4 5 0
-3 4 -2 1 -3 0
3
-2 2 -1
1 1 -5
8
2 4 -2 -2 0 1 -2 3
-2 -5 1 -2 2 1 2 5
4
0 -2 3 2
-1 4 -3 4
6
2 5 -1 1 -5 1
3 -3 -4 2 2 1
6
-2 3 2 5 5 1
-4 -5 -5 -3 -2 -2
6
1 5 -5 -2 -4 1
1 4 -5 1 3 4
3
4 1 -1
-3 -2 2
5
-2 2 1 -2 -4
-4 -1 5 5 -1
6
5 -1 -5 4 2 4
5 1 2 -1 -2 2
6
-5 0 -1 0 4 -4
-5 1 5 -4 -3 5
4
3 2 -3 1
0 -4 2 -3
4
2 3 1 -2
4 5 5 4
7
-1 0 2 5 4 3 -1
-1 2 -4 -5 -4 4 5
4
3 2 4 4
2 0 -4 5
8
-3 2 4 1 -1 2 -3 3
4 -4 -4 0 3 0 2 1
8
-5 -2 -1 0 0 -2 -1 -5
-3 1 4 -4 5 -5 1 -3
3
-4 5 0
0 1 -3
9
-1 -5 3 3 2 4 1 3 -4
2 1 -5 0 0 -4 5 -4 5
4
3 4 -4 -5
-3 -3 3 -1
9
-2 0 -1 4 0 -2 1 -3 -3
3 0 1 4 -5 5 -4 -1 -5
6
-1 -5 -1 2 -2 -4
-5 0 2 4 3 3
8
-2 0 -4 -2 -3 2 -3 0
5 -5 -1 -1 1 0 -1 1
7
-1 2 4 -5 3 -3 3
1 -4 3 1 -5 -1 2
3
2 0 -1
3 -5 0
5
-5 3 5 5 5
2 2 5 -1 -2
6
5 -3 5 -3 4 -4
1 -5 -3 -3 0 -3
4
-3 -4 1 0
4 0 4 -3
9
-1 -2 4 3 5 0 1 2 -2
-3 -2 -5 -2 -5 -4 -2 5 1
6
-3 3 -2 -1 -5 3
2 0 -2 3 -5 1
3
-5 1 3
3 -4 3
9
-2 -1 3 4 -4 -2 2 -4 3
1 -5 -3 -1 1 2 4 5 -1
3
-3 -2 -1
-1 3 -3
3
-4 -1 -4
-4 5 0
9
3 3 3 5 -4 3 0 -2 -1
-5 1 4 -1 -5 -5 -5 3 -5
7
-1 -3 -4 -3 -3 5 -4
2 1 -1 -1 5 -4 4
7
1 4 0 0 2 -2 3
4 -2 3 -3 -5 -2 3
5
-1 1 1 5 2
-2 5 0 -1 3
3
-3 1 3
-3 4 -2
5
2 -4 -2 -2 -4
5 3 -4 -2 -2
9
-3 -5 3 -3 2 -2 2 0 -2
3 -2 -4 0 3 4 1 -1 -3
9
2 -2 -4 -1 0 -2 3 3 -1
-2 4 1 1 5 -3 4 5 4
8
-3 5 -4 5 1 -4 -2 1
-2 5 -4 -3 -3 -1 5 4
7
4 -5 -2 4 3 -3 0
1 -1 5 -5 -4 5 -1
6
-2 3 -2 1 1 0
-2 -2 -3 3 -4 1
7
-5 4 -3 4 2 2 2
2 -3 3 -2 5 -2 0
3
5 0 -4
-1 5 -3
3
2 -1 3
-1 -4 -2
3
-1 0 -5
0 -1 -1
4
2 3 -2 5
1 2 -4 4
8
4 -4 -2 -2 2 -5 2 0
-2 2 5 -4 -2 5 -5 -2
6
5 2 -2 2 -1 3
-5 3 -3 -2 2 1
8
5 3 1 2 2 0 -1 0
-4 -3 -4 3 -3 -1 -1 5
7
-2 -1 -2 -5 -2 -4 -1
-1 1 3 -2 -4 -5 0
4
3 0 2 0
-5 -1 -1 -2
6
3 -4 1 0 -5 1
5 -1 4 -3 -3 5
3
3 3 4
-2 -5 2
4
4 2 -1 -2
5 -5 -3 -2
5
-2 -1 -4 5 -5
-2 0 2 -3 -3
4
0 5 -1 -1
0 2 5 -4
4
-2 -1 -5 -5
-5 -1 4 -4
4
3 5 1 2
4 -4 -4 -3
5
1 -4 -2 -2 -1
0 2 -3 4 -2
3
1 -4 -1
5 -2 2
8
-3 -4 0 -2 3 4 0 -4
-1 -2 -5 2 -5 2 -1 -3
5
-2 -4 -5 4 0
-1 -4 -5 -4 5
6
5 -4 5 -4 -1 1
2 -2 -1 1 1 -2
4
-2 4 -1 -1
0 5 -2 1
3
-5 -3 -2
0 -1 -3
7
-1 -3 -2 0 -3 0 1
1 -1 -5 2 4 0 -5
8
2 -4 0 -2 4 5 -2 5
4 0 -2 4 -1 4 -2 0
9
-5 3 1 -1 -4 0 4 -2 -3
3 -4 -3 0 3 -5 4 2 -4
3
-3 0 -1
-3 5 -3
6
-5 -5 -2 1 3 5
1 4 5 -2 4 2
8
0 -1 -3 1 2 -4 -3 4
-5 3 5 1 5 1 -4 -2
8
-2 5 2 -5 2 5 0 2
-1 2 -5 -4 1 4 0 2
7
5 -4 1 -2 4 1 5
-1 -5 2 5 5 2 5
7
-3 4 -4 -5 5 3 -1
-1 1 -2 -5 -5 4 -4
9
-4 -2 -5 3 5 3 5 -2 -3
1 4 -3 -4 5 -5 2 1 1
6
1 -4 0 -5 -4 0
-1 5 -5 1 1 -4
5
4 1 4 -1 -3
-3 -2 1 3 4
7
4 -4 -2 -1 -3 -5 -4
-4 -2 -1 -4 -2 3 0
5
-1 2 -3 1 4
-4 4 2 -3 -2
8
-1 -5 4 3 5 3 -4 0
3 0 -4 1 1 0 -1 -3
4
-1 4 3 4
3 0 4 2
`

func solveCase(a, b []int64) int {
	n := len(a)
	prefix1, prefix2, prefixAll := int64(0), int64(0), int64(0)
	last1 := map[int64]int{0: 0}
	last2 := map[int64]int{0: 0}
	lastAll := map[int64]int{0: 0}
	lastPair := map[[2]int64]int{{0, 0}: 0}
	dp := make([]int, n+1)
	for i := 1; i <= n; i++ {
		prefix1 += a[i-1]
		prefix2 += b[i-1]
		prefixAll += a[i-1] + b[i-1]
		best := dp[i-1]
		if j, ok := last1[prefix1]; ok && best < dp[j]+1 {
			best = dp[j] + 1
		}
		if j, ok := last2[prefix2]; ok && best < dp[j]+1 {
			best = dp[j] + 1
		}
		if j, ok := lastAll[prefixAll]; ok && best < dp[j]+1 {
			best = dp[j] + 1
		}
		if j, ok := lastPair[[2]int64{prefix1, prefix2}]; ok && best < dp[j]+2 {
			best = dp[j] + 2
		}
		dp[i] = best
		last1[prefix1] = i
		last2[prefix2] = i
		lastAll[prefixAll] = i
		lastPair[[2]int64{prefix1, prefix2}] = i
	}
	return dp[n]
}

func runCandidate(bin string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = bytes.NewReader(input)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("execution failed: %v", err)
	}
	return strings.TrimSpace(string(out)), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	scan := bufio.NewScanner(strings.NewReader(testcasesERaw))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		scan.Scan()
		n, _ := strconv.Atoi(scan.Text())
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			scan.Scan()
			v, _ := strconv.ParseInt(scan.Text(), 10, 64)
			a[i] = v
		}
		b := make([]int64, n)
		for i := 0; i < n; i++ {
			scan.Scan()
			v, _ := strconv.ParseInt(scan.Text(), 10, 64)
			b[i] = v
		}
		expected := solveCase(a, b)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i, v := range a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
		for i, v := range b {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
		out, err := runCandidate(os.Args[1], []byte(sb.String()))
		if err != nil {
			fmt.Printf("case %d failed: %v\n", caseIdx+1, err)
			os.Exit(1)
		}
		got, err := strconv.Atoi(strings.TrimSpace(out))
		if err != nil || got != expected {
			fmt.Printf("case %d failed: expected %d got %s\n", caseIdx+1, expected, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
