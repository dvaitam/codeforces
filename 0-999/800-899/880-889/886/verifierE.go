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

const MOD int64 = 1_000_000_007

type testCase struct {
	input    string
	expected int64
}

func expectedE(n, k int) int64 {
	if k >= n {
		return 0
	}
	fact := make([]int64, n+1)
	fact[0] = 1
	for i := 1; i <= n; i++ {
		fact[i] = fact[i-1] * int64(i) % MOD
	}
	inv := make([]int64, n)
	inv[1] = 1
	for i := 2; i < n; i++ {
		inv[i] = MOD - (MOD/int64(i))*inv[int(MOD%int64(i))]%MOD
	}
	w := make([]int64, k)
	if k > 0 {
		w[0] = 1
	}
	dpSum := int64(1)
	for i := 1; i < n; i++ {
		total := int64(0)
		for j := 0; j < k; j++ {
			total += w[j]
		}
		total %= MOD
		new0 := total * inv[i] % MOD
		for j := k - 1; j >= 1; j-- {
			w[j] = w[j-1]
		}
		if k > 0 {
			w[0] = new0
		}
		dpSum += w[0]
		if dpSum >= MOD {
			dpSum -= MOD
		}
	}
	good := fact[n-1] * dpSum % MOD
	bad := (fact[n] - good) % MOD
	if bad < 0 {
		bad += MOD
	}
	return bad
}

func generateCase(rng *rand.Rand) testCase {
	n := rng.Intn(10) + 1
	k := rng.Intn(10) + 1
	input := fmt.Sprintf("%d %d\n", n, k)
	return testCase{input: input, expected: expectedE(n, k)}
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCase(rng)
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
		var got int64
		if _, err := fmt.Sscan(out, &got); err != nil {
			fmt.Fprintf(os.Stderr, "case %d: cannot parse output: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d: expected %d got %d\ninput:\n%s", i+1, tc.expected, got, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
