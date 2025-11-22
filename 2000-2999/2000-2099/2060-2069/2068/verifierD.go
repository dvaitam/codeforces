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
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	candidatePath := os.Args[1]

	inputData, err := io.ReadAll(os.Stdin)
	if err != nil {
		fail("failed to read input: %v", err)
	}
	n, weights, err := parseInput(inputData)
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
	refCodes, err := parseCodes(refOut, n)
	if err != nil {
		fail("failed to parse reference output: %v\noutput:\n%s", err, refOut)
	}
	if err := validateCodes(refCodes); err != nil {
		fail("reference produced invalid codes: %v", err)
	}
	optCost := computeCost(weights, refCodes)

	candOut, candErr, err := runProgram(candBin, inputData)
	if err != nil {
		fail("candidate runtime error: %v\n%s", err, candErr)
	}
	candCodes, err := parseCodes(candOut, n)
	if err != nil {
		fail("invalid candidate output: %v\noutput:\n%s", err, candOut)
	}
	if err := validateCodes(candCodes); err != nil {
		fail("invalid codes: %v", err)
	}
	candCost := computeCost(weights, candCodes)

	if candCost != optCost {
		fail("wrong expected transmission time: expected %d got %d (scaled by 10000)", optCost, candCost)
	}

	fmt.Println("OK")
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func parseInput(data []byte) (int, []int64, error) {
	reader := bufio.NewReader(bytes.NewReader(data))
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return 0, nil, err
	}
	weights := make([]int64, n)
	for i := 0; i < n; i++ {
		var s string
		if _, err := fmt.Fscan(reader, &s); err != nil {
			return 0, nil, err
		}
		weights[i] = int64(parseFreq(s))
	}
	return n, weights, nil
}

func parseFreq(s string) int {
	if strings.IndexByte(s, '.') == -1 {
		v, _ := strconv.Atoi(s)
		return v * 10000
	}
	parts := strings.SplitN(s, ".", 2)
	whole, _ := strconv.Atoi(parts[0])
	frac := parts[1]
	for len(frac) < 4 {
		frac += "0"
	}
	if len(frac) > 4 {
		frac = frac[:4]
	}
	f, _ := strconv.Atoi(frac)
	return whole*10000 + f
}

func buildReference() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier directory")
	}
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "2068D.go")

	tmp, err := os.CreateTemp("", "2068D-ref-*")
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
		tmp, err := os.CreateTemp("", "2068D-cand-*")
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

func parseCodes(out string, n int) ([]string, error) {
	tokens := strings.Fields(out)
	if len(tokens) != n {
		return nil, fmt.Errorf("expected %d codes, got %d", n, len(tokens))
	}
	codes := make([]string, n)
	for i, c := range tokens {
		codes[i] = c
	}
	return codes, nil
}

type trieNode struct {
	child [2]*trieNode
	end   bool
}

func validateCodes(codes []string) error {
	root := &trieNode{}
	for idx, code := range codes {
		if len(code) == 0 {
			return fmt.Errorf("code %d is empty", idx+1)
		}
		cur := root
		for i := 0; i < len(code); i++ {
			var id int
			if code[i] == '.' {
				id = 0
			} else if code[i] == '-' {
				id = 1
			} else {
				return fmt.Errorf("code %d contains invalid character %q", idx+1, code[i])
			}
			if cur.end {
				return fmt.Errorf("code %d extends another code (prefix violation)", idx+1)
			}
			if cur.child[id] == nil {
				cur.child[id] = &trieNode{}
			}
			cur = cur.child[id]
		}
		if cur.end {
			return fmt.Errorf("duplicate code at position %d", idx+1)
		}
		if cur.child[0] != nil || cur.child[1] != nil {
			return fmt.Errorf("code %d is a prefix of another code", idx+1)
		}
		cur.end = true
	}
	return nil
}

func computeCost(weights []int64, codes []string) int64 {
	var cost int64
	for i, code := range codes {
		var t int64
		for j := 0; j < len(code); j++ {
			if code[j] == '.' {
				t += 1
			} else {
				t += 2
			}
		}
		cost += weights[i] * t
	}
	return cost
}
