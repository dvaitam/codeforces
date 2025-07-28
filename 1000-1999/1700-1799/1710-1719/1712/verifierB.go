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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b int) int {
	return a / gcd(a, b) * b
}

func bruteMaxSum(n int) int {
	used := make([]bool, n)
	maxSum := 0
	var dfs func(int, int)
	dfs = func(pos, sum int) {
		if pos == n {
			if sum > maxSum {
				maxSum = sum
			}
			return
		}
		for i := 1; i <= n; i++ {
			if !used[i-1] {
				used[i-1] = true
				dfs(pos+1, sum+lcm(pos+1, i))
				used[i-1] = false
			}
		}
	}
	dfs(0, 0)
	return maxSum
}

func generateCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(7) + 1
	maxSum := bruteMaxSum(n)
	input := fmt.Sprintf("1\n%d\n", n)
	return input, maxSum
}

func parsePermutation(out string, n int) ([]int, error) {
	fields := strings.Fields(out)
	if len(fields) != n {
		return nil, fmt.Errorf("expected %d numbers got %d", n, len(fields))
	}
	perm := make([]int, n)
	used := make([]bool, n)
	for i, f := range fields {
		val, err := strconv.Atoi(f)
		if err != nil {
			return nil, fmt.Errorf("bad integer: %v", err)
		}
		if val < 1 || val > n {
			return nil, fmt.Errorf("value out of range")
		}
		if used[val-1] {
			return nil, fmt.Errorf("duplicate value")
		}
		used[val-1] = true
		perm[i] = val
	}
	return perm, nil
}

func permSum(perm []int) int {
	sum := 0
	for i, v := range perm {
		sum += lcm(i+1, v)
	}
	return sum
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, expSum := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		nLines := strings.Split(strings.TrimSpace(in), "\n")
		n, _ := strconv.Atoi(strings.TrimSpace(nLines[1]))
		perm, err := parsePermutation(out, n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if permSum(perm) != expSum {
			fmt.Fprintf(os.Stderr, "case %d failed: wrong sum\ninput:\n%soutput:\n%s", i+1, in, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
