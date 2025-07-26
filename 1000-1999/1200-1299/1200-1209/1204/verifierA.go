package main

import (
	"bytes"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	s string
}

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func solveA(s string) string {
	n := new(big.Int)
	n.SetString(s, 2)
	pow := big.NewInt(1)
	cnt := 0
	for pow.Cmp(n) < 0 {
		cnt++
		pow.Lsh(pow, 2)
	}
	return fmt.Sprintf("%d", cnt)
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(1))
	tests := make([]testCase, 0, 100)
	fixed := []string{"0", "1", "10", "11", "101", "1111", "100000", "101010", "11111111", "100000000"}
	for _, f := range fixed {
		tests = append(tests, testCase{s: f})
	}
	for len(tests) < 100 {
		length := rng.Intn(100) + 1
		if length == 1 {
			if rng.Intn(2) == 0 {
				tests = append(tests, testCase{s: "0"})
				continue
			}
		}
		sb := make([]byte, length)
		for i := 0; i < length; i++ {
			b := byte('0' + rng.Intn(2))
			if i == 0 && b == '0' && length > 1 {
				b = '1'
			}
			sb[i] = b
		}
		tests = append(tests, testCase{s: string(sb)})
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		input := fmt.Sprintf("%s\n", t.s)
		expect := solveA(t.s)
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: execution failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != expect {
			fmt.Printf("test %d failed: expected %q got %q\n", i+1, expect, strings.TrimSpace(out))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
