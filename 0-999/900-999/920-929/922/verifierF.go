package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCaseF struct {
	N int
	K int
}

func expectedF(N, K int) string {
	A := make([]int, N+2)
	s := 0
	n := 0
	for i := 1; i <= N; i++ {
		for j := i * 2; j <= N; j += i {
			A[j]++
		}
		s += A[i]
		if s >= K {
			n = i
			break
		}
	}
	if s < K {
		return "No"
	}
	s -= K
	k := n
	removed := make([]bool, n+2)
	for i := 2; i <= n && s > 0; i++ {
		cost := A[i] + n/i - 1
		if s >= cost {
			s -= cost
			removed[i] = true
			k--
		}
	}
	var sb strings.Builder
	sb.WriteString("Yes\n")
	sb.WriteString(fmt.Sprintf("%d\n", k))
	for i := 1; i <= n; i++ {
		if !removed[i] {
			sb.WriteString(fmt.Sprintf("%d ", i))
		}
	}
	sb.WriteByte('\n')
	return strings.TrimSpace(sb.String())
}

func genTestsF() []testCaseF {
	rand.Seed(6)
	tests := make([]testCaseF, 0, 100)
	for len(tests) < 100 {
		N := rand.Intn(100) + 2
		maxPairs := N * (N - 1) / 2
		K := rand.Intn(maxPairs + 1)
		if rand.Intn(5) == 0 { // sometimes impossible
			K = maxPairs + rand.Intn(100) + 1
		}
		tests = append(tests, testCaseF{N: N, K: K})
	}
	return tests
}

func runCase(bin string, tc testCaseF) error {
	input := fmt.Sprintf("%d %d\n", tc.N, tc.K)
	cmd := exec.Command(bin)
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	expect := expectedF(tc.N, tc.K)
	if got != expect {
		return fmt.Errorf("N=%d K=%d expected:\n%s\n\ngot:\n%s", tc.N, tc.K, expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsF()
	for i, tc := range tests {
		if err := runCase(bin, tc); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
