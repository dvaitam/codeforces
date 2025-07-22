package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type testCaseD struct {
	n   int
	h   []int64
	ans int
}

func solveD(tc testCaseD) int {
	n := tc.n
	p := make([]int64, n+1)
	for i := 0; i < n; i++ {
		p[i+1] = p[i] + tc.h[i]
	}
	dp := make([]int, n+1)
	best := make([]int64, n+1)
	const inf int64 = 1<<63 - 1
	best[0] = 0
	for i := 1; i <= n; i++ {
		dp[i] = 0
		best[i] = inf
		for j := i - 1; j >= 0; j-- {
			s := p[i] - p[j]
			if s < best[j] {
				continue
			}
			seg := dp[j] + 1
			if seg > dp[i] || (seg == dp[i] && s < best[i]) {
				dp[i] = seg
				best[i] = s
			}
		}
	}
	return n - dp[n]
}

func genCaseD(rng *rand.Rand) testCaseD {
	n := rng.Intn(6) + 1
	h := make([]int64, n)
	for i := range h {
		h[i] = int64(rng.Intn(10) + 1)
	}
	tc := testCaseD{n: n, h: h}
	tc.ans = solveD(tc)
	return tc
}

func runCaseD(bin string, tc testCaseD) error {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", tc.n)
	for i, v := range tc.h {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(strings.TrimSpace(out.String()))
	if len(fields) != 1 {
		return fmt.Errorf("expected 1 number got %d", len(fields))
	}
	val, err := strconv.Atoi(fields[0])
	if err != nil {
		return fmt.Errorf("invalid integer: %v", err)
	}
	if val != tc.ans {
		return fmt.Errorf("expected %d got %d", tc.ans, val)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := genCaseD(rng)
		if err := runCaseD(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
