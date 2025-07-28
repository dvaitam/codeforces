package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCaseA struct {
	n int
	k int
	b []int
}

func expectedA(n, k int, b []int) string {
	steps := k
	if steps > n {
		steps = n
	}
	s := 0
	for i := 0; i < steps; i++ {
		idx := n - s - 1
		if idx < 0 {
			idx %= n
			idx += n
		}
		idx %= n
		x := b[idx]
		if x < 1 || x > n {
			return "No"
		}
		s = (s + x) % n
	}
	return "Yes"
}

func genTestsA() []testCaseA {
	rand.Seed(1)
	tests := make([]testCaseA, 0, 100)
	for len(tests) < 100 {
		n := rand.Intn(8) + 1
		k := rand.Intn(12) + 1
		b := make([]int, n)
		for i := range b {
			if rand.Intn(4) == 0 {
				b[i] = rand.Intn(n+3) + 1
			} else {
				b[i] = rand.Intn(n) + 1
			}
		}
		tests = append(tests, testCaseA{n: n, k: k, b: b})
	}
	return tests
}

func runCase(bin string, tc testCaseA) error {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.k))
	for i, v := range tc.b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')

	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	expect := expectedA(tc.n, tc.k, tc.b)
	if got != expect {
		return fmt.Errorf("expected %s got %s (n=%d k=%d b=%v)", expect, got, tc.n, tc.k, tc.b)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsA()
	for i, tc := range tests {
		if err := runCase(bin, tc); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
