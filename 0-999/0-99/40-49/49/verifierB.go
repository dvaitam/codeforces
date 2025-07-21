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

type testCase struct {
	a int
	b int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		input := fmt.Sprintf("%d %d\n", t.a, t.b)
		expect := solveB(t.a, t.b)
		out, err := runBinary(bin, input)
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

func solveB(a, b int) string {
	sa := strconv.Itoa(a)
	sb := strconv.Itoa(b)
	la, lb := len(sa), len(sb)
	L := la
	if lb > L {
		L = lb
	}
	ai := make([]int, L)
	bi := make([]int, L)
	for i := 0; i < la; i++ {
		ai[i] = int(sa[la-1-i] - '0')
	}
	for i := 0; i < lb; i++ {
		bi[i] = int(sb[lb-1-i] - '0')
	}
	maxDigit := 0
	sumMax := 0
	for i := 0; i < L; i++ {
		if ai[i] > maxDigit {
			maxDigit = ai[i]
		}
		if bi[i] > maxDigit {
			maxDigit = bi[i]
		}
		s := ai[i] + bi[i]
		if s > sumMax {
			sumMax = s
		}
	}
	minBase := maxDigit + 1
	if minBase < 2 {
		minBase = 2
	}
	ans := L
	for p := minBase; p <= sumMax; p++ {
		carry := 0
		length := 0
		for i := 0; i < L; i++ {
			s := ai[i] + bi[i] + carry
			carry = s / p
			length++
		}
		for carry > 0 {
			carry /= p
			length++
		}
		if length > ans {
			ans = length
		}
	}
	return fmt.Sprintf("%d\n", ans)
}

func generateTests() []testCase {
	rand.Seed(2)
	tests := make([]testCase, 0, 100)
	fixed := []testCase{
		{1, 1}, {9, 9}, {15, 15}, {123, 456}, {999, 1}, {500, 500},
	}
	tests = append(tests, fixed...)
	for len(tests) < 100 {
		a := rand.Intn(1000) + 1
		b := rand.Intn(1000) + 1
		tests = append(tests, testCase{a, b})
	}
	return tests
}
