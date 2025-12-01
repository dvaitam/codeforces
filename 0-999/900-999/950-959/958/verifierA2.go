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

const (
	refSourceA2  = "958A2.go"
	randomTrials = 150
)

type testInstance struct {
	input  string
	N, M   int
	first  []string
	second []string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA2.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference(refSourceA2)
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	tests := deterministicCases()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < randomTrials; i++ {
		tests = append(tests, randomCase(rng))
	}

	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fail("reference failed on case %d: %v\ninput:\n%s", idx+1, err, tc.input)
		}
		if err := validateAnswer(refOut, tc); err != nil {
			fail("reference produced invalid output on case %d: %v\ninput:\n%s\noutput:\n%s", idx+1, err, tc.input, refOut)
		}
		got, err := runCandidate(candidate, tc.input)
		if err != nil {
			fail("candidate failed on case %d: %v\ninput:\n%s", idx+1, err, tc.input)
		}
		if err := validateAnswer(got, tc); err != nil {
			fail("candidate invalid on case %d: %v\ninput:\n%s\noutput:\n%s", idx+1, err, tc.input, got)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference(src string) (string, error) {
	tmp, err := os.CreateTemp("", "958A2-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	source := filepath.Join(".", filepath.Clean(src))
	cmd := exec.Command("go", "build", "-o", tmp.Name(), source)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	return runCommand(cmd, input)
}

func runCandidate(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	return runCommand(cmd, input)
}

func runCommand(cmd *exec.Cmd, input string) (string, error) {
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

func validateAnswer(out string, tc testInstance) error {
	fields := strings.Fields(out)
	if len(fields) != 2 {
		return fmt.Errorf("expected exactly two integers, got %d tokens", len(fields))
	}
	i, err := strconv.Atoi(fields[0])
	if err != nil {
		return fmt.Errorf("invalid row index %q: %v", fields[0], err)
	}
	j, err := strconv.Atoi(fields[1])
	if err != nil {
		return fmt.Errorf("invalid column index %q: %v", fields[1], err)
	}
	maxPos := tc.N - tc.M + 1
	if i < 1 || i > maxPos {
		return fmt.Errorf("row index %d out of range [1,%d]", i, maxPos)
	}
	if j < 1 || j > maxPos {
		return fmt.Errorf("column index %d out of range [1,%d]", j, maxPos)
	}
	row := i - 1
	col := j - 1
	for r := 0; r < tc.M; r++ {
		rowFirst := tc.first[row+r]
		rowSecond := tc.second[r]
		for c := 0; c < tc.M; c++ {
			if rowFirst[c] != rowSecond[col+c] {
				return fmt.Errorf("mismatch at local (%d,%d)", r+1, c+1)
			}
		}
	}
	return nil
}

func deterministicCases() []testInstance {
	var tests []testInstance

	tests = append(tests, buildInstance(1, 1, 0, 0, [][]byte{{'h'}}, constantFiller('h'), constantFiller('h')))

	tests = append(tests, buildInstance(3, 2, 1, 1,
		[][]byte{
			[]byte("ab"),
			[]byte("cd"),
		},
		constantFiller('x'),
		constantFiller('y'),
	))

	tests = append(tests, buildInstance(5, 3, 2, 0,
		[][]byte{
			[]byte("qwe"),
			[]byte("rty"),
			[]byte("uio"),
		},
		constantFiller('a'),
		constantFiller('b'),
	))

	tests = append(tests, buildInstance(200, 200, 0, 0, patternedBlock(200, 'a'), constantFiller('z'), constantFiller('y')))

	largeBlock := make([][]byte, 200)
	for r := range largeBlock {
		row := make([]byte, 200)
		for c := range row {
			row[c] = byte('a' + (r+c)%26)
		}
		largeBlock[r] = row
	}
	tests = append(tests, buildInstance(2000, 200, 100, 400, largeBlock, constantFiller('m'), constantFiller('n')))

	seeds := []int64{7, 1337, 2024}
	for _, seed := range seeds {
		rng := rand.New(rand.NewSource(seed))
		tests = append(tests, randomCase(rng))
	}

	return tests
}

func randomCase(rng *rand.Rand) testInstance {
	M := rng.Intn(200) + 1
	N := rng.Intn(2000-M+1) + M
	block := randomBlock(rng, M)
	rowStart := rng.Intn(N - M + 1)
	colStart := rng.Intn(N - M + 1)
	fill1 := randomFiller(rng)
	fill2 := randomFiller(rng)
	return buildInstance(N, M, rowStart, colStart, block, fill1, fill2)
}

func buildInstance(N, M, rowStart, colStart int, block [][]byte, fillFirst, fillSecond func() byte) testInstance {
	if fillFirst == nil {
		fillFirst = constantFiller('a')
	}
	if fillSecond == nil {
		fillSecond = fillFirst
	}
	first := make([][]byte, N)
	for i := 0; i < N; i++ {
		row := make([]byte, M)
		for j := 0; j < M; j++ {
			row[j] = fillFirst()
		}
		first[i] = row
	}
	for r := 0; r < M; r++ {
		copy(first[rowStart+r], block[r])
	}

	second := make([][]byte, M)
	for i := 0; i < M; i++ {
		row := make([]byte, N)
		for j := 0; j < N; j++ {
			row[j] = fillSecond()
		}
		second[i] = row
	}
	for r := 0; r < M; r++ {
		copy(second[r][colStart:colStart+M], block[r])
	}

	firstRows := make([]string, N)
	for i := range first {
		firstRows[i] = string(first[i])
	}
	secondRows := make([]string, M)
	for i := range second {
		secondRows[i] = string(second[i])
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", N, M)
	for _, row := range firstRows {
		sb.WriteString(row)
		sb.WriteByte('\n')
	}
	for _, row := range secondRows {
		sb.WriteString(row)
		sb.WriteByte('\n')
	}

	return testInstance{
		input:  sb.String(),
		N:      N,
		M:      M,
		first:  firstRows,
		second: secondRows,
	}
}

func patternedBlock(size int, start byte) [][]byte {
	block := make([][]byte, size)
	for r := 0; r < size; r++ {
		row := make([]byte, size)
		for c := 0; c < size; c++ {
			row[c] = byte(int(start) + (r+c)%26)
		}
		block[r] = row
	}
	return block
}

func randomBlock(rng *rand.Rand, size int) [][]byte {
	block := make([][]byte, size)
	for r := 0; r < size; r++ {
		row := make([]byte, size)
		for c := 0; c < size; c++ {
			row[c] = byte('a' + rng.Intn(26))
		}
		block[r] = row
	}
	return block
}

func constantFiller(ch byte) func() byte {
	return func() byte { return ch }
}

func randomFiller(rng *rand.Rand) func() byte {
	return func() byte {
		return byte('a' + rng.Intn(26))
	}
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
