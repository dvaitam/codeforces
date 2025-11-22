package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const refSource = "2096H.go"

type testCase struct {
	name string
	n    int
	m    int
	seg  [][2]int
}

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Fprintln(os.Stderr, "usage: go run verifierH.go /path/to/candidate")
		os.Exit(1)
	}
	candPath := os.Args[len(os.Args)-1]

	refBin, cleanupRef, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference solution:", err)
		os.Exit(1)
	}
	defer cleanupRef()

	candBin, cleanupCand, err := prepareCandidate(candPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to prepare candidate:", err)
		os.Exit(1)
	}
	defer cleanupCand()

	tests := buildTests()
	input := buildInput(tests)

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\n", err)
		os.Exit(1)
	}
	refAns, err := parseOutputs(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference produced invalid output: %v\n%s", err, refOut)
		os.Exit(1)
	}

	candOut, err := runProgram(candBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\ninput preview:\n%s\n", err, previewInput(input))
		os.Exit(1)
	}
	candAns, err := parseOutputs(candOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate produced invalid output: %v\noutput:\n%s\ninput preview:\n%s\n", err, candOut, previewInput(input))
		os.Exit(1)
	}

	for i := range tests {
		if refAns[i] != candAns[i] {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s): expected %d got %d\n", i+1, tests[i].name, refAns[i], candAns[i])
			fmt.Fprintln(os.Stderr, previewInput(buildInput([]testCase{tests[i]})))
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier directory")
	}
	dir := filepath.Dir(file)

	tmpDir, err := os.MkdirTemp("", "ref2096H-")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "ref")

	cmd := exec.Command("go", "build", "-o", binPath, refSource)
	cmd.Dir = dir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("reference build failed: %v\n%s", err, stderr.String())
	}

	cleanup := func() { _ = os.RemoveAll(tmpDir) }
	return binPath, cleanup, nil
}

func prepareCandidate(path string) (string, func(), error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", nil, err
	}
	if filepath.Ext(abs) != ".go" {
		return abs, func() {}, nil
	}

	tmpDir, err := os.MkdirTemp("", "cand2096H-")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "candidate")

	cmd := exec.Command("go", "build", "-o", binPath, abs)
	cmd.Dir = filepath.Dir(abs)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("candidate build failed: %v\n%s", err, stderr.String())
	}

	cleanup := func() { _ = os.RemoveAll(tmpDir) }
	return binPath, cleanup, nil
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseOutputs(out string, expected int) ([]uint64, error) {
	tokens := strings.Fields(out)
	if len(tokens) != expected {
		return nil, fmt.Errorf("expected %d integers, got %d", expected, len(tokens))
	}
	ans := make([]uint64, expected)
	for i, tok := range tokens {
		v, err := strconv.ParseUint(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q: %v", tok, err)
		}
		ans[i] = v
	}
	return ans, nil
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
		for _, p := range tc.seg {
			sb.WriteString(fmt.Sprintf("%d %d\n", p[0], p[1]))
		}
	}
	return sb.String()
}

func buildTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	tests := []testCase{}
	tests = append(tests, sampleTests()...)
	tests = append(tests, singlePoint(), fullRange(), alternatingSingles(), repeatedHeavy())
	tests = append(tests, exhaustiveSmall()...)
	tests = append(tests, randomTests(rng, 6, 2, 4, 1, 4, "rand_tiny")...)
	tests = append(tests, randomTests(rng, 10, 5, 8, 2, 10, "rand_small")...)
	tests = append(tests, randomTests(rng, 8, 9, 14, 10, 60, "rand_mid")...)
	tests = append(tests, randomTests(rng, 4, 15, 18, 80, 160, "rand_big")...)
	tests = append(tests, structuredLarge()...)

	return tests
}

func sampleTests() []testCase {
	return []testCase{
		{name: "sample_like_1", n: 2, m: 2, seg: [][2]int{{0, 2}, {1, 3}}},
		{name: "sample_like_2", n: 3, m: 3, seg: [][2]int{{0, 5}, {2, 6}, {1, 7}}},
		{name: "sample_like_3", n: 5, m: 5, seg: [][2]int{{3, 12}, {0, 31}, {7, 14}, {9, 22}, {16, 20}}},
		{name: "sample_like_4", n: 1, m: 1, seg: [][2]int{{0, 1}}},
	}
}

