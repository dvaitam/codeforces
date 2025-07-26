package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type test struct {
	in  string
	out string
}

func solveB(n int64, a []int64) string {
	bestIdx := 1
	bestCount := n / a[0]
	maxTransport := bestCount * a[0]
	for i := 1; i < len(a); i++ {
		cnt := n / a[i]
		transported := cnt * a[i]
		if transported > maxTransport {
			maxTransport = transported
			bestIdx = i + 1
			bestCount = cnt
		}
	}
	return fmt.Sprintf("%d %d", bestIdx, bestCount)
}

func generateTests() []test {
	rand.Seed(2)
	tests := make([]test, 0, 100)
	for i := 0; i < 100; i++ {
		n := rand.Int63n(1_000_000) + 1
		k := rand.Intn(10) + 1
		a := make([]int64, k)
		for j := range a {
			a[j] = rand.Int63n(1_000_000) + 1
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, k)
		for j := 0; j < k; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", a[j])
		}
		sb.WriteByte('\n')
		tests = append(tests, test{in: sb.String(), out: solveB(n, a)})
	}
	return tests
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return stderr.String(), err
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		out, err := run(bin, t.in)
		if err != nil {
			fmt.Printf("Test %d failed to run: %v\n", i+1, err)
			fmt.Print(out)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != t.out {
			fmt.Printf("Test %d failed. Expected %q, got %q. Input:\n%s", i+1, t.out, out, t.in)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed!\n", len(tests))
}
