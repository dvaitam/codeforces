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

type pair struct{ n, k int64 }

func feasible(n, k int64) bool {
	if k == 1 || n == 1 {
		return false
	}
	allowed := []int{}
	for i := 2; int64(i) <= n; i++ {
		if k%int64(i) == 0 {
			allowed = append(allowed, i)
		}
	}
	if len(allowed) == 0 {
		return false
	}
	dp := make([]bool, n+1)
	dp[0] = true
	for i := 0; i <= int(n); i++ {
		if dp[i] {
			for _, l := range allowed {
				if i+l <= int(n) {
					dp[i+l] = true
				}
			}
		}
	}
	return dp[n]
}

func genCase(rng *rand.Rand) (string, string) {
	t := rng.Intn(5) + 1
	cases := make([]pair, t)
	for i := 0; i < t; i++ {
		cases[i] = pair{int64(rng.Intn(8) + 1), int64(rng.Intn(10) + 1)}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	for _, c := range cases {
		fmt.Fprintf(&sb, "%d %d\n", c.n, c.k)
	}
	input := sb.String()
	var out strings.Builder
	for _, c := range cases {
		if feasible(c.n, c.k) {
			out.WriteString("YES\n")
		} else {
			out.WriteString("NO\n")
		}
	}
	expected := out.String()
	return input, expected
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if strings.TrimSpace(out.String()) != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %s got %s", expected, out.String())
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
