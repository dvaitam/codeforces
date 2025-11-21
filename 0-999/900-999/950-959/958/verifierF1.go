package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const referenceSolutionRel = "0-999/900-999/950-959/958/958F1.go"

var referenceSolutionPath string

func init() {
	referenceSolutionPath = referenceSolutionRel
	if _, file, _, ok := runtime.Caller(0); ok {
		dir := filepath.Dir(file)
		candidate := filepath.Join(dir, "958F1.go")
		if _, err := os.Stat(candidate); err == nil {
			referenceSolutionPath = candidate
			return
		}
	}
	if abs, err := filepath.Abs(referenceSolutionRel); err == nil {
		if _, err := os.Stat(abs); err == nil {
			referenceSolutionPath = abs
		}
	}
}

type testCase struct {
	n, m int
	col  []int
	need []int
}

func inputString(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
	for i, v := range tc.col {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	for i, v := range tc.need {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func solve(tc testCase) string {
	total := 0
	for _, v := range tc.need {
		total += v
	}
	if total == 0 {
		return "NO"
	}
	for start := 0; start+total <= tc.n; start++ {
		cnt := make([]int, tc.m)
		for j := start; j < start+total; j++ {
			cnt[tc.col[j]-1]++
		}
		match := true
		for i := 0; i < tc.m; i++ {
			if cnt[i] != tc.need[i] {
				match = false
				break
			}
		}
		if match {
			return "YES"
		}
	}
	return "NO"
}

func runProgram(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseAnswer(out string) (string, error) {
	reader := strings.NewReader(out)
	var word string
	if _, err := fmt.Fscan(reader, &word); err != nil {
		return "", fmt.Errorf("failed to read verdict: %v\nfull output:\n%s", err, out)
	}
	word = strings.ToUpper(word)
	if word != "YES" && word != "NO" {
		return "", fmt.Errorf("invalid verdict %q", word)
	}
	return word, nil
}

func buildReferenceBinary() (string, func(), error) {
	tmpDir, err := os.MkdirTemp("", "958F1-ref")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "ref_958F1")
	cmd := exec.Command("go", "build", "-o", binPath, referenceSolutionPath)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, out.String())
	}
	cleanup := func() { os.RemoveAll(tmpDir) }
	return binPath, cleanup, nil
}

func randomCase(rng *rand.Rand, maxN int) testCase {
	n := rng.Intn(maxN-1) + 2
	m := rng.Intn(n) + 1
	col := make([]int, n)
	for i := range col {
		col[i] = rng.Intn(m) + 1
	}
	need := make([]int, m)
	total := rng.Intn(n) + 1
	for total > 0 {
		idx := rng.Intn(m)
		need[idx]++
		total--
	}
	return testCase{n: n, m: m, col: col, need: need}
}

func edgeCases() []testCase {
	return []testCase{
		{
			n:    1,
			m:    1,
			col:  []int{1},
			need: []int{1},
		},
		{
			n:    2,
			m:    1,
			col:  []int{1, 1},
			need: []int{2},
		},
		{
			n:    3,
			m:    2,
			col:  []int{1, 2, 1},
			need: []int{1, 1},
		},
	}
}

func genTests() []testCase {
	rng := rand.New(rand.NewSource(20250311))
	tests := edgeCases()
	for i := 0; i < 80; i++ {
		tests = append(tests, randomCase(rng, 10))
	}
	for i := 0; i < 80; i++ {
		tests = append(tests, randomCase(rng, 100))
	}
	for i := 0; i < 40; i++ {
		tests = append(tests, randomCase(rng, 100))
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierF1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	if bin == "--" {
		fmt.Println("usage: go run verifierF1.go /path/to/binary")
		os.Exit(1)
	}
	refBin, cleanup, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	tests := genTests()
	for i, tc := range tests {
		in := inputString(tc)
		expected := solve(tc)

		refOut, err := runProgram(refBin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\ninput:\n%soutput:\n%s\n", i+1, err, in, refOut)
			os.Exit(1)
		}
		refVerdict, err := parseAnswer(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference parse error on test %d: %v\ninput:\n%soutput:\n%s\n", i+1, err, in, refOut)
			os.Exit(1)
		}
		if refVerdict != expected {
			fmt.Fprintf(os.Stderr, "reference mismatch on test %d: expected %s got %s\ninput:\n%soutput:\n%s\n", i+1, expected, refVerdict, in, refOut)
			os.Exit(1)
		}

		out, runErr := runProgram(bin, in)
		if runErr != nil {
			fmt.Fprintf(os.Stderr, "test %d runtime error: %v\ninput:\n%soutput:\n%s\n", i+1, runErr, in, out)
			os.Exit(1)
		}
		ans, err := parseAnswer(out)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d parse error: %v\ninput:\n%soutput:\n%s\n", i+1, err, in, out)
			os.Exit(1)
		}
		if ans != expected {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %s got %s\ninput:\n%soutput:\n%s\n", i+1, expected, ans, in, out)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
