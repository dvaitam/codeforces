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

type testCase struct {
	a, b uint64
}

func gcd(a, b uint64) uint64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func expected(a, b uint64) uint64 {
	var cnt uint64
	for a != 1 || b != 1 {
		if a > b {
			k := (a - 1) / b
			cnt += k
			a -= k * b
		} else {
			k := (b - 1) / a
			cnt += k
			b -= k * a
		}
	}
	return cnt + 1
}

func runCase(bin string, tc testCase) error {
	input := fmt.Sprintf("%d %d\n", tc.a, tc.b)
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
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got uint64
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	exp := expected(tc.a, tc.b)
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func generateCases(rng *rand.Rand) []testCase {
	cases := []testCase{{1, 1}, {3, 2}, {1, 2}}
	for len(cases) < 100 {
		a := rng.Uint64()%1000000000000000000 + 1
		b := rng.Uint64()%1000000000000000000 + 1
		if gcd(a, b) != 1 {
			continue
		}
		cases = append(cases, testCase{a, b})
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := generateCases(rng)
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput: %d %d\n", i+1, err, tc.a, tc.b)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
