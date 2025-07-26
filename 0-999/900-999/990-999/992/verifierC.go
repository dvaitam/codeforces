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

const mod int64 = 1000000007

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

func modPow(a, b int64) int64 {
	res := int64(1)
	a %= mod
	for b > 0 {
		if b&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		b >>= 1
	}
	return res
}

func computeC(x, k int64) int64 {
	if x == 0 {
		return 0
	}
	pow2k := modPow(2, k)
	pow2k1 := pow2k * 2 % mod
	ans := (pow2k1*(x%mod)%mod - pow2k + 1) % mod
	if ans < 0 {
		ans += mod
	}
	return ans
}

type testCaseC struct {
	input    string
	expected int64
}

func generateCaseC(rng *rand.Rand) testCaseC {
	x := rng.Int63n(1_000_000_000_000)
	k := rng.Int63n(1_000_000_000_000)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", x, k)
	return testCaseC{input: sb.String(), expected: computeC(x, k)}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		tc := generateCaseC(rng)
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, tc.input)
			os.Exit(1)
		}
		var val int64
		if _, err := fmt.Sscan(out, &val); err != nil || val%mod != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d: expected %d got %s\ninput:\n%s", i, tc.expected, out, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
