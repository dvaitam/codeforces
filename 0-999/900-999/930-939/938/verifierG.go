package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Edge struct {
	u, v int
	w    int
}

type Query struct {
	typ int
	x   int
	y   int
	w   int
}

type Test struct {
	n     int
	edges []Edge
	qs    []Query
}

func (t Test) Input() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", t.n, len(t.edges)))
	for _, e := range t.edges {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", e.u, e.v, e.w))
	}
	sb.WriteString(strconv.Itoa(len(t.qs)))
	sb.WriteByte('\n')
	for _, q := range t.qs {
		if q.typ == 1 {
			sb.WriteString(fmt.Sprintf("1 %d %d %d\n", q.x, q.y, q.w))
		} else if q.typ == 2 {
			sb.WriteString(fmt.Sprintf("2 %d %d\n", q.x, q.y))
		} else {
			sb.WriteString(fmt.Sprintf("3 %d %d\n", q.x, q.y))
		}
	}
	return sb.String()
}

func runExe(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func buildRef() (string, error) {
	ref := "./refG.bin"
	cmd := exec.Command("go", "build", "-o", ref, "938G.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

func genTests() []Test {
	rand.Seed(6)
	tests := make([]Test, 0, 101)
	for i := 0; i < 100; i++ {
		n := rand.Intn(3) + 2
		edges := make([]Edge, 0, n-1)
		for j := 1; j < n; j++ {
			edges = append(edges, Edge{j, j + 1, rand.Intn(20)})
		}
		qs := make([]Query, 0, rand.Intn(4)+1)
		for j := 0; j < cap(qs); j++ {
			x := rand.Intn(n) + 1
			y := rand.Intn(n) + 1
			for y == x {
				y = rand.Intn(n) + 1
			}
			if x > y {
				x, y = y, x
			}
			qs = append(qs, Query{typ: 3, x: x, y: y})
		}
		tests = append(tests, Test{n, edges, qs})
	}
	tests = append(tests, Test{2, []Edge{{1, 2, 1}}, []Query{{typ: 3, x: 1, y: 2}}})
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierG.go /path/to/binary")
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
		exp, err := runExe(ref, tc.Input())
		if err != nil {
			fmt.Printf("reference runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runExe(bin, tc.Input())
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(exp) != strings.TrimSpace(got) {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%sGot:\n%s\n", i+1, tc.Input(), exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
