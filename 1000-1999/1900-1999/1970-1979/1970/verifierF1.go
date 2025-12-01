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

const refSource = "./1970F1.go"

type testCase struct {
	name  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF1.go /path/to/binary")
		os.Exit(1)
	}

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}

	tests := generateTests()
	candidate := os.Args[1]

	for idx, tc := range tests {
		refOut, err := runBinary(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		refEvents, refFinal, err := parseOutput(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candEvents, candFinal, err := parseOutput(candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, candOut)
			os.Exit(1)
		}

		if candFinal != refFinal {
			fmt.Fprintf(os.Stderr, "wrong final score on test %d (%s): expected %s got %s\ninput:\n%s\n", idx+1, tc.name, refFinal, candFinal, tc.input)
			os.Exit(1)
		}

		if len(candEvents) != len(refEvents) {
			fmt.Fprintf(os.Stderr, "wrong number of goal events on test %d (%s): expected %d got %d\ninput:\n%s\n", idx+1, tc.name, len(refEvents), len(candEvents), tc.input)
			os.Exit(1)
		}
		for i := range refEvents {
			if candEvents[i] != refEvents[i] {
				fmt.Fprintf(os.Stderr, "goal event mismatch on test %d (%s) at event %d: expected %s, got %s\ninput:\n%s\n", idx+1, tc.name, i+1, refEvents[i], candEvents[i], tc.input)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "1970F1-ref-*")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmp.Name())
	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return tmp.Name(), nil
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return stdout.String() + stderr.String(), err
	}
	return stdout.String(), nil
}

func parseOutput(out string) ([]string, string, error) {
	lines := filterNonEmpty(strings.Split(out, "\n"))
	if len(lines) == 0 {
		return nil, "", fmt.Errorf("empty output")
	}
	final := lines[len(lines)-1]
	if !strings.HasPrefix(final, "FINAL SCORE:") {
		return nil, "", fmt.Errorf("missing final score line")
	}
	events := lines[:len(lines)-1]
	for _, e := range events {
		if !strings.Contains(e, "GOAL") {
			return nil, "", fmt.Errorf("invalid goal event format: %s", e)
		}
	}
	return events, final, nil
}

func filterNonEmpty(lines []string) []string {
	var res []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			res = append(res, line)
		}
	}
	return res
}

func generateTests() []testCase {
	tests := []testCase{
		buildCase("simple", `
3 3
.. .. ..
.Q .. ..
R0 .. BG
5
R0 C .Q
R0 T
.Q U
R0 C .Q
R0 T
`),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	dirs := []string{"U", "D", "L", "R"}
	entities := []string{".Q", "R0", "B0"}

	for i := 0; i < 20; i++ {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", 3, 3)
		fmt.Fprintf(&sb, ".Q .. BG\nRG R0 ..\n.. B0 ..\n")
		T := rng.Intn(10) + 1
		fmt.Fprintf(&sb, "%d\n", T)
		for j := 0; j < T; j++ {
			entity := entities[rng.Intn(len(entities))]
			action := dirs[rng.Intn(len(dirs))]
			fmt.Fprintf(&sb, "%s %s\n", entity, action)
		}
		tests = append(tests, testCase{name: fmt.Sprintf("random-%d", i+1), input: sb.String()})
	}
	return tests
}

func buildCase(name, input string) testCase {
	return testCase{name: name, input: strings.TrimSpace(input) + "\n"}
}
