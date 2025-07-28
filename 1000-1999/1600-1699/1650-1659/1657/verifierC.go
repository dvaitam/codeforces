package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const base uint64 = 911382323

type testCaseC struct {
	s string
}

func genTestsC() []testCaseC {
	rng := rand.New(rand.NewSource(44))
	tests := []testCaseC{
		{"()"}, {"(("}, {"())("}, {"(()())"},
	}
	letters := []byte{'(', ')'}
	for len(tests) < 100 {
		n := rng.Intn(12) + 1
		b := make([]byte, n)
		for i := 0; i < n; i++ {
			b[i] = letters[rng.Intn(2)]
		}
		tests = append(tests, testCaseC{string(b)})
	}
	return tests
}

func solveC(tc testCaseC) string {
	b := []byte(tc.s)
	n := len(b)
	pow := make([]uint64, n+1)
	pref := make([]uint64, n+1)
	pow[0] = 1
	for i := 0; i < n; i++ {
		pref[i+1] = pref[i]*base + uint64(b[i])
		pow[i+1] = pow[i] * base
	}
	rev := make([]byte, n)
	for i := 0; i < n; i++ {
		rev[i] = b[n-1-i]
	}
	rpref := make([]uint64, n+1)
	for i := 0; i < n; i++ {
		rpref[i+1] = rpref[i]*base + uint64(rev[i])
	}
	isPal := func(l, r int) bool {
		hf := pref[r+1] - pref[l]*pow[r-l+1]
		hr := rpref[n-l] - rpref[n-1-r]*pow[r-l+1]
		return hf == hr
	}
	i := 0
	ops := 0
	for i < n {
		bal := 0
		minBal := 0
		found := false
		for j := i; j < n; j++ {
			if b[j] == '(' {
				bal++
			} else {
				bal--
			}
			if bal < minBal {
				minBal = bal
			}
			if bal == 0 && minBal >= 0 {
				ops++
				i = j + 1
				found = true
				break
			}
			if j-i+1 >= 2 && isPal(i, j) {
				ops++
				i = j + 1
				found = true
				break
			}
		}
		if !found {
			break
		}
	}
	return fmt.Sprintf("%d %d", ops, n-i)
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("runtime error: %v", err)
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	tests := genTestsC()
	for i, tc := range tests {
		input := fmt.Sprintf("%d\n%s\n", len(tc.s), tc.s)
		exp := solveC(tc)
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: %v\noutput:\n%s", i+1, err, got)
			os.Exit(1)
		}
		got = strings.TrimSpace(got)
		if got != exp {
			fmt.Printf("test %d failed: expected %q got %q\ninput:\n%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
