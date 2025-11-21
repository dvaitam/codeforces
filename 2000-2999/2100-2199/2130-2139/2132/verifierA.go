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

const refSource = "2000-2999/2100-2199/2130-2139/2132/2132A.go"

func buildBinary(path string) (string, func(), error) {
	if !strings.HasSuffix(path, ".go") {
		return path, func() {}, nil
	}
	tmp, err := os.CreateTemp("", "cf-2132A-*")
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

func parseOutput(out string, expected int) ([]string, error) {
	lines := strings.Fields(out)
	if len(lines) != expected {
		return nil, fmt.Errorf("expected %d outputs, got %d", expected, len(lines))
	}
	return lines, nil
}

func fixedTests() []string {
	return []string{
		"4\n2\not\n2\nad\nDV\n3\nefo\n7\nrdcoecs\nDVDVDVD\n3\naca\n4\nbbaa\nDVDV\n3\nbiz\n4\nabon\nVVDD\n",
		"2\n1\na\n1\nb\nV\n1\nz\n1\ny\nD\n",
	}
}

func randStringLetters(rng *rand.Rand, n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = 'a' + byte(rng.Intn(26))
	}
	return string(b)
}

func randomCase(rng *rand.Rand) string {
	// n, m in [1,10]
	n := rng.Intn(10) + 1
	m := rng.Intn(10) + 1
	a := randStringLetters(rng, n)
	b := randStringLetters(rng, m)
	cBytes := make([]byte, m)
	for i := 0; i < m; i++ {
		if rng.Intn(2) == 0 {
			cBytes[i] = 'V'
		} else {
			cBytes[i] = 'D'
		}
	}
	return fmt.Sprintf("%d\n%s\n%d\n%s\n%s\n", n, a, m, b, string(cBytes))
}

func randomInput(rng *rand.Rand) string {
	t := rng.Intn(25) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		sb.WriteString(randomCase(rng))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
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
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 40; i++ {
		tests = append(tests, randomInput(rng))
	}

	for idx, input := range tests {
		var t int
		if _, err := fmt.Fscan(strings.NewReader(input), &t); err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse t for case %d: %v\n", idx+1, err)
			os.Exit(1)
		}

		expOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on case %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		expLines, err := parseOutput(expOut, t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output parse failed on case %d: %v\noutput:\n%s", idx+1, err, expOut)
			os.Exit(1)
		}

		gotOut, err := runProgram(candBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on case %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		gotLines, err := parseOutput(gotOut, t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output parse failed on case %d: %v\noutput:\n%s", idx+1, err, gotOut)
			os.Exit(1)
		}

		for i := 0; i < t; i++ {
			if expLines[i] != gotLines[i] {
				fmt.Fprintf(os.Stderr, "mismatch on case %d test %d: expected %s got %s\ninput:\n%s", idx+1, i+1, expLines[i], gotLines[i], input)
				os.Exit(1)
			}
		}
	}

	fmt.Println("All tests passed")
}