func singlePoint() testCase {
	return testCase{
		name: "single_point",
		n:    3,
		m:    3,
		seg: [][2]int{
			{5, 5},
			{0, 0},
			{7, 7},
		},
	}
}

func fullRange() testCase {
	return testCase{
		name: "full_range",
		n:    4,
		m:    5,
		seg: [][2]int{
			{0, 31},
			{0, 31},
			{0, 31},
			{0, 31},
		},
	}
}

func alternatingSingles() testCase {
	return testCase{
		name: "alternating_points",
		n:    8,
		m:    4,
		seg: func() [][2]int {
			res := make([][2]int, 0, 8)
			vals := []int{0, 1, 3, 2, 5, 7, 8, 10}
			for _, v := range vals {
				res = append(res, [2]int{v, v})
			}
			return res
		}(),
	}
}

func repeatedHeavy() testCase {
	m := 10
	size := 1 << m
	return testCase{
		name: "heavy_repeats",
		n:    30000,
		m:    m,
		seg: func() [][2]int {
			ranges := [][2]int{{0, size - 1}, {100, 200}, {777, 999}}
			res := make([][2]int, 0, 30000)
			for i := 0; i < 30000; i++ {
				res = append(res, ranges[i%len(ranges)])
			}
			return res
		}(),
	}
}

func exhaustiveSmall() []testCase {
	var res []testCase
	// n up to 3, m up to 3
	for m := 1; m <= 3; m++ {
		limit := 1 << m
		for n := 1; n <= 3; n++ {
			// iterate over small set of interval endpoints
			points := make([]int, 0, limit)
			for i := 0; i < limit; i++ {
				points = append(points, i)
			}
			total := limit * limit
			for mask := 0; mask < total; mask++ {
				if len(res) >= 20 {
					break
				}
				l := mask / limit
				r := mask % limit
				if l > r {
					continue
				}
				seg := make([][2]int, n)
				for i := 0; i < n; i++ {
					seg[i] = [2]int{points[l], points[r]}
				}
				res = append(res, testCase{
					name: fmt.Sprintf("exhaust_m%d_n%d_l%d_r%d", m, n, l, r),
					n:    n,
					m:    m,
					seg:  seg,
				})
			}
		}
	}
	return res
}

func randomTests(rng *rand.Rand, count, minM, maxM, minN, maxN int, tag string) []testCase {
	res := make([]testCase, 0, count)
	for i := 0; i < count; i++ {
		m := rng.Intn(maxM-minM+1) + minM
		n := rng.Intn(maxN-minN+1) + minN
		maxVal := 1 << m
		seg := make([][2]int, n)
		for j := 0; j < n; j++ {
			l := rng.Intn(maxVal)
			r := rng.Intn(maxVal)
			if l > r {
				l, r = r, l
			}
			seg[j] = [2]int{l, r}
		}
		res = append(res, testCase{
			name: fmt.Sprintf("%s_%d", tag, i+1),
			n:    n,
			m:    m,
			seg:  seg,
		})
	}
	return res
}

func structuredLarge() []testCase {
	m := 18
	size := 1 << m
	return []testCase{
		{
			name: "wide_cover",
			n:    120,
			m:    m,
			seg: func() [][2]int {
				res := make([][2]int, 0, 120)
				step := size / 60
				for i := 0; i < 60; i++ {
					l := i * step
					r := l + step - 1
					if r >= size {
						r = size - 1
					}
					res = append(res, [2]int{l, r})
					res = append(res, [2]int{0, r})
				}
				return res
			}(),
		},
		{
			name: "max_n_repeats",
			n:    140000,
			m:    m,
			seg: func() [][2]int {
				res := make([][2]int, 0, 140000)
				pool := [][2]int{{0, size - 1}, {size / 4, size/4 + 1000}, {size / 2, size/2 + 500}, {size - 2000, size - 1}}
				for len(res) < 140000 {
					res = append(res, pool[len(res)%len(pool)])
				}
				return res
			}(),
		},
	}
}

func previewInput(input string) string {
	lines := strings.Split(input, "\n")
	if len(lines) > 10 {
		return strings.Join(lines[:10], "\n") + "\n..."
	}
	return input
}
