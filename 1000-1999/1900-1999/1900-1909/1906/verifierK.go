package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const (
	refSource = "1000-1999/1900-1999/1900-1909/1906/1906K.go"
	mod       = "998244353"
)

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierK.go /path/to/binary")
		os.Exit(1)
	}
	candidate := args[0]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	for idx, input := range tests {
		exp, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		got, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		if normalize(exp) != normalize(got) {
			fmt.Fprintf(os.Stderr, "test %d mismatch:\ninput:\n%s\nexpected: %s\ngot: %s", idx+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func buildReference() (string, error) {
	outPath := "./ref_1906K.bin"
	cmd := exec.Command("go", "build", "-o", outPath, refSource)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return outPath, nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstdout:\n%s\nstderr:\n%s", err, stdout.String(), stderr.String())
	}
	return stdout.String(), nil
}

func normalize(s string) string {
	return strings.TrimSpace(s)
}

func buildTests() []string {
	var tests []string

	add := func(arr []int) {
		tests = append(tests, formatTest(arr))
	}

	add([]int{1, 1})
	add([]int{1, 2})
	add([]int{2, 4, 8})
	add([]int{5, 5, 5, 5})
	add([]int{7, 7, 7, 7, 7})
	add([]int{1, 3, 5, 7, 9, 11})

	// Large identical values
	add(repeatValue(100000, 1))
	// Sequential values up to near limit
	add(sequentialValues(99999, 1))
	// Alternating high/low
	add(alternatingValues(50000, 1, 100000))

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// small random tests
	for i := 0; i < 50; i++ {
		n := rng.Intn(7) + 2
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rng.Intn(16)
		}
		add(arr)
	}

	// medium random tests
	for i := 0; i < 30; i++ {
		n := rng.Intn(198) + 2
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rng.Intn(100000) + 1
		}
		add(arr)
	}

	// larger random stress tests
	for i := 0; i < 10; i++ {
		n := rng.Intn(1500) + 500
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			if rng.Intn(5) == 0 {
				arr[j] = 1
			} else {
				arr[j] = rng.Intn(100000) + 1
			}
		}
		add(arr)
	}

	// One more huge random test
	huge := make([]int, 100000)
	for i := range huge {
		if rng.Intn(7) == 0 {
			huge[i] = 1
		} else {
			huge[i] = rng.Intn(100000) + 1
		}
	}
	add(huge)

	return tests
}

func formatTest(arr []int) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(arr)))
	sb.WriteByte('\n')
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func repeatValue(n, val int) []int {
	arr := make([]int, n)
	for i := range arr {
		arr[i] = val
	}
	return arr
}

func sequentialValues(n, start int) []int {
	arr := make([]int, n)
	for i := range arr {
		arr[i] = start + i
	}
	return arr
}

func alternatingValues(n, a, b int) []int {
	arr := make([]int, n)
	for i := range arr {
		if i%2 == 0 {
			arr[i] = a
		} else {
			arr[i] = b
		}
	}
	return arr
}
