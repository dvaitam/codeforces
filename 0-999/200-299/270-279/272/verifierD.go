package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
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
	xs := make([]int, 2*n)
	for i := 0; i < n; i++ {
		fmt.Fscan(r, &xs[i])
	}
	for i := 0; i < n; i++ {
		fmt.Fscan(r, &xs[n+i])
	}
	var m int
	fmt.Fscan(r, &m)
	sort.Ints(xs)
	fact := make([]int, 2*n+1)
	fact[0] = 1 % m
	for i := 1; i <= 2*n; i++ {
		fact[i] = int((int64(fact[i-1]) * int64(i)) % int64(m))
	}
	ans := 1 % m
	for i := 0; i < 2*n; {
		j := i + 1
		for j < 2*n && xs[j] == xs[i] {
			j++
		}
		cnt := j - i
		ans = int((int64(ans) * int64(fact[cnt])) % int64(m))
		i = j
	}
	return fmt.Sprintf("%d\n", ans)
}

func generateTests() []test {
	rand.Seed(45)
	var tests []test
	fixed := []struct {
		a   []int
		b   []int
		mod int
	}{
		{[]int{1}, []int{1}, 2},
		{[]int{1, 2}, []int{2, 1}, 10},
		{[]int{1, 1}, []int{1, 1}, 7},
	}
	for _, f := range fixed {
		var sb strings.Builder
		n := len(f.a)
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i, v := range f.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteString("\n")
		for i, v := range f.b {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteString("\n")
		sb.WriteString(fmt.Sprintf("%d\n", f.mod))
		inp := sb.String()
		tests = append(tests, test{inp, solve(inp)})
	}
	for len(tests) < 100 {
		n := rand.Intn(5) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		a := make([]int, n)
		b := make([]int, n)
		for i := 0; i < n; i++ {
			a[i] = rand.Intn(5)
			b[i] = rand.Intn(5)
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", a[i]))
		}
		sb.WriteString("\n")
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", b[i]))
		}
		sb.WriteString("\n")
		mod := rand.Intn(100) + 2
		sb.WriteString(fmt.Sprintf("%d\n", mod))
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
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
