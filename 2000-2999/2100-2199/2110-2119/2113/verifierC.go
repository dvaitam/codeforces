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

const refSource = "2000-2999/2100-2199/2110-2119/2113/2113C.go"

type testBatch struct {
	text    string
	answers int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/candidate")
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
	tmp, err := os.CreateTemp("", "2113C-ref-*")
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
	for len(tests) < 150 {
		tests = append(tests, randomBatch(rng))
	}
	return tests
}

func sampleTest() testBatch {
	text := "" +
		"3\n" +
		"2 3 1\n" +
		"#.#\n" +
		"g.g\n" +
		"2 3 2\n" +
		"#.#\n" +
		"g.g\n" +
		"3 4 2\n" +
		".gg.\n" +
		"..#g\n" +
		"##.\n"
	return testBatch{text: text, answers: 3}
}

func fixedTests() []testBatch {
	// Tight k, wide k, sparse gold, all gold unreachable behind stones, and multiple empty cells.
	text1 := "" +
		"4\n" +
		"1 5 1\n" +
		".g#g.\n" +
		"4 4 1\n" +
		"....\n" +
		".g#.\n" +
		"g#gg\n" +
		"#..g\n" +
		"3 3 3\n" +
		"g#g\n" +
		"#.#\n" +
		"g#g\n" +
		"2 2 2\n" +
		"..\n" +
		"..\n"

	text2 := "" +
		"2\n" +
		"5 5 2\n" +
		"g.g.g\n" +
		".###.\n" +
		"g.g.g\n" +
		".###.\n" +
		"g.g.g\n" +
		"3 6 1\n" +
		"#g#g#g\n" +
		".#.#.#\n" +
		"#g#g#g\n"

	return []testBatch{
		{text: text1, answers: 4},
		{text: text2, answers: 2},
	}
}

func randomBatch(rng *rand.Rand) testBatch {
	t := rng.Intn(6) + 1

	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	total := 0
	for i := 0; i < t; i++ {
		n := rng.Intn(20) + 1
		m := rng.Intn(20) + 1
		k := rng.Intn(10) + 1
		total++

		fmt.Fprintf(&sb, "%d %d %d\n", n, m, k)
		emptyMade := false
		for r := 0; r < n; r++ {
			row := make([]byte, m)
			for c := 0; c < m; c++ {
				if !emptyMade && r == n-1 && c == m-1 {
					row[c] = '.'
					emptyMade = true
					continue
				}
				switch v := rng.Intn(100); {
				case v < 20:
					row[c] = '#'
				case v < 60:
					row[c] = 'g'
				default:
					row[c] = '.'
					emptyMade = true
				}
			}
			sb.Write(row)
			sb.WriteByte('\n')
		}
	}

	return testBatch{text: sb.String(), answers: total}
}
