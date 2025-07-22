package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func parseBase(s string, base int) int {
	v := 0
	for _, c := range s {
		var d int
		if c >= '0' && c <= '9' {
			d = int(c - '0')
		} else {
			d = int(c-'A') + 10
		}
		v = v*base + d
	}
	return v
}

func toBase(n, base int) string {
	if n == 0 {
		return "0"
	}
	var digits []rune
	for n > 0 {
		rem := n % base
		var c rune
		if rem < 10 {
			c = rune('0' + rem)
		} else {
			c = rune('A' + rem - 10)
		}
		digits = append(digits, c)
		n /= base
	}
	for i, j := 0, len(digits)-1; i < j; i, j = i+1, j-1 {
		digits[i], digits[j] = digits[j], digits[i]
	}
	return string(digits)
}

type test struct {
	n     int
	radix int
	nums  []string
}

func randNum(radix int) string {
	length := rand.Intn(5) + 1
	digits := make([]rune, length)
	for i := range digits {
		d := rand.Intn(radix)
		if d < 10 {
			digits[i] = rune('0' + d)
		} else {
			digits[i] = rune('A' + d - 10)
		}
	}
	return string(digits)
}

func generateTests() []test {
	rand.Seed(42)
	var tests []test
	for len(tests) < 100 {
		n := rand.Intn(10) + 1
		radix := rand.Intn(35) + 2
		nums := make([]string, n)
		for i := range nums {
			nums[i] = randNum(radix)
		}
		tests = append(tests, test{n: n, radix: radix, nums: nums})
	}
	return tests
}

func expected(t test) string {
	sum := 0
	for _, s := range t.nums {
		sum += parseBase(s, t.radix)
	}
	return toBase(sum, t.radix)
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n%d\n", t.n, t.radix))
		for _, s := range t.nums {
			sb.WriteString(fmt.Sprintf("%s\n", s))
		}
		input := sb.String()
		exp := expected(t)
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("Runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if out != exp {
			fmt.Printf("Wrong answer on test %d\nInput:\n%sExpected:%s\nGot:%s\n", i+1, input, exp, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
