package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const mod = 998244353

var fact []int
var invFact []int

func modPow(a, e int) int {
	res := 1
	for e > 0 {
		if e&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		e >>= 1
	}
	return res
}

func precomputeFactorials(limit int) {
	fact = make([]int, limit+1)
	invFact = make([]int, limit+1)
	fact[0] = 1
	for i := 1; i <= limit; i++ {
		fact[i] = fact[i-1] * i % mod
	}
	invFact[limit] = modPow(fact[limit], mod-2)
	for i := limit; i >= 1; i-- {
		invFact[i-1] = invFact[i] * i % mod
	}
}

func comb(n, k int) int {
	if k < 0 || k > n {
		return 0
	}
	return fact[n] * invFact[k] % mod * invFact[n-k] % mod
}

func countForK(arr []int, n, k int) int {
	prefix := make([]int, n+1)
	for i := range prefix {
		prefix[i] = -1
	}
	assign := func(idx, value int) bool {
		if value < 0 || value > k {
			return false
		}
		if prefix[idx] == -1 {
			prefix[idx] = value
			return true
		}
		return prefix[idx] == value
	}
	if !assign(0, 0) || !assign(n, k) {
		return 0
	}
	for idx := 1; idx <= n; idx++ {
		val := arr[idx-1]
		if val == -1 {
			continue
		}
		if val <= k {
			if !assign(idx-1, val-1) || !assign(idx, val) {
				return 0
			}
		} else {
			required := idx + k - val
			if !assign(idx-1, required) || !assign(idx, required) {
				return 0
			}
		}
	}
	if prefix[0] == -1 || prefix[n] == -1 {
		return 0
	}
	result := 1
	lastIdx := 0
	lastVal := prefix[0]
	for i := 1; i <= n; i++ {
		if prefix[i] == -1 {
			continue
		}
		length := i - lastIdx
		delta := prefix[i] - lastVal
		if delta < 0 || delta > length {
			return 0
		}
		result = result * comb(length, delta) % mod
		lastIdx = i
		lastVal = prefix[i]
	}
	if lastIdx != n {
		return 0
	}
	return result
}

func solveRef(input string) (string, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return "", err
	}
	results := make([]string, t)
	for tc := 0; tc < t; tc++ {
		var n int
		fmt.Fscan(reader, &n)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
		}
		identityPossible := true
		for i := 0; i < n; i++ {
			if arr[i] != -1 && arr[i] != i+1 {
				identityPossible = false
			}
		}
		total := 0
		for k := 1; k <= n-1; k++ {
			total += countForK(arr, n, k)
			if total >= mod {
				total -= mod
			}
		}
		if identityPossible {
			deduction := n - 2
			if deduction < 0 {
				deduction = 0
			}
			total = (total - deduction) % mod
			if total < 0 {
				total += mod
			}
		}
		results[tc] = fmt.Sprintf("%d", total)
	}
	return strings.Join(results, "\n"), nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
	}
	if stderr.Len() > 0 {
		return "", fmt.Errorf("unexpected stderr output: %s", stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

type testCase struct {
	name  string
	input string
}

func makeCase(name string, arrays [][]int) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(arrays))
	for _, arr := range arrays {
		fmt.Fprintf(&sb, "%d\n", len(arr))
		for i, v := range arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", v)
		}
		sb.WriteByte('\n')
	}
	return testCase{name: name, input: sb.String()}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	for idx := 0; idx < 30; idx++ {
		tcCount := rng.Intn(3) + 1
		arrays := make([][]int, tcCount)
		for i := 0; i < tcCount; i++ {
			n := rng.Intn(6) + 2
			arr := make([]int, n)
			for j := 0; j < n; j++ {
				if rng.Intn(3) == 0 {
					arr[j] = -1
				} else {
					arr[j] = rng.Intn(n) + 1
				}
			}
			arrays[i] = arr
		}
		tests = append(tests, makeCase(fmt.Sprintf("rand_%d", idx+1), arrays))
	}
	return tests
}

func handcraftedTests() []testCase {
	return []testCase{
		makeCase("example", [][]int{
			{-1, -1, -1, -1, -1},
			{1, 2, 3, 4, 5},
			{-1, -1, -1, 2, -1},
		}),
		makeCase("identity", [][]int{
			{1, 2, 3, 4},
		}),
	}
}

func main() {
	precomputeFactorials(3000)
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierF1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := append(handcraftedTests(), randomTests()...)
	for idx, tc := range tests {
		expect, err := solveRef(tc.input)
		if err != nil {
			fmt.Printf("failed to compute reference for %s: %v\n", tc.name, err)
			os.Exit(1)
		}
		out, err := runCandidate(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d (%s) runtime error: %v\ninput:\n%s\n", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		expLines := strings.Fields(expect)
		gotLines := strings.Fields(out)
		if len(expLines) != len(gotLines) {
			fmt.Printf("test %d (%s) mismatch in outputs count\ninput:\n%s\nexpect:\n%s\nactual:\n%s\n", idx+1, tc.name, tc.input, expect, out)
			os.Exit(1)
		}
		for i := range expLines {
			if expLines[i] != gotLines[i] {
				fmt.Printf("test %d (%s) mismatch at output %d\ninput:\n%s\nexpect:\n%s\nactual:\n%s\n", idx+1, tc.name, i+1, tc.input, expect, out)
				os.Exit(1)
			}
		}
	}
	fmt.Println("all tests passed")
}
