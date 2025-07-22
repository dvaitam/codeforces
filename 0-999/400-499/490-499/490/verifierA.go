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

type caseA struct {
	n   int
	arr []int
}

func generateTests() []caseA {
	r := rand.New(rand.NewSource(42))
	var tests []caseA
	tests = append(tests, caseA{3, []int{1, 2, 3}})
	tests = append(tests, caseA{1, []int{1}})
	tests = append(tests, caseA{4, []int{2, 2, 2, 2}})
	for len(tests) < 120 {
		n := r.Intn(20) + 1
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			arr[i] = r.Intn(3) + 1
		}
		tests = append(tests, caseA{n, arr})
	}
	return tests
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func verify(t caseA, out string) error {
	parts := strings.Fields(out)
	if len(parts) == 0 {
		return fmt.Errorf("no output")
	}
	w, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("invalid integer: %v", parts[0])
	}
	c1, c2, c3 := 0, 0, 0
	for _, v := range t.arr {
		switch v {
		case 1:
			c1++
		case 2:
			c2++
		case 3:
			c3++
		}
	}
	expected := c1
	if c2 < expected {
		expected = c2
	}
	if c3 < expected {
		expected = c3
	}
	if w != expected {
		return fmt.Errorf("expected %d teams, got %d", expected, w)
	}
	used := make([]bool, t.n+1)
	idx := 1
	for i := 0; i < w; i++ {
		if idx+2 >= len(parts) {
			return fmt.Errorf("missing team line %d", i+1)
		}
		a, err1 := strconv.Atoi(parts[idx])
		b, err2 := strconv.Atoi(parts[idx+1])
		c, err3 := strconv.Atoi(parts[idx+2])
		if err1 != nil || err2 != nil || err3 != nil {
			return fmt.Errorf("invalid team indices")
		}
		idx += 3
		vals := []int{a, b, c}
		seen := make(map[int]bool)
		types := make(map[int]bool)
		for _, v := range vals {
			if v < 1 || v > t.n {
				return fmt.Errorf("index out of range")
			}
			if used[v] {
				return fmt.Errorf("index used twice")
			}
			used[v] = true
			seen[v] = true
			types[t.arr[v-1]] = true
		}
		if len(types) != 3 {
			return fmt.Errorf("team %d does not contain all skills", i+1)
		}
	}
	if idx != len(parts) {
		return fmt.Errorf("extra output")
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		var sb strings.Builder
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte('\n')
		for j, v := range tc.arr {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		out, err := run(bin, sb.String())
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := verify(tc, out); err != nil {
			fmt.Printf("wrong answer on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
