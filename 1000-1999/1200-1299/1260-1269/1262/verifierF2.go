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

const mod = 998244353

func solve(k int, h []int) int {
	n := len(h)
	size := 2*n + 5
	dp := make([]int, size)
	off := n + 2
	dp[off] = 1
	for i := 0; i < n; i++ {
		next := make([]int, size)
		if h[i] == h[(i+1)%n] {
			for d := 0; d < size; d++ {
				if dp[d] != 0 {
					next[d] = (next[d] + dp[d]*k) % mod
				}
			}
		} else {
			for d := 0; d < size; d++ {
				val := dp[d]
				if val == 0 {
					continue
				}
				next[d] = (next[d] + val*(k-2)) % mod
				if d+1 < size {
					next[d+1] = (next[d+1] + val) % mod
				}
				if d-1 >= 0 {
					next[d-1] = (next[d-1] + val) % mod
				}
			}
		}
		dp = next
	}
	ans := 0
	for d := off + 1; d < size; d++ {
		ans = (ans + dp[d]) % mod
	}
	return ans
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

func runCase(bin string, n, k int, h []int) error {
	input := fmt.Sprintf("%d %d\n", n, k)
	for i, v := range h {
		if i > 0 {
			input += " "
		}
		input += fmt.Sprintf("%d", v)
	}
	input += "\n"
	expect := solve(k, h)
	out, err := run(bin, input)
	if err != nil {
		return err
	}
	var got int
	if _, err := fmt.Sscan(out, &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != expect {
		return fmt.Errorf("expected %d got %d", expect, got)
	}
	return nil
}

func genCase(rng *rand.Rand) (int, int, []int) {
	n := rng.Intn(10) + 1
	k := rng.Intn(6) + 1
	h := make([]int, n)
	for i := range h {
		h[i] = rng.Intn(k) + 1
	}
	return n, k, h
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF2.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	n, k, h := 3, 3, []int{1, 2, 3}
	if err := runCase(bin, n, k, h); err != nil {
		fmt.Fprintf(os.Stderr, "case 1 failed: %v\n", err)
		os.Exit(1)
	}
	for i := 1; i < 100; i++ {
		n, k, h := genCase(rng)
		if err := runCase(bin, n, k, h); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
