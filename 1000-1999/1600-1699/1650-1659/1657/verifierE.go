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

type testCaseE struct{ n, k int }

func genTestsE() []testCaseE {
	rng := rand.New(rand.NewSource(46))
	tests := []testCaseE{{2, 1}, {3, 2}, {4, 3}}
	for len(tests) < 100 {
		n := rng.Intn(8) + 2
		k := rng.Intn(8) + 1
		tests = append(tests, testCaseE{n, k})
	}
	return tests
}

func compileRef() (string, error) {
	exe := filepath.Join(os.TempDir(), "1657E_ref")
	cmd := exec.Command("go", "build", "-o", exe, "1657E.go")
	cmd.Dir = "."
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("compile reference failed: %v\n%s", err, out)
	}
	return exe, nil
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("runtime error: %v", err)
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref, err := compileRef()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer os.Remove(ref)
	rand.Seed(time.Now().UnixNano())
	tests := genTestsE()
	for i, tc := range tests {
		input := fmt.Sprintf("%d %d\n", tc.n, tc.k)
		exp, err := runBinary(ref, input)
		if err != nil {
			fmt.Printf("reference failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: %v\noutput:\n%s", i+1, err, got)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(exp) {
			fmt.Printf("test %d failed: expected %q got %q\ninput:\n%s", i+1, strings.TrimSpace(exp), strings.TrimSpace(got), input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
