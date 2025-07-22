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
		n := rand.Intn(5) + 2
		lines := []string{fmt.Sprintf("%d", n)}
		for i := 2; i <= n; i++ {
			p := rand.Intn(i-1) + 1
			w := rand.Intn(10) + 1
			lines = append(lines, fmt.Sprintf("%d %d %d", p, i, w))
		}
		q := rand.Intn(3) + 1
		lines = append(lines, fmt.Sprintf("%d", q))
		for i := 0; i < q; i++ {
			u := rand.Intn(n) + 1
			v := rand.Intn(n) + 1
			lines = append(lines, fmt.Sprintf("%d %d", u, v))
		}
		tests = append(tests, strings.Join(lines, "\n")+"\n")
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go <binary>")
		os.Exit(1)
	}
	candPath, err := prepareBinary(os.Args[1], "candD")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	refSrc := filepath.Join(baseDir(), "494D.go")
	refPath, err := prepareBinary(refSrc, "refD")
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
