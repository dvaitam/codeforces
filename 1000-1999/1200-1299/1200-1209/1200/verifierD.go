package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Test struct {
	input string
}

func buildRef() (string, error) {
	ref := "refD.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1200D.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

func runExe(path, input string) (string, error) {
	if !strings.Contains(path, "/") {
		path = "./" + path
	}
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func randCase(rng *rand.Rand) Test {
	n := rng.Intn(6) + 1
	k := rng.Intn(n) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if rng.Intn(2) == 0 {
				sb.WriteByte('B')
			} else {
				sb.WriteByte('W')
			}
		}
		sb.WriteByte('\n')
	}
	return Test{sb.String()}
}

func genTests() []Test {
	rng := rand.New(rand.NewSource(3))
	tests := make([]Test, 0, 105)
	for i := 0; i < 100; i++ {
		tests = append(tests, randCase(rng))
	}
	tests = append(tests, Test{"1 1\nB\n"})
	tests = append(tests, Test{"1 1\nW\n"})
	tests = append(tests, Test{"2 1\nBW\nWB\n"})
	tests = append(tests, Test{"2 2\nBB\nBB\n"})
	tests = append(tests, Test{"3 1\nWWW\nWWW\nWWW\n"})
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)
	tests := genTests()
	for i, tc := range tests {
		exp, err := runExe(ref, tc.input)
		if err != nil {
			fmt.Printf("reference runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runExe(bin, tc.input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if exp != got {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%s\nGot:\n%s\n", i+1, tc.input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
