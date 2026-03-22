package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

// Embedded solver for 823E
func solve823E(K int) int64 {
	mod := int64(1000000007)

	dp := make([]int64, K+2)
	dp[0] = 1
	dp[1] = 1

	for k := 2; k <= K; k++ {
		limit := K - k + 1
		t := make([]int64, limit+2)

		for c := 0; c <= limit+1; c++ {
			var sum int64 = 0
			for i := 0; i <= c; i++ {
				val := (dp[i] * dp[c-i]) % mod
				sum = (sum + val) % mod
			}
			t[c] = sum
		}

		nextDp := make([]int64, limit+1)
		for c := 0; c <= limit; c++ {
			var res int64 = 0
			if c > 0 {
				res = (res + t[c-1]) % mod
			}
			term2 := (int64(2*c+1) * t[c]) % mod
			res = (res + term2) % mod

			term3 := (int64(c+1) * int64(c)) % mod
			term3 = (term3 * t[c+1]) % mod
			res = (res + term3) % mod

			nextDp[c] = res
		}
		dp = nextDp
	}

	return dp[1]
}

func runBinary(bin, input string) (string, error) {
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

type Case struct{ input string }

func genCases() []Case {
	rng := rand.New(rand.NewSource(8234))
	cases := make([]Case, 100)
	for i := range cases {
		k := rng.Intn(20) + 1
		cases[i] = Case{fmt.Sprintf("%d\n", k)}
	}
	return cases
}

func runCase(bin string, c Case) error {
	var K int
	fmt.Sscan(c.input, &K)
	expected := fmt.Sprintf("%d", solve823E(K))

	got, err := runBinary(bin, c.input)
	if err != nil {
		return err
	}
	if strings.TrimSpace(expected) != strings.TrimSpace(got) {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	if bin == "--" {
		if len(os.Args) < 3 {
			fmt.Println("usage: go run verifierE.go /path/to/binary")
			os.Exit(1)
		}
		bin = os.Args[2]
	}
	cases := genCases()
	for i, c := range cases {
		if err := runCase(bin, c); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, c.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
