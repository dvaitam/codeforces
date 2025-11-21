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

// solver replicates the reference logic to compute the minimum number of operations.
func solver(a []int) int {
	n := len(a)
	if n == 0 {
		return 0
	}
	ans := 0
	i := 0
	for i < n-1 {
		ans++
		j := i
		if a[j+1] > a[j] {
			for j+1 < n && a[j+1] > a[j] {
				j++
			}
		} else {
			for j+1 < n && a[j+1] < a[j] {
				j++
			}
		}
		i = j
	}
	ans++
	return ans
}

type testCase struct {
	n int
	a []int
}

func genCase(rng *rand.Rand) testCase {
	n := rng.Intn(40) + 2
	a := make([]int, n)
	a[0] = rng.Intn(2000000001) - 1000000000
	for i := 1; i < n; i++ {
		for {
			val := rng.Intn(2000000001) - 1000000000
			if val != a[i-1] {
				a[i] = val
				break
			}
		}
	}
	return testCase{n: n, a: a}
}

func buildInput(cases []testCase) (string, []int) {
	var sb strings.Builder
	fmt.Fprintln(&sb, len(cases))
	expected := make([]int, len(cases))
	for idx, tc := range cases {
		fmt.Fprintln(&sb, tc.n)
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprint(&sb, v)
		}
		sb.WriteByte('\n')
		expected[idx] = solver(tc.a)
	}
	return sb.String(), expected
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/2081B_binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var cases []testCase
	for i := 0; i < 200; i++ {
		cases = append(cases, genCase(rng))
	}
	input, expected := buildInput(cases)
	output, err := runCandidate(bin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to run candidate: %v\n", err)
		os.Exit(1)
	}
	lines := strings.Fields(output)
	if len(lines) < len(expected) {
		fmt.Fprintf(os.Stderr, "not enough outputs: got %d expected %d\n", len(lines), len(expected))
		os.Exit(1)
	}
	for i, exp := range expected {
		val, err := strconv.Atoi(lines[i])
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid integer on case %d: %q\n", i+1, lines[i])
			os.Exit(1)
		}
		if val != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\narray: %v\n", i+1, exp, val, cases[i].a)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
