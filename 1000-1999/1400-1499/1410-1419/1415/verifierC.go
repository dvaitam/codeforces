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

func run(bin string, input string) (string, error) {
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

func solveC(n, p, k int, s string, x, y int64) int64 {
	a := []byte(s)
	dp := make([]int64, n)
	for i := n - 1; i >= 0; i-- {
		add := int64(0)
		if a[i] == '0' {
			add = 1
		}
		if i+k < n {
			dp[i] = add + dp[i+k]
		} else {
			dp[i] = add
		}
	}
	ans := int64(1 << 62)
	for d := 0; p-1+d < n; d++ {
		idx := p - 1 + d
		cost := int64(d)*y + dp[idx]*x
		if cost < ans {
			ans = cost
		}
	}
	return ans
}

func runCase(bin string, n, p, k int, s string, x, y int64) error {
	input := fmt.Sprintf("1\n%d %d %d\n%s\n%d %d\n", n, p, k, s, x, y)
	expect := fmt.Sprintf("%d", solveC(n, p, k, s, x, y))
	out, err := run(bin, input)
	if err != nil {
		return err
	}
	if out != expect {
		return fmt.Errorf("expected %s got %s", expect, out)
	}
	return nil
}

func randString(rng *rand.Rand, n int) string {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			b[i] = '0'
		} else {
			b[i] = '1'
		}
	}
	return string(b)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	total := 0
	// edge cases
	if err := runCase(bin, 1, 1, 1, "0", 1, 1); err != nil {
		fmt.Fprintf(os.Stderr, "case %d failed: %v\n", total+1, err)
		os.Exit(1)
	}
	total++
	if err := runCase(bin, 5, 2, 2, "11001", 3, 4); err != nil {
		fmt.Fprintf(os.Stderr, "case %d failed: %v\n", total+1, err)
		os.Exit(1)
	}
	total++
	for total < 100 {
		n := rng.Intn(20) + 1
		p := rng.Intn(n) + 1
		k := rng.Intn(n) + 1
		s := randString(rng, n)
		x := int64(rng.Intn(10) + 1)
		y := int64(rng.Intn(10) + 1)
		if err := runCase(bin, n, p, k, s, x, y); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", total+1, err)
			os.Exit(1)
		}
		total++
	}
	fmt.Printf("All %d tests passed\n", total)
}
