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

const target = "CODEFORCES"

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	for idx, s := range tests {
		input := s + "\n"
		refOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		candOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(refOut) != strings.TrimSpace(candOut) {
			fmt.Fprintf(os.Stderr, "test %d mismatch\ninput:%s\nexpected:%s\ngot:%s", idx+1, s, refOut, candOut)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func buildReference() (string, error) {
	const refName = "./ref_538A.bin"
	cmd := exec.Command("go", "build", "-o", refName, "538A.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return refName, nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstdout:\n%s\nstderr:\n%s", err, stdout.String(), stderr.String())
	}
	return stdout.String(), nil
}

func buildTests() []string {
	tests := []string{
		target + target,
		target[:5] + "XYZ" + target[5:],
		"XYZ" + target[1:],
		target[:len(target)-1] + "Z",
		"ABCDE",
	}

	alphabet := "CODEFORCESXYZ"
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 500 {
		length := rng.Intn(20) + 1
		var sb strings.Builder
		for i := 0; i < length; i++ {
			sb.WriteByte(alphabet[rng.Intn(len(alphabet))])
		}
		if sb.String() == target {
			continue
		}
		tests = append(tests, sb.String())
	}
	return tests
}
