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

func solveC(n int, a []int64, s, f int) string {
	k := f - s
	prefix := make([]int64, n+k+1)
	for i := 1; i <= n+k; i++ {
		prefix[i] = prefix[i-1] + a[(i-1)%n]
	}
	bestSum := int64(-1)
	bestTime := 1
	for start := 1; start <= n; start++ {
		sum := prefix[start+k-1] - prefix[start-1]
		time := (s-start+n)%n + 1
		if sum > bestSum || (sum == bestSum && time < bestTime) {
			bestSum = sum
			bestTime = time
		}
	}
	return fmt.Sprintf("%d", bestTime)
}

func generateTests() []test {
	rand.Seed(3)
	tests := make([]test, 0, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(9) + 2 // 2..10
		a := make([]int64, n)
		for j := range a {
			a[j] = rand.Int63n(100) + 1
		}
		k := rand.Intn(n-1) + 1
		s := rand.Intn(n-k+1) + 1
		f := s + k
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", a[j])
		}
		sb.WriteByte('\n')
		fmt.Fprintf(&sb, "%d %d\n", s, f)
		tests = append(tests, test{in: sb.String(), out: solveC(n, a, s, f)})
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
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
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
