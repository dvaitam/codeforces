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

func solveCase(n int, h int64, a []int64) string {
	pos := 0
	for i := 1; i < n; i++ {
		if a[i] < a[pos] {
			pos = i
		}
	}
	b := append([]int64(nil), a...)
	sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })
	val1 := b[n-1] + b[n-2] - b[0] - b[1]
	max12 := b[n-1] + b[n-2]
	if b[0]+b[n-1]+h > max12 {
		max12 = b[0] + b[n-1] + h
	}
	min01 := b[1] + b[2]
	if b[0]+b[1]+h < min01 {
		min01 = b[0] + b[1] + h
	}
	val2 := max12 - min01
	if val1 < val2 {
		pos = -1
	}
	res := val1
	if val2 < res {
		res = val2
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", res))
	for i := 0; i < n; i++ {
		if i == pos {
			sb.WriteString("2 ")
		} else {
			sb.WriteString("1 ")
		}
	}
	return strings.TrimSpace(sb.String())
}

type test struct{ input, expected string }

func generateTests() []test {
	rng := rand.New(rand.NewSource(43))
	var tests []test
	// some fixed cases
	fixed := []struct {
		n int
		h int64
		a []int64
	}{
		{2, 1, []int64{1, 2}},
		{3, 0, []int64{1, 2, 3}},
	}
	for _, f := range fixed {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", f.n, f.h))
		for i, v := range f.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		inp := sb.String()
		tests = append(tests, test{inp, solveCase(f.n, f.h, f.a)})
	}
	for len(tests) < 100 {
		n := rng.Intn(8) + 2
		h := int64(rng.Intn(5))
		a := make([]int64, n)
		for i := range a {
			a[i] = int64(rng.Intn(10))
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, h))
		for i, v := range a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		inp := sb.String()
		tests = append(tests, test{inp, solveCase(n, h, a)})
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
