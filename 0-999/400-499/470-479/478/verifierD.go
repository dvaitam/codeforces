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

const mod = 1000000007

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

func intSqrt(x int) int {
	lo, hi := 0, x
	for lo <= hi {
		mid := (lo + hi) / 2
		if mid*mid <= x {
			lo = mid + 1
		} else {
			hi = mid - 1
		}
	}
	return hi
}

func expected(r, g int) string {
	total := r + g
	h := int((-1 + intSqrt(1+8*total)) / 2)
	S := h * (h + 1) / 2
	small, big := r, g
	if small > big {
		small, big = big, small
	}
	low := S - big
	if low < 0 {
		low = 0
	}
	dp := make([]int, small+1)
	dp[0] = 1
	for i := 1; i <= h; i++ {
		if i > small {
			continue
		}
		for s := small; s >= i; s-- {
			dp[s] += dp[s-i]
			if dp[s] >= mod {
				dp[s] -= mod
			}
		}
	}
	ans := 0
	for s := low; s <= small; s++ {
		ans += dp[s]
		if ans >= mod {
			ans -= mod
		}
	}
	return fmt.Sprintf("%d", ans)
}

func genCase(rng *rand.Rand) (string, string) {
	r := rng.Intn(200000 + 1)
	g := rng.Intn(200000 + 1)
	input := fmt.Sprintf("%d %d\n", r, g)
	return input, expected(r, g)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
