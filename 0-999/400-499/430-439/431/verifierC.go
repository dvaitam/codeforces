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

func solveC(n, k, d int) string {
	dp := make([][2]int, n+1)
	dp[0][0] = 1
	for s := 1; s <= n; s++ {
		var without, with int
		for i := 1; i <= k && i <= s; i++ {
			if i < d {
				without = (without + dp[s-i][0]) % mod
				with = (with + dp[s-i][1]) % mod
			} else {
				with = (with + dp[s-i][0] + dp[s-i][1]) % mod
			}
		}
		dp[s][0] = without
		dp[s][1] = with
	}
	return fmt.Sprintf("%d", dp[n][1])
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	k := rng.Intn(10) + 1
	d := rng.Intn(k) + 1
	input := fmt.Sprintf("%d %d %d\n", n, k, d)
	return input, solveC(n, k, d)
}

func runCase(bin, input, expected string) error {
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
	outStr := strings.TrimSpace(out.String())
	expStr := strings.TrimSpace(expected)
	if outStr != expStr {
		return fmt.Errorf("expected %q got %q", expStr, outStr)
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

	cases := []struct{ in, out string }{}
	for i := 0; i < 102; i++ {
		in, out := generateCase(rng)
		cases = append(cases, struct{ in, out string }{in, out})
	}

	for i, tc := range cases {
		if err := runCase(bin, tc.in, tc.out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
