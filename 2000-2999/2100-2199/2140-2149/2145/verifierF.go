package main

import (
	"bytes"
	"container/list"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	n int
	m int64
	a []int64
	b []int64
}

func runBinary(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 1, m: 3, a: []int64{2}, b: []int64{0}},
		{n: 2, m: 5, a: []int64{2, 3}, b: []int64{0, 1}},
		{n: 3, m: 7, a: []int64{2, 3, 4}, b: []int64{1, 2, 3}},
		{n: 2, m: 4, a: []int64{3, 4}, b: []int64{1, 2}},
	}
}

func randomTest(rng *rand.Rand) testCase {
	n := rng.Intn(3) + 1
	m := int64(rng.Intn(200) + 1)
	a := make([]int64, n)
	b := make([]int64, n)
	for i := 0; i < n; i++ {
		a[i] = int64(rng.Intn(8) + 2)
		b[i] = int64(rng.Intn(int(a[i])))
	}
	return testCase{n: n, m: m, a: a, b: b}
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
		for i := 0; i < tc.n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(tc.a[i], 10))
		}
		sb.WriteByte('\n')
		for i := 0; i < tc.n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(tc.b[i], 10))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseOutput(out string, t int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != t {
		return nil, fmt.Errorf("expected %d outputs, got %d", t, len(fields))
	}
	res := make([]int64, t)
	for i, f := range fields {
		val, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", f)
		}
		res[i] = val
	}
	return res, nil
}

func bruteForce(tc testCase) int64 {
	type state struct {
		pos   int64
		phase int
	}
	n := tc.n
	maxPos := tc.m
	if maxPos < 0 {
		maxPos = 0
	}
	visited := make([][]bool, n)
	for i := 0; i < n; i++ {
		visited[i] = make([]bool, maxPos+1)
	}
	queue := list.New()
	queue.PushBack(struct {
		state
		turn int64
	}{state{0, 0}, 0})
	visited[0][0] = true

	for queue.Len() > 0 {
		front := queue.Front().Value.(struct {
			state
			turn int64
		})
		queue.Remove(queue.Front())
		pos := front.pos
		phase := front.phase
		turn := front.turn
		if pos == tc.m {
			return turn
		}
		for _, delta := range []int64{0, 1} {
			newPos := pos + delta
			if newPos > tc.m {
				continue
			}
			if newPos%tc.a[phase] == tc.b[phase] {
				continue
			}
			newPhase := (phase + 1) % n
			if !visited[newPhase][newPos] {
				visited[newPhase][newPos] = true
				queue.PushBack(struct {
					state
					turn int64
				}{state{newPos, newPhase}, turn + 1})
			}
		}
	}
	return -1
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	tests := deterministicTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomTest(rng))
	}

	input := buildInput(tests)
	out, err := runBinary(target, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	results, err := parseOutput(out, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	for i, tc := range tests {
		exp := bruteForce(tc)
		if results[i] != exp {
			fmt.Fprintf(os.Stderr, "test %d mismatch: expected %d got %d\nn=%d m=%d a=%v b=%v\n", i+1, exp, results[i], tc.n, tc.m, tc.a, tc.b)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
