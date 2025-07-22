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

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func checkOutput(n, k int, output string) error {
	fields := strings.Fields(strings.TrimSpace(output))
	if len(fields) != 2*n {
		return fmt.Errorf("expected %d numbers got %d", 2*n, len(fields))
	}
	used := make([]bool, 2*n+1)
	a := make([]int, 2*n)
	for i, f := range fields {
		val, err := strconv.Atoi(f)
		if err != nil {
			return fmt.Errorf("bad number %q", f)
		}
		if val < 1 || val > 2*n {
			return fmt.Errorf("value out of range: %d", val)
		}
		if used[val] {
			return fmt.Errorf("duplicate value: %d", val)
		}
		used[val] = true
		a[i] = val
	}
	sumPair := 0
	sumDiff := 0
	for i := 0; i < n; i++ {
		sumPair += abs(a[2*i] - a[2*i+1])
		sumDiff += a[2*i] - a[2*i+1]
	}
	if sumPair-abs(sumDiff) != 2*k {
		return fmt.Errorf("equation not satisfied")
	}
	return nil
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func generateCase(rng *rand.Rand) (string, int, int) {
	n := rng.Intn(30) + 1
	k := rng.Intn(n/2 + 1)
	input := fmt.Sprintf("%d %d\n", n, k)
	return input, n, k
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, n, k := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if err := checkOutput(n, k, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s", i+1, err, in, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
