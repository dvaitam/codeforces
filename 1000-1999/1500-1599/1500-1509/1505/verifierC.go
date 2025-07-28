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

type testCase struct {
	input  string
	expect string
}

func solve(s string) string {
	isPal := true
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		if s[i] != s[j] {
			isPal = false
			break
		}
	}
	if isPal {
		return "YES"
	}
	return "NO"
}

func genTests() []testCase {
	r := rand.New(rand.NewSource(1))
	tests := make([]testCase, 100)
	for i := range tests {
		l := r.Intn(10) + 1
		b := make([]byte, l)
		for j := 0; j < l; j++ {
			b[j] = byte('A' + r.Intn(26))
		}
		s := string(b)
		tests[i].input = fmt.Sprintf("%s\n", s)
		tests[i].expect = solve(s)
	}
	return tests
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return out.String() + errBuf.String(), fmt.Errorf("runtime error: %v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	tests := genTests()
	for i, tc := range tests {
		got, err := run(bin, tc.input)
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			fmt.Print(tc.input)
			return
		}
		if got != tc.expect {
			fmt.Printf("case %d failed: expected %s got %s\n", i+1, tc.expect, got)
			fmt.Print(tc.input)
			return
		}
	}
	fmt.Println("All tests passed")
}
