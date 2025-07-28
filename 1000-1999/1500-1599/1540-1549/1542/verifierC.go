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

const mod int64 = 1e9 + 7

type testCase struct {
	input    string
	expected int64
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b, limit int64) int64 {
	g := gcd(a, b)
	res := a / g * b
	if res > limit {
		return limit + 1
	}
	return res
}

func solve(n int64) int64 {
	ans := int64(0)
	l := int64(1)
	for m := int64(2); l <= n; m++ {
		nl := lcm(l, m, n)
		cnt := n/l - n/nl
		ans = (ans + cnt*m) % mod
		l = nl
	}
	return ans % mod
}

func generateCase(rng *rand.Rand) testCase {
	n := rng.Int63n(1_000_000_000_000_0000) + 1 // up to 1e16
	input := fmt.Sprintf("1\n%d\n", n)
	return testCase{input: input, expected: solve(n)}
}

func runCase(bin string, tc testCase) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(tc.input)
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
	if got%mod != tc.expected%mod {
		return fmt.Errorf("expected %d got %d", tc.expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	cases := []testCase{
		{input: "1\n1\n", expected: solve(1)},
		{input: "1\n4\n", expected: solve(4)},
	}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCase(rng))
	}

	for idx, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", idx+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
