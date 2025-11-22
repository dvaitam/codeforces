package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	input string
	ns    []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refSrc, err := locateReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	refBin, err := buildReference(refSrc)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for idx, tc := range tests {
		// Run and validate reference to ensure test is feasible and validator is sane.
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\nInput:\n%s", idx+1, err, tc.input)
			os.Exit(1)
		}
		if err := validateOutput(tc, refOut); err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d: %v\nInput:\n%sOutput:\n%s", idx+1, err, tc.input, refOut)
			os.Exit(1)
		}

		gotOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\nInput:\n%s", idx+1, err, tc.input)
			os.Exit(1)
		}
		if err := validateOutput(tc, gotOut); err != nil {
			fmt.Fprintf(os.Stderr, "candidate produced invalid output on test %d: %v\nInput:\n%sOutput:\n%s", idx+1, err, tc.input, gotOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func locateReference() (string, error) {
	candidates := []string{
		"2049C.go",
		filepath.Join("2000-2999", "2000-2099", "2040-2049", "2049", "2049C.go"),
	}
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("could not find 2049C.go")
}

func buildReference(src string) (string, error) {
	outPath := filepath.Join(os.TempDir(), fmt.Sprintf("ref2049C_%d.bin", time.Now().UnixNano()))
	cmd := exec.Command("go", "build", "-o", outPath, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return outPath, nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstderr:\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func validateOutput(tc testCase, out string) error {
	ns := tc.ns
	vals, err := parseOutputs(out, ns)
	if err != nil {
		return err
	}

	reader := strings.NewReader(tc.input)
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return fmt.Errorf("failed to read t: %v", err)
	}
	offset := 0
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		var n, x, y int
		if _, err := fmt.Fscan(reader, &n, &x, &y); err != nil {
			return fmt.Errorf("failed to read case %d header: %v", caseIdx+1, err)
		}
		ans := vals[offset : offset+n]
		offset += n
		if err := checkCase(n, x, y, ans); err != nil {
			return fmt.Errorf("case %d invalid: %v", caseIdx+1, err)
		}
	}
	if offset != len(vals) {
		return fmt.Errorf("extra output values detected")
	}
	return nil
}

func parseOutputs(out string, ns []int) ([]int64, error) {
	fields := strings.Fields(out)
	total := 0
	for _, n := range ns {
		total += n
	}
	if len(fields) != total {
		return nil, fmt.Errorf("expected %d integers, got %d", total, len(fields))
	}
	res := make([]int64, total)
	for i, f := range fields {
		v, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", f)
		}
		res[i] = v
	}
	return res, nil
}

func checkCase(n, x, y int, ans []int64) error {
	if len(ans) != n {
		return fmt.Errorf("expected %d numbers, got %d", n, len(ans))
	}
	adj := make([][]int, n)
	addEdge := func(u, v int) {
		for _, nei := range adj[u] {
			if nei == v {
				return
			}
		}
		adj[u] = append(adj[u], v)
	}
	for i := 0; i < n; i++ {
		addEdge(i, (i+1)%n)
		addEdge((i+1)%n, i)
	}
	x--
	y--
	addEdge(x, y)
	addEdge(y, x)

	for i := 0; i < n; i++ {
		if ans[i] < 0 {
			return fmt.Errorf("node %d has negative value %d", i+1, ans[i])
		}
		neis := adj[i]
		seen := [4]bool{}
		for _, v := range neis {
			val := ans[v]
			if val >= 0 && val < 4 {
				seen[val] = true
			}
		}
		expected := int64(0)
		for j := 0; j < 4; j++ {
			if !seen[j] {
				expected = int64(j)
				break
			}
		}
		if ans[i] != expected {
			return fmt.Errorf("node %d: expected %d, got %d", i+1, expected, ans[i])
		}
	}
	return nil
}

func generateTests() []testCase {
	tests := make([]testCase, 0, 50)
	tests = append(tests, sampleTest())
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 40; i++ {
		tests = append(tests, randomTest(rng))
	}
	return tests
}

func sampleTest() testCase {
	// The sample from the statement (t=7).
	input := "7\n5 1 3\n4 2 4\n6 3 5\n7 3 6\n3 2 3\n5 1 5\n6 2 5\n"
	ns := []int{5, 4, 6, 7, 3, 5, 6}
	return testCase{input: input, ns: ns}
}

func randomTest(rng *rand.Rand) testCase {
	t := rng.Intn(4) + 1 // 1..4 test cases
	ns := make([]int, t)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		n := rng.Intn(20) + 3 // 3..22
		ns[i] = n
		x := rng.Intn(n) + 1
		y := rng.Intn(n-1) + 1
		if y >= x {
			y++
		}
		if x > y {
			x, y = y, x
		}
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, x, y))
	}
	return testCase{input: sb.String(), ns: ns}
}
