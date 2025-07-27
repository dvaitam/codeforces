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

type testCaseB struct {
	n, m int
	k    int64
	a, b []int
}

func genTests() []testCaseB {
	rng := rand.New(rand.NewSource(1))
	tests := make([]testCaseB, 100)
	for i := range tests {
		n := rng.Intn(6) + 1
		m := rng.Intn(6) + 1
		a := make([]int, n)
		b := make([]int, m)
		for j := range a {
			if rng.Intn(2) == 0 {
				a[j] = 0
			} else {
				a[j] = 1
			}
		}
		for j := range b {
			if rng.Intn(2) == 0 {
				b[j] = 0
			} else {
				b[j] = 1
			}
		}
		k := int64(rng.Intn(n*m+3) + 1)
		tests[i] = testCaseB{n, m, k, a, b}
	}
	// add edge cases
	tests = append(tests,
		testCaseB{1, 1, 1, []int{1}, []int{1}},
		testCaseB{2, 2, 4, []int{1, 1}, []int{1, 1}},
		testCaseB{3, 3, 2, []int{1, 0, 1}, []int{0, 1, 1}},
		testCaseB{4, 1, 3, []int{1, 1, 1, 1}, []int{1}},
		testCaseB{1, 4, 2, []int{1}, []int{1, 1, 1, 1}},
	)
	return tests
}

func getFreq(arr []int) []int {
	n := len(arr)
	freq := make([]int, n+1)
	count := 0
	for _, v := range arr {
		if v == 1 {
			count++
		} else if count > 0 {
			freq[count]++
			count = 0
		}
	}
	if count > 0 {
		freq[count]++
	}
	return freq
}

func buildPrefix(freq []int) ([]int64, []int64) {
	n := len(freq) - 1
	cnt := make([]int64, n+2)
	sum := make([]int64, n+2)
	for i := n; i >= 1; i-- {
		cnt[i] = cnt[i+1] + int64(freq[i])
		sum[i] = sum[i+1] + int64(freq[i]*i)
	}
	return cnt, sum
}

func val(cnt, sum []int64, length int) int64 {
	if length >= len(cnt) || length <= 0 {
		return 0
	}
	return sum[length] - int64(length-1)*cnt[length]
}

func expected(tc testCaseB) int64 {
	freqA := getFreq(tc.a)
	freqB := getFreq(tc.b)
	cntA, sumA := buildPrefix(freqA)
	cntB, sumB := buildPrefix(freqB)
	var ans int64
	k := tc.k
	for p := int64(1); p*p <= k; p++ {
		if k%p != 0 {
			continue
		}
		q := k / p
		if int(p) <= tc.n && int(q) <= tc.m {
			ans += val(cntA, sumA, int(p)) * val(cntB, sumB, int(q))
		}
		if p != q {
			if int(q) <= tc.n && int(p) <= tc.m {
				ans += val(cntA, sumA, int(q)) * val(cntB, sumB, int(p))
			}
		}
	}
	return ans
}

func run(bin string, input string) (string, error) {
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
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func verifyCase(bin string, tc testCaseB) error {
	var input bytes.Buffer
	fmt.Fprintf(&input, "%d %d %d\n", tc.n, tc.m, tc.k)
	for i, v := range tc.a {
		if i > 0 {
			input.WriteByte(' ')
		}
		fmt.Fprint(&input, v)
	}
	input.WriteByte('\n')
	for i, v := range tc.b {
		if i > 0 {
			input.WriteByte(' ')
		}
		fmt.Fprint(&input, v)
	}
	input.WriteByte('\n')
	out, err := run(bin, input.String())
	if err != nil {
		return err
	}
	valOut, err := strconv.ParseInt(strings.TrimSpace(out), 10, 64)
	if err != nil {
		return fmt.Errorf("non-integer output %q", out)
	}
	exp := expected(tc)
	if valOut != exp {
		return fmt.Errorf("expected %d got %d", exp, valOut)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()
	for i, tc := range tests {
		if err := verifyCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput:\n%v\n", i+1, err, tc)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
