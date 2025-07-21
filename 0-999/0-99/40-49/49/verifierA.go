package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	q string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		expect := solveA(t.q)
		out, err := runBinary(bin, t.q)
		if err != nil {
			fmt.Printf("test %d: execution failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(expect) {
			fmt.Printf("test %d failed: expected %q got %q\n", i+1, strings.TrimSpace(expect), strings.TrimSpace(out))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func solveA(q string) string {
	q = strings.TrimRight(q, "\r\n")
	var last byte
	for i := len(q) - 1; i >= 0; i-- {
		c := q[i]
		if (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') {
			last = c
			break
		}
	}
	last = byte(strings.ToUpper(string(last))[0])
	vowels := "AEIOUY"
	ans := "NO"
	if strings.ContainsRune(vowels, rune(last)) {
		ans = "YES"
	}
	return ans + "\n"
}

func generateTests() []testCase {
	rand.Seed(1)
	tests := make([]testCase, 0, 100)
	fixed := []string{
		"Q?\n",
		"A?\n",
		"z?\n",
		"Hello?\n",
		"abcde?\n",
		"y?\n",
		"Why??\n",
		"Sly fox?\n",
	}
	for _, f := range fixed {
		tests = append(tests, testCase{q: f})
	}
	for len(tests) < 100 {
		l := rand.Intn(50) + 1
		var sb strings.Builder
		for i := 0; i < l-1; i++ {
			if rand.Intn(10) == 0 {
				sb.WriteByte(' ')
			} else {
				x := rand.Intn(52)
				if x < 26 {
					sb.WriteByte(byte('a' + x))
				} else {
					sb.WriteByte(byte('A' + x - 26))
				}
			}
		}
		sb.WriteByte('?')
		sb.WriteByte('\n')
		tests = append(tests, testCase{q: sb.String()})
	}
	return tests
}
