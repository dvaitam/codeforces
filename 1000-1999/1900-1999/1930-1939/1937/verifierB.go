package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func buildRef() (string, error) {
	ref := "refB.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1937B.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return ref, nil
}

func runBinary(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type testCase struct {
	n     int
	s1    string
	s2    string
	input string
}

func generateCase(rng *rand.Rand) testCase {
	n := rng.Intn(20) + 2
	b1 := make([]byte, n)
	b2 := make([]byte, n)
	for i := 0; i < n; i++ {
		b1[i] = byte('0' + rng.Intn(2))
		b2[i] = byte('0' + rng.Intn(2))
	}
	input := fmt.Sprintf("1\n%d\n%s\n%s\n", n, string(b1), string(b2))
	return testCase{n: n, s1: string(b1), s2: string(b2), input: input}
}

func runCase(candidate, ref string, tc testCase) error {
	want, err := runBinary(ref, tc.input)
	if err != nil {
		return fmt.Errorf("reference failed: %v", err)
	}
	got, err := runBinary(candidate, tc.input)
	if err != nil {
		return err
	}
	if strings.TrimSpace(got) != strings.TrimSpace(want) {
		return fmt.Errorf("expected:\n%s\ngot:\n%s", want, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCase(rng)
		if err := runCase(candidate, ref, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
