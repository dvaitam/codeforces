package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func abs64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func expectedF(mice []int64, holes []int64, cap []int) int64 {
	n := len(mice)
	m := len(holes)
	const INF int64 = math.MaxInt64 / 4
	dp := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		dp[i] = INF
	}
	for j := 0; j < m; j++ {
		prefix := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			diff := abs64(mice[i-1] - holes[j])
			prefix[i] = prefix[i-1] + diff
		}
		newDP := make([]int64, n+1)
		for i := 0; i <= n; i++ {
			newDP[i] = INF
		}
		for i := 0; i <= n; i++ {
			if dp[i] == INF {
				continue
			}
			for t := 0; t <= cap[j] && i+t <= n; t++ {
				cost := dp[i] + prefix[i+t] - prefix[i]
				if cost < newDP[i+t] {
					newDP[i+t] = cost
				}
			}
		}
		dp = newDP
	}
	if dp[n] >= INF {
		return -1
	}
	return dp[n]
}

func generateCase(rng *rand.Rand) (string, int64) {
	n := rng.Intn(4) + 1
	m := rng.Intn(4) + 1
	mice := make([]int64, n)
	for i := 0; i < n; i++ {
		mice[i] = int64(rng.Intn(21) - 10)
	}
	holes := make([]int64, m)
	capc := make([]int, m)
	for i := 0; i < m; i++ {
		holes[i] = int64(rng.Intn(21) - 10)
		capc[i] = rng.Intn(n) + 1
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", mice[i]))
	}
	sb.WriteByte('\n')
	for i := 0; i < m; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", holes[i], capc[i]))
	}
	expect := expectedF(append([]int64(nil), mice...), append([]int64(nil), holes...), append([]int(nil), capc...))
	return sb.String(), expect
}

func runCase(bin string, input string, expect int64) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
	}
	resStr := strings.TrimSpace(out.String())
	got, err := strconv.ParseInt(resStr, 10, 64)
	if err != nil {
		return fmt.Errorf("bad output %q", resStr)
	}
	if got != expect {
		return fmt.Errorf("expected %d got %d", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
