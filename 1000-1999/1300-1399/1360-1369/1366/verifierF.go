package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Edge struct {
	u, v int
	w    int64
}

const INF int64 = math.MinInt64 / 4
const MOD int64 = 1_000_000_007

func solveCase(n, m, q int, edges []Edge) int64 {
	dp := make([][]int64, q+1)
	for i := range dp {
		dp[i] = make([]int64, n+1)
		for j := 1; j <= n; j++ {
			dp[i][j] = INF
		}
	}
	dp[0][1] = 0
	ans := int64(0)
	for t := 1; t <= q; t++ {
		for j := 1; j <= n; j++ {
			dp[t][j] = INF
		}
		for _, e := range edges {
			if dp[t-1][e.u] != INF {
				val := dp[t-1][e.u] + e.w
				if val > dp[t][e.v] {
					dp[t][e.v] = val
				}
			}
			if dp[t-1][e.v] != INF {
				val := dp[t-1][e.v] + e.w
				if val > dp[t][e.u] {
					dp[t][e.u] = val
				}
			}
		}
		best := INF
		for i := 1; i <= n; i++ {
			if dp[t][i] > best {
				best = dp[t][i]
			}
		}
		if best == INF {
			best = 0
		}
		ans = (ans + ((best%MOD)+MOD)%MOD) % MOD
	}
	return ans % MOD
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesF.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		if len(fields) < 3 {
			fmt.Printf("invalid test %d\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(fields[0])
		m, _ := strconv.Atoi(fields[1])
		q, _ := strconv.Atoi(fields[2])
		if len(fields) != 3+3*m {
			fmt.Printf("test %d wrong count\n", idx)
			os.Exit(1)
		}
		edges := make([]Edge, m)
		pos := 3
		for i := 0; i < m; i++ {
			u, _ := strconv.Atoi(fields[pos])
			v, _ := strconv.Atoi(fields[pos+1])
			w, _ := strconv.ParseInt(fields[pos+2], 10, 64)
			pos += 3
			edges[i] = Edge{u, v, w}
		}
		expected := solveCase(n, m, q, edges)

		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d %d %d\n", n, m, q))
		for i, e := range edges {
			input.WriteString(fmt.Sprintf("%d %d %d", e.u, e.v, e.w))
			if i < m-1 {
				input.WriteByte('\n')
			}
		}
		input.WriteByte('\n')

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input.String())
		var out bytes.Buffer
		var errBuf bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &errBuf
		if err := cmd.Run(); err != nil {
			fmt.Printf("Test %d: runtime error: %v\n%s", idx, err, errBuf.String())
			os.Exit(1)
		}
		outStr := strings.TrimSpace(out.String())
		if outStr != fmt.Sprintf("%d", expected) {
			fmt.Printf("Test %d failed: expected %d got %s\n", idx, expected, outStr)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
