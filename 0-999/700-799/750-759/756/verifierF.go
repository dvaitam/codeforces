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

func randNum(r *rand.Rand) string {
	return strconv.Itoa(r.Intn(99) + 1)
}

func randTerm(r *rand.Rand, depth int) string {
	typ := r.Intn(3)
	if depth <= 0 && typ > 0 {
		typ = 0
	}
	switch typ {
	case 0:
		return randNum(r)
	case 1:
		l := r.Intn(50) + 1
		r2 := l + r.Intn(50)
		return fmt.Sprintf("%d-%d", l, r2)
	default:
		return randNum(r) + "(" + randExpr(r, depth-1) + ")"
	}
}

func randExpr(r *rand.Rand, depth int) string {
	terms := r.Intn(3) + 1
	res := randTerm(r, depth)
	for i := 1; i < terms; i++ {
		res += "+" + randTerm(r, depth)
	}
	return res
}

func generateTests() []string {
	r := rand.New(rand.NewSource(47))
	tests := make([]string, 0, 20)
	for i := 0; i < 20; i++ {
		expr := randExpr(r, 2)
		tests = append(tests, expr+"\n")
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: verifierF.go <binary>")
		os.Exit(1)
	}
	candidate := os.Args[1]
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	ref := filepath.Join(dir, "756F.go")

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
