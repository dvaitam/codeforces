package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCase struct {
	n int
	m int
	s string
}

func solve(n, m int, s string) int {
	w := make([][]int, m)
	for i := range w {
		w[i] = make([]int, m)
	}
	for i := 1; i < n; i++ {
		a := int(s[i-1] - 'a')
		b := int(s[i] - 'a')
		if a == b || a >= m || b >= m {
			continue
		}
		w[a][b]++
		w[b][a]++
	}
	total := make([]int, m)
	for i := 0; i < m; i++ {
		sum := 0
		for j := 0; j < m; j++ {
			sum += w[i][j]
		}
		total[i] = sum
	}
	maxMask := 1 << m
	sumW := make([][]int, m)
	for j := 0; j < m; j++ {
		arr := make([]int, maxMask)
		for mask := 1; mask < maxMask; mask++ {
			lb := mask & -mask
			k := bits.TrailingZeros(uint(lb))
			arr[mask] = arr[mask^lb] + w[j][k]
		}
		sumW[j] = arr
	}
	cross := make([]int, maxMask)
	for mask := 1; mask < maxMask; mask++ {
		lb := mask & -mask
		j := bits.TrailingZeros(uint(lb))
		prev := mask ^ lb
		cross[mask] = cross[prev] + total[j] - 2*sumW[j][prev]
	}
	const inf int = int(1e18)
	dp := make([]int, maxMask)
	for i := 1; i < maxMask; i++ {
		dp[i] = inf
	}
	for mask := 0; mask < maxMask; mask++ {
		for j := 0; j < m; j++ {
			if mask&(1<<j) == 0 {
				next := mask | (1 << j)
				val := dp[mask] + cross[next]
				if val < dp[next] {
					dp[next] = val
				}
			}
		}
	}
	return dp[maxMask-1]
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 100)
	for i := 0; i < 100; i++ {
		m := rng.Intn(4) + 1
		n := rng.Intn(10) + 1
		b := make([]byte, n)
		for j := 0; j < n; j++ {
			b[j] = byte('a' + rng.Intn(m))
		}
		tests = append(tests, testCase{n: n, m: m, s: string(b)})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		input := fmt.Sprintf("%d %d\n%s\n", t.n, t.m, t.s)
		want := fmt.Sprintf("%d", solve(t.n, t.m, t.s))
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("case %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("case %d failed: expected %s got %s\ninput:%s", i+1, want, got, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
