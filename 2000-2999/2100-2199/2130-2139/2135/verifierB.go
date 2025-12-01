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
)

const refSource = "./2135B.go"

type testCase struct {
	anchors [][2]int64
	ansX    int64
	ansY    int64
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("usage: go run verifierB.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := args[len(args)-1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	input := buildInput(tests)

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference failed: %v\n", err)
		os.Exit(1)
	}
	exp, err := parseOutput(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\n%s", err, refOut)
		os.Exit(1)
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n", err)
		os.Exit(1)
	}
	got, err := parseOutput(candOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse candidate output: %v\n%s", err, candOut)
		os.Exit(1)
	}

	for i := range tests {
		if exp[i][0] != got[i][0] || exp[i][1] != got[i][1] {
			tc := tests[i]
			fmt.Fprintf(os.Stderr, "wrong answer on case %d: expected (%d %d), got (%d %d)\n", i+1, tc.ansX, tc.ansY, got[i][0], got[i][1])
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2135B-ref-*")
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

func buildTests() []testCase {
	var tcs []testCase
	add := func(x, y int64, anchors [][2]int64) {
		tcs = append(tcs, testCase{anchors: anchors, ansX: x, ansY: y})
	}

	add(0, 0, [][2]int64{{0, 0}})
	add(5, -3, [][2]int64{{1, 1}, {-1, -1}})
	add(-100, 100, [][2]int64{{-100, 99}, {-101, 100}, {-99, 100}})
	add(1_000_000_000, -1_000_000_000, [][2]int64{{0, 0}, {1_000_000_000, 0}, {0, -1_000_000_000}})

	rng := rand.New(rand.NewSource(2135))
	for len(tcs) < 35 {
		n := rng.Intn(100) + 1
		coords := make(map[[2]int64]struct{})
		anchors := make([][2]int64, 0, n)
		for len(anchors) < n {
			x := rng.Int63n(2_000_000_001) - 1_000_000_000
			y := rng.Int63n(2_000_000_001) - 1_000_000_000
			p := [2]int64{x, y}
			if _, ok := coords[p]; ok {
				continue
			}
			coords[p] = struct{}{}
			anchors = append(anchors, p)
		}
		ansX := rng.Int63n(2_000_000_001) - 1_000_000_000
		ansY := rng.Int63n(2_000_000_001) - 1_000_000_000
		add(ansX, ansY, anchors)
	}

	return tcs
}

func buildInput(tcs []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tcs)))
	sb.WriteByte('\n')
	for _, tc := range tcs {
		sb.WriteString(strconv.Itoa(len(tc.anchors)))
		sb.WriteByte('\n')
		for _, p := range tc.anchors {
			sb.WriteString(fmt.Sprintf("%d %d\n", p[0], p[1]))
		}
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.ansX, tc.ansY))
	}
	return sb.String()
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
	return strings.TrimSpace(stdout.String()), nil
}

func parseOutput(out string, t int) ([][2]int64, error) {
	tokens := strings.Fields(out)
	if len(tokens) != 2*t {
		return nil, fmt.Errorf("expected %d tokens, got %d", 2*t, len(tokens))
	}
	res := make([][2]int64, t)
	for i := 0; i < t; i++ {
		x, err1 := strconv.ParseInt(tokens[2*i], 10, 64)
		y, err2 := strconv.ParseInt(tokens[2*i+1], 10, 64)
		if err1 != nil || err2 != nil {
			return nil, fmt.Errorf("invalid integer around position %d", 2*i+1)
		}
		res[i] = [2]int64{x, y}
	}
	return res, nil
}
