package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCaseA struct {
	n   int
	arr []int
}

func genTests() []testCaseA {
	rng := rand.New(rand.NewSource(1))
	tests := make([]testCaseA, 100)
	for i := range tests {
		n := rng.Intn(10) + 1
		arr := make([]int, n)
		for j := range arr {
			arr[j] = rng.Intn(100) + 1
		}
		tests[i] = testCaseA{n, arr}
	}
	// additional edge cases
	tests = append(tests,
		testCaseA{1, []int{1}},
		testCaseA{1, []int{2}},
		testCaseA{2, []int{1, 3}},
		testCaseA{2, []int{2, 2}},
		testCaseA{3, []int{1, 1, 1}},
	)
	return tests
}

func hasEvenSubset(a []int) bool {
	even := false
	oddCount := 0
	for _, v := range a {
		if v%2 == 0 {
			even = true
		} else {
			oddCount++
		}
	}
	if even {
		return true
	}
	return oddCount >= 2
}

func run(bin string, input string) (string, error) {
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
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return out.String(), nil
}

func verifyCase(bin string, tc testCaseA) error {
	var input bytes.Buffer
	fmt.Fprintln(&input, 1)
	fmt.Fprintln(&input, tc.n)
	for i, v := range tc.arr {
		if i > 0 {
			input.WriteByte(' ')
		}
		fmt.Fprint(&input, v)
	}
	input.WriteByte('\n')
	output, err := run(bin, input.String())
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(strings.NewReader(strings.TrimSpace(output)))
	if !scanner.Scan() {
		return fmt.Errorf("no output")
	}
	line := strings.TrimSpace(scanner.Text())
	if line == "-1" {
		if hasEvenSubset(tc.arr) {
			return fmt.Errorf("expected subset but got -1")
		}
		if scanner.Scan() {
			return fmt.Errorf("extra output after -1")
		}
		return nil
	}
	k, err := strconv.Atoi(line)
	if err != nil || k <= 0 || k > tc.n {
		return fmt.Errorf("invalid k: %s", line)
	}
	if !scanner.Scan() {
		return fmt.Errorf("missing index line")
	}
	idxFields := strings.Fields(scanner.Text())
	if len(idxFields) != k {
		return fmt.Errorf("expected %d indices got %d", k, len(idxFields))
	}
	used := make(map[int]bool)
	sum := 0
	for _, f := range idxFields {
		idx, err := strconv.Atoi(f)
		if err != nil || idx < 1 || idx > tc.n || used[idx] {
			return fmt.Errorf("invalid index %s", f)
		}
		used[idx] = true
		sum += tc.arr[idx-1]
	}
	if sum%2 != 0 {
		return fmt.Errorf("subset sum %d not even", sum)
	}
	if scanner.Scan() {
		return fmt.Errorf("extra output after indices")
	}
	if !hasEvenSubset(tc.arr) {
		return fmt.Errorf("should have output -1")
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()
	for i, tc := range tests {
		if err := verifyCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput:\n%v\n", i+1, err, tc)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
