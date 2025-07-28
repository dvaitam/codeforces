package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCaseA struct {
	n   int
	k   int64
	arr []int64
}

func generateTests() []testCaseA {
	rand.Seed(42)
	tests := make([]testCaseA, 100)
	for i := range tests {
		n := rand.Intn(6) + 2 // 2..7
		arr := make([]int64, n)
		vals := map[int64]bool{}
		for j := 0; j < n; j++ {
			for {
				v := rand.Int63n(201) - 100 // -100..100
				if !vals[v] {
					vals[v] = true
					arr[j] = v
					break
				}
			}
		}
		k := rand.Int63n(201) - 100
		tests[i] = testCaseA{n: n, k: k, arr: arr}
	}
	return tests
}

func solveCase(t testCaseA) string {
	gcd := func(a, b int64) int64 {
		for b != 0 {
			a, b = b, a%b
		}
		if a < 0 {
			return -a
		}
		return a
	}
	d := int64(0)
	for i := 1; i < t.n; i++ {
		diff := t.arr[i] - t.arr[0]
		if diff < 0 {
			diff = -diff
		}
		d = gcd(d, diff)
	}
	if (t.k-t.arr[0])%d == 0 {
		return "YES"
	}
	return "NO"
}

func buildInput(tests []testCaseA) string {
	var b strings.Builder
	fmt.Fprintln(&b, len(tests))
	for _, t := range tests {
		fmt.Fprintf(&b, "%d %d\n", t.n, t.k)
		for i, v := range t.arr {
			if i > 0 {
				b.WriteByte(' ')
			}
			fmt.Fprintf(&b, "%d", v)
		}
		b.WriteByte('\n')
	}
	return b.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	tests := generateTests()
	input := buildInput(tests)

	cmd := exec.Command(binary)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "execution failed:", err)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(&out)
	outputs := []string{}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			outputs = append(outputs, strings.ToUpper(line))
		}
	}
	if len(outputs) != len(tests) {
		fmt.Fprintf(os.Stderr, "expected %d lines of output, got %d\n", len(tests), len(outputs))
		os.Exit(1)
	}
	for i, t := range tests {
		exp := solveCase(t)
		if outputs[i] != exp {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %s got %s\n", i+1, exp, outputs[i])
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
