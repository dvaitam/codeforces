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
	"time"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for caseNum := 0; caseNum < 100; caseNum++ {
		input := generateInput(rng)

		refOut, err := runProgram(refBin, []byte(input))
		if err != nil {
			fail("reference execution failed on case %d: %v", caseNum+1, err)
		}
		refAns, err := parseAnswers(refOut)
		if err != nil {
			fail("failed to parse reference output on case %d: %v", caseNum+1, err)
		}

		candOut, err := runProgram(candidate, []byte(input))
		if err != nil {
			fail("candidate execution failed on case %d: %v", caseNum+1, err)
		}
		candAns, err := parseAnswers(candOut)
		if err != nil {
			fail("failed to parse candidate output on case %d: %v", caseNum+1, err)
		}

		if len(refAns) != len(candAns) {
			fail("case %d: expected %d answers, got %d", caseNum+1, len(refAns), len(candAns))
		}
		for i := range refAns {
			if refAns[i] != candAns[i] {
				fail("case %d test %d: expected %d got %d\ninput:\n%s", caseNum+1, i+1, refAns[i], candAns[i], input)
			}
		}
	}

	fmt.Println("OK")
}

func generateInput(rng *rand.Rand) string {
	t := rng.Intn(10) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	for i := 0; i < t; i++ {
		// n can be from 1 to 10^18
		var n int64
		switch rng.Intn(4) {
		case 0:
			n = rng.Int63n(10) + 1
		case 1:
			n = rng.Int63n(100) + 1
		case 2:
			n = rng.Int63n(1000000) + 1
		case 3:
			n = rng.Int63n(1000000000000000000) + 1
		}
		fmt.Fprintf(&sb, "%d\n", n)
	}
	return sb.String()
}

func buildReference() (string, error) {
	refSource := os.Getenv("REFERENCE_SOURCE_PATH")
	if refSource == "" {
		return "", fmt.Errorf("REFERENCE_SOURCE_PATH not set")
	}

	tmp := filepath.Join(os.TempDir(), "2043A-ref")
	cmd := exec.Command("go", "build", "-o", tmp, refSource)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp)
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp, nil
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
