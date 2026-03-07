package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const MOD = 1000000007

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expected(n, k, m int) int64 {
	// Bitmask DP: state = (mask of chosen planets, last chosen planet, both 0-indexed)
	// Movement rule: from planet x (0-indexed), can go to y (0-indexed) if y <= x+m and y not visited.
	type state struct{ mask, last int }
	dp := make(map[state]int64, n)
	for s := 0; s < n; s++ {
		dp[state{1 << s, s}] = 1
	}
	for step := 1; step < k; step++ {
		ndp := make(map[state]int64, len(dp))
		for st, ways := range dp {
			for y := 0; y < n; y++ {
				if st.mask>>y&1 != 0 {
					continue
				}
				if y <= st.last+m {
					ns := state{st.mask | (1 << y), y}
					ndp[ns] = (ndp[ns] + ways) % MOD
				}
			}
		}
		dp = ndp
	}
	var ans int64
	for _, ways := range dp {
		ans = (ans + ways) % MOD
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, int64) {
	n := rng.Intn(20) + 1
	k := rng.Intn(n) + 1
	if k > 12 {
		k = 12
	}
	m := rng.Intn(4) + 1
	input := fmt.Sprintf("%d %d %d\n", n, k, m)
	return input, expected(n, k, m)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	if strings.HasSuffix(bin, ".go") {
		tmp, err := os.CreateTemp("", "verifierF1-bin-*")
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to create temp file: %v\n", err)
			os.Exit(1)
		}
		tmp.Close()
		defer os.Remove(tmp.Name())
		out, err := exec.Command("go", "build", "-o", tmp.Name(), bin).CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "compile error: %v\n%s", err, out)
			os.Exit(1)
		}
		bin = tmp.Name()
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	fixed := [][3]int{
		{3, 3, 1},
		{4, 2, 2},
	}
	idx := 0
	for ; idx < len(fixed); idx++ {
		n := fixed[idx][0]
		k := fixed[idx][1]
		m := fixed[idx][2]
		inp := fmt.Sprintf("%d %d %d\n", n, k, m)
		exp := strconv.FormatInt(expected(n, k, m), 10)
		out, err := run(bin, inp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", idx+1, err, inp)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", idx+1, exp, out, inp)
			os.Exit(1)
		}
	}
	for ; idx < 100; idx++ {
		inp, expVal := generateCase(rng)
		out, err := run(bin, inp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", idx+1, err, inp)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strconv.FormatInt(expVal, 10) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\ninput:%s", idx+1, expVal, out, inp)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
