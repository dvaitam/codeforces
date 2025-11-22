package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type test struct {
	t      int
	cases  []int
	rawInp string
}

func modMul(a, b, mod uint64) uint64 {
	hi, lo := bits.Mul64(a, b)
	_, rem := bits.Div64(hi, lo, mod)
	return rem
}

func modPow(a, e, mod uint64) uint64 {
	res := uint64(1)
	for e > 0 {
		if e&1 == 1 {
			res = modMul(res, a, mod)
		}
		a = modMul(a, a, mod)
		e >>= 1
	}
	return res
}

func isPrime(n uint64) bool {
	if n < 2 {
		return false
	}
	small := []uint64{2, 3, 5, 7, 11}
	for _, p := range small {
		if n == p {
			return true
		}
		if n%p == 0 {
			return false
		}
	}

	d := n - 1
	s := 0
	for d&1 == 0 {
		d >>= 1
		s++
	}

	witnesses := []uint64{2, 3, 5, 7, 11}
	for _, a := range witnesses {
		if a >= n {
			continue
		}
		x := modPow(a, d, n)
		if x == 1 || x == n-1 {
			continue
		}
		composite := true
		for r := 1; r < s; r++ {
			x = modMul(x, x, n)
			if x == n-1 {
				composite = false
				break
			}
		}
		if composite {
			return false
		}
	}
	return true
}

func parseAndVerify(output string, tcases []int) error {
	reader := bufio.NewReader(strings.NewReader(output))
	caseIdx := 0
	for _, n := range tcases {
		perm := make([]int, n)
		seen := make([]bool, n+1)
		for i := 0; i < n; i++ {
			if _, err := fmt.Fscan(reader, &perm[i]); err != nil {
				return fmt.Errorf("test %d: insufficient numbers (expected %d)", caseIdx+1, n)
			}
			if perm[i] < 1 || perm[i] > n {
				return fmt.Errorf("test %d: value %d out of range [1,%d]", caseIdx+1, perm[i], n)
			}
			if seen[perm[i]] {
				return fmt.Errorf("test %d: duplicate value %d", caseIdx+1, perm[i])
			}
			seen[perm[i]] = true
		}

		req := n/3 - 1
		if req < 0 {
			req = 0
		}
		sum := int64(0)
		primes := 0
		for i := 0; i < n; i++ {
			sum += int64(perm[i])
			ci := (sum + int64(i)) / int64(i+1) // ceil(sum/(i+1))
			if isPrime(uint64(ci)) {
				primes++
			}
		}
		if primes < req {
			return fmt.Errorf("test %d: prime count %d below requirement %d", caseIdx+1, primes, req)
		}
		caseIdx++
	}

	var extra int
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return fmt.Errorf("extra output detected after processing %d test cases", caseIdx)
	}
	return nil
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(2089))
	var tests []test

	samples := [][]int{{2}, {3, 5}, {5}, {10}}
	for _, arr := range samples {
		t := len(arr)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", t)
		for _, v := range arr {
			fmt.Fprintf(&sb, "%d\n", v)
		}
		tests = append(tests, test{t: t, cases: arr, rawInp: sb.String()})
	}

	for len(tests) < 120 {
		t := rng.Intn(5) + 1
		cases := make([]int, t)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", t)
		for i := 0; i < t; i++ {
			n := rng.Intn(100000-2+1) + 2
			cases[i] = n
			fmt.Fprintf(&sb, "%d\n", n)
		}
		tests = append(tests, test{t: t, cases: cases, rawInp: sb.String()})
	}

	return tests
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
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()

	for i, tcase := range tests {
		got, err := runBinary(bin, tcase.rawInp)
		if err != nil {
			fmt.Printf("Runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := parseAndVerify(got, tcase.cases); err != nil {
			fmt.Printf("Wrong answer on test %d: %v\nInput:\n%s\nOutput:\n%s\n", i+1, err, tcase.rawInp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
