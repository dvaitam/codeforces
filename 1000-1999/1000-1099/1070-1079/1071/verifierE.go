package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type drop struct {
	t int
	x int
	y int
}

type testCase struct {
	n      int
	w, h   int
	e1, e2 int
	drops  []drop
}

func buildReferenceBinary() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("unable to locate verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "ref-1071E-")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "ref_1071E")
	cmd := exec.Command("go", "build", "-o", binPath, "1071E.go")
	cmd.Dir = dir
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, out.String())
	}
	cleanup := func() { os.RemoveAll(tmpDir) }
	return binPath, cleanup, nil
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
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return out.String(), nil
}

func parseAnswer(output string) (float64, error) {
	fields := strings.Fields(output)
	if len(fields) == 0 {
		return 0, fmt.Errorf("empty output")
	}
	val, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse float from %q: %v", fields[0], err)
	}
	return val, nil
}

func answersMatch(expected, got float64) bool {
	if expected < 0 && got < 0 {
		return true
	}
	if expected < 0 || got < 0 {
		return false
	}
	diff := math.Abs(expected - got)
	den := math.Max(1.0, math.Abs(expected))
	return diff/den <= 1e-4+1e-7
}

func inputString(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, tc.w, tc.h))
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.e1, tc.e2))
	for _, d := range tc.drops {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", d.t, d.x, d.y))
	}
	return sb.String()
}

func deterministicTests() []testCase {
	return []testCase{
		{
			n: 1, w: 5, h: 4,
			e1: 2, e2: 3,
			drops: []drop{
				{t: 1, x: 2, y: 2},
			},
		},
		{
			n: 2, w: 6, h: 5,
			e1: 1, e2: 5,
			drops: []drop{
				{t: 2, x: 3, y: 1},
				{t: 4, x: 4, y: 3},
			},
		},
		{
			n: 3, w: 10, h: 6,
			e1: 5, e2: 4,
			drops: []drop{
				{t: 1, x: 5, y: 1},
				{t: 2, x: 2, y: 3},
				{t: 5, x: 8, y: 2},
			},
		},
	}
}

func randomTest(rng *rand.Rand, maxN int) testCase {
	n := rng.Intn(maxN-1) + 1
	w := rng.Intn(30) + 2
	h := rng.Intn(30) + 2
	e1 := rng.Intn(w + 1)
	e2 := rng.Intn(w + 1)
	drops := make([]drop, n)
	tCur := 0
	for i := 0; i < n; i++ {
		tCur += rng.Intn(5) + 1
		x := rng.Intn(w + 1)
		y := rng.Intn(h-1) + 1 // ensure 0 < y < h
		drops[i] = drop{t: tCur, x: x, y: y}
	}
	return testCase{
		n:     n,
		w:     w,
		h:     h,
		e1:    e1,
		e2:    e2,
		drops: drops,
	}
}

func genTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := deterministicTests()
	for i := 0; i < 60; i++ {
		tests = append(tests, randomTest(rng, 5))
	}
	for i := 0; i < 60; i++ {
		tests = append(tests, randomTest(rng, 20))
	}
	for i := 0; i < 40; i++ {
		tests = append(tests, randomTest(rng, 80))
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[len(os.Args)-1]
	if target == "--" {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}

	refBin, cleanup, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := genTests()
	for idx, tc := range tests {
		input := inputString(tc)
		refOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		expVal, err := parseAnswer(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference parse error on test %d: %v\ninput:\n%soutput:\n%s", idx+1, err, input, refOut)
			os.Exit(1)
		}

		out, err := runProgram(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d runtime error: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		gotVal, err := parseAnswer(out)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d parse error: %v\ninput:\n%soutput:\n%s", idx+1, err, input, out)
			os.Exit(1)
		}
		if !answersMatch(expVal, gotVal) {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %.10f got %.10f\ninput:\n%sreference output:\n%suser output:\n%s", idx+1, expVal, gotVal, input, refOut, out)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
