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

const refSource = "./2161A.go"

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/candidate")
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
		want, err := parseOutput(wantOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d: %v\ninput:\n%soutput:\n%s", i+1, err, input, wantOut)
			os.Exit(1)
		}

		gotOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		got, err := parseOutput(gotOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d: %v\ninput:\n%soutput:\n%s", i+1, err, input, gotOut)
			os.Exit(1)
		}

		if len(want) != len(got) {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %d answers got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
				i+1, len(want), len(got), input, wantOut, gotOut)
			os.Exit(1)
		}
		for idx := range want {
			if want[idx] != got[idx] {
				fmt.Fprintf(os.Stderr, "test %d case %d mismatch: expected %d got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
					i+1, idx+1, want[idx], got[idx], input, wantOut, gotOut)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2161A-ref-*")
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

func parseOutput(out string) ([]int, error) {
	sc := bufio.NewScanner(strings.NewReader(out))
	sc.Split(bufio.ScanWords)
	var res []int
	for sc.Scan() {
		val, err := strconv.Atoi(sc.Text())
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

type testCase struct {
	r0, x, d, n int
	rounds      string
}

func buildTest(cases []testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, cs := range cases {
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", cs.r0, cs.x, cs.d, cs.n))
		sb.WriteString(cs.rounds)
		sb.WriteByte('\n')
	}
	return sb.String()
}

var rng = rand.New(rand.NewSource(20240605))

func randomRounds(n int) string {
	var sb strings.Builder
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			sb.WriteByte('1')
		} else {
			sb.WriteByte('2')
		}
	}
	return sb.String()
}

func randomCase() testCase {
	r0 := rng.Intn(1_000_000_000)
	x := rng.Intn(1_000_000_000-1) + 1
	d := rng.Intn(999) + 1
	n := rng.Intn(1000) + 1
	rounds := randomRounds(n)
	return testCase{r0: r0, x: x, d: d, n: n, rounds: rounds}
}

func generateTests() []string {
	var tests []string
	// fixed edge cases
	tests = append(tests, buildTest([]testCase{
		{r0: 0, x: 1, d: 1, n: 1, rounds: "1"},
		{r0: 0, x: 1, d: 1, n: 1, rounds: "2"},
	}))

	tests = append(tests, buildTest([]testCase{
		{r0: 100, x: 100, d: 10, n: 5, rounds: "11111"},
		{r0: 100, x: 1000, d: 10, n: 5, rounds: "22222"},
	}))

	tests = append(tests, buildTest([]testCase{
		{r0: 200, x: 150, d: 20, n: 10, rounds: "1212121212"},
	}))

	tests = append(tests, buildTest([]testCase{
		{r0: 500, x: 600, d: 100, n: 15, rounds: "211211211211211"},
		{r0: 1000, x: 1200, d: 50, n: 12, rounds: "121122112211"},
	}))

	tests = append(tests, buildTest([]testCase{
		randomCase(),
		randomCase(),
	}))

	for i := 0; i < 5; i++ {
		tests = append(tests, buildTest([]testCase{
			randomCase(),
			randomCase(),
			randomCase(),
		}))
	}

	return tests
}
