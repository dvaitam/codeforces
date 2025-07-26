package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCase struct {
	n int
	v []int
}

func generateTests() []testCase {
	rnd := rand.New(rand.NewSource(1))
	tests := make([]testCase, 100)
	for i := range tests {
		n := rnd.Intn(50) + 2
		p := rnd.Perm(n)
		v := make([]int, n)
		for j := 0; j < n; j++ {
			v[j] = p[j] + 1
		}
		for j := 0; j < n; j++ {
			if v[j] == j+1 {
				if j+1 < n {
					v[j], v[j+1] = v[j+1], v[j]
				} else {
					v[j], v[0] = v[0], v[j]
				}
			}
		}
		zeros := rnd.Intn(n-1) + 2
		if zeros > n {
			zeros = n
		}
		perm := rnd.Perm(n)
		for k := 0; k < zeros; k++ {
			v[perm[k]] = 0
		}
		tests[i] = testCase{n, v}
	}
	return tests
}

func solveExpected(tc testCase) []int {
	n := tc.n
	v := append([]int(nil), tc.v...)
	used := make([]bool, n)
	for i := 0; i < n; i++ {
		if v[i] != 0 {
			used[v[i]-1] = true
		}
	}
	var givers, missing []int
	for i := 0; i < n; i++ {
		if v[i] == 0 {
			givers = append(givers, i)
		}
		if !used[i] {
			missing = append(missing, i)
		}
	}
	k := len(givers)
	for i := 0; i < k; i++ {
		j := (i + 1) % k
		v[givers[i]] = missing[j] + 1
	}
	return v
}

func runBinary(bin string, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return out.String(), fmt.Errorf("timeout")
	}
	if err != nil {
		return out.String(), fmt.Errorf("exec error: %v", err)
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for idx, tc := range tests {
		input := fmt.Sprintf("%d\n", tc.n)
		for i, x := range tc.v {
			if i > 0 {
				input += " "
			}
			input += fmt.Sprintf("%d", x)
		}
		input += "\n"
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		fields := strings.Fields(out)
		if len(fields) != tc.n {
			fmt.Printf("test %d: expected %d numbers got %d\n", idx+1, tc.n, len(fields))
			os.Exit(1)
		}
		res := make([]int, tc.n)
		for i, f := range fields {
			if _, err := fmt.Sscan(f, &res[i]); err != nil {
				fmt.Printf("test %d: invalid number %q\n", idx+1, f)
				os.Exit(1)
			}
		}
		expected := solveExpected(tc)
		for i := 0; i < tc.n; i++ {
			if res[i] < 1 || res[i] > tc.n {
				fmt.Printf("test %d: value out of range\n", idx+1)
				os.Exit(1)
			}
		}
		// check permutation
		seen := make(map[int]bool)
		for i, v := range res {
			if tc.v[i] != 0 && v != tc.v[i] {
				fmt.Printf("test %d: value %d expected fixed %d\n", idx+1, v, tc.v[i])
				os.Exit(1)
			}
			if v == i+1 {
				fmt.Printf("test %d: self gift\n", idx+1)
				os.Exit(1)
			}
			if seen[v] {
				fmt.Printf("test %d: duplicate value %d\n", idx+1, v)
				os.Exit(1)
			}
			seen[v] = true
		}
		// expected arrangement from algorithm
		for i, v := range expected {
			if res[i] != v {
				fmt.Printf("test %d: expected %v got %v\n", idx+1, expected, res)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
