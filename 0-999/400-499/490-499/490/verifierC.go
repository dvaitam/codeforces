package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type caseC struct {
	s     string
	a     int64
	b     int64
	ok    bool
	left  string
	right string
}

func modStr(s string, mod int64) int64 {
	var res int64
	for _, ch := range s {
		res = (res*10 + int64(ch-'0')) % mod
	}
	return res
}

func solveCaseC(s string, a, b int64) (bool, string, string) {
	n := len(s)
	prefix := make([]int64, n)
	for i := 0; i < n; i++ {
		d := int64(s[i] - '0')
		if i == 0 {
			prefix[i] = d % a
		} else {
			prefix[i] = (prefix[i-1]*10 + d) % a
		}
	}
	suffix := make([]int64, n+1)
	var mult int64 = 1
	for i := n - 1; i >= 0; i-- {
		d := int64(s[i] - '0')
		suffix[i] = (d*mult + suffix[i+1]) % b
		mult = (mult * 10) % b
	}
	for i := 1; i < n; i++ {
		if prefix[i-1] == 0 && suffix[i] == 0 && s[i] != '0' {
			return true, s[:i], s[i:]
		}
	}
	return false, "", ""
}

func generateTests() []caseC {
	r := rand.New(rand.NewSource(44))
	var tests []caseC
	fixed := []struct {
		s    string
		a, b int64
	}{
		{"35", 5, 7},
		{"100", 10, 2},
		{"999", 3, 9},
	}
	for _, f := range fixed {
		ok, l, r := solveCaseC(f.s, f.a, f.b)
		tests = append(tests, caseC{f.s, f.a, f.b, ok, l, r})
	}
	for len(tests) < 120 {
		n := r.Intn(18) + 2
		var sb strings.Builder
		sb.WriteByte(byte('1' + r.Intn(9)))
		for i := 1; i < n; i++ {
			sb.WriteByte(byte('0' + r.Intn(10)))
		}
		s := sb.String()
		a := int64(r.Intn(1000) + 1)
		b := int64(r.Intn(1000) + 1)
		ok, l, rp := solveCaseC(s, a, b)
		tests = append(tests, caseC{s, a, b, ok, l, rp})
	}
	return tests
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func verify(tc caseC, out string) error {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) == 0 {
		return fmt.Errorf("no output")
	}
	if lines[0] == "NO" {
		if tc.ok {
			return fmt.Errorf("should be YES")
		}
		if len(lines) != 1 {
			return fmt.Errorf("extra output after NO")
		}
		return nil
	}
	if lines[0] != "YES" {
		return fmt.Errorf("first line must be YES or NO")
	}
	if len(lines) < 3 {
		return fmt.Errorf("missing parts")
	}
	left := strings.TrimSpace(lines[1])
	right := strings.TrimSpace(lines[2])
	if left+right != tc.s {
		return fmt.Errorf("concatenation mismatch")
	}
	if len(left) > 1 && left[0] == '0' {
		return fmt.Errorf("left leading zero")
	}
	if len(right) > 1 && right[0] == '0' {
		return fmt.Errorf("right leading zero")
	}
	if modStr(left, tc.a) != 0 || modStr(right, tc.b) != 0 {
		return fmt.Errorf("parts not divisible")
	}
	if !tc.ok {
		return fmt.Errorf("should be NO")
	}
	if len(lines) != 3 {
		return fmt.Errorf("extra output")
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		input := fmt.Sprintf("%s\n%d %d\n", tc.s, tc.a, tc.b)
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := verify(tc, out); err != nil {
			fmt.Printf("wrong answer on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
