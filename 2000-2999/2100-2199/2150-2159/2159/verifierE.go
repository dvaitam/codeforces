package main

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const mod int64 = 1000000007

type query struct {
	n int
	k int
}

const (
	nMax      = 300000
	kMax      = 600000
	blockSize = 1024
)

type solver struct {
	a      int64
	b      int64
	c      int64
	invC   int64
	powC   []int64
	invK   []int64
	small  [][]uint32
	blocks [][]uint32
}

func modPow(base, exp int64) int64 {
	res := int64(1)
	b := base % mod
	e := exp
	for e > 0 {
		if e&1 == 1 {
			res = res * b % mod
		}
		b = b * b % mod
		e >>= 1
	}
	return res
}

func newSolver(a, b, c int64) *solver {
	s := &solver{
		a: a % mod,
		b: b % mod,
		c: c % mod,
	}
	s.invC = modPow(s.c, mod-2)
	s.powC = buildPowC(s.c)
	s.invK = buildInvK()
	s.small = s.buildSmallPolys()
	s.blocks = s.buildBlockPrefixes()
	return s
}

func buildPowC(c int64) []int64 {
	pow := make([]int64, nMax+1)
	pow[0] = 1
	for i := 1; i <= nMax; i++ {
		pow[i] = pow[i-1] * c % mod
	}
	return pow
}

func buildInvK() []int64 {
	inv := make([]int64, kMax+2)
	inv[1] = 1
	for i := 2; i < len(inv); i++ {
		inv[i] = (mod - (mod/int64(i))*inv[int(mod%int64(i))]%mod) % mod
	}
	return inv
}

func (s *solver) nextCoeff(m, k int, prev1, prev2 int64) int64 {
	t1 := int64(m-k+1) % mod
	if t1 < 0 {
		t1 += mod
	}
	t2 := int64(2*m-k+2) % mod
	if t2 < 0 {
		t2 += mod
	}
	val := (t1*s.b%mod*prev1%mod + t2*s.a%mod*prev2%mod) % mod
	val = val * s.invC % mod
	val = val * s.invK[k] % mod
	return val
}

func (s *solver) buildSmallPolys() [][]uint32 {
	size := blockSize
	if size > nMax+1 {
		size = nMax + 1
	}
	res := make([][]uint32, size)
	for m := 0; m < size; m++ {
		deg := 2 * m
		if deg > kMax {
			deg = kMax
		}
		coeff := make([]uint32, deg+1)
		c0 := s.powC[m]
		coeff[0] = uint32(c0)
		if deg >= 1 && m > 0 {
			c1 := int64(m) * s.b % mod * s.powC[m-1] % mod
			coeff[1] = uint32(c1)
			prev2 := c0
			prev1 := c1
			for k := 2; k <= deg; k++ {
				cur := s.nextCoeff(m, k, prev1, prev2)
				coeff[k] = uint32(cur)
				prev2 = prev1
				prev1 = cur
			}
		}
		res[m] = coeff
	}
	return res
}

func (s *solver) buildBlockPrefixes() [][]uint32 {
	maxQ := nMax / blockSize
	res := make([][]uint32, maxQ+1)
	for q := 0; q <= maxQ; q++ {
		m := q * blockSize
		deg := 2 * m
		if deg > kMax {
			deg = kMax
		}
		pref := make([]uint32, deg+1)
		c0 := s.powC[m]
		sum := c0 % mod
		pref[0] = uint32(sum)
		if deg >= 1 && m > 0 {
			c1 := int64(m) * s.b % mod * s.powC[m-1] % mod
			sum = (sum + c1) % mod
			pref[1] = uint32(sum)
			prev2 := c0
			prev1 := c1
			for k := 2; k <= deg; k++ {
				cur := s.nextCoeff(m, k, prev1, prev2)
				sum = (sum + cur) % mod
				pref[k] = uint32(sum)
				prev2 = prev1
				prev1 = cur
			}
		}
		res[q] = pref
	}
	return res
}

func (s *solver) solveQuery(n, k int) int64 {
	q := n / blockSize
	r := n - q*blockSize
	smallPoly := s.small[r]
	limit := len(smallPoly) - 1
	if limit > k {
		limit = k
	}
	pref := s.blocks[q]
	degBig := len(pref) - 1
	totalBig := int64(pref[degBig])
	ans := int64(0)
	for i := 0; i <= limit; i++ {
		t := k - i
		var prefVal int64
		if t >= degBig {
			prefVal = totalBig
		} else {
			prefVal = int64(pref[t])
		}
		contrib := int64(smallPoly[i]) * prefVal % mod
		ans += contrib
		if ans >= mod {
			ans -= mod
		}
	}
	return ans
}

