package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	refSourceB = "2000-2999/2000-2099/2030-2039/2039/2039B.go"
	sumLimit   = 300000
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[len(os.Args)-1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	input := buildInput(tests)

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference failed: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}
	candOut, err := runCandidate(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate failed: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}

	refAns, err := parseOutputs(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse reference output: %v\n", err)
		os.Exit(1)
	}
	candAns, err := parseOutputs(candOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse candidate output: %v\n", err)
		os.Exit(1)
	}

	for i, s := range tests {
		exp := refAns[i]
		ans := candAns[i]
		if exp == "-1" {
			if ans != "-1" {
				fmt.Fprintf(os.Stderr, "test %d: expected -1 but got %q\n", i+1, ans)
				os.Exit(1)
			}
			continue
		}
		if ans == "-1" {
			fmt.Fprintf(os.Stderr, "test %d: candidate printed -1 but solution exists\n", i+1)
			os.Exit(1)
		}
		if !strings.Contains(s, ans) {
			fmt.Fprintf(os.Stderr, "test %d: output %q is not a substring of %q\n", i+1, ans, s)
			os.Exit(1)
		}
		if !hasEvenDistinct(ans) {
			fmt.Fprintf(os.Stderr, "test %d: output %q has odd number of distinct substrings\n", i+1, ans)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2039B-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSourceB))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func commandFor(path string) *exec.Cmd {
	switch filepath.Ext(path) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseOutputs(output string, t int) ([]string, error) {
	tokens := strings.Fields(output)
	if len(tokens) < t {
		return nil, fmt.Errorf("expected %d outputs, got %d", t, len(tokens))
	}
	if len(tokens) > t {
		return nil, fmt.Errorf("extra output detected starting at token %q", tokens[t])
	}
	return tokens, nil
}

type samState struct {
	next [26]int
	link int
	len  int
}

func hasEvenDistinct(s string) bool {
	if len(s) == 0 {
		return false
	}
	sam := newSAM(len(s))
	for i := 0; i < len(s); i++ {
		c := int(s[i] - 'a')
		if c < 0 || c >= 26 {
			return false
		}
		sam.extend(c)
	}
	var total int64
	for i := 1; i < len(sam.states); i++ {
		total += int64(sam.states[i].len - sam.states[sam.states[i].link].len)
	}
	return total%2 == 0
}

type suffixAutomaton struct {
	states []samState
	last   int
}

func newSAM(maxLen int) *suffixAutomaton {
	st := make([]samState, 1, 2*maxLen)
	for i := range st[0].next {
		st[0].next[i] = -1
	}
	st[0].link = -1
	return &suffixAutomaton{
		states: st,
		last:   0,
	}
}

func (sam *suffixAutomaton) extend(c int) {
	cur := len(sam.states)
	sam.states = append(sam.states, samState{})
	for i := range sam.states[cur].next {
		sam.states[cur].next[i] = -1
	}
	sam.states[cur].len = sam.states[sam.last].len + 1
	p := sam.last
	for p != -1 && sam.states[p].next[c] == -1 {
		sam.states[p].next[c] = cur
		p = sam.states[p].link
	}
	if p == -1 {
		sam.states[cur].link = 0
	} else {
		q := sam.states[p].next[c]
		if sam.states[p].len+1 == sam.states[q].len {
			sam.states[cur].link = q
		} else {
			clone := len(sam.states)
			sam.states = append(sam.states, sam.states[q])
			sam.states[clone].len = sam.states[p].len + 1
			for p != -1 && sam.states[p].next[c] == q {
				sam.states[p].next[c] = clone
				p = sam.states[p].link
			}
			sam.states[q].link = clone
			sam.states[cur].link = clone
		}
	}
	sam.last = cur
}

func generateTests() []string {
	rng := rand.New(rand.NewSource(2039))
	var tests []string
	total := 0
	add := func(s string) {
		if len(s) == 0 {
			return
		}
		if total+len(s) > sumLimit {
			return
		}
		tests = append(tests, s)
		total += len(s)
	}

	seed := []string{
		"dcabaa",
		"a",
		"acayouknowwhocodeforcesbangladesh",
		"bbbbbb",
		"abcabcabc",
	}
	for _, s := range seed {
		add(s)
	}

	for total < sumLimit {
		n := rng.Intn(200) + 1
		start := make([]byte, n)
		for i := 0; i < n; i++ {
			start[i] = byte('a' + rng.Intn(26))
		}
		add(string(start))
		if len(tests) > 1000 {
			break
		}
	}
	return tests
}

func buildInput(tests []string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(tests))
	for _, s := range tests {
		fmt.Fprintln(&b, s)
	}
	return b.String()
}
