package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
	"unicode"
)

type operation struct {
	t int
	c int
}

type testCase struct {
	name string
	n    int
	a    int
	b    int
	s    []int
}

func canSort(s []int) bool {
	n := len(s)
	stack := make([]int, 0, n)
	need := 1
	for i := n - 1; i >= 0; i-- {
		stack = append(stack, s[i])
		for len(stack) > 0 && stack[len(stack)-1] == need {
			stack = stack[:len(stack)-1]
			need++
		}
	}
	return need == n+1
}

func simulate(tc testCase, ops []operation) error {
	dirty := append([]int(nil), tc.s...)
	inter := make([]int, 0, tc.n)
	dryer := make([]int, 0, tc.n)
	for idx, op := range ops {
		switch op.t {
		case 1:
			if op.c < 1 || op.c > tc.a {
				return fmt.Errorf("operation %d: move count %d exceeds limit a=%d", idx+1, op.c, tc.a)
			}
			if len(dirty) < op.c {
				return fmt.Errorf("operation %d: trying to move %d from dirty but only %d left", idx+1, op.c, len(dirty))
			}
			move := dirty[len(dirty)-op.c:]
			dirty = dirty[:len(dirty)-op.c]
			inter = append(inter, move...)
		case 2:
			if op.c < 1 || op.c > tc.b {
				return fmt.Errorf("operation %d: move count %d exceeds limit b=%d", idx+1, op.c, tc.b)
			}
			if len(inter) < op.c {
				return fmt.Errorf("operation %d: trying to move %d from intermediate but only %d available", idx+1, op.c, len(inter))
			}
			move := inter[len(inter)-op.c:]
			inter = inter[:len(inter)-op.c]
			dryer = append(dryer, move...)
		default:
			return fmt.Errorf("operation %d: invalid type %d", idx+1, op.t)
		}
	}
	if len(dirty) != 0 {
		return fmt.Errorf("dirty stack not empty after operations")
	}
	if len(inter) != 0 {
		return fmt.Errorf("intermediate stack not empty after operations")
	}
	if len(dryer) != tc.n {
		return fmt.Errorf("dryer has %d plates, expected %d", len(dryer), tc.n)
	}
	for i := 1; i < len(dryer); i++ {
		if dryer[i] <= dryer[i-1] {
			return fmt.Errorf("dryer is not strictly increasing (positions %d and %d)", i, i+1)
		}
	}
	return nil
}

func hasExtraTokens(r *bufio.Reader) bool {
	for {
		ch, _, err := r.ReadRune()
		if err != nil {
			return false
		}
		if !unicode.IsSpace(ch) {
			return true
		}
	}
}

func parseOutput(tc testCase, output string, expectPossible bool) error {
	reader := bufio.NewReader(strings.NewReader(output))
	var status string
	if _, err := fmt.Fscan(reader, &status); err != nil {
		if strings.TrimSpace(output) == "" {
			return fmt.Errorf("no output produced")
		}
		return fmt.Errorf("failed to read status: %v", err)
	}
	statusUpper := strings.ToUpper(status)
	switch statusUpper {
	case "NO":
		if expectPossible {
			return fmt.Errorf("expected YES but got NO")
		}
		if hasExtraTokens(reader) {
			return fmt.Errorf("extra data after NO")
		}
		return nil
	case "YES":
		if !expectPossible {
			return fmt.Errorf("expected NO but got YES")
		}
		var k int
		if _, err := fmt.Fscan(reader, &k); err != nil {
			return fmt.Errorf("failed to read number of operations: %v", err)
		}
		if k < 0 {
			return fmt.Errorf("negative number of operations")
		}
		ops := make([]operation, k)
		for i := 0; i < k; i++ {
			if _, err := fmt.Fscan(reader, &ops[i].t, &ops[i].c); err != nil {
				return fmt.Errorf("failed to read operation %d: %v", i+1, err)
			}
		}
		if hasExtraTokens(reader) {
			return fmt.Errorf("extra data after listed operations")
		}
		if err := simulate(tc, ops); err != nil {
			return fmt.Errorf("invalid operations: %v", err)
		}
		return nil
	default:
		return fmt.Errorf("unexpected status %q", status)
	}
}

func formatInput(tc testCase) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d %d %d\n", tc.n, tc.a, tc.b)
	for i := 0; i < tc.n; i++ {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(&b, "%d", tc.s[i])
	}
	b.WriteByte('\n')
	return b.String()
}

func runCandidate(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
	}
	if stderr.Len() != 0 {
		// In case the solution prints debugging info to stderr, treat as error.
		return "", fmt.Errorf("unexpected stderr output: %s", stderr.String())
	}
	return stdout.String(), nil
}

func randomPermutation(rng *rand.Rand, n int) []int {
	perm := make([]int, n)
	for i := 0; i < n; i++ {
		perm[i] = i + 1
	}
	for i := n - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		perm[i], perm[j] = perm[j], perm[i]
	}
	return perm
}

func handcraftedTests() []testCase {
	return []testCase{
		{name: "n1", n: 1, a: 1, b: 1, s: []int{1}},
		{name: "already_sorted", n: 4, a: 2, b: 2, s: []int{1, 2, 3, 4}},
		{name: "reverse_order", n: 4, a: 2, b: 2, s: []int{4, 3, 2, 1}},
		{name: "impossible_small", n: 3, a: 1, b: 1, s: []int{2, 1, 3}},
		{name: "a1b3", n: 5, a: 1, b: 3, s: []int{5, 1, 2, 3, 4}},
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	// small cases
	for i := 0; i < 80; i++ {
		n := rng.Intn(6) + 1
		a := rng.Intn(n) + 1
		b := rng.Intn(n) + 1
		tests = append(tests, testCase{
			name: fmt.Sprintf("small_%d", i+1),
			n:    n,
			a:    a,
			b:    b,
			s:    randomPermutation(rng, n),
		})
	}
	// medium cases
	for i := 0; i < 40; i++ {
		n := rng.Intn(20) + 6
		a := rng.Intn(n) + 1
		b := rng.Intn(n) + 1
		tests = append(tests, testCase{
			name: fmt.Sprintf("medium_%d", i+1),
			n:    n,
			a:    a,
			b:    b,
			s:    randomPermutation(rng, n),
		})
	}
	// larger cases to test performance
	for i := 0; i < 10; i++ {
		n := rng.Intn(200) + 100
		a := rng.Intn(n) + 1
		b := rng.Intn(n) + 1
		tests = append(tests, testCase{
			name: fmt.Sprintf("large_%d", i+1),
			n:    n,
			a:    a,
			b:    b,
			s:    randomPermutation(rng, n),
		})
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := append(handcraftedTests(), randomTests()...)
	for idx, tc := range tests {
		input := formatInput(tc)
		expect := canSort(tc.s)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("test %d (%s) runtime error: %v\n", idx+1, tc.name, err)
			os.Exit(1)
		}
		if err := parseOutput(tc, out, expect); err != nil {
			fmt.Printf("test %d (%s) failed: %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, input, out)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
