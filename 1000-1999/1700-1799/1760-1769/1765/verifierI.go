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

const refSource = "1000-1999/1700-1799/1760-1769/1765/1765I.go"

type testCase struct {
	name  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierI.go /path/to/binary")
		os.Exit(1)
	}

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()

	candidate := os.Args[1]
	for idx, tc := range tests {
		refOut, err := runBinary(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		refVal, err := parseInt(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candVal, err := parseInt(candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, candOut)
			os.Exit(1)
		}

		if candVal != refVal {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s): expected %d, got %d\ninput:\n%s\n", idx+1, tc.name, refVal, candVal, tc.input)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "1765I-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmp.Name())
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

func parseInt(out string) (int, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return 0, fmt.Errorf("empty output")
	}
	val, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q: %v", fields[0], err)
	}
	return val, nil
}

func generateTests() []testCase {
	tests := []testCase{
		buildCase("simple-no-pieces", 1, 1, 5, 5, nil),
		buildCase("knight-block", 2, 2, 3, 3, []piece{{'N', 3, 3}}),
		buildCase("rook-wall", 10, 4, 12, 5, []piece{{'R', 11, 4}}),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		xs := rng.Intn(100) + 1
		ys := rng.Intn(8) + 1
		xt := xs + rng.Intn(20) - 10
		if xt < 1 {
			xt = 1
		}
		yt := rng.Intn(8) + 1
		n := rng.Intn(6)
		pieces := make([]piece, n)
		types := []byte{'K', 'Q', 'R', 'B', 'N'}
		for j := 0; j < n; j++ {
			pieces[j] = piece{
				t: types[rng.Intn(len(types))],
				x: xs + rng.Intn(20) - 10,
				y: rng.Intn(8) + 1,
			}
			if pieces[j].x < 1 {
				pieces[j].x = 1
			}
		}
		tests = append(tests, buildCase(fmt.Sprintf("random-%d", i+1), xs, ys, xt, yt, pieces))
	}
	return tests
}

type piece struct {
	t byte
	x int
	y int
}

func buildCase(name string, xs, ys, xt, yt int, pieces []piece) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n%d %d\n%d\n", xs, ys, xt, yt, len(pieces))
	for _, p := range pieces {
		fmt.Fprintf(&sb, "%c %d %d\n", p.t, p.x, p.y)
	}
	return testCase{name: name, input: sb.String()}
}
