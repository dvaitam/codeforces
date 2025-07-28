package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCaseC struct {
	s string
	r string
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func expectedC(s, r string) string {
	n := len(s)
	a := 0
	b := 0
	for i := 0; i < n; i++ {
		if s[i] != r[i] {
			a++
		}
		if s[i] != r[n-1-i] {
			b++
		}
	}
	ans1 := 2*a - a%2
	ans2 := 2*b - (1 - b%2)
	if ans2 < 0 {
		ans2 = 1
	}
	return fmt.Sprint(min(ans1, ans2))
}

func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = byte('a' + rand.Intn(3))
	}
	return string(b)
}

func genTestsC() []testCaseC {
	rand.Seed(3)
	tests := make([]testCaseC, 0, 100)
	for len(tests) < 100 {
		n := rand.Intn(10) + 1
		s := randString(n)
		r := randString(n)
		tests = append(tests, testCaseC{s: s, r: r})
	}
	return tests
}

func runCase(bin string, tc testCaseC) error {
	input := fmt.Sprintf("1\n%d\n%s\n%s\n", len(tc.s), tc.s, tc.r)
	cmd := exec.Command(bin)
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	expect := expectedC(tc.s, tc.r)
	if got != expect {
		return fmt.Errorf("s=%s r=%s expected %s got %s", tc.s, tc.r, expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsC()
	for i, tc := range tests {
		if err := runCase(bin, tc); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
