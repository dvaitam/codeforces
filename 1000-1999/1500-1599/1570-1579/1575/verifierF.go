package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const mod = 1000000007

func modpow(a, b, m int64) int64 {
	a %= m
	if a < 0 {
		a += m
	}
	res := int64(1)
	for b > 0 {
		if b&1 == 1 {
			res = res * a % m
		}
		a = a * a % m
		b >>= 1
	}
	return res
}

func modinv(a, m int64) int64 {
	return modpow(a, m-2, m)
}

// modBruteForceF computes f(b) mod p using Gaussian elimination on the Markov chain.
// Only works for small state spaces (k^n <= ~2000).
func modBruteForceF(b []int, k int) int64 {
	n := len(b)
	if n == 0 {
		return 0
	}
	allEq := true
	for i := 1; i < n; i++ {
		if b[i] != b[0] {
			allEq = false
			break
		}
	}
	if allEq {
		return 0
	}

	totalStates := 1
	for i := 0; i < n; i++ {
		totalStates *= k
	}

	encode := func(arr []int) int {
		v := 0
		for i := 0; i < n; i++ {
			v = v*k + arr[i]
		}
		return v
	}

	decode := func(v int) []int {
		arr := make([]int, n)
		for i := n - 1; i >= 0; i-- {
			arr[i] = v % k
			v /= k
		}
		return arr
	}

	isAbsorbing := func(v int) bool {
		arr := decode(v)
		for i := 1; i < n; i++ {
			if arr[i] != arr[0] {
				return false
			}
		}
		return true
	}

	var nonAbs []int
	stateIdx := make(map[int]int)
	for s := 0; s < totalStates; s++ {
		if !isAbsorbing(s) {
			stateIdx[s] = len(nonAbs)
			nonAbs = append(nonAbs, s)
		}
	}
	sz := len(nonAbs)
	if sz == 0 {
		return 0
	}

	p := int64(mod)
	invNK := modinv(int64(n)*int64(k), p)

	A := make([][]int64, sz)
	rhs := make([]int64, sz)
	for r := 0; r < sz; r++ {
		A[r] = make([]int64, sz)
		A[r][r] = 1
		rhs[r] = 1

		s := nonAbs[r]
		arr := decode(s)
		for i := 0; i < n; i++ {
			for j := 0; j < k; j++ {
				old := arr[i]
				arr[i] = j
				ns := encode(arr)
				arr[i] = old
				if !isAbsorbing(ns) {
					col := stateIdx[ns]
					A[r][col] = (A[r][col] - invNK + p) % p
				}
			}
		}
	}

	for col := 0; col < sz; col++ {
		pivot := -1
		for r := col; r < sz; r++ {
			if A[r][col] != 0 {
				pivot = r
				break
			}
		}
		if pivot == -1 {
			continue
		}
		A[col], A[pivot] = A[pivot], A[col]
		rhs[col], rhs[pivot] = rhs[pivot], rhs[col]

		inv := modinv(A[col][col], p)
		for j := col; j < sz; j++ {
			A[col][j] = A[col][j] * inv % p
		}
		rhs[col] = rhs[col] * inv % p

		for r := 0; r < sz; r++ {
			if r == col || A[r][col] == 0 {
				continue
			}
			factor := A[r][col]
			for j := col; j < sz; j++ {
				A[r][j] = (A[r][j] - factor*A[col][j]%p + p) % p
			}
			rhs[r] = (rhs[r] - factor*rhs[col]%p + p) % p
		}
	}

	ef := make(map[int]int64)
	for r := 0; r < sz; r++ {
		ef[nonAbs[r]] = rhs[r]
	}

	startState := encode(b)
	if isAbsorbing(startState) {
		return 0
	}
	return ef[startState]
}

// modBruteForceExpected computes E[f(a)] mod p by enumerating all replacements of -1 entries.
func modBruteForceExpected(a []int, k int) int64 {
	n := len(a)
	var freePos []int
	for i := 0; i < n; i++ {
		if a[i] == -1 {
			freePos = append(freePos, i)
		}
	}
	numFree := len(freePos)
	totalCombinations := 1
	for i := 0; i < numFree; i++ {
		totalCombinations *= k
	}

	p := int64(mod)
	sum := int64(0)
	arr := make([]int, n)
	copy(arr, a)
	for combo := 0; combo < totalCombinations; combo++ {
		v := combo
		for _, pos := range freePos {
			arr[pos] = v % k
			v /= k
		}
		sum = (sum + modBruteForceF(arr, k)) % p
	}
	invTotal := modinv(int64(totalCombinations), p)
	return sum * invTotal % p
}

func runBinary(bin, input string) (string, error) {
	if !strings.Contains(bin, "/") {
		bin = "./" + bin
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("timeout")
	}
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, out)
	}
	return strings.TrimSpace(string(out)), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: go run verifierF.go /path/to/binary\n")
		os.Exit(1)
	}
	cand := os.Args[1]

	type testCase struct {
		input    string
		expected int64
	}

	var tests []testCase

	// k=1: answer is always 0
	tests = append(tests, testCase{"1 1\n0 \n", 0})
	tests = append(tests, testCase{"2 1\n0 0 \n", 0})
	tests = append(tests, testCase{"1 1\n-1 \n", 0})

	// Generate small brute-forceable cases: keep k^n <= 500 and numFree*k small
	rng := rand.New(rand.NewSource(42))
	for i := 0; i < 15; i++ {
		// n in [1,3], k in [2,6], but ensure k^n <= 500
		n := rng.Intn(3) + 1
		maxK := 6
		if n == 3 {
			maxK = 4 // 4^3=64
		}
		k := rng.Intn(maxK-1) + 2

		a := make([]int, n)
		numFree := 0
		for j := 0; j < n; j++ {
			if rng.Intn(3) == 0 && numFree < 2 {
				a[j] = -1
				numFree++
			} else {
				a[j] = rng.Intn(k)
			}
		}

		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
		for j := 0; j < n; j++ {
			sb.WriteString(fmt.Sprintf("%d ", a[j]))
		}
		sb.WriteByte('\n')
		input := sb.String()

		exp := modBruteForceExpected(a, k)
		tests = append(tests, testCase{input, exp})
	}

	for i, tc := range tests {
		got, err := runBinary(cand, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed to run: %v\n", i+1, err)
			os.Exit(1)
		}
		gotVal, err := strconv.ParseInt(strings.TrimSpace(got), 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: invalid output %q\n", i+1, got)
			os.Exit(1)
		}
		if gotVal != tc.expected {
			fmt.Fprintf(os.Stderr, "test %d failed:\ninput:\n%sexpected: %d\ngot: %d\n", i+1, tc.input, tc.expected, gotVal)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
