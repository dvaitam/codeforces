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

func reachableBoundaries(s string) []int {
	n := len(s)
	leftOne := make([]int, n)
	last := 0
	for i := 0; i < n; i++ {
		if s[i] == '1' {
			last = i + 1
		}
		leftOne[i] = last
	}
	const inf = int(1e9)
	minLeft := make([]int, n+1)
	for i := range minLeft {
		minLeft[i] = inf
	}
	for i := n - 1; i >= 0; i-- {
		minLeft[i] = minLeft[i+1]
		if s[i] == '0' {
			if leftOne[i] < minLeft[i] {
				minLeft[i] = leftOne[i]
			}
		}
	}
	res := make([]int, 0)
	for j := 0; j < n-1; j++ {
		if s[j] == '0' && minLeft[j+1] > j+1 {
			res = append(res, j+1)
		}
	}
	return res
}

func canEqual(a, b string) bool {
	ba := reachableBoundaries(a)
	bb := reachableBoundaries(b)
	i, j := 0, 0
	for i < len(ba) && j < len(bb) {
		if ba[i] == bb[j] {
			return true
		}
		if ba[i] < bb[j] {
			i++
		} else {
			j++
		}
	}
	return false
}

func solveCase(a, b string) string {
	if canEqual(a, b) {
		return "YES"
	} else {
		return "NO"
	}
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(43))
	var tests []test
	for len(tests) < 100 {
		n := rng.Intn(9) + 2 // length 2..10
		a := make([]byte, n)
		b := make([]byte, n)
		a[0] = '0'
		b[0] = '0'
		a[n-1] = '1'
		b[n-1] = '1'
		for i := 1; i < n-1; i++ {
			a[i] = byte('0' + rng.Intn(2))
			b[i] = byte('0' + rng.Intn(2))
		}
		sa := string(a)
		sb := string(b)
		input := fmt.Sprintf("1\n%s\n%s\n", sa, sb)
		tests = append(tests, test{input, solveCase(sa, sb)})
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
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
