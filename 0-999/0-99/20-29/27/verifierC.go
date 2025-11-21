package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func buildReference() (string, error) {
	ref := "./refC.bin"
	cmd := exec.Command("go", "build", "-o", ref, "27C.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return ref, nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
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

func parseInputSequence(input string) ([]int, error) {
	in := bufio.NewReader(strings.NewReader(input))
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return nil, fmt.Errorf("failed to read n: %v", err)
	}
	if n < 1 || n > 100000 {
		return nil, fmt.Errorf("n out of range: %d", n)
	}
	seq := make([]int, n)
	for i := 0; i < n; i++ {
		if _, err := fmt.Fscan(in, &seq[i]); err != nil {
			return nil, fmt.Errorf("failed to read a[%d]: %v", i, err)
		}
	}
	return seq, nil
}

func parseReferenceLength(out string) (int, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return 0, fmt.Errorf("empty reference output")
	}
	val, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, fmt.Errorf("invalid reference length %q: %v", fields[0], err)
	}
	if val < 0 {
		return 0, fmt.Errorf("negative reference length %d", val)
	}
	return val, nil
}

type answer struct {
	k       int
	indices []int
}

func parseCandidateAnswer(out string, n int) (answer, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return answer{}, fmt.Errorf("empty output")
	}
	k, err := strconv.Atoi(fields[0])
	if err != nil {
		return answer{}, fmt.Errorf("invalid length %q: %v", fields[0], err)
	}
	if k < 0 {
		return answer{}, fmt.Errorf("length cannot be negative: %d", k)
	}
	if k > n {
		return answer{}, fmt.Errorf("length %d exceeds n=%d", k, n)
	}
	if k == 0 {
		if len(fields) != 1 {
			return answer{}, fmt.Errorf("expected only length for 0 answer, got %d numbers", len(fields))
		}
		return answer{k: 0}, nil
	}
	if len(fields) != k+1 {
		return answer{}, fmt.Errorf("expected %d indices, got %d", k, len(fields)-1)
	}
	if k < 3 {
		return answer{}, fmt.Errorf("length %d is too small to be unordered", k)
	}
	indices := make([]int, k)
	for i := 0; i < k; i++ {
		idx, err := strconv.Atoi(fields[i+1])
		if err != nil {
			return answer{}, fmt.Errorf("invalid index %q: %v", fields[i+1], err)
		}
		indices[i] = idx
	}
	return answer{k: k, indices: indices}, nil
}

func validateSubsequence(seq []int, indices []int) error {
	n := len(seq)
	prev := 0
	values := make([]int, len(indices))
	for i, idx := range indices {
		if idx < 1 || idx > n {
			return fmt.Errorf("index %d out of range", idx)
		}
		if i > 0 && idx <= prev {
			return fmt.Errorf("indices must be strictly increasing")
		}
		values[i] = seq[idx-1]
		prev = idx
	}
	if isOrdered(values) {
		return fmt.Errorf("reported subsequence %v is ordered", values)
	}
	return nil
}

func isOrdered(seq []int) bool {
	nonDec := true
	nonInc := true
	for i := 1; i < len(seq); i++ {
		if seq[i] < seq[i-1] {
			nonDec = false
		}
		if seq[i] > seq[i-1] {
			nonInc = false
		}
	}
	return nonDec || nonInc
}

func deterministicCases() []string {
	return []string{
		"1\n5\n",
		"2\n1 1\n",
		"3\n1 2 3\n",
		"3\n3 2 1\n",
		"3\n1 3 2\n",
		"5\n5 5 5 5 5\n",
		"6\n1 2 3 2 2 2\n",
		"6\n10 9 8 9 7 6\n",
		"7\n-1 -1 -1 0 -1 -2 -3\n",
		"8\n1 4 4 4 5 3 3 3\n",
	}
}

func generateRandomCase(rng *rand.Rand) string {
	n := rng.Intn(50) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		val := rng.Intn(2000001) - 1000000
		sb.WriteString(strconv.Itoa(val))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func verifyCase(candidate, ref, input string) error {
	seq, err := parseInputSequence(input)
	if err != nil {
		return err
	}
	refOut, err := runProgram(ref, input)
	if err != nil {
		return fmt.Errorf("reference error: %v", err)
	}
	expectedLen, err := parseReferenceLength(refOut)
	if err != nil {
		return fmt.Errorf("invalid reference output: %v", err)
	}
	candOut, err := runProgram(candidate, input)
	if err != nil {
		return err
	}
	ans, err := parseCandidateAnswer(candOut, len(seq))
	if err != nil {
		return err
	}
	if ans.k != expectedLen {
		return fmt.Errorf("expected length %d, got %d", expectedLen, ans.k)
	}
	if ans.k > 0 {
		if err := validateSubsequence(seq, ans.indices); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	ref, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := deterministicCases()
	for i := 0; i < 200; i++ {
		tests = append(tests, generateRandomCase(rng))
	}

	for idx, input := range tests {
		if err := verifyCase(candidate, ref, input); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
