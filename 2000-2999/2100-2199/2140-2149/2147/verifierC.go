package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	input  string
	expect string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
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
	exp := strings.TrimSpace(expect)
	act := strings.ToLower(strings.TrimSpace(actual))
	if exp == "YES" {
		if act != "yes" {
			return fmt.Errorf("expected YES but got %s", actual)
		}
	} else {
		if act != "no" {
			return fmt.Errorf("expected NO but got %s", actual)
		}
	}
	return nil
}

func genTests() []testCase {
	rand.Seed(42)
	tests := []testCase{
		makeCase("0101"),
		makeCase("00100"),
		makeCase("01010"),
		makeCase("0110"),
		breakCase("010"),
	}
	for i := 0; i < 200; i++ {
		n := rand.Intn(50) + 1
		s := randomString(n)
		tests = append(tests, makeCase(s))
	}
	return tests
}

func randomString(n int) string {
	var sb strings.Builder
	for i := 0; i < n; i++ {
		if rand.Intn(2) == 0 {
			sb.WriteByte('0')
		} else {
			sb.WriteByte('1')
		}
	}
	return sb.String()
}

func makeCase(s string) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d\n%s\n", len(s), s)
	return testCase{
		input:  sb.String(),
		expect: solveReference(s),
	}
}

func breakCase(s string) testCase {
	return makeCase(s)
}

func solveReference(s string) string {
	n := len(s)
	forced := make([]bool, n)
	adj := make([][]int, n)
	for i := 0; i < n; i++ {
		if s[i] != '0' {
			continue
		}
		boundary := i == 0 || i == n-1
		leftZero := i > 0 && s[i-1] == '0'
		rightZero := i+1 < n && s[i+1] == '0'
		if !(boundary || leftZero || rightZero) {
			forced[i] = true
		}
	}
	for i := 0; i+2 < n; i++ {
		if s[i] == '0' && s[i+1] == '1' && s[i+2] == '0' {
			adj[i] = append(adj[i], i+2)
			adj[i+2] = append(adj[i+2], i)
		}
	}
	for i := 0; i < n; i++ {
		if s[i] == '0' && forced[i] && len(adj[i]) == 0 {
			return "NO"
		}
	}
	visited := make([]bool, n)
	for i := 0; i < n; i++ {
		if s[i] != '0' || len(adj[i]) == 0 || visited[i] {
			continue
		}
		var stack []int
		stack = append(stack, i)
		visited[i] = true
		var comp []int
		for len(stack) > 0 {
			u := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			comp = append(comp, u)
			for _, v := range adj[u] {
				if !visited[v] {
					visited[v] = true
					stack = append(stack, v)
				}
			}
		}
		if !componentPossible(comp, forced) {
			return "NO"
		}
	}
	return "YES"
}

func componentPossible(nodes []int, forced []bool) bool {
	k := len(nodes)
	dpFree := true
	dpTaken := false
	for idx := 0; idx < k; idx++ {
		pos := nodes[idx]
		forcedNode := forced[pos]
		nextFree := false
		nextTaken := false
		if dpTaken {
			nextFree = true
		}
		if dpFree {
			if !forcedNode {
				nextFree = true
			}
			if idx+1 < k {
				nextTaken = true
			}
		}
		dpFree, dpTaken = nextFree, nextTaken
	}
	return dpFree
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
