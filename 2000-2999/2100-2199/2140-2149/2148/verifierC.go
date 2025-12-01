package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const refSource = "./2148C.go"

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for i, input := range tests {
		wantOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		gotOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		want, err := parseOutput(wantOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d: %v\ninput:\n%soutput:\n%s", i+1, err, input, wantOut)
			os.Exit(1)
		}
		got, err := parseOutput(gotOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d: %v\ninput:\n%soutput:\n%s", i+1, err, input, gotOut)
			os.Exit(1)
		}
		if len(want) != len(got) {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %d results got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
				i+1, len(want), len(got), input, wantOut, gotOut)
			os.Exit(1)
		}
		for j := range want {
			if want[j] != got[j] {
				fmt.Fprintf(os.Stderr, "test %d failed on case %d: expected %d got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
					i+1, j+1, want[j], got[j], input, wantOut, gotOut)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2148C-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	cmd := exec.Command("go", "build", "-o", tmp.Name(), refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("go build failed: %v\n%s", err, stderr.String())
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
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func parseOutput(out string) ([]int64, error) {
	sc := bufio.NewScanner(strings.NewReader(out))
	sc.Split(bufio.ScanWords)
	var res []int64
	for sc.Scan() {
		val, err := strconv.ParseInt(sc.Text(), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", sc.Text())
		}
		res = append(res, val)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func generateTests() []string {
	var tests []string
	tests = append(tests, buildTest([]testCase{{
		n: 1,
		m: 5,
		segments: []segment{
			{time: 2, side: 1},
		},
	}}))

	tests = append(tests, buildTest([]testCase{{
		n: 3,
		m: 20,
		segments: []segment{
			{time: 4, side: 0},
			{time: 9, side: 1},
			{time: 15, side: 0},
		},
	}}))

	tests = append(tests, buildTest([]testCase{
		randomCase(5, 50),
		randomCase(4, 80),
	}))

	for i := 0; i < 5; i++ {
		tests = append(tests, buildTest([]testCase{
			randomCase(10+i*5, int64(100+50*i)),
		}))
	}

	return tests
}

type segment struct {
	time int64
	side int64
}

type testCase struct {
	n        int
	m        int64
	segments []segment
}

func buildTest(cases []testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, tc := range cases {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
		for _, seg := range tc.segments {
			sb.WriteString(fmt.Sprintf("%d %d\n", seg.time, seg.side))
		}
	}
	return sb.String()
}

var rng = rand.New(rand.NewSource(20240605))

func randomCase(maxN int, maxM int64) testCase {
	if maxN < 1 {
		maxN = 1
	}
	n := rng.Intn(maxN) + 1
	if maxM < 1 {
		maxM = 1
	}
	m := rng.Int63n(maxM) + int64(n)
	segments := make([]segment, n)
	var current int64
	for i := 0; i < n; i++ {
		step := rng.Int63n((m/int64(n))+2) + 1
		current += step
		if current >= m {
			current = m - int64(n-i)
		}
		if current < 0 {
			current = int64(i + 1)
		}
		segments[i] = segment{
			time: current,
			side: int64(rng.Intn(2)),
		}
	}
	if segments[len(segments)-1].time >= m {
		segments[len(segments)-1].time = m - 1
		if segments[len(segments)-1].time < 0 {
			segments[len(segments)-1].time = 0
		}
	}
	for i := 1; i < len(segments); i++ {
		if segments[i].time <= segments[i-1].time {
			segments[i].time = segments[i-1].time + 1
		}
		if segments[i].time >= m {
			segments[i].time = m - 1
		}
	}
	return testCase{
		n:        n,
		m:        m,
		segments: segments,
	}
}
