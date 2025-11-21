package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type event struct {
	t int
	p byte
}

type testCase struct {
	n    int
	c    int
	d    int
	evs  []event
	tEnd int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refPath, err := buildReference()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer os.Remove(refPath)

	tests := buildTests()

	for i, tc := range tests {
		input := buildInput(tc)

		refOut, err := runProgram(refPath, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\n%s", i+1, err, refOut)
			os.Exit(1)
		}
		expected, err := parseOutput(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output parse error on test %d: %v\n%s", i+1, err, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\n%s", i+1, err, candOut)
			os.Exit(1)
		}
		got, err := parseOutput(candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output parse error on test %d: %v\n%s", i+1, err, candOut)
			os.Exit(1)
		}

		if expected != got {
			fmt.Fprintf(os.Stderr, "Mismatch on test %d: expected %d, got %d\nInput:\n%s\n", i+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "1120F-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	cmd := exec.Command("go", "build", "-o", tmp.Name(), "1120F.go")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runProgram(target, input string) (string, error) {
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
		return stdout.String() + stderr.String(), err
	}
	return stdout.String(), nil
}

func buildTests() []testCase {
	var tests []testCase
	sumN := 0
	add := func(tc testCase) {
		if sumN+tc.n > 100000 {
			return
		}
		tests = append(tests, tc)
		sumN += tc.n
	}

	add(testCase{
		n: 5, c: 1, d: 4,
		evs: []event{
			{0, 'P'}, {1, 'W'}, {3, 'P'}, {5, 'P'}, {8, 'P'},
		},
		tEnd: 10,
	})
	add(testCase{
		n: 10, c: 10, d: 94,
		evs: []event{
			{17, 'W'}, {20, 'W'}, {28, 'W'}, {48, 'W'}, {51, 'P'},
			{52, 'W'}, {56, 'W'}, {62, 'P'}, {75, 'P'}, {78, 'P'},
		},
		tEnd: 87,
	})

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 50 && sumN < 100000 {
		n := rng.Intn(2000) + 1
		c := rng.Intn(100) + 1
		d := rng.Intn(100_000_000) + 1
		evs := make([]event, n)
		cur := 0
		for i := 0; i < n; i++ {
			cur += rng.Intn(10) + 1
			if cur > 1_000_000 {
				cur = 1_000_000 - (n - i)
			}
			evs[i] = event{t: cur, p: []byte{'W', 'P'}[rng.Intn(2)]}
		}
		tEnd := cur + rng.Intn(10) + 1
		if tEnd > 1_000_000 {
			tEnd = 1_000_000
		}
		if tEnd <= evs[len(evs)-1].t {
			tEnd = evs[len(evs)-1].t + 1
			if tEnd > 1_000_000 {
				tEnd = 1_000_000
			}
		}
		add(testCase{n: n, c: c, d: d, evs: evs, tEnd: tEnd})
	}

	return tests
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, tc.c, tc.d))
	for _, ev := range tc.evs {
		sb.WriteString(fmt.Sprintf("%d %c\n", ev.t, ev.p))
	}
	sb.WriteString(strconv.Itoa(tc.tEnd))
	sb.WriteByte('\n')
	return sb.String()
}

func parseOutput(out string) (int64, error) {
	fields := strings.Fields(out)
	if len(fields) != 1 {
		return 0, fmt.Errorf("expected 1 output, got %d", len(fields))
	}
	val, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q", fields[0])
	}
	return val, nil
}
