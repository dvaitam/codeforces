package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// refSource points to the local reference solution to avoid GOPATH resolution.
const refSource = "2092D.go"

type caseData struct {
	n    int
	s    string
	id   int
	name string
}

type testCase struct {
	name  string
	input string
	cases []caseData
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	candidate := os.Args[1]

	for idx, tc := range tests {
		out, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, out)
			os.Exit(1)
		}
		if err := validateOutput(tc, out); err != nil {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, out)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2092D-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return tmp.Name(), nil
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
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

func validateOutput(tc testCase, out string) error {
	reader := bufio.NewReader(strings.NewReader(out))
	var tOut int
	if _, err := fmt.Fscan(reader, &tOut); err != nil {
		return fmt.Errorf("failed to read t: %v", err)
	}
	if tOut != len(tc.cases) {
		return fmt.Errorf("expected t=%d, got %d", len(tc.cases), tOut)
	}

	for idx, c := range tc.cases {
		var m int
		if _, err := fmt.Fscan(reader, &m); err != nil {
			return fmt.Errorf("case %d: failed to read operations count: %v", idx+1, err)
		}
		if m == -1 {
			if hasSolution(c.s) {
				return fmt.Errorf("case %d: claimed impossible but a solution exists", idx+1)
			}
			continue
		}
		if m < 0 {
			return fmt.Errorf("case %d: negative operations count %d", idx+1, m)
		}
		if m > 2*c.n {
			return fmt.Errorf("case %d: operations %d exceed limit %d", idx+1, m, 2*c.n)
		}
		str := []byte(c.s)
		for op := 0; op < m; op++ {
			var pos int
			if _, err := fmt.Fscan(reader, &pos); err != nil {
				return fmt.Errorf("case %d: failed to read operation %d: %v", idx+1, op+1, err)
			}
			if pos <= 0 || pos >= len(str) {
				return fmt.Errorf("case %d: operation %d index %d out of bounds (length %d)", idx+1, op+1, pos, len(str))
			}
			if str[pos-1] == str[pos] {
				return fmt.Errorf("case %d: operation %d chooses equal neighbors %c %c", idx+1, op+1, str[pos-1], str[pos])
			}
			ins := thirdChar(str[pos-1], str[pos])
			str = append(str[:pos], append([]byte{ins}, str[pos:]...)...)
		}
		if !balanced(str) {
			return fmt.Errorf("case %d: final string not balanced (%d, %d, %d)", idx+1, count(str, 'L'), count(str, 'I'), count(str, 'T'))
		}
	}

	var extra string
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return fmt.Errorf("unexpected extra output: %s", extra)
	}
	return nil
}

func hasSolution(s string) bool {
	distinct := map[byte]bool{}
	for i := 0; i < len(s); i++ {
		distinct[s[i]] = true
	}
	return len(distinct) >= 2
}

func balanced(s []byte) bool {
	return count(s, 'L') == count(s, 'I') && count(s, 'L') == count(s, 'T')
}

func count(s []byte, ch byte) int {
	cnt := 0
	for _, v := range s {
		if v == ch {
			cnt++
		}
	}
	return cnt
}

func thirdChar(a, b byte) byte {
	if a == b {
		return 0
	}
	if (a == 'L' && b == 'I') || (a == 'I' && b == 'L') {
		return 'T'
	}
	if (a == 'L' && b == 'T') || (a == 'T' && b == 'L') {
		return 'I'
	}
	return 'L'
}

func generateTests() []testCase {
	tests := []testCase{
		buildCase("single-impossible", []string{"L"}),
		buildCase("already-balanced", []string{"LIT"}),
		buildCase("two-chars-easy", []string{"LI", "LLII"}),
		buildCase("sample-like", []string{"TILII"}),
	}
	rng := rand.New(rand.NewSource(20920228))
	for i := 0; i < 60; i++ {
		cases := rng.Intn(4) + 1
		strs := make([]string, 0, cases)
		for j := 0; j < cases; j++ {
			n := rng.Intn(100) + 1
			var sb strings.Builder
			for k := 0; k < n; k++ {
				sb.WriteByte("LIT"[rng.Intn(3)])
			}
			strs = append(strs, sb.String())
		}
		tests = append(tests, buildCase(fmt.Sprintf("random-%d", i+1), strs))
	}
	return tests
}

func buildCase(name string, strs []string) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(strs))
	cases := make([]caseData, len(strs))
	for i, s := range strs {
		fmt.Fprintf(&sb, "%d\n%s\n", len(s), s)
		cases[i] = caseData{n: len(s), s: s, id: i + 1, name: name}
	}
	return testCase{name: name, input: sb.String(), cases: cases}
}
