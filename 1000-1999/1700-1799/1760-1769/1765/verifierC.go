package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const MOD_REF int64 = 998244353

type testCase struct {
	name  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}

	tests := generateTests()
	candidate := os.Args[1]

	for idx, tc := range tests {
		refVal := solveReference(tc.input)

		candOut, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candVal, err := parseInt(candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, candOut)
			os.Exit(1)
		}

		if candVal != refVal {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s): expected %d, got %d\ninput:\n%s\n", idx+1, tc.name, refVal, candVal, tc.input)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

// Embedded solver from the ACCEPTED solution.
func solveReference(input string) int64 {
	var n, k int
	fmt.Sscanf(input, "%d %d", &n, &k)

	fact := make([]int64, 4*n+1)
	invFact := make([]int64, 4*n+1)
	fact[0] = 1
	invFact[0] = 1
	for i := 1; i <= 4*n; i++ {
		fact[i] = (fact[i-1] * int64(i)) % MOD_REF
	}
	invFact[4*n] = power(fact[4*n], MOD_REF-2)
	for i := 4*n - 1; i >= 1; i-- {
		invFact[i] = (invFact[i+1] * int64(i+1)) % MOD_REF
	}

	nCr := func(n, r int) int64 {
		if r < 0 || r > n {
			return 0
		}
		num := fact[n]
		den := (invFact[r] * invFact[n-r]) % MOD_REF
		return (num * den) % MOD_REF
	}

	F := make([]int64, 4*n+1)
	F2 := make([]int64, 4*n+1)
	F3 := make([]int64, 4*n+1)
	F4 := make([]int64, 4*n+1)
	S := make([]int64, 4*n+1)

	F[n] = 1
	F2[2*n] = 1
	F3[3*n] = 1
	F4[4*n] = 1
	S[4*n] = 1

	for x := n - 1; x >= 1; x-- {
		c := nCr(n, x)
		c2 := (c * c) % MOD_REF
		c3 := (c2 * c) % MOD_REF
		c4 := (c3 * c) % MOD_REF

		c4_3 := (4 * c) % MOD_REF
		c6_2 := (6 * c2) % MOD_REF
		c4_1 := (4 * c3) % MOD_REF

		c3_2 := (3 * c) % MOD_REF
		c3_1 := (3 * c2) % MOD_REF

		c2_1 := (2 * c) % MOD_REF

		for i := 4 * n; i >= 0; i-- {
			new_F4 := F4[i]
			if i-x >= 0 {
				new_F4 = (new_F4 + F3[i-x]*c4_3) % MOD_REF
			}
			if i-2*x >= 0 {
				new_F4 = (new_F4 + F2[i-2*x]*c6_2) % MOD_REF
			}
			if i-3*x >= 0 {
				new_F4 = (new_F4 + F[i-3*x]*c4_1) % MOD_REF
			}
			if i == 4*x {
				new_F4 = (new_F4 + c4) % MOD_REF
			}
			F4[i] = new_F4
		}

		for i := 3 * n; i >= 0; i-- {
			new_F3 := F3[i]
			if i-x >= 0 {
				new_F3 = (new_F3 + F2[i-x]*c3_2) % MOD_REF
			}
			if i-2*x >= 0 {
				new_F3 = (new_F3 + F[i-2*x]*c3_1) % MOD_REF
			}
			if i == 3*x {
				new_F3 = (new_F3 + c3) % MOD_REF
			}
			F3[i] = new_F3
		}

		for i := 2 * n; i >= 0; i-- {
			new_F2 := F2[i]
			if i-x >= 0 {
				new_F2 = (new_F2 + F[i-x]*c2_1) % MOD_REF
			}
			if i == 2*x {
				new_F2 = (new_F2 + c2) % MOD_REF
			}
			F2[i] = new_F2
		}

		for i := n; i >= 0; i-- {
			new_F := F[i]
			if i == x {
				new_F = (new_F + c) % MOD_REF
			}
			F[i] = new_F
		}

		for i := 0; i <= 4*n; i++ {
			S[i] = (S[i] + F4[i]) % MOD_REF
		}
	}

	ans := int64(0)
	for w := 0; w <= 4*n-1; w++ {
		term1 := (S[w] * modInverse(nCr(4*n, w))) % MOD_REF
		num := (int64(n) - term1 + MOD_REF) % MOD_REF
		den := int64(4*n - w)
		E_w := (num * modInverse(den)) % MOD_REF

		count := int64(0)
		if w < k {
			count = 1
		} else if w == k {
			count = int64(4*n - k)
		}

		ans = (ans + count*E_w) % MOD_REF
	}

	return ans
}

func power(base, exp int64) int64 {
	res := int64(1)
	base %= MOD_REF
	for exp > 0 {
		if exp%2 == 1 {
			res = (res * base) % MOD_REF
		}
		base = (base * base) % MOD_REF
		exp /= 2
	}
	return res
}

func modInverse(n int64) int64 {
	return power(n, MOD_REF-2)
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return stdout.String() + stderr.String(), err
	}
	return stdout.String(), nil
}

func parseInt(out string) (int64, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return 0, fmt.Errorf("empty output")
	}
	val, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q: %v", fields[0], err)
	}
	return val, nil
}

func generateTests() []testCase {
	tests := []testCase{
		buildCase("n1-k1", 1, 1),
		buildCase("n2-k1", 2, 1),
		buildCase("n2-k3", 2, 3),
		buildCase("n3-k4", 3, 4),
		buildCase("n5-k20", 5, 20),
	}

	rng := rand.New(rand.NewSource(1765))
	for i := 0; i < 40; i++ {
		n := rng.Intn(30) + 1
		k := rng.Intn(4*n) + 1
		tests = append(tests, buildCase(fmt.Sprintf("random-%d", i+1), n, k))
	}
	return tests
}

func buildCase(name string, n, k int) testCase {
	return testCase{
		name:  name,
		input: fmt.Sprintf("%d %d\n", n, k),
	}
}
