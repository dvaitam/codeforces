package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Query struct{ a, b, c int }

type Test struct {
	n       int
	parents []int // length n-1 for nodes 2..n
	q       int
	queries []Query
}

func (tc Test) Input() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.q))
	for i := 2; i <= tc.n; i++ {
		sb.WriteString(fmt.Sprintf("%d ", tc.parents[i-2]))
	}
	sb.WriteByte('\n')
	for _, qu := range tc.queries {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", qu.a, qu.b, qu.c))
	}
	return sb.String()
}

func runExe(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func buildRef() (string, error) {
	ref := "./refD.bin"
	cmd := exec.Command("go", "build", "-o", ref, "832D.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

func genTree(n int) []int {
	p := make([]int, n-1)
	for i := 2; i <= n; i++ {
		p[i-2] = rand.Intn(i-1) + 1
	}
	return p
}

func genTests() []Test {
	rand.Seed(time.Now().UnixNano())
	tests := make([]Test, 0, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(7) + 3
		parents := genTree(n)
		q := rand.Intn(5) + 1
		queries := make([]Query, q)
		for j := 0; j < q; j++ {
			a := rand.Intn(n) + 1
			b := rand.Intn(n) + 1
			c := rand.Intn(n) + 1
			queries[j] = Query{a, b, c}
		}
		tests = append(tests, Test{n: n, parents: parents, q: q, queries: queries})
	}
	// simple tree
	tests = append(tests, Test{n: 3, parents: []int{1, 1}, q: 1, queries: []Query{{1, 2, 3}}})
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := genTests()
	for i, tc := range tests {
		input := tc.Input()
		exp, err := runExe(ref, input)
		if err != nil {
			fmt.Printf("reference runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runExe(bin, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(exp) != strings.TrimSpace(got) {
			fmt.Printf("Test %d failed\nInput:%sExpected:%sGot:%s\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
