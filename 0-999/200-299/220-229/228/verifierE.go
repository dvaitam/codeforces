package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func compileRef() (string, error) {
	src := filepath.Join(".", "228E.go")
	out := filepath.Join(os.TempDir(), "refE.bin")
	cmd := exec.Command("go", "build", "-o", out, src)
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return out, nil
}

func genTests() []string {
	r := rand.New(rand.NewSource(42))
	tests := []string{}
	for len(tests) < 100 {
		n := r.Intn(3) + 2 // 2..4
		maxM := n * (n - 1) / 2
		m := r.Intn(maxM + 1)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, m)
		edges := make(map[[2]int]bool)
		for i := 0; i < m; i++ {
			var a, b int
			for {
				a = r.Intn(n) + 1
				b = r.Intn(n) + 1
				if a != b && !edges[[2]int{a, b}] && !edges[[2]int{b, a}] {
					break
				}
			}
			edges[[2]int{a, b}] = true
			c := r.Intn(2)
			fmt.Fprintf(&sb, "%d %d %d\n", a, b, c)
		}
		tests = append(tests, sb.String())
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref, err := compileRef()
	if err != nil {
		fmt.Println("failed to compile reference:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)
	tests := genTests()
	for i, in := range tests {
		exp, err := runBinary(ref, in)
		if err != nil {
			fmt.Printf("reference error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		out, err := runBinary(bin, in)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if out != exp {
			fmt.Printf("test %d failed:\ninput:\n%sexpected=%s got=%s\n", i+1, in, exp, out)
			os.Exit(1)
		}
	}
	fmt.Printf("ok %d tests\n", len(tests))
}
