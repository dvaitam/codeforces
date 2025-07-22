package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func baseDir() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Dir(file)
}

func prepareBinary(path, tag string) (string, error) {
	if strings.HasSuffix(path, ".go") {
		bin := filepath.Join(os.TempDir(), tag)
		cmd := exec.Command("go", "build", "-o", bin, path)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("build %s: %v\n%s", path, err, out)
		}
		return bin, nil
	}
	return path, nil
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func genTests() []string {
	rand.Seed(1)
	var tests []string
	for t := 0; t < 100; t++ {
		n := rand.Intn(10) + 1
		m := rand.Intn(5) + 1
		k := rand.Intn(n) + 1
		lines := []string{fmt.Sprintf("%d %d %d", n, m, k)}
		for i := 0; i < m; i++ {
			a := rand.Intn(n) + 1
			c := a + rand.Intn(n-a+1)
			b := rand.Intn(n) + 1
			d := b + rand.Intn(n-b+1)
			lines = append(lines, fmt.Sprintf("%d %d %d %d", a, b, c, d))
		}
		tests = append(tests, strings.Join(lines, "\n")+"\n")
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go <binary>")
		os.Exit(1)
	}
	candPath, err := prepareBinary(os.Args[1], "candE")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	refSrc := filepath.Join(baseDir(), "494E.go")
	refPath, err := prepareBinary(refSrc, "refE")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	tests := genTests()
	for i, input := range tests {
		exp, err := runBinary(refPath, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runBinary(candPath, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(exp) != strings.TrimSpace(got) {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%sGot:\n%s", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
