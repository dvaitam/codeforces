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
	"time"
)

type testCase struct{ a, l int64 }

func compileRef() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	exe := filepath.Join(os.TempDir(), "ref109E")
	cmd := exec.Command("go", "build", "-o", exe, filepath.Join(dir, "109E.go"))
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return exe, nil
}

func generateCase(rng *rand.Rand) testCase {
	a := rng.Int63n(1000000) + 1
	l := rng.Int63n(1000) + 1
	return testCase{a, l}
}

func runCase(bin, ref string, tc testCase) error {
	input := fmt.Sprintf("%d %d\n", tc.a, tc.l)
	run := func(path string) (string, error) {
		cmd := exec.Command(path)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		err := cmd.Run()
		return strings.TrimSpace(out.String()), err
	}
	expected, err := run(ref)
	if err != nil {
		return fmt.Errorf("ref error: %v", err)
	}
	got, err := run(bin)
	if err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, got)
	}
	if expected != got {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref, err := compileRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference solution:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCase{{1, 1}, {4, 3}}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, ref, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput: %v\n", i+1, err, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
