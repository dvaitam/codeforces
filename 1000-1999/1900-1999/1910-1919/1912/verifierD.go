package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type testInput struct {
	text string
}

type pair struct {
	b int64
	n int64
}

type answer struct {
	a int
	k int64
}

type factor struct {
	p int64
	e int64
}

func sieve(n int) []int {
	if n < 2 {
		return nil
	}
	isComp := make([]bool, n+1)
	pr := make([]int, 0)
	for i := 2; i*i <= n; i++ {
		if !isComp[i] {
			for j := i * i; j <= n; j += i {
				isComp[j] = true
			}
		}
	}
	for i := 2; i <= n; i++ {
		if !isComp[i] {
			pr = append(pr, i)
		}
	}
	return pr
}

func factorize(x int64, primes []int) []factor {
	res := make([]factor, 0)
	n := x
	for _, p := range primes {
		pp := int64(p)
		if pp*pp > n {
			break
		}
		if n%pp == 0 {
			cnt := int64(0)
			for n%pp == 0 {
				n /= pp
				cnt++
			}
			res = append(res, factor{p: pp, e: cnt})
		}
	}
	if n > 1 {
		res = append(res, factor{p: n, e: 1})
	}
	return res
}

func phi(n int64, primes []int) int64 {
	res := n
	for _, fe := range factorize(n, primes) {
		res = res / fe.p * (fe.p - 1)
	}
	return res
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func powmod(a, e, mod int64) int64 {
	if mod == 1 {
		return 0
	}
	res := int64(1)
	base := ((a % mod) + mod) % mod
	for e > 0 {
		if e&1 == 1 {
			res = (res * base) % mod
		}
		base = (base * base) % mod
		e >>= 1
	}
	return res
}

func solveOne(b, n int64, primes []int) answer {
	best := answer{a: 0, k: math.MaxInt64}

	fn := factorize(n, primes)
	possible1 := true
	k1 := int64(0)
	for _, fe := range fn {
		cnt := int64(0)
		tb := b
		for tb%fe.p == 0 {
			tb /= fe.p
			cnt++
		}
		if cnt == 0 {
			possible1 = false
			break
		}
		need := (fe.e + cnt - 1) / cnt
		if need > k1 {
			k1 = need
		}
	}
	if possible1 {
		best = answer{a: 1, k: k1}
	}

	if gcd(b, n) == 1 {
		ph := phi(n, primes)
		r := ph
		for _, fe := range factorize(ph, primes) {
			for r%fe.p == 0 && powmod(b, r/fe.p, n) == 1%n {
				r /= fe.p
			}
		}
		if r < best.k || (r == best.k && 2 < best.a) {
			best = answer{a: 2, k: r}
		}
		target := (n - 1) % n
		k3 := int64(-1)
		if n == 2 {
			k3 = r
		} else if r%2 == 0 && powmod(b, r/2, n) == target {
			k3 = r / 2
		}
		if k3 != -1 && (k3 < best.k || (k3 == best.k && 3 < best.a)) {
			best = answer{a: 3, k: k3}
		}
	}

	if best.a == 0 {
		return answer{a: 0, k: -1}
	}
	return best
}

func solveAll(input string) ([]answer, error) {
	in := bufio.NewReader(strings.NewReader(input))
	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return nil, fmt.Errorf("bad input: %w", err)
	}
	items := make([]pair, t)
	mx := int64(2)
	for i := 0; i < t; i++ {
		if _, err := fmt.Fscan(in, &items[i].b, &items[i].n); err != nil {
			return nil, fmt.Errorf("bad input at test %d: %w", i+1, err)
		}
		if items[i].b > mx {
			mx = items[i].b
		}
		if items[i].n > mx {
			mx = items[i].n
		}
	}
	primes := sieve(int(mx))
	ans := make([]answer, t)
	for i := range items {
		ans[i] = solveOne(items[i].b, items[i].n, primes)
	}
	return ans, nil
}

func parseCandidateOutput(out string, t int) ([]answer, error) {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != t {
		return nil, fmt.Errorf("expected %d lines, got %d", t, len(lines))
	}
	ans := make([]answer, t)
	for i, ln := range lines {
		fields := strings.Fields(strings.TrimSpace(ln))
		if len(fields) == 1 {
			if fields[0] != "0" {
				return nil, fmt.Errorf("line %d: single token must be 0", i+1)
			}
			ans[i] = answer{a: 0, k: -1}
			continue
		}
		if len(fields) != 2 {
			return nil, fmt.Errorf("line %d: expected either 0 or 'a k'", i+1)
		}
		a, err := strconv.Atoi(fields[0])
		if err != nil || a < 1 || a > 3 {
			return nil, fmt.Errorf("line %d: invalid a", i+1)
		}
		k, err := strconv.ParseInt(fields[1], 10, 64)
		if err != nil || k < 0 {
			return nil, fmt.Errorf("line %d: invalid k", i+1)
		}
		ans[i] = answer{a: a, k: k}
	}
	return ans, nil
}

func commandForPath(path string) *exec.Cmd {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	case ".js":
		return exec.Command("node", path)
	default:
		return exec.Command(path)
	}
}

func runBinary(path, input string) (string, error) {
	cmd := commandForPath(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("%v\n%s", err, errBuf.String())
	}
	return out.String(), nil
}

func fixedTests() []testInput {
	return []testInput{
		{text: "3\n10 3\n10 11\n10 4\n"},
		{text: "2\n8 2\n7 5\n"},
	}
}

func randomTest(rng *rand.Rand) testInput {
	t := rng.Intn(20) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		b := rng.Intn(2000) + 2
		n := rng.Intn(2000) + 2
		sb.WriteString(fmt.Sprintf("%d %d\n", b, n))
	}
	return testInput{text: sb.String()}
}

func edgeTests() []testInput {
	cases := []struct {
		b int
		n int
	}{
		{2, 2},
		{1000000, 1000000},
		{999983, 999979},
		{1000000, 2},
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, c := range cases {
		sb.WriteString(fmt.Sprintf("%d %d\n", c.b, c.n))
	}
	return []testInput{{text: sb.String()}}
}

func generateTests() []testInput {
	tests := fixedTests()
	tests = append(tests, edgeTests()...)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 60 {
		tests = append(tests, randomTest(rng))
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests := generateTests()
	for idx, input := range tests {
		expected, err := solveAll(input.text)
		if err != nil {
			fmt.Fprintf(os.Stderr, "internal solver failed on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		gotRaw, err := runBinary(bin, input.text)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%s\n", idx+1, err, preview(input.text))
			os.Exit(1)
		}
		got, err := parseCandidateOutput(gotRaw, len(expected))
		if err != nil {
			fmt.Fprintf(os.Stderr, "bad output on test %d: %v\ninput:\n%s\noutput:\n%s\n", idx+1, err, preview(input.text), gotRaw)
			os.Exit(1)
		}
		for i := range expected {
			if got[i] != expected[i] {
				fmt.Fprintf(os.Stderr, "mismatch on test %d case %d\ninput:\n%s\nexpected: %v\ngot: %v\n", idx+1, i+1, preview(input.text), expected[i], got[i])
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func preview(s string) string {
	if len(s) <= 400 {
		return s
	}
	return s[:400] + "...\n"
}
