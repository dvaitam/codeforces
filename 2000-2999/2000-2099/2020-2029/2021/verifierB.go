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

const refSource = "2000-2999/2000-2099/2020-2029/2021/2021B.go"

func main() {
	if len(os.Args) != 2 {
		fail("usage: go run verifierB.go /path/to/candidate")
	}
	candidate := os.Args[1]

	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fail("failed to read input: %v", err)
	}

	testCases, err := parseInput(input)
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
	refAnswers, err := parseReference(refOut, len(testCases))
	if err != nil {
		fail("failed to parse reference output: %v", err)
	}

	candOut, err := runProgram(commandFor(candidate), input)
	if err != nil {
		fail("candidate execution failed: %v", err)
	}
	if err := checkCandidate(candOut, testCases, refAnswers); err != nil {
		fail("%v", err)
	}

	fmt.Println("OK")
}

type testCase struct {
	n int
	x int64
	a []int64
}

func parseInput(data []byte) ([]testCase, error) {
	reader := bufio.NewReader(bytes.NewReader(data))
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return nil, err
	}
	tc := make([]testCase, t)
	for i := 0; i < t; i++ {
		var n int
		var x int64
		if _, err := fmt.Fscan(reader, &n, &x); err != nil {
			return nil, err
		}
		tc[i].n = n
		tc[i].x = x
		tc[i].a = make([]int64, n)
		for j := 0; j < n; j++ {
			if _, err := fmt.Fscan(reader, &tc[i].a[j]); err != nil {
				return nil, err
			}
		}
	}
	return tc, nil
}

func parseReference(out string, t int) ([]int64, error) {
	reader := bufio.NewReader(strings.NewReader(out))
	ans := make([]int64, t)
	for i := 0; i < t; i++ {
		token, err := readToken(reader)
		if err != nil {
			return nil, fmt.Errorf("reference output ended early at test %d: %v", i+1, err)
		}
		val, err := strconv.ParseInt(token, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("reference output invalid integer %q", token)
		}
		ans[i] = val
	}
	if extra, err := readToken(reader); err != io.EOF {
		if err == nil {
			return nil, fmt.Errorf("reference output has extra token %q", extra)
		}
		return nil, err
	}
	return ans, nil
}

func checkCandidate(out string, tests []testCase, refAnswers []int64) error {
	reader := bufio.NewReader(strings.NewReader(out))
	for idx, tc := range tests {
		token, err := readToken(reader)
		if err != nil {
			return fmt.Errorf("candidate output ended early at test %d: %v", idx+1, err)
		}
		val, err := strconv.ParseInt(token, 10, 64)
		if err != nil {
			return fmt.Errorf("candidate output invalid integer %q at test %d", token, idx+1)
		}
		if val < refAnswers[idx] {
			return fmt.Errorf("test %d: reported mex %d smaller than optimal %d", idx+1, val, refAnswers[idx])
		}
		if err := validateMex(tc, val); err != nil {
			return fmt.Errorf("test %d: %v", idx+1, err)
		}
		if val > refAnswers[idx] {
			return fmt.Errorf("test %d: reported mex %d exceeds optimal %d", idx+1, val, refAnswers[idx])
		}
	}
	if extra, err := readToken(reader); err != io.EOF {
		if err == nil {
			return fmt.Errorf("candidate output has extra token %q", extra)
		}
		return err
	}
	return nil
}

func validateMex(tc testCase, mex int64) error {
	if mex < 0 {
		return fmt.Errorf("mex cannot be negative")
	}
	if mex > int64(tc.n) {
		return fmt.Errorf("mex %d exceeds array size %d", mex, tc.n)
	}
	buckets := make(map[int64][]int64)
	for _, v := range tc.a {
		r := v % tc.x
		buckets[r] = append(buckets[r], v)
	}
	for r := range buckets {
		sort.Slice(buckets[r], func(i, j int) bool { return buckets[r][i] < buckets[r][j] })
	}
	pointers := make(map[int64]int)
	for cur := int64(0); cur < mex; cur++ {
		r := cur % tc.x
		list := buckets[r]
		ptr := pointers[r]
		if ptr < len(list) && list[ptr] <= cur {
			pointers[r] = ptr + 1
			continue
		}
		return fmt.Errorf("value %d cannot be formed", cur)
	}
	return nil
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2021B-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, buf.String())
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
