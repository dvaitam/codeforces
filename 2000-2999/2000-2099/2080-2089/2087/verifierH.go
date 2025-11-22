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

type query struct {
	typ   int
	s     int
	piles []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	targetPath := os.Args[1]

	refBin, refCleanup, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer refCleanup()

	candBin, candCleanup, err := prepareCandidate(targetPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to prepare candidate: %v\n", err)
		os.Exit(1)
	}
	defer candCleanup()

	seed := time.Now().UnixNano()
	tests := generateTests(seed)

	for i, input := range tests {
		expRaw, err := runBinary(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		expAns, err := parseOutputs(expRaw)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d: %v\noutput:\n%s\n", i+1, err, expRaw)
			os.Exit(1)
		}

		actRaw, err := runBinary(candBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		actAns, err := parseOutputs(actRaw)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d: %v\ninput:\n%soutput:\n%s\n", i+1, err, input, actRaw)
			os.Exit(1)
		}

		if len(actAns) != len(expAns) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: expected %d answers, got %d\ninput:\n%sreference:\n%v\ncandidate:\n%v\n",
				i+1, len(expAns), len(actAns), input, expAns, actAns)
			os.Exit(1)
		}
		for j := range expAns {
			if actAns[j] != expAns[j] {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d at position %d: expected %q, got %q\ninput:\n%sreference:\n%v\ncandidate:\n%v\n",
					i+1, j+1, expAns[j], actAns[j], input, expAns, actAns)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed (seed %d).\n", len(tests), seed)
}

func buildReference() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("unable to determine verifier location")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "verifier-2087H-ref-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "ref2087H")
	cmd := exec.Command("go", "build", "-o", outPath, "2087H.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("go build failed: %v\n%s", err, out)
	}
	cleanup := func() {
		_ = os.RemoveAll(tmpDir)
	}
	return outPath, cleanup, nil
}

func prepareCandidate(path string) (string, func(), error) {
	if !strings.HasSuffix(path, ".go") {
		return path, func() {}, nil
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", nil, err
	}
	dir := filepath.Dir(absPath)
	tmpDir, err := os.MkdirTemp("", "verifier-2087H-cand-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "candidate2087H")
	cmd := exec.Command("go", "build", "-o", outPath, absPath)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build candidate: %v\n%s", err, out)
	}
	cleanup := func() {
		_ = os.RemoveAll(tmpDir)
	}
	return outPath, cleanup, nil
}

func runBinary(binPath, input string) (string, error) {
	cmd := exec.Command(binPath)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func parseOutputs(out string) ([]string, error) {
	fields := strings.Fields(out)
	res := make([]string, len(fields))
	for i, f := range fields {
		if f != "First" && f != "Second" {
			return nil, fmt.Errorf("invalid token %q at position %d", f, i+1)
		}
		res[i] = f
	}
	return res, nil
}

func buildInput(qs []query) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(qs)))
	sb.WriteByte('\n')
	for _, q := range qs {
		if q.typ == 1 {
			sb.WriteString("1 ")
			sb.WriteString(strconv.Itoa(q.s))
		} else {
			sb.WriteString("2 ")
			sb.WriteString(strconv.Itoa(len(q.piles)))
			for _, v := range q.piles {
				sb.WriteByte(' ')
				sb.WriteString(strconv.Itoa(v))
			}
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func deterministicTests() []string {
	var tests []string

	tests = append(tests, buildInput([]query{
		{typ: 2, piles: []int{1}},
	}))

	tests = append(tests, buildInput([]query{
		{typ: 1, s: 2},
		{typ: 2, piles: []int{2}},
		{typ: 1, s: 2},
		{typ: 2, piles: []int{2}},
		{typ: 2, piles: []int{2, 2, 2}},
	}))

	tests = append(tests, buildInput([]query{
		{typ: 1, s: 5},
		{typ: 1, s: 10},
		{typ: 2, piles: []int{5, 10}},
		{typ: 1, s: 5},
		{typ: 2, piles: []int{6, 7, 8}},
		{typ: 1, s: 10},
		{typ: 2, piles: []int{11}},
	}))

	tests = append(tests, buildInput([]query{
		{typ: 1, s: 300000},
		{typ: 1, s: 1},
		{typ: 2, piles: []int{1, 2, 300000}},
		{typ: 2, piles: []int{300000, 300000, 300000}},
	}))

	return tests
}

func randomQueries(rng *rand.Rand, q int) []query {
	res := make([]query, q)
	type2Seen := false
	for i := 0; i < q; i++ {
		if i == q-1 && !type2Seen {
			res[i] = randomType2(rng)
			type2Seen = true
			continue
		}
		if rng.Intn(3) == 0 {
			res[i] = randomType2(rng)
			type2Seen = true
		} else {
			res[i] = query{typ: 1, s: rng.Intn(300000) + 1}
		}
	}
	if !type2Seen {
		res[0] = randomType2(rng)
	}
	return res
}

func randomType2(rng *rand.Rand) query {
	k := rng.Intn(3) + 1
	piles := make([]int, k)
	for i := 0; i < k; i++ {
		piles[i] = rng.Intn(300000) + 1
	}
	return query{typ: 2, piles: piles}
}

func generateTests(seed int64) []string {
	rng := rand.New(rand.NewSource(seed))
	tests := deterministicTests()

	for i := 0; i < 40; i++ {
		q := rng.Intn(200) + 1
		tests = append(tests, buildInput(randomQueries(rng, q)))
	}

	for i := 0; i < 3; i++ {
		q := rng.Intn(1500) + 500
		tests = append(tests, buildInput(randomQueries(rng, q)))
	}

	return tests
}
