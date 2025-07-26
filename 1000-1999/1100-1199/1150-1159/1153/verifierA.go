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
	input    string
	expected string
}

func generateTests() []Test {
	rng := rand.New(rand.NewSource(42))
	tests := make([]Test, 0, 100)
	for i := 0; i < 100; i++ {
		n := rng.Intn(100) + 1
		t := rng.Intn(100000) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, t))
		starts := make([]int, n)
		d := make([]int, n)
		for j := 0; j < n; j++ {
			starts[j] = rng.Intn(100000) + 1
			d[j] = rng.Intn(100000) + 1
			sb.WriteString(fmt.Sprintf("%d %d\n", starts[j], d[j]))
		}
		input := sb.String()
		bestTime := int64(1 << 62)
		bestIndex := 1
		for j := 0; j < n; j++ {
			arrival := int64(starts[j])
			if arrival < int64(t) {
				delta := int64(t) - arrival
				k := (delta + int64(d[j]) - 1) / int64(d[j])
				arrival += k * int64(d[j])
			}
			if arrival < bestTime {
				bestTime = arrival
				bestIndex = j + 1
			}
		}
		tests = append(tests, Test{input: input, expected: strconv.Itoa(bestIndex)})
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		got, err := run(bin, tc.input)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != tc.expected {
			fmt.Printf("wrong answer on test %d\ninput:\n%sexpected: %s\ngot: %s\n", i+1, tc.input, tc.expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
