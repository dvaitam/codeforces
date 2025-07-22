package main

import (
	"bytes"
	"fmt"
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
	arr := make([]int, n)
	sum := 0
	for i := 0; i < n; i++ {
		fmt.Fscan(r, &arr[i])
		sum += arr[i]
	}
	cnt := 0
	total := n + 1
	for k := 1; k <= 5; k++ {
		if (sum+k)%total != 1 {
			cnt++
		}
	}
	return fmt.Sprintf("%d\n", cnt)
}

func generateTests() []test {
	rand.Seed(42)
	var tests []test
	fixed := []struct {
		n   int
		arr []int
	}{
		{1, []int{1}},
		{2, []int{1, 2}},
		{3, []int{5, 5, 5}},
		{4, []int{1, 1, 1, 1}},
		{5, []int{1, 2, 3, 4, 5}},
	}
	for _, f := range fixed {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", f.n))
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
		n := rand.Intn(100) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", rand.Intn(5)+1))
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
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
