package main

import (
	"bytes"
	"fmt"
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
	src := filepath.Join(".", "228C.go")
	out := filepath.Join(os.TempDir(), "refC.bin")
	cmd := exec.Command("go", "build", "-o", out, src)
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return out, nil
}

func genTests() []string {
	tests := []string{}
	for n := 2; len(tests) < 100 && n <= 3; n++ {
		for m := 2; len(tests) < 100 && m <= 3; m++ {
			cells := n * m
			for mask := 0; len(tests) < 100 && mask < (1<<cells); mask++ {
				var sb strings.Builder
				fmt.Fprintf(&sb, "%d %d\n", n, m)
				for i := 0; i < n; i++ {
					for j := 0; j < m; j++ {
						if mask&(1<<(i*m+j)) != 0 {
							sb.WriteByte('*')
						} else {
							sb.WriteByte('.')
						}
					}
					sb.WriteByte('\n')
				}
				tests = append(tests, sb.String())
			}
		}
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
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
