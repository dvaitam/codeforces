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

const (
	refSourceC   = "./2005C.go"
	randomCases  = 120
	maxTotalSize = 1_000_000
)

type testCase struct {
	n    int
	m    int
	data []string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference(refSourceC)
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	tests := deterministicCases()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < randomCases; i++ {
		tests = append(tests, randomCase(rng))
	}

	input := buildInput(tests)

	expect, err := runProgram(refBin, input)
	if err != nil {
		fail("reference failed: %v", err)
	}
	got, err := runCandidate(candidate, input)
	if err != nil {
		fail("candidate failed: %v", err)
	}
	if normalize(expect) != normalize(got) {
		fail("output mismatch\nexpected:\n%s\ngot:\n%s", expect, got)
	}
	fmt.Printf("All %d cases passed\n", len(tests))
}

func buildReference(src string) (string, error) {
	tmp, err := os.CreateTemp("", "2005C-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(src))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func deterministicCases() []testCase {
	return []testCase{
		{n: 5, m: 2, data: []string{"nn", "aa", "rr", "ee", "kk"}},
		{n: 1, m: 5, data: []string{"narek"}},
		{n: 1, m: 4, data: []string{"nare"}},
		{n: 5, m: 7, data: []string{"nrrarek", "nrnekan", "uuuuuuu", "ppppppp", "nkarekz"}},
		{n: 3, m: 6, data: []string{"narxxx", "eeknar", "zzzzzz"}},
		{n: 4, m: 5, data: []string{"abcde", "narek", "knara", "xxxxx"}},
	}
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(20) + 1
	m := rng.Intn(20) + 1
	data := make([]string, n)
	for i := 0; i < n; i++ {
		data[i] = randomString(rng, m)
	}
	return testCase{n: n, m: m, data: data}
}

func randomString(rng *rand.Rand, m int) string {
	bytes := make([]byte, m)
	for i := 0; i < m; i++ {
		if rng.Intn(4) == 0 {
			switch rng.Intn(5) {
			case 0:
				bytes[i] = 'n'
			case 1:
				bytes[i] = 'a'
			case 2:
				bytes[i] = 'r'
			case 3:
				bytes[i] = 'e'
			case 4:
				bytes[i] = 'k'
			}
		} else {
			bytes[i] = byte('a' + rng.Intn(26))
		}
	}
	return string(bytes)
}

func buildInput(cases []testCase) string {
	var sb strings.Builder
	totalSize := 0
	validCases := make([]testCase, 0, len(cases))
	for _, tc := range cases {
		size := tc.n * tc.m
		if totalSize+size > maxTotalSize {
			break
		}
		totalSize += size
		validCases = append(validCases, tc)
	}
	sb.WriteString(fmt.Sprintf("%d\n", len(validCases)))
	for _, tc := range validCases {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
		for _, s := range tc.data {
			sb.WriteString(s)
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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

func runCandidate(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
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

func normalize(out string) string {
	return strings.TrimSpace(out)
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
