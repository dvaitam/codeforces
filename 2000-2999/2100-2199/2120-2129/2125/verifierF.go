package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// refSource points to the local reference solution to avoid GOPATH resolution.
const refSource = "2125F.go"

type testBatch struct {
	text    string
	answers int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for idx, tc := range tests {
		refOut, err := runProgram(exec.Command(refBin), tc.text)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s\n", idx+1, err, tc.text)
			os.Exit(1)
		}
		refVals, err := parseOutputs(refOut, tc.answers)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not parse reference output on test %d: %v\noutput:\n%s\n", idx+1, err, refOut)
			os.Exit(1)
		}

		candCmd := commandFor(candidate)
		candOut, err := runProgram(candCmd, tc.text)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%s\nstdout/stderr:\n%s\n", idx+1, err, tc.text, candOut)
			os.Exit(1)
		}
		candVals, err := parseOutputs(candOut, tc.answers)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid candidate output on test %d: %v\noutput:\n%s\n", idx+1, err, candOut)
			os.Exit(1)
		}

		for i := 0; i < tc.answers; i++ {
			if refVals[i] != candVals[i] {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d case %d\ninput:\n%s\nexpected: %d\nfound: %d\n", idx+1, i+1, tc.text, refVals[i], candVals[i])
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2125F-ref-*")
	if err != nil {
		return "", err
	}
	tmpPath := tmp.Name()
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmpPath, filepath.Clean(refSource))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmpPath)
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmpPath, nil
}

func commandFor(path string) *exec.Cmd {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	case ".js":
		return exec.Command("node", path)
	default:
		return exec.Command(path)
	}
}

func runProgram(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseOutputs(out string, expected int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d outputs, got %d", expected, len(fields))
	}
	res := make([]int64, expected)
	for i, f := range fields {
		v, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer at position %d: %v", i+1, err)
		}
		res[i] = v
	}
	return res, nil
}

func generateTests() []testBatch {
	tests := []testBatch{sampleTest()}
	tests = append(tests, fixedTests()...)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 140 {
		tests = append(tests, randomBatch(rng))
	}
	return tests
}

func sampleTest() testBatch {
	// Combined from the statement samples.
	text := "" +
		"6\n" +
		"dockerdockerxxxxxx\n" +
		"3\n" +
		"3 3\n" +
		"2 4\n" +
		"1 5\n" +
		"ljglsjfkdieufj\n" +
		"5\n" +
		"1 5\n" +
		"3 3\n" +
		"2 4\n" +
		"3 7\n" +
		"2 9\n" +
		"dockerdockerdockerdockzzdockzz\n" +
		"4\n" +
		"1 1\n" +
		"1 1\n" +
		"4 5\n" +
		"4 5\n" +
		"docker\n" +
		"5\n" +
		"1 1\n" +
		"2 2\n" +
		"3 3\n" +
		"4 4\n" +
		"5 5\n" +
		"ddddddoooooocccccckkkkkkeeeeeerrrrrr\n" +
		"10\n" +
		"1 200\n" +
		"500 600\n" +
		"1 600\n" +
		"6 6\n" +
		"6 6\n" +
		"500 2000\n" +
		"6 400\n" +
		"89 90\n" +
		"4 7\n" +
		"1 10\n" +
		"dockerdockerdockerdockzzdockzz\n" +
		"4\n" +
		"2 2\n" +
		"2 4\n" +
		"5 5\n" +
		"4 5\n"
	return testBatch{text: text, answers: 6}
}

func fixedTests() []testBatch {
	// Edge shapes: short strings, zero or huge acceptable ranges, and many overlapping intervals.
	text1 := "" +
		"5\n" +
		"d\n" +
		"3\n" +
		"1 1\n" +
		"1 2\n" +
		"2 3\n" +
		"docker\n" +
		"2\n" +
		"1 1\n" +
		"1 10\n" +
		"kkkkkk\n" +
		"4\n" +
		"1 1\n" +
		"1 1\n" +
		"2 2\n" +
		"3 3\n" +
		"dockerdocker\n" +
		"3\n" +
		"1 2\n" +
		"2 2\n" +
		"3 3\n" +
		"aaaaaaaaaaaa\n" +
		"1\n" +
		"100 200\n"

	text2 := "" +
		"3\n" +
		"dockerr\n" +
		"5\n" +
		"1 1\n" +
		"2 2\n" +
		"3 3\n" +
		"4 4\n" +
		"5 5\n" +
		"codersdocker\n" +
		"4\n" +
		"1 2\n" +
		"2 3\n" +
		"3 3\n" +
		"4 5\n" +
		"zzzzzzzz\n" +
		"2\n" +
		"1 1\n" +
		"1 2\n"

	return []testBatch{
		{text: text1, answers: 5},
		{text: text2, answers: 3},
	}
}

func randomBatch(rng *rand.Rand) testBatch {
	t := rng.Intn(6) + 1

	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	for i := 0; i < t; i++ {
		strLen := rng.Intn(150) + 1
		s := make([]byte, strLen)
		for j := 0; j < strLen; j++ {
			s[j] = byte('a' + rng.Intn(26))
		}
		sb.Write(s)
		sb.WriteByte('\n')

		n := rng.Intn(60) + 1
		fmt.Fprintf(&sb, "%d\n", n)
		maxC := strLen / 6
		for j := 0; j < n; j++ {
			// allow ranges to go beyond maxC to test clamping logic
			l := rng.Intn(maxC+5) - 2
			r := l + rng.Intn(maxC+5) + rng.Intn(3)
			if l < 1 {
				l = 1
			}
			if r < l {
				r = l
			}
			// occasionally make very large upper bounds
			if rng.Intn(10) == 0 {
				r += rng.Intn(1_000_000_000)
			}
			fmt.Fprintf(&sb, "%d %d\n", l, r)
		}
	}

	return testBatch{text: sb.String(), answers: t}
}
