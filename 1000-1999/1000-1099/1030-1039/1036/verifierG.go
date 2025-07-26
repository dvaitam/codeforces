package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Test struct{ input string }

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
	ref := "./refG.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1036G.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

func genTests() []Test {
	rand.Seed(0)
	tests := make([]Test, 0, 110)
	for i := 0; i < 100; i++ {
		n := rand.Intn(4) + 2
		m := rand.Intn(n*(n-1)/2 + 1)
		edges := make([][2]int, 0, m)
		used := make(map[[2]int]bool)
		for len(edges) < m {
			a := rand.Intn(n) + 1
			b := rand.Intn(n) + 1
			if a == b {
				continue
			}
			if a > b {
				a, b = b, a
			}
			p := [2]int{a, b}
			if !used[p] {
				used[p] = true
				edges = append(edges, p)
			}
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for _, e := range edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
		tests = append(tests, Test{sb.String()})
	}
	tests = append(tests,
		Test{"2 0\n"},
		Test{"2 1\n1 2\n"},
		Test{"3 3\n1 2\n2 3\n1 3\n"},
		Test{"4 0\n"},
		Test{"3 1\n1 2\n"},
	)
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
		exp, err := runExe(ref, tc.input)
		if err != nil {
			fmt.Printf("reference runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runExe(bin, tc.input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(exp) != strings.TrimSpace(got) {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%sGot:\n%s\n", i+1, tc.input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
