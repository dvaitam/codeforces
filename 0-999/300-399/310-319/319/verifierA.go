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

const mod int64 = 1000000007

func compute(x string) int64 {
	n := len(x)
	pow2 := make([]int64, 2*n+1)
	pow2[0] = 1
	for i := 1; i <= 2*n; i++ {
		pow2[i] = (pow2[i-1] * 2) % mod
	}
	var ans int64
	for k := 1; k <= n; k++ {
		ans = (ans * 2) % mod
		if x[n-k] == '1' {
			ans = (ans + pow2[2*k-2]) % mod
		}
	}
	return ans
}

func runCase(bin, x string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(x + "\n")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int64
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	exp := compute(x)
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []string{"0", "1", "00", "01", "10", "11", strings.Repeat("0", 100), strings.Repeat("1", 100)}
	for i := 0; i < 92; i++ {
		n := rng.Intn(100) + 1
		b := make([]byte, n)
		for j := range b {
			if rng.Intn(2) == 1 {
				b[j] = '1'
			} else {
				b[j] = '0'
			}
		}
		cases = append(cases, string(b))
	}
	for i, x := range cases {
		if err := runCase(bin, x); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s\n", i+1, err, x)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
