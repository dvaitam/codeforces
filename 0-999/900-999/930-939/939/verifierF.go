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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func solveF(n int, intervals [][2]int) string {
	maxT := 2 * n
	allowed := make([]bool, maxT+1)
	for _, p := range intervals {
		for t := p[0]; t <= p[1] && t <= maxT; t++ {
			allowed[t] = true
		}
	}
	inf := int(1e9)
	offset := n
	curA := make([]int, 2*n+1)
	curB := make([]int, 2*n+1)
	for i := range curA {
		curA[i] = inf
		curB[i] = inf
	}
	curA[offset] = 0
	for t := 0; t < maxT; t++ {
		nextA := make([]int, 2*n+1)
		nextB := make([]int, 2*n+1)
		for i := range nextA {
			nextA[i] = inf
			nextB[i] = inf
		}
		for d := 0; d <= 2*n; d++ {
			if curA[d] < inf {
				if d+1 <= 2*n {
					nextA[d+1] = min(nextA[d+1], curA[d])
				}
				if allowed[t] && d-1 >= 0 {
					nextB[d-1] = min(nextB[d-1], curA[d]+1)
				}
			}
			if curB[d] < inf {
				if d-1 >= 0 {
					nextB[d-1] = min(nextB[d-1], curB[d])
				}
				if allowed[t] && d+1 <= 2*n {
					nextA[d+1] = min(nextA[d+1], curB[d]+1)
				}
			}
		}
		curA = nextA
		curB = nextB
	}
	ans := min(curA[offset], curB[offset])
	if ans >= inf {
		return "Hungry"
	}
	return fmt.Sprintf("Full\n%d", ans)
}

func generateTests() []test {
	rand.Seed(6)
	tests := make([]test, 0, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(5) + 1
		k := rand.Intn(4) + 1
		intervals := make([][2]int, k)
		for j := 0; j < k; j++ {
			l := rand.Intn(2*n + 1)
			r := l + rand.Intn(2*n-l+1)
			intervals[j] = [2]int{l, r}
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, k)
		for _, p := range intervals {
			fmt.Fprintf(&sb, "%d %d\n", p[0], p[1])
		}
		tests = append(tests, test{in: sb.String(), out: solveF(n, intervals)})
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
		fmt.Println("Usage: go run verifierF.go /path/to/binary")
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
