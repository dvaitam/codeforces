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

type testCase struct {
	n      int
	arr    []int
	expect int
}

func expectedCount(arr []int) int {
	if len(arr) == 0 {
		return 0
	}
	minVal, maxVal := arr[0], arr[0]
	for _, v := range arr {
		if v < minVal {
			minVal = v
		}
		if v > maxVal {
			maxVal = v
		}
	}
	cnt := 0
	for _, v := range arr {
		if v > minVal && v < maxVal {
			cnt++
		}
	}
	return cnt
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 100)
	// deterministic edge cases
	edges := [][]int{
		{0},
		{1},
		{1, 1},
		{0, 1},
		{5, 1, 5},
		{1, 2, 3, 4, 5},
		{5, 4, 3, 2, 1},
	}
	for _, arr := range edges {
		tests = append(tests, testCase{n: len(arr), arr: append([]int(nil), arr...), expect: expectedCount(arr)})
	}
	for len(tests) < 100 {
		n := rng.Intn(10) + 1
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			arr[i] = rng.Intn(100)
		}
		tests = append(tests, testCase{n: n, arr: arr, expect: expectedCount(arr)})
	}
	return tests
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d\n", tc.n))
		for j, v := range tc.arr {
			if j > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(v))
		}
		input.WriteByte('\n')
		expected := strconv.Itoa(tc.expect)
		got, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input.String())
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expected, got, input.String())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
