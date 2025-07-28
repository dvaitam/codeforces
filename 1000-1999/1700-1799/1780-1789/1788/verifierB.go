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

func buildSolution() (string, error) {
	exe := "./_verifier_solB"
	cmd := exec.Command("go", "build", "-o", exe, "1788B.go")
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return exe, nil
}

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	timer := time.AfterFunc(2*time.Second, func() { cmd.Process.Kill() })
	err := cmd.Run()
	timer.Stop()
	if err != nil {
		return out.String(), fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type testCase struct {
	n     int64
	input string
}

func genTest(rng *rand.Rand) testCase {
	n := rng.Int63n(1_000_000_000) + 1
	input := fmt.Sprintf("1\n%d\n", n)
	return testCase{n: n, input: input}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	solExe, err := buildSolution()
	if err != nil {
		fmt.Println("failed to build reference solution:", err)
		return
	}
	defer os.Remove(solExe)

	fixed := []testCase{{n: 1, input: "1\n1\n"}, {n: 9, input: "1\n9\n"}, {n: 10, input: "1\n10\n"}, {n: 99, input: "1\n99\n"}}
	tests := make([]testCase, 0, 100)
	tests = append(tests, fixed...)
	for len(tests) < 100 {
		tests = append(tests, genTest(rand.New(rand.NewSource(rand.Int63()))))
	}
	for i, tc := range tests {
		exp, err := runBinary(solExe, tc.input)
		if err != nil {
			fmt.Printf("reference solution failed on test %d: %v\n", i+1, err)
			return
		}
		out, err := runBinary(bin, tc.input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			return
		}
		if strings.TrimSpace(exp) != strings.TrimSpace(out) {
			fmt.Printf("wrong answer on test %d\nexpected:\n%s\nactual:\n%s\n", i+1, exp, out)
			return
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
