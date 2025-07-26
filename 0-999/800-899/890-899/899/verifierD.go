package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type test struct {
	input    string
	expected string
}

func countPairs(n, sum int64) int64 {
	if sum < 2 || sum > 2*n {
		return 0
	}
	iMin := int64(1)
	if sum-n > iMin {
		iMin = sum - n
	}
	iMax := n
	if sum-1 < iMax {
		iMax = sum - 1
	}
	maxI := (sum - 1) / 2
	if iMax > maxI {
		iMax = maxI
	}
	if iMin > iMax {
		return 0
	}
	return iMax - iMin + 1
}

func solve(input string) string {
	reader := strings.NewReader(input)
	var n int64
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return ""
	}
	pow10 := [10]int64{1, 10, 100, 1000, 10000, 100000, 1000000, 10000000, 100000000, 1000000000}
	for k := 9; k >= 0; k-- {
		m := pow10[k]
		r := m - 1
		var total int64
		for s := r; s <= 2*n; s += m {
			if s < 2 {
				continue
			}
			total += countPairs(n, s)
		}
		if total > 0 {
			return fmt.Sprintf("%d", total)
		}
	}
	return "0"
}

func generateTests() []test {
	rand.Seed(8994)
	var tests []test
	fixed := []string{"10\n", "50\n", "100\n"}
	for _, f := range fixed {
		tests = append(tests, test{f, solve(f)})
	}
	for len(tests) < 100 {
		n := rand.Int63n(1000000) + 1
		inp := fmt.Sprintf("%d\n", n)
		tests = append(tests, test{inp, solve(inp)})
	}
	return tests
}

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", bin)
	} else {
		cmd = exec.CommandContext(ctx, bin)
	}
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("time limit")
	}
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, out)
	}
	return strings.TrimSpace(string(out)), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		out, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if out != strings.TrimSpace(t.expected) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\ninput:\n%sexpected:%s\n got:%s\n", i+1, t.input, t.expected, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
