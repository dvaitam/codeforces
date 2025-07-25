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

type Test struct {
	n, m int
}

func (t Test) Input() string {
	return fmt.Sprintf("1\n%d %d\n", t.n, t.m)
}

func runExe(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func genTests() []Test {
	rand.Seed(1)
	tests := make([]Test, 0, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(8) + 1
		m := rand.Intn(60) + 1
		tests = append(tests, Test{n: n, m: m})
	}
	return tests
}

func possible(n, m int) bool {
	if m <= n-1 {
		return false
	}
	if n%2 == 0 && m%2 == 1 {
		return false
	}
	return true
}

func check(tc Test, out string) error {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return fmt.Errorf("empty output")
	}
	poss := possible(tc.n, tc.m)
	if fields[0] == "No" {
		if poss {
			return fmt.Errorf("expected Yes but got No")
		}
		if len(fields) != 1 {
			return fmt.Errorf("unexpected tokens after No")
		}
		return nil
	}
	if fields[0] != "Yes" {
		return fmt.Errorf("output should start with Yes or No")
	}
	if !poss {
		return fmt.Errorf("expected No but got Yes")
	}
	if len(fields)-1 != tc.n {
		return fmt.Errorf("expected %d numbers, got %d", tc.n, len(fields)-1)
	}
	seq := make([]int, tc.n)
	sum := 0
	for i := 0; i < tc.n; i++ {
		v, err := strconv.Atoi(fields[i+1])
		if err != nil {
			return fmt.Errorf("invalid number")
		}
		if v <= 0 {
			return fmt.Errorf("values must be positive")
		}
		seq[i] = v
		sum += v
	}
	if sum != tc.m {
		return fmt.Errorf("sum mismatch")
	}
	for i := 0; i < tc.n; i++ {
		p := 0
		for j := 0; j < tc.n; j++ {
			if seq[j] < seq[i] {
				p ^= seq[j]
			}
		}
		if p != 0 {
			return fmt.Errorf("xor condition failed at index %d", i)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests := genTests()
	for i, tc := range tests {
		out, err := runExe(bin, tc.Input())
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := check(tc, out); err != nil {
			fmt.Printf("Test %d failed: %v\nInput:%sOutput:%s\n", i+1, err, tc.Input(), out)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
