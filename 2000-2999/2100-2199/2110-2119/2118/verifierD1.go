package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

const refSource = "2000-2999/2100-2199/2110-2119/2118/2118D1.go"

type testCase struct {
	n, k   int
	pos    []int64
	del    []int
	q      int
	starts []int64
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2118D1-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	return tmp.Name(), nil
}

func runProgram(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		abs, err := filepath.Abs(bin)
		if err != nil {
			return "", err
		}
		cmd = exec.Command("go", "run", abs)
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
	return stdout.String(), nil
}

func deterministicTests() []testCase {
	return []testCase{
		{
			n: 2, k: 2,
			pos:    []int64{1, 4},
			del:    []int{1, 0},
			q:      3,
			starts: []int64{1, 2, 3},
		},
		{
			n: 1, k: 3,
			pos:    []int64{5},
			del:    []int{2},
			q:      2,
			starts: []int64{1, 6},
		},
		{
			n: 3, k: 2,
			pos:    []int64{2, 5, 9},
			del:    []int{1, 0, 1},
			q:      4,
			starts: []int64{2, 3, 4, 10},
		},
	}
}

func randomTest(rng *rand.Rand, n, k, q int) testCase {
	pos := make([]int64, n)
	current := int64(rng.Intn(10) + 1)
	for i := 0; i < n; i++ {
		current += int64(rng.Intn(10) + 1)
		pos[i] = current
	}
	del := make([]int, n)
	for i := 0; i < n; i++ {
		del[i] = rng.Intn(k)
	}
	starts := make([]int64, q)
	for i := 0; i < q; i++ {
		if rng.Intn(4) == 0 {
			starts[i] = pos[rng.Intn(n)]
		} else {
			offset := int64(rng.Intn(20)) - 10
			base := pos[rng.Intn(n)]
			if base+offset < 1 {
				offset = 1 - base
			}
			starts[i] = base + offset
		}
	}
	// make sure positions sorted
	sort.Slice(pos, func(i, j int) bool { return pos[i] < pos[j] })
	return testCase{n: n, k: k, pos: pos, del: del, q: q, starts: starts}
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.k))
		for i, v := range tc.pos {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
		for i, v := range tc.del {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		sb.WriteString(strconv.Itoa(tc.q))
		sb.WriteByte('\n')
		for i, v := range tc.starts {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseOutput(out string, tests []testCase) ([]string, error) {
	lines := strings.Fields(out)
	res := make([]string, 0, len(lines))
	res = append(res, lines...)
	expect := 0
	for _, tc := range tests {
		expect += tc.q
	}
	if len(res) != expect {
		return nil, fmt.Errorf("expected %d answers, got %d", expect, len(res))
	}
	return res, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD1.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := deterministicTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	total := 0
	for _, tc := range tests {
		total += tc.n + tc.q
	}

	for len(tests) < 35 && total < 480 {
		n := rng.Intn(15) + 1
		k := rng.Intn(30) + 1
		q := rng.Intn(15) + 1
		if total+n+q > 500 {
			break
		}
		tests = append(tests, randomTest(rng, n, k, q))
		total += n + q
	}
	for len(tests) < 45 && total < 500 {
		n := rng.Intn(60) + 20
		k := rng.Intn(500) + 1
		q := rng.Intn(20) + 1
		if total+n+q > 500 {
			break
		}
		tests = append(tests, randomTest(rng, n, k, q))
		total += n + q
	}

	input := buildInput(tests)

	wantOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\n", err)
		os.Exit(1)
	}
	want, err := parseOutput(wantOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\n", err)
		os.Exit(1)
	}

	gotOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n", err)
		os.Exit(1)
	}
	got, err := parseOutput(gotOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse candidate output: %v\n", err)
		os.Exit(1)
	}

	if len(want) != len(got) {
		fmt.Fprintf(os.Stderr, "answer count mismatch: expected %d got %d\n", len(want), len(got))
		os.Exit(1)
	}

	for i := range want {
		if strings.ToUpper(want[i]) != strings.ToUpper(got[i]) {
			fmt.Fprintf(os.Stderr, "mismatch at answer %d: expected %s got %s\n", i+1, want[i], got[i])
			os.Exit(1)
		}
	}

	fmt.Printf("All %d test cases passed\n", len(tests))
}
