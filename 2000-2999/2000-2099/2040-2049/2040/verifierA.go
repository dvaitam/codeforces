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
)

const refSource = "./2040A.go"

type testCase struct {
	name string
	ns   []int
	ks   []int
	arrs [][]int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	tests := buildTests()
	for idx, tc := range tests {
		input := buildInput(tc)

		refOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, input, refOut)
			os.Exit(1)
		}
		if err := validateOutput(tc, refOut); err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, input, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, input, candOut)
			os.Exit(1)
		}
		if err := validateOutput(tc, candOut); err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, input, candOut)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-2040A-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref2040A.bin")
	cmd := exec.Command("go", "build", "-o", binPath, refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	cleanup := func() { _ = os.RemoveAll(dir) }
	return binPath, cleanup, nil
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
	err := cmd.Run()
	return out.String(), err
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tc.ns)))
	for idx, n := range tc.ns {
		sb.WriteString(fmt.Sprintf("%d %d\n", n, tc.ks[idx]))
		for i, v := range tc.arrs[idx] {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func validateOutput(tc testCase, output string) error {
	lines := splitNonEmptyLines(output)
	expectedLines := 0
	for range tc.ns {
		expectedLines++ // YES/NO
		// add one more line if there is a winning answer
	}

	idx := 0
	for caseIdx := range tc.ns {
		if idx >= len(lines) {
			return fmt.Errorf("case %d: missing YES/NO line", caseIdx+1)
		}
		state := strings.ToUpper(lines[idx])
		idx++
		if state != "YES" && state != "NO" {
			return fmt.Errorf("case %d: expected YES/NO, got %q", caseIdx+1, lines[idx-1])
		}
		hasWinning := findWinningIndex(tc.arrs[caseIdx], tc.ks[caseIdx])
		if state == "NO" {
			if hasWinning {
				return fmt.Errorf("case %d: answer NO but winning index exists", caseIdx+1)
			}
			continue
		}
		if !hasWinning {
			return fmt.Errorf("case %d: answer YES but no winning index exists", caseIdx+1)
		}
		if idx >= len(lines) {
			return fmt.Errorf("case %d: missing index after YES", caseIdx+1)
		}
		chosen, err := strconv.Atoi(lines[idx])
		idx++
		if err != nil {
			return fmt.Errorf("case %d: invalid integer index %q", caseIdx+1, lines[idx-1])
		}
		if chosen < 1 || chosen > tc.ns[caseIdx] {
			return fmt.Errorf("case %d: index %d out of range 1..%d", caseIdx+1, chosen, tc.ns[caseIdx])
		}
		if !isWinning(tc.arrs[caseIdx], tc.ks[caseIdx], chosen-1) {
			return fmt.Errorf("case %d: index %d is not a winning choice", caseIdx+1, chosen)
		}
		expectedLines++ // account for YES + index
	}
	if idx != len(lines) {
		return fmt.Errorf("expected %d lines, got %d", idx, len(lines))
	}
	return nil
}

func splitNonEmptyLines(output string) []string {
	raw := strings.Split(output, "\n")
	lines := make([]string, 0, len(raw))
	for _, line := range raw {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			lines = append(lines, trimmed)
		}
	}
	return lines
}

func findWinningIndex(arr []int, k int) bool {
	for i := range arr {
		if isWinning(arr, k, i) {
			return true
		}
	}
	return false
}

func isWinning(arr []int, k int, idx int) bool {
	for j := range arr {
		if j == idx {
			continue
		}
		if abs(arr[idx]-arr[j])%k == 0 {
			return false
		}
	}
	return true
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func buildTests() []testCase {
	tests := []testCase{
		{
			name: "sample",
			ns:   []int{3, 4, 5},
			ks:   []int{2, 2, 3},
			arrs: [][]int{
				{1, 2, 3},
				{1, 2, 4, 5},
				{10, 7, 3, 4, 5},
			},
		},
		{
			name: "single_element",
			ns:   []int{1},
			ks:   []int{1},
			arrs: [][]int{{5}},
		},
	}

	rng := rand.New(rand.NewSource(123456789))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomTest(rng, i))
	}
	return tests
}

func randomTest(rng *rand.Rand, idx int) testCase {
	t := rng.Intn(5) + 1
	ns := make([]int, t)
	ks := make([]int, t)
	arrs := make([][]int, t)
	for i := 0; i < t; i++ {
		n := rng.Intn(100) + 1
		k := rng.Intn(100) + 1
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rng.Intn(100) + 1
		}
		ns[i] = n
		ks[i] = k
		arrs[i] = arr
	}
	return testCase{
		name: fmt.Sprintf("random_%d", idx+1),
		ns:   ns,
		ks:   ks,
		arrs: arrs,
	}
}
