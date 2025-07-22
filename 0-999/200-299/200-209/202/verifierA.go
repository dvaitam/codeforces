package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func generateTests() []string {
	r := rand.New(rand.NewSource(42))
	tests := make([]string, 100)
	for i := 0; i < 100; i++ {
		length := r.Intn(10) + 1 // 1..10
		var sb strings.Builder
		for j := 0; j < length; j++ {
			sb.WriteByte(byte('a' + r.Intn(26)))
		}
		tests[i] = sb.String() + "\n"
	}
	return tests
}

func solve(input string) string {
	s := strings.TrimSpace(input)
	n := len(s)
	best := ""
	for mask := 1; mask < (1 << n); mask++ {
		var t []byte
		for i := 0; i < n; i++ {
			if mask&(1<<i) != 0 {
				t = append(t, s[i])
			}
		}
		ok := true
		for i, j := 0, len(t)-1; i < j; i, j = i+1, j-1 {
			if t[i] != t[j] {
				ok = false
				break
			}
		}
		if ok {
			str := string(t)
			if str > best {
				best = str
			}
		}
	}
	return best
}

func runBinary(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, test := range tests {
		expected := solve(test)
		actual, err := runBinary(bin, test)
		if err != nil {
			fmt.Printf("Test %d: execution failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if actual != expected {
			fmt.Printf("Test %d failed.\nInput: %sExpected: %s\nGot: %s\n", i+1, test, expected, actual)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
