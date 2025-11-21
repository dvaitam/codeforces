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

const refSource = "2000-2999/2000-2099/2070-2079/2074/2074F.go"

func buildBinary(path string) (string, func(), error) {
	if !strings.HasSuffix(path, ".go") {
		return path, func() {}, nil
	}
	tmp, err := os.CreateTemp("", "cf-2074F-*")
	if err != nil {
		return "", func() {}, err
	}
	tmp.Close()
	os.Remove(tmp.Name())

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Base(path))
	cmd.Dir = filepath.Dir(path)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", func() {}, fmt.Errorf("go build error: %v\n%s", err, stderr.String())
	}
	return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func fixedTests() []string {
	return []string{
		"5\n0 1 1 2\n0 2 0 2\n1 3 1 3\n0 2 1 5\n9 98 244 353\n",
		"3\n0 1 0 1\n0 1 0 2\n2 4 3 5\n",
	}
}

func randomSegment(rng *rand.Rand) (int, int) {
	l := rng.Intn(1_000_000)
	r := l + 1 + rng.Intn(1_000_000-l)
	return l, r
}

func generateRandomInput(rng *rand.Rand) string {
	t := rng.Intn(25) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		l1, r1 := randomSegment(rng)
		l2, r2 := randomSegment(rng)
		// occasionally swap to ensure robustness
		if rng.Intn(5) == 0 {
			l1, r1 = r1, l1
		}
		if rng.Intn(5) == 0 {
			l2, r2 = r2, l2
		}
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", l1, r1, l2, r2))
	}
	return sb.String()
}

func powerOfTwoCase(k int) string {
	size := 1 << k
	return fmt.Sprintf("1\n0 %d 0 %d\n", size, size)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	candPath, err := filepath.Abs(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to resolve candidate path: %v\n", err)
		os.Exit(1)
	}
	refPath, err := filepath.Abs(refSource)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to resolve reference path: %v\n", err)
		os.Exit(1)
	}

	refBin, cleanupRef, err := buildBinary(refPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer cleanupRef()

	candBin, cleanupCand, err := buildBinary(candPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build candidate: %v\n", err)
		os.Exit(1)
	}
	defer cleanupCand()

	tests := fixedTests()
	for k := 0; k <= 20; k += 4 {
		tests = append(tests, powerOfTwoCase(k))
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 40; i++ {
		tests = append(tests, generateRandomInput(rng))
	}

	for idx, input := range tests {
		exp, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on case %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		got, err := runProgram(candBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on case %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		if exp != got {
			fmt.Fprintf(os.Stderr, "mismatch on case %d\nexpected:\n%s\n\ngot:\n%s\ninput:\n%s", idx+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
