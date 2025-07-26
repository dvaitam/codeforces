package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
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

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func unionSize(nums []int) int {
	m := len(nums)
	ans := 0
	for mask := 1; mask < 1<<m; mask++ {
		g := 0
		bits := 0
		for i := 0; i < m; i++ {
			if mask>>i&1 == 1 {
				if g == 0 {
					g = nums[i]
				} else {
					g = gcd(g, nums[i])
				}
				bits++
			}
		}
		if bits%2 == 1 {
			ans += g
		} else {
			ans -= g
		}
	}
	return ans
}

func bruteForce(n, k int) int {
	values := make([]int, n-2)
	for i := range values {
		values[i] = i + 3
	}
	choose := make([]int, k)
	best := math.MaxInt32
	var dfs func(start, idx int)
	dfs = func(start, idx int) {
		if idx == k {
			val := unionSize(choose)
			if val < best {
				best = val
			}
			return
		}
		for i := start; i <= len(values)-(k-idx); i++ {
			choose[idx] = values[i]
			dfs(i+1, idx+1)
		}
	}
	dfs(0, 0)
	return best
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 4
	k := rng.Intn(min(n-2, 4)) + 1
	input := fmt.Sprintf("%d %d\n", n, k)
	expect := fmt.Sprintf("%d", bruteForce(n, k))
	return input, expect
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