func buildReferenceBinary() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to determine current file path")
	}
	dir := filepath.Dir(file)
	refPath := filepath.Join(dir, "2159E_ref.bin")
	cmd := exec.Command("go", "build", "-o", refPath, "2159E.go")
	cmd.Dir = dir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		_ = os.Remove(refPath)
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	return refPath, nil
}

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		if stderr.Len() > 0 {
			return "", fmt.Errorf("%v: %s", err, strings.TrimSpace(stderr.String()))
		}
		return "", err
	}
	return stdout.String(), nil
}

func parseOutputs(out string, expected int) ([]int64, error) {
	reader := strings.NewReader(out)
	results := make([]int64, 0, expected)
	for len(results) < expected {
		var val int64
		if _, err := fmt.Fscan(reader, &val); err != nil {
			if err == io.EOF {
				return nil, fmt.Errorf("expected %d answers, got %d", expected, len(results))
			}
			return nil, fmt.Errorf("failed to parse output: %v", err)
		}
		val %= mod
		if val < 0 {
			val += mod
		}
		results = append(results, val)
	}
	var extra int64
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return nil, fmt.Errorf("extra output detected (starts with %d)", extra)
	} else if err != io.EOF {
		return nil, fmt.Errorf("failed checking extra output: %v", err)
	}
	return results, nil
}

func parseQueryCount(input string) (int, error) {
	reader := strings.NewReader(input)
	var a, b, c int
	if _, err := fmt.Fscan(reader, &a, &b, &c); err != nil {
		return 0, fmt.Errorf("failed to read coefficients: %v", err)
	}
	var q int
	if _, err := fmt.Fscan(reader, &q); err != nil {
		return 0, fmt.Errorf("failed to read q: %v", err)
	}
	return q, nil
}

func encodeQueries(a, b, c int, actual []query) (string, error) {
	s := newSolver(int64(a), int64(b), int64(c))
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", a, b, c))
	sb.WriteString(fmt.Sprintf("%d\n", len(actual)))
	ansPrev := int64(0)
	for _, q := range actual {
		nPrime := q.n ^ int(ansPrev)
		kPrime := q.k ^ int(ansPrev)
		sb.WriteString(fmt.Sprintf("%d %d\n", nPrime, kPrime))
		ansPrev = s.solveQuery(q.n, q.k)
	}
	return sb.String(), nil
}

func randomQueries(rng *rand.Rand, q int, maxN int) []query {
	res := make([]query, q)
	for i := 0; i < q; i++ {
		n := rng.Intn(maxN + 1)
		kMax := 2 * n
		var k int
		if kMax > 0 {
			k = rng.Intn(kMax + 1)
		}
		res[i] = query{n: n, k: k}
	}
	return res
}

func generateTests() ([]string, error) {
	var tests []string

	rng := rand.New(rand.NewSource(2024))

	// Deterministic small test with multiple queries
	baseQueries := []query{
		{n: 1, k: 0},
		{n: 2, k: 2},
		{n: 5, k: 3},
	}
	input, err := encodeQueries(3, 4, 5, baseQueries)
	if err != nil {
		return nil, err
	}
	tests = append(tests, input)

	// Single small query with n=0
	queries := []query{{n: 0, k: 0}}
	input, err = encodeQueries(5, 7, 11, queries)
	if err != nil {
		return nil, err
	}
	tests = append(tests, input)

	// Random medium test
	actual := randomQueries(rng, 5, 2000)
	input, err = encodeQueries(1234567, 7654321, 1111111, actual)
	if err != nil {
		return nil, err
	}
	tests = append(tests, input)

	// Another random test with larger q
	actual = randomQueries(rng, 8, 5000)
	input, err = encodeQueries(98765432, 13579111, 24680247, actual)
	if err != nil {
		return nil, err
	}
	tests = append(tests, input)

	// Stress test with big n values
	actual = []query{
		{n: 300000, k: 600000},
		{n: 250000, k: 123456},
		{n: 1, k: 1},
	}
	input, err = encodeQueries(999999937, 12345, 67890, actual)
	if err != nil {
		return nil, err
	}
	tests = append(tests, input)

	return tests, nil
}

func compareOutputs(expected, actual []int64) error {
	if len(expected) != len(actual) {
		return fmt.Errorf("expected %d answers, got %d", len(expected), len(actual))
	}
	for i := range expected {
		if expected[i] != actual[i] {
			return fmt.Errorf("mismatch at query %d: expected %d, got %d", i+1, expected[i], actual[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	refPath, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refPath)

	tests, err := generateTests()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to generate tests: %v\n", err)
		os.Exit(1)
	}

	for idx, input := range tests {
		qCount, err := parseQueryCount(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: failed to parse query count: %v\n", idx+1, err)
			os.Exit(1)
		}
		refOut, err := runProgram(refPath, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		expected, err := parseOutputs(refOut, qCount)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}

		userOut, err := runProgram(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", idx+1, err)
			os.Exit(1)
		}
		actual, err := parseOutputs(userOut, qCount)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: invalid output: %v\ninput:\n%s\noutput:\n%s\n", idx+1, err, input, userOut)
			os.Exit(1)
		}
		if err := compareOutputs(expected, actual); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}
