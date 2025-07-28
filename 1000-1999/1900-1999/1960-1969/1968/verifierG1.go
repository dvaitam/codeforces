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

func zAlgorithm(s string) []int {
	n := len(s)
	z := make([]int, n)
	l, r := 0, 0
	for i := 1; i < n; i++ {
		if i <= r {
			if r-i+1 < z[i-l] {
				z[i] = r - i + 1
			} else {
				z[i] = z[i-l]
			}
		}
		for i+z[i] < n && s[z[i]] == s[i+z[i]] {
			z[i]++
		}
		if i+z[i]-1 > r {
			l = i
			r = i + z[i] - 1
		}
	}
	z[0] = n
	return z
}

func can(s string, z []int, k, d int) bool {
	n := len(s)
	if d == 0 {
		return true
	}
	if k*d > n {
		return false
	}
	occ := make([]bool, n+2)
	limit := n - d + 1
	for i := 1; i <= limit; i++ {
		if z[i-1] >= d {
			occ[i] = true
		}
	}
	next := make([]int, n+2)
	next[n+1] = n + 1
	for i := n; i >= 1; i-- {
		if occ[i] {
			next[i] = i
		} else {
			next[i] = next[i+1]
		}
	}
	pos := 1
	for i := 1; i < k; i++ {
		target := pos + d
		if target > limit {
			return false
		}
		pos = next[target]
		if pos == n+1 {
			return false
		}
	}
	return pos <= limit
}

func solve(n, k int, s string) int {
	z := zAlgorithm(s)
	lo, hi := 0, n/k
	for lo < hi {
		mid := (lo + hi + 1) / 2
		if can(s, z, k, mid) {
			lo = mid
		} else {
			hi = mid - 1
		}
	}
	return lo
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	k := rng.Intn(n) + 1
	letters := []byte("abcdefghijklmnopqrstuvwxyz")
	sb := make([]byte, n)
	for i := 0; i < n; i++ {
		sb[i] = letters[rng.Intn(len(letters))]
	}
	s := string(sb)
	input := fmt.Sprintf("1\n%d %d %d\n%s\n", n, k, k, s)
	expect := fmt.Sprint(solve(n, k, s))
	return input, expect
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	const cases = 100
	for i := 0; i < cases; i++ {
		inp, exp := genCase(rng)
		got, err := run(bin, inp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed\nexpected: %s\ngot: %s\ninput:\n%s", i+1, exp, got, inp)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", cases)
}
