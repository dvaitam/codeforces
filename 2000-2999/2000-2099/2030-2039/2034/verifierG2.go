package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

type interval struct {
	l int
	r int
}

type testCase struct {
	n    int
	intv []interval
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG2.go /path/to/binary")
		os.Exit(1)
	}
	candidatePath := os.Args[1]

	inputData, err := io.ReadAll(os.Stdin)
	if err != nil {
		fail("failed to read input: %v", err)
	}
	tests, err := parseInput(inputData)
	if err != nil {
		fail("failed to parse input: %v", err)
	}

	refBin, cleanupRef, err := buildReference()
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer cleanupRef()

	candBin, cleanupCand, err := prepareCandidate(candidatePath)
	if err != nil {
		fail("failed to prepare candidate: %v", err)
	}
	defer cleanupCand()

	refOut, refErr, err := runProgram(refBin, inputData)
	if err != nil {
		fail("reference runtime error: %v\n%s", err, refErr)
	}
	refSolutions, err := parseOutput(refOut, tests, false)
	if err != nil {
		fail("failed to parse reference output: %v\noutput:\n%s", err, refOut)
	}

	candOut, candErr, err := runProgram(candBin, inputData)
	if err != nil {
		fail("candidate runtime error: %v\n%s", err, candErr)
	}
	candSolutions, err := parseOutput(candOut, tests, true)
	if err != nil {
		fail("invalid candidate output: %v\noutput:\n%s", err, candOut)
	}

	for i, tc := range tests {
		refK := refSolutions[i].k
		sol := candSolutions[i]
		if sol.k != refK {
			fail("test case %d: wrong minimal colors, expected %d got %d", i+1, refK, sol.k)
		}
		if len(sol.colors) != tc.n {
			fail("test case %d: expected %d colors, got %d", i+1, tc.n, len(sol.colors))
		}
		if err := validateColoring(tc, sol.k, sol.colors); err != nil {
			fail("test case %d: invalid coloring: %v", i+1, err)
		}
	}

	fmt.Println("OK")
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func parseInput(data []byte) ([]testCase, error) {
	reader := bufio.NewReader(bytes.NewReader(data))
	var T int
	if _, err := fmt.Fscan(reader, &T); err != nil {
		return nil, err
	}
	tests := make([]testCase, T)
	for i := 0; i < T; i++ {
		var n int
		if _, err := fmt.Fscan(reader, &n); err != nil {
			return nil, err
		}
		tc := testCase{n: n, intv: make([]interval, n)}
		for j := 0; j < n; j++ {
			if _, err := fmt.Fscan(reader, &tc.intv[j].l, &tc.intv[j].r); err != nil {
				return nil, err
			}
		}
		tests[i] = tc
	}
	return tests, nil
}

func buildReference() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier directory")
	}
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "2034G2.go")

	tmp, err := os.CreateTemp("", "2034G2-ref-*")
	if err != nil {
		return "", nil, err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), src)
	cmd.Dir = dir
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", nil, fmt.Errorf("%v\n%s", err, out.String())
	}
	cleanup := func() {
		os.Remove(tmp.Name())
	}
	return tmp.Name(), cleanup, nil
}

func prepareCandidate(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		abs, err := filepath.Abs(path)
		if err != nil {
			return "", nil, err
		}
		tmp, err := os.CreateTemp("", "2034G2-cand-*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		cmd := exec.Command("go", "build", "-o", tmp.Name(), abs)
		cmd.Dir = filepath.Dir(abs)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("%v\n%s", err, out.String())
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func runProgram(path string, input []byte) (string, string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	return out.String(), errBuf.String(), err
}

type solution struct {
	k      int
	colors []int
}

func parseOutput(out string, tests []testCase, needColors bool) ([]solution, error) {
	tokens := strings.Fields(out)
	pos := 0
	res := make([]solution, len(tests))
	for i, tc := range tests {
		if pos >= len(tokens) {
			return nil, fmt.Errorf("test case %d: missing k", i+1)
		}
		k, err := strconv.Atoi(tokens[pos])
		if err != nil {
			return nil, fmt.Errorf("test case %d: invalid k %q", i+1, tokens[pos])
		}
		pos++
		if k <= 0 {
			return nil, fmt.Errorf("test case %d: k must be positive", i+1)
		}
		if pos+tc.n > len(tokens) {
			return nil, fmt.Errorf("test case %d: not enough color values", i+1)
		}
		colors := make([]int, tc.n)
		for j := 0; j < tc.n; j++ {
			v, err := strconv.Atoi(tokens[pos+j])
			if err != nil {
				return nil, fmt.Errorf("test case %d: invalid color %q", i+1, tokens[pos+j])
			}
			colors[j] = v
		}
		pos += tc.n
		if !needColors {
			colors = nil
		}
		res[i] = solution{k: k, colors: colors}
	}
	if pos != len(tokens) {
		return nil, fmt.Errorf("extra tokens at the end of output")
	}
	return res, nil
}

type event struct {
	pos   int
	color int
	delta int
}

func validateColoring(tc testCase, k int, colors []int) error {
	for i, c := range colors {
		if c < 1 || c > k {
			return fmt.Errorf("color %d at position %d out of range [1,%d]", c, i+1, k)
		}
	}
	events := make([]event, 0, 2*tc.n)
	for i, iv := range tc.intv {
		c := colors[i]
		events = append(events, event{pos: iv.l, color: c, delta: 1})
		events = append(events, event{pos: iv.r + 1, color: c, delta: -1})
	}
	sort.Slice(events, func(i, j int) bool {
		return events[i].pos < events[j].pos
	})
	if len(events) == 0 {
		return fmt.Errorf("no intervals provided")
	}
	counts := make([]int, k+1)
	ones := 0
	total := 0
	prevPos := events[0].pos
	idx := 0
	for idx < len(events) {
		curPos := events[idx].pos
		if curPos > prevPos && total > 0 && ones == 0 {
			return fmt.Errorf("no uniquely colored interval in segment [%d,%d]", prevPos, curPos-1)
		}
		for idx < len(events) && events[idx].pos == curPos {
			ev := events[idx]
			old := counts[ev.color]
			if old == 1 {
				ones--
			}
			counts[ev.color] = old + ev.delta
			if counts[ev.color] < 0 {
				return fmt.Errorf("negative coverage for color %d at position %d", ev.color, curPos)
			}
			if counts[ev.color] == 1 {
				ones++
			}
			total += ev.delta
			idx++
		}
		prevPos = curPos
	}
	if total > 0 && ones == 0 {
		return fmt.Errorf("no uniquely colored interval after last event")
	}
	return nil
}
