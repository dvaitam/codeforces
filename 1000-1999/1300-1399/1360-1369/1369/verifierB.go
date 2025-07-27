package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solveB(s string) string {
	n := len(s)
	prefix := 0
	for prefix < n && s[prefix] == '0' {
		prefix++
	}
	suffix := n - 1
	for suffix >= 0 && s[suffix] == '1' {
		suffix--
	}
	if prefix > suffix {
		return s
	}
	return s[:prefix] + "0" + s[suffix+1:]
}

func runBinary(binPath string, input string) (string, error) {
	cmd := exec.Command(binPath)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func randomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		if rand.Intn(2) == 0 {
			b[i] = '0'
		} else {
			b[i] = '1'
		}
	}
	return string(b)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(2)
	const tests = 100
	for t := 0; t < tests; t++ {
		n := rand.Intn(20) + 1
		s := randomString(n)
		input := fmt.Sprintf("1\n%d\n%s\n", n, s)
		expected := solveB(s) + "\n"
		output, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", t+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(output) != strings.TrimSpace(expected) {
			fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", t+1, input, expected, output)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
