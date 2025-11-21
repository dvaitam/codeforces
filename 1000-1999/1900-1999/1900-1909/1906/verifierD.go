package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const ref1906D = "1000-1999/1900-1999/1900-1909/1906/1906D.go"

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to read input:", err)
		os.Exit(1)
	}
	q, err := queryCount(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	refBin, cleanup, err := buildReference(ref1906D)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer cleanup()

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\n", err)
		os.Exit(1)
	}
	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n", err)
		os.Exit(1)
	}

	refVals, err := parseFloatLines(refOut)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to parse reference output:", err)
		os.Exit(1)
	}
	candVals, err := parseFloatLines(candOut)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to parse candidate output:", err)
		os.Exit(1)
	}
	if len(refVals) != q {
		fmt.Fprintf(os.Stderr, "reference output line count mismatch: expected %d, got %d\n", q, len(refVals))
		os.Exit(1)
	}
	if len(candVals) != q {
		fmt.Fprintf(os.Stderr, "candidate output line count mismatch: expected %d, got %d\n", q, len(candVals))
		os.Exit(1)
	}
	for i := 0; i < q; i++ {
		if !closeEnough(refVals[i], candVals[i]) {
			fmt.Fprintf(os.Stderr, "line %d mismatch: expected %.9f, got %.9f\n", i+1, refVals[i], candVals[i])
			fmt.Fprintln(os.Stderr, "reference output:")
			fmt.Fprintln(os.Stderr, refOut)
			fmt.Fprintln(os.Stderr, "candidate output:")
			fmt.Fprintln(os.Stderr, candOut)
			os.Exit(1)
		}
	}

	fmt.Println("Accepted")
}

func queryCount(data []byte) (int, error) {
	reader := bufio.NewReader(bytes.NewReader(data))
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return 0, fmt.Errorf("failed to read n: %v", err)
	}
	for i := 0; i < n; i++ {
		var x, y int
		if _, err := fmt.Fscan(reader, &x, &y); err != nil {
			return 0, fmt.Errorf("failed to read polygon point %d: %v", i+1, err)
		}
	}
	var q int
	if _, err := fmt.Fscan(reader, &q); err != nil {
		return 0, fmt.Errorf("failed to read q: %v", err)
	}
	for i := 0; i < q; i++ {
		var a, b, c, d int
		if _, err := fmt.Fscan(reader, &a, &b, &c, &d); err != nil {
			return 0, fmt.Errorf("failed to read query %d: %v", i+1, err)
		}
	}
	return q, nil
}

func parseFloatLines(out string) ([]float64, error) {
	scanner := bufio.NewScanner(strings.NewReader(out))
	vals := make([]float64, 0)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		v, err := strconv.ParseFloat(line, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid float %q", line)
		}
		vals = append(vals, v)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return vals, nil
}

func closeEnough(a, b float64) bool {
	diff := math.Abs(a - b)
	limit := 1e-6 * math.Max(1.0, math.Abs(a))
	return diff <= limit+1e-9
}

func buildReference(src string) (string, func(), error) {
	dir, err := os.MkdirTemp("", "verifier-1906D-")
	if err != nil {
		return "", nil, err
	}
	bin := filepath.Join(dir, "ref.bin")
	cmd := exec.Command("go", "build", "-o", bin, src)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("go build failed: %v\n%s", err, stderr.String())
	}
	return bin, func() { os.RemoveAll(dir) }, nil
}

func runProgram(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	return out.String(), cmd.Run()
}
