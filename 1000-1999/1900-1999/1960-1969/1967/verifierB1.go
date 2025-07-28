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

type TestCaseB1 struct {
	n int64
	m int64
}

func genCaseB1(rng *rand.Rand) TestCaseB1 {
	n := int64(rng.Intn(100) + 1)
	m := int64(rng.Intn(100) + 1)
	return TestCaseB1{n, m}
}

func expectedB1(n, m int64) int64 {
	ans := int64(0)
	for g := int64(1); g <= m; g++ {
		kmax := (n + g) / (g * g)
		kmin := int64(1)
		if g == 1 {
			kmin = 2
		}
		if kmax >= kmin {
			ans += kmax - kmin + 1
		}
	}
	return ans
}

func runCaseB1(bin string, tc TestCaseB1, expect int64) error {
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierB1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := genCaseB1(rng)
		exp := expectedB1(tc.n, tc.m)
		if err := runCaseB1(bin, tc, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput: n=%d m=%d\n", i+1, err, tc.n, tc.m)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
