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

type testCase struct {
	n   int
	seq []int
}

func solve(n int, seq []int) string {
	arrived := make([]bool, n+1)
	expected := n
	var sb strings.Builder
	for _, x := range seq {
		arrived[x] = true
		first := true
		for expected > 0 && arrived[expected] {
			if !first {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(expected))
			expected--
			first = false
		}
		sb.WriteByte('\n')
	}
	return strings.TrimSpace(sb.String())
}

func (tc testCase) input() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", tc.n))
	for i, v := range tc.seq {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(10) + 1
	perm := rng.Perm(n)
	for i := range perm {
		perm[i]++
	}
	return testCase{n: n, seq: perm}
}

func runProgram(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func runCase(bin string, tc testCase) error {
	in := tc.input()
	expected := solve(tc.n, append([]int(nil), tc.seq...))
	got, err := runProgram(bin, in)
	if err != nil {
		return err
	}
	exp := strings.ReplaceAll(expected, "\r", "")
	got = strings.ReplaceAll(got, "\r", "")
	if err := compareOutput(exp, got); err != nil {
		return fmt.Errorf("comparison failed: %v\nexpected:\n%s\n\ngot:\n%s", err, exp, got)
	}
	return nil
}

func compareOutput(expected, got string) error {
	expLines := strings.Split(expected, "\n")
	gotLines := strings.Split(got, "\n")
	if len(expLines) != len(gotLines) {
		return fmt.Errorf("line count mismatch: expected %d, got %d", len(expLines), len(gotLines))
	}
	for i := range expLines {
		if strings.TrimSpace(expLines[i]) != strings.TrimSpace(gotLines[i]) {
			return fmt.Errorf("line %d mismatch", i+1)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCase{{n: 3, seq: []int{3, 1, 2}}}
	for len(cases) < 100 {
		cases = append(cases, randomCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
