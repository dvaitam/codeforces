package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type testCase struct{ input string }

func compileRef() (string, error) {
	exe := filepath.Join(os.TempDir(), "ref44I")
	cmd := exec.Command("go", "build", "-o", exe, "44I.go")
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return exe, nil
}

func generateRandomCase(rng *rand.Rand) testCase {
	n := rng.Intn(5) + 1
	return testCase{input: fmt.Sprintf("%d\n", n)}
}

func runCase(bin, ref string, tc testCase) error {
	run := func(path string) (string, error) {
		cmd := exec.Command(path)
		cmd.Stdin = strings.NewReader(tc.input)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		err := cmd.Run()
		return strings.TrimSpace(out.String()), err
	}
	exp, err := run(ref)
	if err != nil {
		return fmt.Errorf("ref error: %v", err)
	}
	got, err := run(bin)
	if err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, got)
	}

	expLines := strings.Split(exp, "\n")
	gotLines := strings.Split(got, "\n")

	if len(expLines) != len(gotLines) {
		return fmt.Errorf("expected %d lines, got %d lines\nexpected %q got %q", len(expLines), len(gotLines), exp, got)
	}

	if expLines[0] != gotLines[0] {
		return fmt.Errorf("expected count %q got %q", expLines[0], gotLines[0])
	}

	if len(expLines) > 1 {
		sort.Strings(expLines[1:])
		sort.Strings(gotLines[1:])
	}

	sortedExp := strings.Join(expLines, "\n")
	sortedGot := strings.Join(gotLines, "\n")

	if sortedExp != sortedGot {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierI.go /path/to/binary")
		os.Exit(1)
	}
	ref, err := compileRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	cases := []testCase{{input: "1\n"}}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateRandomCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, ref, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
