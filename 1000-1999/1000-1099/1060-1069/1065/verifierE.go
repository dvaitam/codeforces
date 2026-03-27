package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const MOD = int64(998244353)

func power(base, exp int64) int64 {
	res := int64(1)
	base %= MOD
	for exp > 0 {
		if exp%2 == 1 {
			res = (res * base) % MOD
		}
		base = (base * base) % MOD
		exp /= 2
	}
	return res
}

func oracleSolve(input string) string {
	words := strings.Fields(input)
	idx := 0
	nextInt := func() int64 {
		v := int64(0)
		for _, ch := range words[idx] {
			v = v*10 + int64(ch-'0')
		}
		idx++
		return v
	}
	n := nextInt()
	m := nextInt()
	c := nextInt()

	ans := int64(1)
	prev := int64(0)
	var b int64
	for i := int64(0); i < m; i++ {
		b = nextInt()
		l := b - prev
		term := (power(c, l) + power(c, 2*l)) % MOD
		ans = (ans * term) % MOD
		prev = b
	}

	rem := n - 2*b
	ans = (ans * power(c, rem)) % MOD
	ans = (ans * power(power(2, m), MOD-2)) % MOD

	return fmt.Sprint(ans)
}

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

func generateTests() []string {
	rng := rand.New(rand.NewSource(46))
	tests := make([]string, 0, 100)
	for len(tests) < 100 {
		n := rng.Intn(20) + 2
		m := rng.Intn(n/2) + 1
		A := rng.Intn(10) + 1
		maxVal := n / 2
		if maxVal < 1 {
			maxVal = 1
		}
		vals := make([]int, 0, m)
		last := 0
		ok := true
		for i := 0; i < m; i++ {
			gap := maxVal - last
			if gap <= 0 {
				ok = false
				break
			}
			step := rng.Intn(gap) + 1
			last += step
			vals = append(vals, last)
		}
		if !ok || len(vals) != m {
			continue
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, A))
		for i, v := range vals {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		tests = append(tests, sb.String())
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]
	if cand == "--" && len(os.Args) >= 3 {
		cand = os.Args[2]
	}
	tests := generateTests()
	for i, tc := range tests {
		exp := oracleSolve(tc)
		got, gerr := run(cand, tc)
		if gerr != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\ninput:\n%s", i+1, gerr, tc)
			os.Exit(1)
		}
		if exp != got {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, tc, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
