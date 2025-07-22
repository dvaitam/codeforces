package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

func computeQualifiers(a, b []int) (string, string) {
	n := len(a)
	qa := make([]bool, n)
	qb := make([]bool, n)
	for k := 0; k <= n/2; k++ {
		ta := make([]bool, n)
		tb := make([]bool, n)
		for i := 0; i < k; i++ {
			ta[i] = true
			tb[i] = true
		}
		i, j := k, k
		need := n - 2*k
		for c := 0; c < need; c++ {
			if j >= n || (i < n && a[i] < b[j]) {
				ta[i] = true
				i++
			} else {
				tb[j] = true
				j++
			}
		}
		for idx := 0; idx < n; idx++ {
			if ta[idx] {
				qa[idx] = true
			}
			if tb[idx] {
				qb[idx] = true
			}
		}
	}
	var sbA, sbB strings.Builder
	for i := 0; i < n; i++ {
		if qa[i] {
			sbA.WriteByte('1')
		} else {
			sbA.WriteByte('0')
		}
	}
	for i := 0; i < n; i++ {
		if qb[i] {
			sbB.WriteByte('1')
		} else {
			sbB.WriteByte('0')
		}
	}
	return sbA.String(), sbB.String()
}

func solve(input string) string {
	reader := bufio.NewReader(strings.NewReader(strings.TrimSpace(input)))
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return ""
	}
	a := make([]int, n)
	b := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i], &b[i])
	}
	s1, s2 := computeQualifiers(a, b)
	return fmt.Sprintf("%s\n%s", s1, s2)
}

func runBinary(bin, input string) (string, error) {
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
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

type test struct {
	input    string
	expected string
}

func generateTest(rng *rand.Rand, n int) test {
	values := rng.Perm(2 * n * 5)
	a := make([]int, n)
	b := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = values[i] + 1
		b[i] = values[i+n] + 1
	}
	sort.Ints(a)
	sort.Ints(b)
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(n))
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", a[i], b[i]))
	}
	inp := sb.String()
	return test{inp, solve(inp)}
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(42))
	var tests []test
	fixed := []string{
		"1\n1 2\n",
		"2\n1 4\n2 3\n",
		"3\n1 6\n2 5\n3 4\n",
	}
	for _, f := range fixed {
		tests = append(tests, test{f, solve(f)})
	}
	for len(tests) < 100 {
		n := rng.Intn(8) + 1
		tests = append(tests, generateTest(rng, n))
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		out, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Printf("Runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(t.expected) {
			fmt.Printf("Wrong answer on test %d\nInput:\n%sExpected:\n%s\nGot:\n%s\n", i+1, t.input, t.expected, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
