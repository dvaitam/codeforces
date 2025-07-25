package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func runBinary(path string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func generateTests() []string {
	r := rand.New(rand.NewSource(44))
	tests := make([]string, 0, 20)
	for i := 0; i < 20; i++ {
		m := r.Intn(10) + 1
		var sb strings.Builder
		sb.WriteString(strconv.Itoa(m))
		sb.WriteByte('\n')
		for j := 0; j < m; j++ {
			p := r.Intn(m) + 1
			t := r.Intn(2)
			sb.WriteString(strconv.Itoa(p))
			sb.WriteByte(' ')
			sb.WriteString(strconv.Itoa(t))
			if t == 1 {
				sb.WriteByte(' ')
				x := r.Intn(100) + 1
				sb.WriteString(strconv.Itoa(x))
			}
			if j+1 < m {
				sb.WriteByte('\n')
			}
		}
		sb.WriteByte('\n')
		tests = append(tests, sb.String())
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: verifierC.go <binary>")
		os.Exit(1)
	}
	candidate := os.Args[1]
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	ref := filepath.Join(dir, "756C.go")

	tests := generateTests()
	for i, t := range tests {
		expected, err1 := runBinary(ref, t)
		if err1 != nil {
			fmt.Printf("reference solution failed on test %d: %v\n", i+1, err1)
			os.Exit(1)
		}
		got, err2 := runBinary(candidate, t)
		if err2 != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err2)
			os.Exit(1)
		}
		if expected != got {
			fmt.Printf("test %d failed:\ninput:\n%s\nexpected: %s\ngot: %s\n", i+1, t, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
