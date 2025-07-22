package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type test struct {
	input    string
	expected string
}

func solve(input string) string {
	r := strings.NewReader(strings.TrimSpace(input))
	var n int
	fmt.Fscan(r, &n)
	counts := make([]int64, 33)
	for i := 0; i < n; i++ {
		var x uint
		fmt.Fscan(r, &x)
		counts[bits.OnesCount(x)]++
	}
	var ans int64
	for _, c := range counts {
		ans += c * (c - 1) / 2
	}
	return fmt.Sprintf("%d\n", ans)
}

func generateTests() []test {
	rand.Seed(43)
	var tests []test
	fixed := []struct {
		arr []uint
	}{
		{[]uint{1, 2, 3}},
		{[]uint{1, 1}},
		{[]uint{0}},
		{[]uint{7, 7, 7}},
		{[]uint{5, 10, 5, 10}},
	}
	for _, f := range fixed {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", len(f.arr)))
		for i, v := range f.arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		inp := sb.String()
		tests = append(tests, test{inp, solve(inp)})
	}
	for len(tests) < 100 {
		n := rand.Intn(20) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", rand.Uint32()%1000))
		}
		sb.WriteByte('\n')
		inp := sb.String()
		tests = append(tests, test{inp, solve(inp)})
	}
	return tests
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
		if got != strings.TrimSpace(t.expected) {
			fmt.Printf("Wrong answer on test %d\nInput:\n%sExpected:%sGot:%s\n", i+1, t.input, t.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
