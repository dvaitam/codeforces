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

type Test struct {
	n     int
	edges []Edge
	a     []int
}

func (t Test) Input() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", t.n, len(t.edges)))
	for _, e := range t.edges {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", e.u, e.v, e.w))
	}
	for i, v := range t.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
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
	ref := "./refD.bin"
	cmd := exec.Command("go", "build", "-o", ref, "938D.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

func genTests() []Test {
	rand.Seed(3)
	tests := make([]Test, 0, 101)
	for i := 0; i < 100; i++ {
		n := rand.Intn(4) + 2
		m := rand.Intn(n*(n-1)/2-n+1) + n - 1
		edges := make([]Edge, 0, m)
		used := map[[2]int]bool{}
		for j := 1; j < n; j++ {
			w := rand.Intn(20) + 1
			edges = append(edges, Edge{j, j + 1, w})
			used[[2]int{j, j + 1}] = true
		}
		for len(edges) < m {
			u := rand.Intn(n) + 1
			v := rand.Intn(n) + 1
			if u == v {
				continue
			}
			if u > v {
				u, v = v, u
			}
			if used[[2]int{u, v}] {
				continue
			}
			used[[2]int{u, v}] = true
			w := rand.Intn(20) + 1
			edges = append(edges, Edge{u, v, w})
		}
		a := make([]int, n)
		for j := range a {
			a[j] = rand.Intn(50) + 1
		}
		tests = append(tests, Test{n, edges, a})
	}
	// simple case n=2
	tests = append(tests, Test{2, []Edge{{1, 2, 5}}, []int{3, 4}})
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
