package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const refSource = "2000-2999/2000-2099/2010-2019/2014/2014D.go"

type job struct {
	l int
	r int
}

type testData struct {
	n    int
	d    int
	jobs []job
}

type expect struct {
	nWindows int
	maxStart int
	minStart int
}

type testCase struct {
	name  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}

	refBin, refCleanup, err := buildBinary(refSource)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer refCleanup()

	candBin, candCleanup, err := buildBinary(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to prepare candidate binary: %v\n", err)
		os.Exit(1)
	}
	defer candCleanup()

	tests := generateTests()
	for idx, tc := range tests {
		data, err := parseInput(tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse internal test %d (%s): %v\n", idx+1, tc.name, err)
			os.Exit(1)
		}
		exps := computeExpectations(data)

		refOut, err := runBinary(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		if err := verifyOutput(refOut, exps); err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runBinary(candBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		if err := verifyOutput(candOut, exps); err != nil {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s): %v\ninput:\n%s\ncandidate output:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildBinary(path string) (string, func(), error) {
	cleanPath := filepath.Clean(path)
	if strings.HasSuffix(cleanPath, ".go") {
		tmp, err := os.CreateTemp("", "verifier2014D-*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		cmd := exec.Command("go", "build", "-o", tmp.Name(), cleanPath)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("%v\n%s", err, out.String())
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	abs, err := filepath.Abs(cleanPath)
	if err != nil {
		return "", nil, err
	}
	return abs, func() {}, nil
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return stdout.String() + stderr.String(), fmt.Errorf("%v", err)
	}
	return stdout.String(), nil
}

func parseInput(input string) ([]testData, error) {
	reader := strings.NewReader(input)
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return nil, err
	}
	data := make([]testData, 0, t)
	for i := 0; i < t; i++ {
		var n, d, k int
		if _, err := fmt.Fscan(reader, &n, &d, &k); err != nil {
			return nil, err
		}
		jobs := make([]job, k)
		for j := 0; j < k; j++ {
			if _, err := fmt.Fscan(reader, &jobs[j].l, &jobs[j].r); err != nil {
				return nil, err
			}
		}
		data = append(data, testData{n: n, d: d, jobs: jobs})
	}
	return data, nil
}

func computeExpectations(cases []testData) []expect {
	res := make([]expect, len(cases))
	for i, tc := range cases {
		N := tc.n - tc.d + 1
		if N < 1 {
			N = 1
		}
		diff := make([]int, N+3)
		for _, jb := range tc.jobs {
			start := jb.l - tc.d + 1
			if start < 1 {
				start = 1
			}
			end := jb.r
			if end > N {
				end = N
			}
			if start <= end {
				diff[start]++
				diff[end+1]--
			}
		}
		cur := 0
		maxVal := -1
		minVal := int(1e9)
		maxPos := 1
		minPos := 1
		for pos := 1; pos <= N; pos++ {
			cur += diff[pos]
			if cur > maxVal {
				maxVal = cur
				maxPos = pos
			}
			if cur < minVal {
				minVal = cur
				minPos = pos
			}
		}
		res[i] = expect{nWindows: N, maxStart: maxPos, minStart: minPos}
	}
	return res
}

func verifyOutput(output string, exps []expect) error {
	reader := strings.NewReader(output)
	for idx, exp := range exps {
		var bro, mom int
		if _, err := fmt.Fscan(reader, &bro, &mom); err != nil {
			return fmt.Errorf("case %d: failed to read two integers: %v", idx+1, err)
		}
		if bro < 1 || bro > exp.nWindows {
			return fmt.Errorf("case %d: brother start %d out of range [1,%d]", idx+1, bro, exp.nWindows)
		}
		if mom < 1 || mom > exp.nWindows {
			return fmt.Errorf("case %d: mother start %d out of range [1,%d]", idx+1, mom, exp.nWindows)
		}
		if bro != exp.maxStart {
			return fmt.Errorf("case %d: expected brother at %d but got %d", idx+1, exp.maxStart, bro)
		}
		if mom != exp.minStart {
			return fmt.Errorf("case %d: expected mother at %d but got %d", idx+1, exp.minStart, mom)
		}
	}
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return fmt.Errorf("unexpected extra output token %q", extra)
	}
	return nil
}

func generateTests() []testCase {
	var tests []testCase
	tests = append(tests, testCase{
		name:  "sample",
		input: sampleInput(),
	})
	tests = append(tests, buildTest("simple_cases", []testData{
		{n: 5, d: 2, jobs: []job{{1, 2}, {2, 4}, {3, 5}}},
		{n: 6, d: 3, jobs: []job{{1, 3}, {2, 6}}},
		{n: 4, d: 1, jobs: []job{{2, 2}}},
	}))

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 6; i++ {
		tests = append(tests, buildTest(
			fmt.Sprintf("random_small_%d", i+1),
			randomCases(rng, 4, 8, 20),
		))
	}
	for i := 0; i < 4; i++ {
		tests = append(tests, buildTest(
			fmt.Sprintf("random_mid_%d", i+1),
			randomCases(rng, 5, 20, 200),
		))
	}
	tests = append(tests, buildTest("random_large", randomCases(rng, 3, 40, 100000)))

	return tests
}

func sampleInput() string {
	return `6
2 1 1
1 2
4 1 2
1 2
2 4
7 2 3
1 2
1 3
6 7
5 1 2
1 2
3 5
9 2 1
2 8
9 2 4
7 9
4 8
1 3
2 3
`
}

func buildTest(name string, cases []testData) testCase {
	var b strings.Builder
	fmt.Fprintln(&b, len(cases))
	for _, tc := range cases {
		fmt.Fprintf(&b, "%d %d %d\n", tc.n, tc.d, len(tc.jobs))
		for _, jb := range tc.jobs {
			fmt.Fprintf(&b, "%d %d\n", jb.l, jb.r)
		}
	}
	return testCase{name: name, input: b.String()}
}

func randomCases(rng *rand.Rand, tMin, tMax, nMax int) []testData {
	if tMin > tMax {
		tMin = tMax
	}
	t := rng.Intn(tMax-tMin+1) + tMin
	res := make([]testData, 0, t)
	for i := 0; i < t; i++ {
		if nMax < 1 {
			nMax = 1
		}
		n := rng.Intn(nMax) + 1
		d := rng.Intn(n) + 1
		k := rng.Intn(n) + 1
		jobs := make([]job, k)
		for j := 0; j < k; j++ {
			l := rng.Intn(n) + 1
			r := rng.Intn(n-l+1) + l
			jobs[j] = job{l: l, r: r}
		}
		res = append(res, testData{n: n, d: d, jobs: jobs})
	}
	return res
}
