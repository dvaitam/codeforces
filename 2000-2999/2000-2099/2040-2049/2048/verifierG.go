package main

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var refSource = getRefSource()

func getRefSource() string {
	if v := os.Getenv("REFERENCE_SOURCE_PATH"); v != "" {
		return v
	}
	return "2000-2999/2000-2099/2040-2049/2048/2048G.go"
}

func generateInput() []byte {
	rng := rand.New(rand.NewSource(2048))
	type tc struct{ n, m, v int }
	var cases []tc

	// Fixed test cases from the problem
	cases = append(cases, tc{2, 2, 2})
	cases = append(cases, tc{2, 3, 4})
	cases = append(cases, tc{11, 45, 14})

	// Small exhaustive cases
	for n := 1; n <= 5; n++ {
		for v := 1; v <= 6; v++ {
			if n*v > 1000000 {
				continue
			}
			for _, m := range []int{1, 2, 3, 5, 10, 100, 1000000000} {
				cases = append(cases, tc{n, m, v})
			}
		}
	}

	// Random cases
	for i := 0; i < 50; i++ {
		n := rng.Intn(10) + 1
		v := rng.Intn(100) + 1
		if n*v > 1000000 {
			v = 1000000 / n
			if v < 1 {
				v = 1
			}
		}
		m := rng.Intn(1000000000) + 1
		cases = append(cases, tc{n, m, v})
	}

	var buf bytes.Buffer
	fmt.Fprintln(&buf, len(cases))
	for _, c := range cases {
		fmt.Fprintf(&buf, "%d %d %d\n", c.n, c.m, c.v)
	}
	return buf.Bytes()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	input := generateInput()

	refBin, err := buildReference()
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fail("reference execution failed: %v", err)
	}
	refAnswers, err := parseAnswers(refOut)
	if err != nil {
		fail("failed to parse reference output: %v", err)
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fail("candidate execution failed: %v", err)
	}
	candAnswers, err := parseAnswers(candOut)
	if err != nil {
		fail("failed to parse candidate output: %v", err)
	}

	if len(refAnswers) != len(candAnswers) {
		fail("expected %d answers, got %d", len(refAnswers), len(candAnswers))
	}
	for i := range refAnswers {
		if refAnswers[i] != candAnswers[i] {
			fail("test %d: expected %d got %d", i+1, refAnswers[i], candAnswers[i])
		}
	}

	fmt.Println("OK")
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2048G-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runProgram(path string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = bytes.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func parseAnswers(out string) ([]int64, error) {
	reader := strings.NewReader(out)
	var res []int64
	for {
		var val int64
		_, err := fmt.Fscan(reader, &val)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("failed to parse integer: %v", err)
		}
		res = append(res, val)
	}
	return res, nil
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
