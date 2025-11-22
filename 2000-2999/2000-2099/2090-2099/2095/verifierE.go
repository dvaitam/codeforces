package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type test struct {
	input, expected string
}

func prepareBinary(path, tag string) (string, func(), error) {
	if !strings.HasSuffix(path, ".go") {
		return path, func() {}, nil
	}
	bin := filepath.Join(os.TempDir(), tag)
	cmd := exec.Command("go", "build", "-o", bin, path)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", func() {}, fmt.Errorf("build %s: %v\n%s", path, err, out)
	}
	cleanup := func() { os.Remove(bin) }
	return bin, cleanup, nil
}

func runBinary(path, input string, timeout time.Duration) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	timer := time.AfterFunc(timeout, func() { cmd.Process.Kill() })
	err := cmd.Run()
	timer.Stop()
	if err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func formatInput(n int, p, k int64, a []int64) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, p, k))
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func xorCount(n int, p, k int64, a []int64) int64 {
	sq := make([]int64, n)
	for i := 0; i < n; i++ {
		sq[i] = a[i] * a[i]
	}
	var ans int64
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			x := (a[i] ^ a[j]) % p
			y := (sq[i] ^ sq[j]) % p
			if (x*y)%p == k {
				ans++
			}
		}
	}
	return ans
}

func buildExpected(n int, p, k int64, a []int64) string {
	return fmt.Sprintf("%d", xorCount(n, p, k, a))
}

func uniqueNumbers(rng *rand.Rand, n int, p int64) []int64 {
	used := make(map[int64]struct{}, n)
	res := make([]int64, 0, n)
	for len(res) < n {
		v := rng.Int63n(p)
		if _, ok := used[v]; ok {
			continue
		}
		used[v] = struct{}{}
		res = append(res, v)
	}
	return res
}

func randomPrime(rng *rand.Rand) int64 {
	primes := []int64{
		2, 3, 5, 7, 11, 13, 17, 19, 23, 29,
		31, 37, 97, 131, 997, 1229, 50021, 99991, 199999, 499979, 899999, 999983,
	}
	return primes[rng.Intn(len(primes))]
}

func generateTests() []test {
	var tests []test
	// samples
	tests = append(tests, test{
		input:    "3 3 2\n0 1 2\n",
		expected: "1",
	})
	tests = append(tests, test{
		input:    "6 11 2\n1 3 5 6 7 8\n",
		expected: "3",
	})

	// deterministic small cases
	tests = append(tests, func() test {
		n := 2
		p := int64(2)
		a := []int64{0, 1}
		k := int64(0)
		return test{
			input:    formatInput(n, p, k, a),
			expected: buildExpected(n, p, k, a),
		}
	}())

	rng := rand.New(rand.NewSource(2095))
	for i := 0; i < 40; i++ {
		n := rng.Intn(25) + 2
		p := randomPrime(rng)
		k := rng.Int63n(p)
		a := uniqueNumbers(rng, n, p)
		tests = append(tests, test{
			input:    formatInput(n, p, k, a),
			expected: buildExpected(n, p, k, a),
		})
	}

	// larger stress case
	n := 300
	p := int64(999983)
	k := int64(12345 % int(p))
	a := uniqueNumbers(rand.New(rand.NewSource(4242)), n, p)
	tests = append(tests, test{
		input:    formatInput(n, p, k, a),
		expected: buildExpected(n, p, k, a),
	})

	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}

	candBin, candCleanup, err := prepareBinary(os.Args[1], "cand2095E")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer candCleanup()

	_, file, _, _ := runtime.Caller(0)
	refPath := filepath.Join(filepath.Dir(file), "2095E.go")
	refBin, refCleanup, err := prepareBinary(refPath, "ref2095E")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer refCleanup()

	tests := generateTests()
	for i, tc := range tests {
		exp, err := runBinary(refBin, tc.input, 6*time.Second)
		if err != nil {
			fmt.Printf("Reference failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		// cross-check with internal expected to catch ref mismatch early
		if exp != strings.TrimSpace(tc.expected) {
			fmt.Printf("Reference mismatch on test %d\nInput:\n%s\nRef:%s\nCalc:%s\n", i+1, tc.input, exp, tc.expected)
			os.Exit(1)
		}
		got, err := runBinary(candBin, tc.input, 6*time.Second)
		if err != nil {
			fmt.Printf("Candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if exp != got {
			fmt.Printf("Wrong answer on test %d\nInput:\n%s\nExpected:\n%s\nGot:\n%s\n", i+1, tc.input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
