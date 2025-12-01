package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type test struct{ input, expected string }

func solveCase(s string) string {
	size := 0
	need := -1
	added := 0
	valid := true
	for i := 0; i < len(s) && valid; i++ {
		switch s[i] {
		case '+':
			size++
			added++
		case '-':
			size--
			if added > 0 {
				added--
			}
			if need != -1 && size < need {
				need = -1
			}
		case '1':
			if need != -1 {
				valid = false
				break
			}
			added = 0
		case '0':
			if size < 2 || added == 0 {
				valid = false
				break
			}
			need = size
		}
	}
	if valid {
		return "YES"
	} else {
		return "NO"
	}
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(44))
	var tests []test
	ops := []byte{'+', '-', '0', '1'}
	for len(tests) < 100 {
		n := rng.Intn(15) + 1
		var sb strings.Builder
		for i := 0; i < n; i++ {
			sb.WriteByte(ops[rng.Intn(4)])
		}
		s := sb.String()
		input := fmt.Sprintf("1\n%s\n", s)
		tests = append(tests, test{input, solveCase(s)})
	}
	return tests
}

func runBinary(bin, input string) (string, error) {
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
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		got, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Printf("Runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(t.expected) {
			fmt.Printf("Wrong answer on test %d\nInput:%sExpected:%s\nGot:%s\n", i+1, t.input, t.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
