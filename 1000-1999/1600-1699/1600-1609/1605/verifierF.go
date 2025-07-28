package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type testCaseF struct {
	n int
	k int
	m int
}

func generateTestsF() []testCaseF {
	r := rand.New(rand.NewSource(1))
	primes := []int{100000007, 100000037, 100000039, 100000049}
	tests := []testCaseF{}
	for len(tests) < 120 {
		n := r.Intn(4) + 1
		k := r.Intn(3) + 1
		m := primes[r.Intn(len(primes))]
		tests = append(tests, testCaseF{n, k, m})
	}
	return tests
}

func nextPerm(a []int) bool {
	i := len(a) - 2
	for i >= 0 && a[i] >= a[i+1] {
		i--
	}
	if i < 0 {
		return false
	}
	j := len(a) - 1
	for a[j] <= a[i] {
		j--
	}
	a[i], a[j] = a[j], a[i]
	for l, r := i+1, len(a)-1; l < r; l, r = l+1, r-1 {
		a[l], a[r] = a[r], a[l]
	}
	return true
}

func isPalindORme(arr []int) bool {
	n := len(arr)
	pre := 0
	suf := 0
	for i := 0; i < n; i++ {
		pre |= arr[i]
		suf |= arr[n-1-i]
		if pre != suf {
			return false
		}
	}
	return true
}

func isGood(arr []int) bool {
	sort.Ints(arr)
	tmp := make([]int, len(arr))
	copy(tmp, arr)
	for {
		if isPalindORme(tmp) {
			return true
		}
		if !nextPerm(tmp) {
			break
		}
	}
	return false
}

func bruteForce(n, k, m int) int {
	limit := 1 << k
	arr := make([]int, n)
	count := 0
	var dfs func(int)
	dfs = func(pos int) {
		if pos == n {
			tmp := make([]int, n)
			copy(tmp, arr)
			if isGood(tmp) {
				count++
			}
			return
		}
		for v := 0; v < limit; v++ {
			arr[pos] = v
			dfs(pos + 1)
		}
	}
	dfs(0)
	return count % m
}

func runCandidate(bin, input string) (string, error) {
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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTestsF()
	for i, tc := range tests {
		input := fmt.Sprintf("%d %d %d\n", tc.n, tc.k, tc.m)
		exp := fmt.Sprintf("%d", bruteForce(tc.n, tc.k, tc.m))
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != exp {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\ninput:%sexpected:%s\ngot:%s\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
