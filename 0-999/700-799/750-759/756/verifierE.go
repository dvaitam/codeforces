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
	r := rand.New(rand.NewSource(46))
	tests := make([]string, 0, 20)
	for i := 0; i < 20; i++ {
		n := r.Intn(4) + 2 // at least 2
		a := make([]int64, n-1)
		for j := 0; j < n-1; j++ {
			a[j] = int64(r.Intn(3) + 1)
		}
		b := make([]int, n)
		for j := 0; j < n; j++ {
			b[j] = r.Intn(3)
		}
		m := int64(r.Intn(1000))
		var sb strings.Builder
		sb.WriteString(strconv.Itoa(n))
		sb.WriteByte('\n')
		for j, v := range a {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
		for j, v := range b {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		sb.WriteString(strconv.FormatInt(m, 10))
		sb.WriteByte('\n')
		tests = append(tests, sb.String())
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: verifierE.go <binary>")
		os.Exit(1)
	}
	candidate := os.Args[1]
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	ref := filepath.Join(dir, "756E.go")

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
