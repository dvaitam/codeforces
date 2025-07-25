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

const MOD = 1000000007

func runProgram(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

var memo map[[3]int]int

func ways(w, h, k int) int {
	if k == 0 {
		return 1
	}
	if w <= 0 || h <= 0 {
		return 0
	}
	key := [3]int{w, h, k}
	if v, ok := memo[key]; ok {
		return v
	}
	var total int64
	for nw := 1; nw <= w-2; nw++ {
		for nh := 1; nh <= h-2; nh++ {
			cnt := (w - nw - 1) * (h - nh - 1)
			if cnt <= 0 {
				continue
			}
			inner := ways(nw, nh, k-1)
			total += int64(cnt) * int64(inner)
		}
	}
	memo[key] = int(total % MOD)
	return memo[key]
}

func solve(n, m, k int) int {
	memo = make(map[[3]int]int)
	return ways(n, m, k)
}

func randomTest(rng *rand.Rand) (int, int, int) {
	n := rng.Intn(8) + 3
	m := rng.Intn(8) + 3
	k := rng.Intn(3) + 1
	return n, m, k
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(42))
	for i := 0; i < 100; i++ {
		n, m, k := randomTest(rng)
		input := fmt.Sprintf("%d %d %d\n", n, m, k)
		expected := solve(n, m, k)
		out, err := runProgram(bin, input)
		if err != nil {
			fmt.Printf("Test %d runtime error: %v\n", i+1, err)
			return
		}
		out = strings.TrimSpace(out)
		if out != fmt.Sprintf("%d", expected) {
			fmt.Printf("Test %d failed: expected %d got %s\n", i+1, expected, out)
			return
		}
	}
	fmt.Println("All tests passed")
}
