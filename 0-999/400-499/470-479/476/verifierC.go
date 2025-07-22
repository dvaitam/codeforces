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

func runCandidate(bin, input string) (string, error) {
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

const mod int64 = 1000000007

func expected(a, b int64) int64 {
	inv2 := (mod + 1) / 2
	aMod := a % mod
	bMod := b % mod
	sumK := aMod * ((aMod + 1) % mod) % mod * inv2 % mod
	t := (bMod*sumK%mod + aMod) % mod
	sumR := ((bMod - 1 + mod) % mod) * bMod % mod * inv2 % mod
	ans := t * sumR % mod
	return ans
}

type testCase struct {
	a int64
	b int64
}

func generateRandomCase(rng *rand.Rand) testCase {
	a := rng.Int63n(10000000) + 1
	b := rng.Int63n(10000000) + 1
	return testCase{a, b}
}

func runCase(bin string, tc testCase) error {
	input := fmt.Sprintf("%d %d\n", tc.a, tc.b)
	exp := expected(tc.a, tc.b)
	out, err := runCandidate(bin, input)
	if err != nil {
		return err
	}
	var got int64
	if _, err := fmt.Sscan(out, &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCase{
		{1, 1},
		{2, 3},
	}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateRandomCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
