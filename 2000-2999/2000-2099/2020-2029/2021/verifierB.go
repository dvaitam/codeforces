package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// solveReference is the correct solver for 2021B, embedded directly.
func solveReference(input []byte) string {
	reader := bufio.NewReader(bytes.NewReader(input))
	var out bytes.Buffer
	writer := bufio.NewWriter(&out)

	var t int
	fmt.Fscan(reader, &t)

	for i := 0; i < t; i++ {
		var n, x int
		fmt.Fscan(reader, &n, &x)

		M := x
		if n+1 < M {
			M = n + 1
		}

		lists := make([][]int, M)

		for j := 0; j < n; j++ {
			var a int
			fmt.Fscan(reader, &a)
			r := a % x
			if r < M {
				lists[r] = append(lists[r], a)
			}
		}

		for j := 0; j < M; j++ {
			sort.Slice(lists[j], func(p, q int) bool {
				return lists[j][p] > lists[j][q]
			})
		}

		ans := 0
		for v := 0; v <= n; v++ {
			r := v % x
			if len(lists[r]) == 0 {
				ans = v
				break
			}
			last := lists[r][len(lists[r])-1]
			lists[r] = lists[r][:len(lists[r])-1]
			if last > v {
				ans = v
				break
			}
		}
		fmt.Fprintln(writer, ans)
	}
	writer.Flush()
	return out.String()
}

func main() {
	if len(os.Args) != 2 {
		fail("usage: go run verifierB.go /path/to/candidate")
	}
	candidate := os.Args[1]

	tests := generateTests()

	for i, input := range tests {
		refOut := solveReference(input)
		refAnswers, err := parseReference(refOut)
		if err != nil {
			fail("failed to parse reference output on test %d: %v", i+1, err)
		}

		testCases, err := parseInput(input)
		if err != nil {
			fail("failed to parse input on test %d: %v", i+1, err)
		}

		candOut, err := runProgram(commandFor(candidate), input)
		if err != nil {
			fail("candidate execution failed on test %d: %v", i+1, err)
		}
		if err := checkCandidate(candOut, testCases, refAnswers); err != nil {
			fail("test %d: %v", i+1, err)
		}
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

func parseReference(out string) ([]int64, error) {
	reader := bufio.NewReader(strings.NewReader(out))
	var ans []int64
	for {
		token, err := readToken(reader)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		val, err := strconv.ParseInt(token, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("reference output invalid integer %q", token)
		}
		ans = append(ans, val)
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

func generateTests() [][]byte {
	rng := rand.New(rand.NewSource(2021))
	var tests [][]byte

	// Sample test
	tests = append(tests, []byte("4\n5 3\n0 1 2 3 4\n5 3\n0 1 1 1 1\n5 3\n4 3 0 1 2\n5 3\n3 0 1 2 5\n"))

	// Small edge cases
	tests = append(tests, []byte("1\n1 1\n0\n"))
	tests = append(tests, []byte("1\n1 2\n1\n"))
	tests = append(tests, []byte("1\n3 2\n0 1 2\n"))

	// Random tests
	for iter := 0; iter < 30; iter++ {
		t := rng.Intn(5) + 1
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", t)
		for c := 0; c < t; c++ {
			n := rng.Intn(50) + 1
			x := rng.Intn(50) + 1
			fmt.Fprintf(&sb, "%d %d\n", n, x)
			for j := 0; j < n; j++ {
				if j > 0 {
					sb.WriteByte(' ')
				}
				fmt.Fprintf(&sb, "%d", rng.Intn(n+1))
			}
			sb.WriteByte('\n')
		}
		tests = append(tests, []byte(sb.String()))
	}

	// Larger random test
	{
		var sb strings.Builder
		t := 3
		fmt.Fprintf(&sb, "%d\n", t)
		for c := 0; c < t; c++ {
			n := 500 + rng.Intn(500)
			x := rng.Intn(500) + 1
			fmt.Fprintf(&sb, "%d %d\n", n, x)
			for j := 0; j < n; j++ {
				if j > 0 {
					sb.WriteByte(' ')
				}
				fmt.Fprintf(&sb, "%d", rng.Intn(n+1))
			}
			sb.WriteByte('\n')
		}
		tests = append(tests, []byte(sb.String()))
	}

	return tests
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
