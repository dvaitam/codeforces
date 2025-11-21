package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const refSource = "2000-2999/2000-2099/2030-2039/2034/2034G1.go"

type interval struct {
	l int64
	r int64
}

type testInput struct {
	n         int
	intervals []interval
}

func main() {
	if len(os.Args) != 2 {
		fail("usage: go run verifierG1.go /path/to/candidate")
	}
	candidate := os.Args[1]

	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fail("failed to read input: %v", err)
	}

	tests, err := parseInput(input)
	if err != nil {
		fail("failed to parse input: %v", err)
	}

	refBin, err := buildReference()
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	refOut, err := runProgram(exec.Command(refBin), input)
	if err != nil {
		fail("reference execution failed: %v", err)
	}
	refK, _, err := parseSolution(refOut, tests, false)
	if err != nil {
		fail("failed to parse reference output: %v", err)
	}

	candOut, err := runProgram(commandFor(candidate), input)
	if err != nil {
		fail("candidate execution failed: %v", err)
	}
	candK, candColors, err := parseSolution(candOut, tests, true)
	if err != nil {
		fail("failed to parse candidate output: %v", err)
	}

	for i, test := range tests {
		if candK[i] != refK[i] {
			fail("test %d: expected minimal colors %d, got %d", i+1, refK[i], candK[i])
		}
		if err := verifyColoring(test, candColors[i], candK[i]); err != nil {
			fail("test %d: %v", i+1, err)
		}
	}

	fmt.Println("OK")
}

func parseInput(data []byte) ([]testInput, error) {
	reader := bufio.NewReader(bytes.NewReader(data))
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return nil, err
	}
	tests := make([]testInput, t)
	for i := 0; i < t; i++ {
		var n int
		if _, err := fmt.Fscan(reader, &n); err != nil {
			return nil, err
		}
		segs := make([]interval, n)
		for j := 0; j < n; j++ {
			var l, r int64
			if _, err := fmt.Fscan(reader, &l, &r); err != nil {
				return nil, err
			}
			segs[j] = interval{l: l, r: r}
		}
		tests[i] = testInput{n: n, intervals: segs}
	}
	return tests, nil
}

func parseSolution(out string, tests []testInput, keepColors bool) ([]int, [][]int, error) {
	reader := bufio.NewReader(strings.NewReader(out))
	ks := make([]int, len(tests))
	var colors [][]int
	if keepColors {
		colors = make([][]int, len(tests))
	}
	for i, test := range tests {
		token, err := readToken(reader)
		if err != nil {
			return nil, nil, fmt.Errorf("missing k for test %d: %v", i+1, err)
		}
		k, err := strconv.Atoi(token)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid k for test %d: %v", i+1, err)
		}
		ks[i] = k
		if keepColors {
			colors[i] = make([]int, test.n)
		}
		for j := 0; j < test.n; j++ {
			token, err = readToken(reader)
			if err != nil {
				return nil, nil, fmt.Errorf("missing color %d for test %d: %v", j+1, i+1, err)
			}
			val, err := strconv.Atoi(token)
			if err != nil {
				return nil, nil, fmt.Errorf("invalid color %d for test %d: %v", j+1, i+1, err)
			}
			if keepColors {
				colors[i][j] = val
			}
		}
	}
	if extra, err := readToken(reader); err != io.EOF {
		if err == nil {
			return nil, nil, fmt.Errorf("extra token %q after outputs", extra)
		}
		return nil, nil, err
	}
	return ks, colors, nil
}

func verifyColoring(test testInput, colors []int, k int) error {
	if k <= 0 {
		return fmt.Errorf("k must be positive, got %d", k)
	}
	if len(colors) != test.n {
		return fmt.Errorf("expected %d colors, got %d", test.n, len(colors))
	}
	for idx, c := range colors {
		if c < 1 || c > k {
			return fmt.Errorf("color at index %d is out of range 1..%d", idx+1, k)
		}
	}

	coords := make([]int64, 0, 2*test.n)
	for _, seg := range test.intervals {
		coords = append(coords, seg.l, seg.r)
	}
	if len(coords) == 0 {
		return nil
	}
	sort.Slice(coords, func(i, j int) bool { return coords[i] < coords[j] })
	pos := 0
	for i := 1; i < len(coords); i++ {
		if coords[i] != coords[pos] {
			pos++
			coords[pos] = coords[i]
		}
	}
	coords = coords[:pos+1]

	index := make(map[int64]int, len(coords))
	for i, v := range coords {
		index[v] = i
	}
	starts := make([][]int, len(coords))
	ends := make([][]int, len(coords))
	for idx, seg := range test.intervals {
		starts[index[seg.l]] = append(starts[index[seg.l]], idx)
		ends[index[seg.r]] = append(ends[index[seg.r]], idx)
	}

	counts := make([]int, k+1)
	activeTotal := 0
	uniqueColors := 0
	inc := func(color int) {
		if counts[color] == 1 {
			uniqueColors--
		}
		counts[color]++
		if counts[color] == 1 {
			uniqueColors++
		}
	}
	dec := func(color int) {
		if counts[color] == 1 {
			uniqueColors--
		}
		counts[color]--
		if counts[color] == 1 {
			uniqueColors++
		}
	}
	check := func(desc string) error {
		if activeTotal > 0 && uniqueColors == 0 {
			return fmt.Errorf("no uniquely colored interval %s", desc)
		}
		return nil
	}

	for idx := 0; idx < len(coords); idx++ {
		if idx > 0 {
			if err := check(fmt.Sprintf("between (%d,%d)", coords[idx-1], coords[idx])); err != nil {
				return err
			}
		}
		for _, segIdx := range starts[idx] {
			inc(colors[segIdx])
			activeTotal++
		}
		if err := check(fmt.Sprintf("at point %d", coords[idx])); err != nil {
			return err
		}
		for _, segIdx := range ends[idx] {
			dec(colors[segIdx])
			activeTotal--
			if activeTotal < 0 {
				return fmt.Errorf("internal error: negative active interval count")
			}
		}
	}
	if activeTotal != 0 {
		return fmt.Errorf("internal error: unmatched intervals")
	}
	return nil
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2034G1-ref-*")
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

func commandFor(path string) *exec.Cmd {
	switch filepath.Ext(path) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func runProgram(cmd *exec.Cmd, input []byte) (string, error) {
	cmd.Stdin = bytes.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		msg := stderr.String()
		if msg == "" {
			msg = stdout.String()
		}
		return "", fmt.Errorf("%v\n%s", err, msg)
	}
	return stdout.String(), nil
}

func readToken(r *bufio.Reader) (string, error) {
	var sb strings.Builder
	for {
		b, err := r.ReadByte()
		if err != nil {
			if err == io.EOF {
				return "", io.EOF
			}
			return "", err
		}
		if !isSpace(b) {
			sb.WriteByte(b)
			break
		}
	}
	for {
		b, err := r.ReadByte()
		if err != nil {
			if err == io.EOF {
				return sb.String(), nil
			}
			return "", err
		}
		if isSpace(b) {
			break
		}
		sb.WriteByte(b)
	}
	return sb.String(), nil
}

func isSpace(b byte) bool {
	return b == ' ' || b == '\n' || b == '\r' || b == '\t' || b == '\v' || b == '\f'
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
