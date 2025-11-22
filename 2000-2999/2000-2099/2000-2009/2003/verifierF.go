package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	refSourceF = "2003F.go"
	refBinaryF = "ref2003F.bin"
)

type testCase struct {
	name string
	n    int
	m    int
	a    []int
	b    []int
	c    []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refPath, err := buildReference()
	if err != nil {
		fmt.Println("failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refPath)

	tests := buildTests()
	for i, tc := range tests {
		input := formatInput(tc)
		expected, err := expectedValue(tc, refPath, input)
		if err != nil {
			fmt.Printf("reference failed on test %d (%s): %v\n", i+1, tc.name, err)
			printInput(input)
			os.Exit(1)
		}

		out, err := runProgram(candidate, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d (%s): %v\n", i+1, tc.name, err)
			printInput(input)
			os.Exit(1)
		}
		got, err := parseOutput(out)
		if err != nil {
			fmt.Printf("failed to parse candidate output on test %d (%s): %v\noutput:\n%s\n", i+1, tc.name, err, out)
			printInput(input)
			os.Exit(1)
		}

		if got != expected {
			fmt.Printf("test %d (%s) failed: expected %d, got %d\n", i+1, tc.name, expected, got)
			printInput(input)
			fmt.Println("Candidate output:")
			fmt.Println(out)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, error) {
	cmd := exec.Command("go", "build", "-o", refBinaryF, refSourceF)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return filepath.Join(".", refBinaryF), nil
}

func runProgram(path string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), err
	}
	return out.String(), nil
}

func parseOutput(out string) (int, error) {
	fields := strings.Fields(out)
	if len(fields) != 1 {
		return 0, fmt.Errorf("expected 1 integer, got %d tokens", len(fields))
	}
	val, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q: %v", fields[0], err)
	}
	return val, nil
}

func formatInput(tc testCase) []byte {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.m)
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	for i, v := range tc.b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	for i, v := range tc.c {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	return []byte(sb.String())
}

func expectedValue(tc testCase, refPath string, input []byte) (int, error) {
	if tc.n <= 10 {
		return bruteForce(tc), nil
	}
	best := -1
	for i := 0; i < 4; i++ {
		out, err := runProgram(refPath, input)
		if err != nil {
			return 0, fmt.Errorf("reference runtime error: %v\n%s", err, out)
		}
		val, err := parseOutput(out)
		if err != nil {
			return 0, fmt.Errorf("invalid reference output: %v\n%s", err, out)
		}
		if val > best {
			best = val
		}
	}
	return best, nil
}

func bruteForce(tc testCase) int {
	if tc.m > tc.n {
		return -1
	}
	used := make([]bool, tc.n+1)
	best := -1
	var dfs func(pos, chosen, lastA, sum int)
	dfs = func(pos, chosen, lastA, sum int) {
		if chosen == tc.m {
			if sum > best {
				best = sum
			}
			return
		}
		remaining := tc.m - chosen
		for i := pos; i <= tc.n-remaining; i++ {
			ai := tc.a[i]
			if ai < lastA {
				continue
			}
			bi := tc.b[i]
			if used[bi] {
				continue
			}
			used[bi] = true
			dfs(i+1, chosen+1, ai, sum+tc.c[i])
			used[bi] = false
		}
	}
	dfs(0, 0, 0, 0)
	return best
}

func buildTests() []testCase {
	tests := make([]testCase, 0, 80)

	// Samples from the statement.
	tests = append(tests, testCase{
		name: "sample1",
		n:    4,
		m:    2,
		a:    []int{2, 3, 4, 2},
		b:    []int{1, 3, 3, 2},
		c:    []int{1, 4, 2, 3},
	})
	tests = append(tests, testCase{
		name: "sample2",
		n:    7,
		m:    3,
		a:    []int{1, 4, 5, 2, 3, 6, 7},
		b:    []int{1, 2, 2, 1, 1, 3, 2},
		c:    []int{1, 5, 6, 7, 3, 2, 4},
	})
	tests = append(tests, testCase{
		name: "sample3",
		n:    5,
		m:    3,
		a:    []int{1, 2, 3, 4, 5},
		b:    []int{1, 1, 2, 1, 2},
		c:    []int{5, 4, 3, 2, 1},
	})

	// Small hand-crafted edge cases.
	tests = append(tests, testCase{
		name: "single-element",
		n:    1,
		m:    1,
		a:    []int{1},
		b:    []int{1},
		c:    []int{7},
	})
	tests = append(tests, testCase{
		name: "no-distinct-b",
		n:    3,
		m:    2,
		a:    []int{1, 2, 3},
		b:    []int{2, 2, 2},
		c:    []int{5, 6, 7},
	})
	tests = append(tests, testCase{
		name: "max-m-equal-n",
		n:    5,
		m:    5,
		a:    []int{1, 1, 2, 2, 3},
		b:    []int{1, 2, 3, 4, 5},
		c:    []int{10, 9, 8, 7, 6},
	})
	tests = append(tests, testCase{
		name: "strict-a-block",
		n:    6,
		m:    3,
		a:    []int{3, 3, 3, 2, 2, 4},
		b:    []int{1, 2, 3, 4, 5, 6},
		c:    []int{5, 4, 3, 2, 1, 6},
	})

	// Random brute-force sized cases (n <= 8).
	smallRnd := rand.New(rand.NewSource(1234567))
	for i := 0; i < 30; i++ {
		n := smallRnd.Intn(8) + 1
		m := smallRnd.Intn(5) + 1
		a := make([]int, n)
		b := make([]int, n)
		c := make([]int, n)
		for j := 0; j < n; j++ {
			a[j] = smallRnd.Intn(n) + 1
			b[j] = smallRnd.Intn(n) + 1
			c[j] = smallRnd.Intn(15) + 1
		}
		tests = append(tests, testCase{
			name: fmt.Sprintf("small-%d", i+1),
			n:    n,
			m:    m,
			a:    a,
			b:    b,
			c:    c,
		})
	}

	// Larger randomized cases for stress testing.
	largeRnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 40; i++ {
		n := largeRnd.Intn(700) + 50
		if i%10 == 0 {
			n = 3000
		}
		m := largeRnd.Intn(5) + 1
		a := make([]int, n)
		b := make([]int, n)
		c := make([]int, n)
		for j := 0; j < n; j++ {
			a[j] = largeRnd.Intn(n) + 1
			b[j] = largeRnd.Intn(n) + 1
			c[j] = largeRnd.Intn(10000) + 1
		}
		tests = append(tests, testCase{
			name: fmt.Sprintf("large-%d", i+1),
			n:    n,
			m:    m,
			a:    a,
			b:    b,
			c:    c,
		})
	}

	return tests
}

func printInput(in []byte) {
	fmt.Println("Input used:")
	fmt.Print(string(in))
}
