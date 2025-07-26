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
	input  string
	output string
}

func solveC(n, k int, arr []int) string {
	mapping := make([]int, 256)
	for i := range mapping {
		mapping[i] = -1
	}
	res := make([]string, n)
	for i := 0; i < n; i++ {
		x := arr[i]
		if mapping[x] == -1 {
			start := x
			for start > 0 && x-(start-1)+1 <= k && mapping[start-1] == -1 {
				start--
			}
			if start > 0 && mapping[start-1] != -1 && x-mapping[start-1]+1 <= k {
				start = mapping[start-1]
			}
			for y := start; y <= x; y++ {
				mapping[y] = start
			}
		}
		res[i] = strconv.Itoa(mapping[x])
	}
	return strings.Join(res, " ") + "\n"
}

func generateTests() []Test {
	rand.Seed(42)
	tests := make([]Test, 0, 100)
	for len(tests) < 100 {
		n := rand.Intn(20) + 1
		k := rand.Intn(10) + 1
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			arr[i] = rand.Intn(256)
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(arr[i]))
		}
		sb.WriteByte('\n')
		input := sb.String()
		output := solveC(n, k, arr)
		tests = append(tests, Test{input: input, output: output})
	}
	return tests
}

func run(binary string, input string) (string, error) {
	cmd := exec.Command(binary)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, stderr.String())
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		got, err := run(binary, t.input)
		if err != nil {
			fmt.Printf("Test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(t.output) {
			fmt.Printf("Test %d failed. Input: %q\nExpected: %qGot: %q\n", i+1, t.input, t.output, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
