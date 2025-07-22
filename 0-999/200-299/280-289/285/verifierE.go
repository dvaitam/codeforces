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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

const MOD = 1000000007

func add(a, b int) int {
	a += b
	if a >= MOD {
		a -= MOD
	}
	return a
}

func sub(a, b int) int {
	a -= b
	if a < 0 {
		a += MOD
	}
	return a
}

func mul(a, b int) int {
	return int((int64(a) * int64(b)) % MOD)
}

func powmod(a, e int) int {
	res := 1
	x := a
	for e > 0 {
		if e&1 == 1 {
			res = mul(res, x)
		}
		x = mul(x, x)
		e >>= 1
	}
	return res
}

func expected(n, k int) string {
	dpCurr := make([][2][2]int, n+1)
	dpCurr[0][0][0] = 1
	for pos := 1; pos <= n; pos++ {
		dpNext := make([][2][2]int, n+1)
		for s := 0; s <= n; s++ {
			for prevA := 0; prevA < 2; prevA++ {
				for prev2A := 0; prev2A < 2; prev2A++ {
					v := dpCurr[s][prevA][prev2A]
					if v == 0 {
						continue
					}
					dpNext[s][0][prevA] = add(dpNext[s][0][prevA], v)
					if pos < n {
						dpNext[s+1][1][prevA] = add(dpNext[s+1][1][prevA], v)
					}
					if pos > 1 && prev2A == 0 {
						dpNext[s+1][0][prevA] = add(dpNext[s+1][0][prevA], v)
					}
				}
			}
		}
		dpCurr = dpNext
	}
	dpEvents := make([]int, n+1)
	for s := 0; s <= n; s++ {
		total := 0
		for prevA := 0; prevA < 2; prevA++ {
			for prev2A := 0; prev2A < 2; prev2A++ {
				total = add(total, dpCurr[s][prevA][prev2A])
			}
		}
		dpEvents[s] = total
	}
	fact := make([]int, n+1)
	invFact := make([]int, n+1)
	fact[0] = 1
	for i := 1; i <= n; i++ {
		fact[i] = mul(fact[i-1], i)
	}
	invFact[n] = powmod(fact[n], MOD-2)
	for i := n; i > 0; i-- {
		invFact[i-1] = mul(invFact[i], i)
	}
	res := 0
	for s := k; s <= n; s++ {
		ds := dpEvents[s]
		if ds == 0 {
			continue
		}
		comb := mul(fact[s], mul(invFact[k], invFact[s-k]))
		term := mul(ds, comb)
		term = mul(term, fact[n-s])
		if (s-k)&1 == 1 {
			res = sub(res, term)
		} else {
			res = add(res, term)
		}
	}
	return fmt.Sprintf("%d", res)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(15) + 1
	k := rng.Intn(n + 1)
	input := fmt.Sprintf("%d %d\n", n, k)
	expect := expected(n, k)
	return input, expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	cases := [][2]int{{1, 0}, {2, 0}, {3, 1}, {4, 2}, {5, 3}}
	for _, c := range cases {
		in := fmt.Sprintf("%d %d\n", c[0], c[1])
		exp := expected(c[0], c[1])
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "edge case failed: %v\ninput:\n%s", err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "edge case failed: expected %s got %s\ninput:\n%s", exp, out, in)
			os.Exit(1)
		}
	}

	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
