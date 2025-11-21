package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const referenceSolutionRel = "2000-2999/2000-2099/2010-2019/2010/2010C2.go"

var referenceSolutionPath string

func init() {
	referenceSolutionPath = referenceSolutionRel
	if _, file, _, ok := runtime.Caller(0); ok {
		dir := filepath.Dir(file)
		candidate := filepath.Join(dir, "2010C2.go")
		if _, err := os.Stat(candidate); err == nil {
			referenceSolutionPath = candidate
			return
		}
	}
	if abs, err := filepath.Abs(referenceSolutionRel); err == nil {
		if _, err := os.Stat(abs); err == nil {
			referenceSolutionPath = abs
		}
	}
}

func runProgram(path, input string) (string, error) {
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
	return out.String(), err
}

func buildReferenceBinary() (string, func(), error) {
	if referenceSolutionPath == "" {
		return "", nil, fmt.Errorf("reference solution path not set")
	}
	if _, err := os.Stat(referenceSolutionPath); err != nil {
		return "", nil, fmt.Errorf("reference solution not found: %v", err)
	}
	tmpDir, err := os.MkdirTemp("", "2010C2-ref")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "ref_2010C2")
	cmd := exec.Command("go", "build", "-o", binPath, referenceSolutionPath)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, out.String())
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return binPath, cleanup, nil
}

func parseOutput(output string) (string, string, error) {
	scanner := bufio.NewScanner(strings.NewReader(output))
	var lines []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}
	if len(lines) == 0 {
		return "", "", fmt.Errorf("empty output")
	}
	verdict := strings.ToUpper(lines[0])
	if verdict != "YES" && verdict != "NO" {
		return "", "", fmt.Errorf("first line must be YES/NO, got %q", lines[0])
	}
	var s string
	if verdict == "YES" {
		if len(lines) < 2 {
			return "", "", fmt.Errorf("YES without string line")
		}
		s = lines[1]
	}
	return verdict, s, nil
}

func checkWitness(t, s string) error {
	n := len(t)
	m := len(s)
	if m < 2 || m >= n {
		return fmt.Errorf("invalid witness length")
	}
	for k := 1; k < m; k++ {
		if s[k:] == s[:m-k] {
			if t == s[:k]+s {
				return nil
			}
		}
	}
	return fmt.Errorf("provided string is not a valid witness")
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC2.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	input, err := os.ReadFile(os.Stdin.Name())
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to read input:", err)
		os.Exit(1)
	}
	t := strings.TrimSpace(string(input))
	if t == "" {
		fmt.Fprintln(os.Stderr, "empty input string")
		os.Exit(1)
	}

	refBin, cleanup, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	refOut, err := runProgram(refBin, t+"\n")
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}
	refVerdict, refWitness, err := parseOutput(refOut)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}

	userOut, err := runProgram(bin, t+"\n")
	if err != nil {
		fmt.Fprintf(os.Stderr, "runtime error: %v\noutput:\n%s\n", err, userOut)
		os.Exit(1)
	}
	userVerdict, userWitness, err := parseOutput(userOut)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse participant output: %v\noutput:\n%s\n", err, userOut)
		os.Exit(1)
	}

	if refVerdict == "NO" {
		if userVerdict != "NO" {
			fmt.Fprintf(os.Stderr, "expected NO but participant answered YES\nreference:\n%s\nparticipant:\n%s\n", refOut, userOut)
			os.Exit(1)
		}
		fmt.Println("All tests passed.")
		return
	}

	if userVerdict != "YES" {
		fmt.Fprintf(os.Stderr, "expected YES but participant answered NO\nreference:\n%s\nparticipant:\n%s\n", refOut, userOut)
		os.Exit(1)
	}

	if err := checkWitness(t, userWitness); err != nil {
		fmt.Fprintf(os.Stderr, "invalid witness: %v\noutput:\n%s\n", err, userOut)
		os.Exit(1)
	}
	fmt.Println("All tests passed.")
}
