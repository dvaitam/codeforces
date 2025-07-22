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

type caseE struct {
	n        int
	patterns []string
	ok       bool
	seq      []int64
}

func dfsFill(pat, lb, cur []byte, pos int, tight bool) bool {
	if pos == len(pat) {
		return true
	}
	if pat[pos] == '?' {
		for d := byte('0'); d <= '9'; d++ {
			if pos == 0 && d == '0' {
				continue
			}
			if tight && d < lb[pos] {
				continue
			}
			cur[pos] = d
			nt := tight && d == lb[pos]
			if dfsFill(pat, lb, cur, pos+1, nt) {
				return true
			}
		}
	} else {
		d := pat[pos]
		if tight && d < lb[pos] {
			return false
		}
		cur[pos] = d
		nt := tight && d == lb[pos]
		if dfsFill(pat, lb, cur, pos+1, nt) {
			return true
		}
	}
	return false
}

func solveCase(patterns []string) (bool, []int64) {
	n := len(patterns)
	res := make([]int64, n)
	var prev int64
	for i, pat := range patterns {
		k := len(pat)
		lb := prev + 1
		lbStr := strconv.FormatInt(lb, 10)
		var cur []byte
		var ok bool
		if k < len(lbStr) {
			ok = false
		} else if k > len(lbStr) {
			cur = make([]byte, k)
			for j := 0; j < k; j++ {
				if pat[j] == '?' {
					if j == 0 {
						cur[j] = '1'
					} else {
						cur[j] = '0'
					}
				} else {
					cur[j] = pat[j]
				}
			}
			ok = true
		} else {
			lbBytes := []byte(lbStr)
			cur = make([]byte, k)
			ok = dfsFill([]byte(pat), lbBytes, cur, 0, true)
		}
		if !ok {
			return false, nil
		}
		x, err := strconv.ParseInt(string(cur), 10, 64)
		if err != nil {
			return false, nil
		}
		res[i] = x
		prev = x
	}
	return true, res
}

func generateTests() []caseE {
	r := rand.New(rand.NewSource(46))
	var tests []caseE
	fixed := [][]string{{"1", "2"}, {"?", "?"}, {"9", "10", "11"}}
	for _, f := range fixed {
		ok, seq := solveCase(f)
		tests = append(tests, caseE{len(f), f, ok, seq})
	}
	for len(tests) < 120 {
		n := r.Intn(5) + 1
		patterns := make([]string, n)
		for i := 0; i < n; i++ {
			l := r.Intn(4) + 1
			var sb strings.Builder
			for j := 0; j < l; j++ {
				if r.Intn(3) == 0 {
					sb.WriteByte('?')
				} else {
					sb.WriteByte(byte('0' + r.Intn(10)))
				}
			}
			if sb.String()[0] == '0' && l > 1 {
				bs := []byte(sb.String())
				bs[0] = '?'
				sb.Reset()
				sb.Write(bs)
			}
			patterns[i] = sb.String()
		}
		ok, seq := solveCase(patterns)
		tests = append(tests, caseE{n, patterns, ok, seq})
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

func matches(pattern string, val int64) bool {
	s := strconv.FormatInt(val, 10)
	if len(s) != len(pattern) {
		return false
	}
	if len(s) > 1 && s[0] == '0' {
		return false
	}
	for i := 0; i < len(s); i++ {
		if pattern[i] != '?' && pattern[i] != s[i] {
			return false
		}
	}
	return true
}

func verify(tc caseE, out string) error {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) == 0 {
		return fmt.Errorf("no output")
	}
	if lines[0] == "NO" {
		if tc.ok {
			return fmt.Errorf("should be YES")
		}
		if len(lines) != 1 {
			return fmt.Errorf("extra output")
		}
		return nil
	}
	if lines[0] != "YES" {
		return fmt.Errorf("first line must be YES or NO")
	}
	if len(lines) != 1+tc.n {
		return fmt.Errorf("expected %d numbers", tc.n)
	}
	seq := make([]int64, tc.n)
	var prev int64 = -1
	for i := 0; i < tc.n; i++ {
		v, err := strconv.ParseInt(strings.TrimSpace(lines[i+1]), 10, 64)
		if err != nil {
			return fmt.Errorf("invalid number")
		}
		if v <= prev {
			return fmt.Errorf("sequence not increasing")
		}
		if !matches(tc.patterns[i], v) {
			return fmt.Errorf("pattern mismatch")
		}
		seq[i] = v
		prev = v
	}
	if !tc.ok {
		return fmt.Errorf("should be NO")
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		var sb strings.Builder
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte('\n')
		for _, p := range tc.patterns {
			sb.WriteString(p)
			sb.WriteByte('\n')
		}
		out, err := run(bin, sb.String())
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
