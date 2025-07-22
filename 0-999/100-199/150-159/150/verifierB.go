package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func powmod(base, exp, mod int64) int64 {
	res := int64(1)
	b := base % mod
	for exp > 0 {
		if exp&1 == 1 {
			res = res * b % mod
		}
		b = b * b % mod
		exp >>= 1
	}
	return res
}

// solve computes the answer for problem B
func solve(n int, m int64, k int) int64 {
	parent := make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i
	}
	var find func(int) int
	find = func(x int) int {
		if parent[x] != x {
			parent[x] = find(parent[x])
		}
		return parent[x]
	}
	union := func(a, b int) {
		ra, rb := find(a), find(b)
		if ra != rb {
			parent[rb] = ra
		}
	}
	half := k / 2
	for i := 0; i <= n-k; i++ {
		for j := 0; j < half; j++ {
			union(i+j, i+k-1-j)
		}
	}
	comps := 0
	for i := 0; i < n; i++ {
		if find(i) == i {
			comps++
		}
	}
	const mod = 1000000007
	return powmod(m, int64(comps), mod)
}

func run(bin, input string) (string, error) {
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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(10) + 1
		mVal := int64(rng.Intn(10) + 1)
		k := rng.Intn(n) + 1
		input := fmt.Sprintf("%d %d %d\n", n, mVal, k)
		expected := fmt.Sprintf("%d", solve(n, mVal, k))
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(expected) != strings.TrimSpace(got) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
