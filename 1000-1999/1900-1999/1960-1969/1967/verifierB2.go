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

type TestCaseB2 struct {
	n int
	m int
}

func genCaseB2(rng *rand.Rand) TestCaseB2 {
	n := rng.Intn(40) + 1
	m := rng.Intn(40) + 1
	return TestCaseB2{n, m}
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func expectedB2(n, m int) int64 {
	count := int64(0)
	for a := 1; a <= n; a++ {
		for b := 1; b <= m; b++ {
			g := gcd(a, b)
			if (b*g)%(a+b) == 0 {
				count++
			}
		}
	}
	return count
}

func runCaseB2(bin string, tc TestCaseB2, expect int64) error {
	input := fmt.Sprintf("1\n%d %d\n", tc.n, tc.m)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int64
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("invalid output: %v", err)
	}
	if got != expect {
		return fmt.Errorf("expected %d got %d", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB2.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := genCaseB2(rng)
		exp := expectedB2(tc.n, tc.m)
		if err := runCaseB2(bin, tc, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput: n=%d m=%d\n", i+1, err, tc.n, tc.m)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
