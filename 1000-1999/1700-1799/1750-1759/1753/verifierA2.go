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
	"strings"
)

type testCase struct {
	n   int
	arr []int
}

type segment struct {
	l, r int
}

type partition struct {
	possible bool
	segs     []segment
}

func runProgram(path string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return out.String(), nil
}

func parseInput(data []byte) ([]testCase, error) {
	reader := bufio.NewReader(bytes.NewReader(data))
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return nil, fmt.Errorf("failed to read t: %v", err)
	}
	cases := make([]testCase, 0, t)
	for ; t > 0; t-- {
		var n int
		if _, err := fmt.Fscan(reader, &n); err != nil {
			return nil, fmt.Errorf("failed to read n: %v", err)
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			if _, err := fmt.Fscan(reader, &arr[i]); err != nil {
				return nil, fmt.Errorf("failed to read array element: %v", err)
			}
		}
		cases = append(cases, testCase{n: n, arr: arr})
	}
	return cases, nil
}

func parseOutput(out string, t int) ([]partition, error) {
	reader := bufio.NewReader(strings.NewReader(out))
	parts := make([]partition, 0, t)
	for i := 0; i < t; i++ {
		var k int
		if _, err := fmt.Fscan(reader, &k); err != nil {
			return nil, fmt.Errorf("failed to read k for case %d: %v", i+1, err)
		}
		if k == -1 {
			parts = append(parts, partition{possible: false})
			continue
		}
		if k < 0 {
			return nil, fmt.Errorf("invalid k=%d for case %d", k, i+1)
		}
		segs := make([]segment, k)
		for j := 0; j < k; j++ {
			if _, err := fmt.Fscan(reader, &segs[j].l, &segs[j].r); err != nil {
				return nil, fmt.Errorf("failed to read segment %d for case %d: %v", j+1, i+1, err)
			}
		}
		parts = append(parts, partition{possible: true, segs: segs})
	}
	return parts, nil
}

func validatePartition(tc testCase, part partition) error {
	if len(part.segs) == 0 {
		return fmt.Errorf("no segments provided")
	}
	prevEnd := 0
	for idx, seg := range part.segs {
		if seg.l > seg.r {
			return fmt.Errorf("segment %d has l>r (%d>%d)", idx+1, seg.l, seg.r)
		}
		if seg.l != prevEnd+1 {
			return fmt.Errorf("segment %d does not continue sequence (expected %d got %d)", idx+1, prevEnd+1, seg.l)
		}
		if seg.r > tc.n {
			return fmt.Errorf("segment %d exceeds array length", idx+1)
		}
		prevEnd = seg.r
	}
	if prevEnd != tc.n {
		return fmt.Errorf("partition ends at %d instead of %d", prevEnd, tc.n)
	}
	total := 0
	for _, seg := range part.segs {
		sign := 1
		sum := 0
		for pos := seg.l; pos <= seg.r; pos++ {
			sum += tc.arr[pos-1] * sign
			sign *= -1
		}
		total += sum
	}
	if total != 0 {
		return fmt.Errorf("total alternating sum %d != 0", total)
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA2.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[len(os.Args)-1]
	if target == "--" {
		fmt.Println("usage: go run verifierA2.go /path/to/binary")
		os.Exit(1)
	}

	inputData, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read input: %v\n", err)
		os.Exit(1)
	}
	tests, err := parseInput(inputData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	_, src, _, _ := runtime.Caller(0)
	baseDir := filepath.Dir(src)
	refPath := filepath.Join(baseDir, "1753A2.go")

	refOut, err := runProgram(refPath, inputData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\n", err)
		os.Exit(1)
	}
	refParts, err := parseOutput(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference output parse error: %v\n", err)
		os.Exit(1)
	}
	for i, part := range refParts {
		if !part.possible {
			continue
		}
		if err := validatePartition(tests[i], part); err != nil {
			fmt.Fprintf(os.Stderr, "reference partition invalid for case %d: %v\n", i+1, err)
			os.Exit(1)
		}
	}

	targetOut, err := runProgram(target, inputData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "target runtime error: %v\n", err)
		os.Exit(1)
	}
	targetParts, err := parseOutput(targetOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "target output parse error: %v\n", err)
		os.Exit(1)
	}

	for i, part := range targetParts {
		if !part.possible {
			if refParts[i].possible {
				fmt.Fprintf(os.Stderr, "case %d: solution exists but target output -1\n", i+1)
				os.Exit(1)
			}
			continue
		}
		if err := validatePartition(tests[i], part); err != nil {
			fmt.Fprintf(os.Stderr, "case %d invalid: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
