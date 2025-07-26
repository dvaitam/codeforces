package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func expectedC(n, m int, a []int64) []int64 {
	sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
	prefix := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		prefix[i] = prefix[i-1] + a[i-1]
	}
	dp := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		if i <= m {
			dp[i] = prefix[i]
		} else {
			dp[i] = prefix[i] + dp[i-m]
		}
	}
	return dp[1:]
}

func generateCaseC(rng *rand.Rand) (int, int, []int64) {
	n := rng.Intn(6) + 1
	m := rng.Intn(n) + 1
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		a[i] = int64(rng.Intn(11))
	}
	return n, m, a
}

func runCaseC(bin string, n, m int, a []int64) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteString("\n")
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(strings.TrimSpace(out.String()))
	if len(fields) != n {
		return fmt.Errorf("expected %d values got %d", n, len(fields))
	}
	expect := expectedC(n, m, append([]int64(nil), a...))
	for i := 0; i < n; i++ {
		var val int64
		if _, err := fmt.Sscan(fields[i], &val); err != nil {
			return fmt.Errorf("failed to parse output: %v", err)
		}
		if val != expect[i] {
			return fmt.Errorf("index %d expected %d got %d", i, expect[i], val)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, m, a := generateCaseC(rng)
		if err := runCaseC(bin, n, m, a); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
