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

type Edge struct{ u, v int }

type Test struct {
	n, m  int
	a, b  int
	edges []Edge
}

func (t Test) Input() string {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d %d %d\n", t.n, t.m, t.a, t.b))
	for _, e := range t.edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e.u, e.v))
	}
	return sb.String()
}

func buildRef() (string, error) {
	ref := "refE.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1259E.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

func runExe(path, input string) (string, error) {
	if !strings.Contains(path, "/") {
		path = "./" + path
	}
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func genGraph(r *rand.Rand, n int) []Edge {
	edgesMap := make(map[[2]int]bool)
	res := make([]Edge, 0)
	for len(res) < n-1 {
		u := r.Intn(n) + 1
		v := r.Intn(n) + 1
		if u == v {
			continue
		}
		a, b := u, v
		if a > b {
			a, b = b, a
		}
		if edgesMap[[2]int{a, b}] {
			continue
		}
		edgesMap[[2]int{a, b}] = true
		res = append(res, Edge{u, v})
	}
	m := n - 1 + r.Intn(n)
	for len(res) < m {
		u := r.Intn(n) + 1
		v := r.Intn(n) + 1
		if u == v {
			continue
		}
		a, b := u, v
		if a > b {
			a, b = b, a
		}
		if edgesMap[[2]int{a, b}] {
			continue
		}
		edgesMap[[2]int{a, b}] = true
		res = append(res, Edge{u, v})
	}
	return res
}

func genTests() []Test {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]Test, 0, 101)
	for i := 0; i < 100; i++ {
		n := r.Intn(6) + 4
		edges := genGraph(r, n)
		a := r.Intn(n) + 1
		b := r.Intn(n) + 1
		for a == b {
			b = r.Intn(n) + 1
		}
		tests = append(tests, Test{n, len(edges), a, b, edges})
	}
	// simple fixed case
	tests = append(tests, Test{4, 3, 1, 2, []Edge{{1, 3}, {2, 3}, {3, 4}}})
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
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
		if exp != got {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%s\nGot:\n%s\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
