package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func compileRef() (string, error) {
	exe := filepath.Join(os.TempDir(), "ref935E")
	cmd := exec.Command("go", "build", "-o", exe, "935E.go")
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return exe, nil
}

func runProg(path, input string) (string, error) {
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

func genExpr(r *rand.Rand, depth int) (string, int) {
	if depth > 3 || r.Intn(2) == 0 {
		d := r.Intn(9) + 1
		return fmt.Sprintf("%d", d), 0
	}
	left, opsL := genExpr(r, depth+1)
	right, opsR := genExpr(r, depth+1)
	return "(" + left + "?" + right + ")", opsL + opsR + 1
}

func generateCase(r *rand.Rand) string {
	expr, ops := genExpr(r, 0)
	P := r.Intn(ops + 1)
	M := ops - P
	return fmt.Sprintf("%s\n%d %d\n", expr, P, M)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref, err := compileRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in := generateCase(rng)
		expect, err := runProg(ref, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		out, err := runProg(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "test %d failed.\nInput:\n%sExpected:\n%s\nGot:\n%s\n", i+1, in, expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
