package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type testCase struct {
	input  string
	expect string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()
	for idx, tc := range tests {
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\ninput:\n%s\n", idx+1, err, tc.input)
			os.Exit(1)
		}
		if err := check(tc.expect, strings.TrimSpace(out)); err != nil {
			fmt.Printf("test %d failed: %v\ninput:\n%s\nexpected:\n%s\nactual:\n%s\n", idx+1, err, tc.input, tc.expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}

func check(expect, actual string) error {
	exp, _ := strconv.Atoi(expect)
	val, err := strconv.Atoi(actual)
	if err != nil {
		return fmt.Errorf("output is not integer: %v", err)
	}
	if val != exp {
		return fmt.Errorf("expected %d but got %d", exp, val)
	}
	return nil
}

func genTests() []testCase {
	rand.Seed(42)
	tests := []testCase{
		makeTest([]int{2, 2, 3, 4, 3, 1}),
		makeTest([]int{1, 2, 3, 4, 5}),
		makeTest([]int{10, 10, 10}),
	}
	for i := 0; i < 200; i++ {
		n := rand.Intn(50) + 1
		arr := make([]int, n)
		for j := range arr {
			arr[j] = rand.Intn(20) + 1
		}
		tests = append(tests, makeTest(arr))
	}
	return tests
}

func makeTest(arr []int) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d\n", len(arr))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	return testCase{
		input:  sb.String(),
		expect: fmt.Sprintf("%d", solveReference(arr)),
	}
}

func solveReference(arr []int) int {
	freq := make(map[int]int)
	for _, x := range arr {
		freq[x]++
	}
	keys := make([]int, 0, len(freq))
	for k := range freq {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	total := 0
	active := 0
	prev := 0
	first := true
	for _, k := range keys {
		if first {
			first = false
		} else if k != prev+1 {
			active = 0
		}
		need := freq[k] - active
		if need > 0 {
			total += need
		}
		active = freq[k]
		prev = k
	}
	return total
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
