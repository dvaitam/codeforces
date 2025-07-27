package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCaseE struct {
	n int
	t int64
	s string
}

func genTestsE() []testCaseE {
	rand.Seed(46)
	tests := make([]testCaseE, 100)
	letters := []byte("abcdefghijklmnopqrstuvwxyz")
	for i := range tests {
		n := rand.Intn(8) + 2
		var sb strings.Builder
		for j := 0; j < n; j++ {
			sb.WriteByte(letters[rand.Intn(26)])
		}
		t := rand.Int63n(2000) - 1000
		tests[i] = testCaseE{n, t, sb.String()}
	}
	tests = append(tests, testCaseE{2, -1, "ba"})
	tests = append(tests, testCaseE{3, -7, "abc"})
	tests = append(tests, testCaseE{7, -475391, "qohshra"})
	return tests
}

func value(ch byte) int64 {
	return 1 << (ch - 'a')
}

func solveE(tc testCaseE) bool {
	vals := make([]int64, tc.n)
	for i := 0; i < tc.n; i++ {
		vals[i] = value(tc.s[i])
	}
	t := tc.t - vals[tc.n-1] + vals[tc.n-2]
	sum := int64(0)
	counts := make([]int64, 26)
	for i := 0; i < tc.n-2; i++ {
		sum += vals[i]
		counts[tc.s[i]-'a']++
	}
	target := t + sum
	if target%2 != 0 {
		return false
	}
	target /= 2
	for b := 25; b >= 0; b-- {
		need := target >> uint(b) & 1
		if need == 1 {
			if counts[b] == 0 {
				return false
			}
			counts[b]--
		}
		if b > 0 {
			counts[b-1] += counts[b] * 2
		}
	}
	return true
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	cmd.Env = append(os.Environ(), "GOMAXPROCS=1")
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("runtime error: %v", err)
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsE()
	for i, tc := range tests {
		input := fmt.Sprintf("%d %d\n%s\n", tc.n, tc.t, tc.s)
		exp := solveE(tc)
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: %v\noutput:\n%s", i+1, err, got)
			os.Exit(1)
		}
		got = strings.TrimSpace(strings.ToLower(got))
		expected := "no"
		if exp {
			expected = "yes"
		}
		if got != expected {
			fmt.Printf("test %d failed: expected %s got %s\ninput:\n%s", i+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
