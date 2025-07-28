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
	"time"
)

type testA struct {
	n int
	a []int
}

func runCandidate(bin, input string) (string, error) {
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
		return out.String(), fmt.Errorf("runtime error: %v", err)
	}
	return out.String(), nil
}

func genTests() []testA {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := []testA{
		{n: 1, a: []int{1}},
		{n: 2, a: []int{1, 2}},
		{n: 2, a: []int{2, 1}},
		{n: 5, a: []int{1, 2, 4, 5, 3}},
	}
	for len(tests) < 100 {
		n := rng.Intn(100) + 1
		perm := rng.Perm(n)
		for i := range perm {
			perm[i]++
		}
		cp := make([]int, n)
		copy(cp, perm)
		tests = append(tests, testA{n: n, a: cp})
	}
	return tests
}

func check(tc testA, out string) error {
	scanner := bufio.NewScanner(strings.NewReader(out))
	scanner.Split(bufio.ScanWords)
	vals := make([]int, 0, tc.n)
	for scanner.Scan() {
		if len(vals) == tc.n {
			return fmt.Errorf("extra output")
		}
		v, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return fmt.Errorf("invalid integer")
		}
		vals = append(vals, v)
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scan error: %v", err)
	}
	if len(vals) != tc.n {
		return fmt.Errorf("expected %d numbers got %d", tc.n, len(vals))
	}
	used := make([]bool, tc.n+1)
	prev := -1 << 60
	for i, v := range vals {
		if v < 1 || v > tc.n || used[v] {
			return fmt.Errorf("invalid permutation")
		}
		used[v] = true
		sum := tc.a[i] + v
		if i > 0 && sum < prev {
			return fmt.Errorf("sum decreases at index %d", i+1)
		}
		prev = sum
	}
	return nil
}

func runCase(bin string, tc testA) error {
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d\n", tc.n)
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	out, err := runCandidate(bin, sb.String())
	if err != nil {
		return fmt.Errorf("%v\noutput:\n%s", err, out)
	}
	return check(tc, strings.TrimSpace(out))
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := args[0]
	tests := genTests()
	for i, tc := range tests {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
