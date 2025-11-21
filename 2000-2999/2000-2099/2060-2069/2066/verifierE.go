package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const maxVolume = 1_000_000

type operation struct {
	op    string
	value int
}

type dataset struct {
	name    string
	initial []int
	queries []operation
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2066E-")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "oracleE")
	cmd := exec.Command("go", "build", "-o", binPath, "2066E.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build oracle: %v\n%s", err, out)
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return binPath, cleanup, nil
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func buildInput(ds dataset) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", len(ds.initial), len(ds.queries)))
	for i, v := range ds.initial {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	for _, op := range ds.queries {
		sb.WriteString(fmt.Sprintf("%s %d\n", op.op, op.value))
	}
	return sb.String()
}

func sampleDatasets() []dataset {
	return []dataset{
		{
			name:    "sample1",
			initial: []int{2, 2, 4, 11},
			queries: []operation{
				{"-", 2},
				{"+", 4},
				{"+", 3},
				{"+", 4},
				{"-", 4},
				{"+", 2},
				{"+", 2},
			},
		},
		{
			name:    "sample2",
			initial: []int{5000, 1000, 400, 400, 100, 99},
			queries: []operation{
				{"+", 1},
				{"-", 5000},
				{"-", 1},
				{"-", 400},
				{"-", 400},
				{"-", 100},
				{"-", 99},
			},
		},
		{
			name:    "small_manual",
			initial: []int{1},
			queries: []operation{
				{"+", 1},
				{"+", 2},
				{"-", 1},
				{"+", 3},
				{"-", 2},
			},
		},
	}
}

func randomDataset(name string, n, q int, rng *rand.Rand) dataset {
	if n < 1 {
		n = 1
	}
	initial := make([]int, n)
	curr := make([]int, n)
	for i := 0; i < n; i++ {
		val := rng.Intn(maxVolume) + 1
		initial[i] = val
		curr[i] = val
	}
	queries := make([]operation, q)
	for i := 0; i < q; i++ {
		if len(curr) == 1 {
			val := rng.Intn(maxVolume) + 1
			queries[i] = operation{"+", val}
			curr = append(curr, val)
			continue
		}
		if rng.Intn(2) == 0 {
			val := rng.Intn(maxVolume) + 1
			queries[i] = operation{"+", val}
			curr = append(curr, val)
		} else {
			idx := rng.Intn(len(curr))
			val := curr[idx]
			queries[i] = operation{"-", val}
			curr[idx] = curr[len(curr)-1]
			curr = curr[:len(curr)-1]
		}
	}
	return dataset{
		name:    name,
		initial: initial,
		queries: queries,
	}
}

func randomDatasets() []dataset {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	dsets := make([]dataset, 0, 5)
	configs := [][2]int{
		{5, 20},
		{10, 50},
		{100, 500},
		{1000, 2000},
		{2000, 4000},
	}
	for i, cfg := range configs {
		dsets = append(dsets, randomDataset(fmt.Sprintf("random_%d", i+1), cfg[0], cfg[1], rng))
	}
	return dsets
}

func stressDatasets() []dataset {
	rng := rand.New(rand.NewSource(time.Now().UnixNano() + 123))
	return []dataset{
		randomDataset("stress_mid", 50000, 50000, rng),
		randomDataset("stress_large", 200000, 200000, rng),
	}
}

func normalizeOutput(out string) []string {
	fields := strings.Fields(out)
	for i := range fields {
		fields[i] = strings.ToLower(fields[i])
	}
	return fields
}

func compareOutputs(expected, actual string, expectedCount int) error {
	exp := normalizeOutput(expected)
	act := normalizeOutput(actual)
	if len(exp) != expectedCount {
		return fmt.Errorf("oracle produced %d tokens, expected %d", len(exp), expectedCount)
	}
	if len(act) != expectedCount {
		return fmt.Errorf("expected %d tokens, got %d", expectedCount, len(act))
	}
	for i := 0; i < expectedCount; i++ {
		if exp[i] != act[i] {
			return fmt.Errorf("mismatch at token %d: expected %s got %s", i+1, exp[i], act[i])
		}
	}
	return nil
}

func verifyDataset(bin string, ds dataset) error {
	input := buildInput(ds)
	expected, err := runBinary(bin, input)
	if err != nil {
		return fmt.Errorf("oracle failed: %v", err)
	}
	actual, err := runBinary(os.Args[1], input)
	if err != nil {
		return fmt.Errorf("target runtime error: %v", err)
	}
	expectedCount := len(ds.queries) + 1
	if err := compareOutputs(expected, actual, expectedCount); err != nil {
		return fmt.Errorf("%v\ninput:\n%s\nexpected:\n%s\nactual:\n%s", err, input, expected, actual)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	datasets := append(sampleDatasets(), randomDatasets()...)
	datasets = append(datasets, stressDatasets()...)

	for idx, ds := range datasets {
		input := buildInput(ds)
		expected, err := runBinary(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on dataset %d (%s): %v\n", idx+1, ds.name, err)
			os.Exit(1)
		}
		actual, err := runBinary(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "dataset %d (%s) runtime error: %v\n", idx+1, ds.name, err)
			os.Exit(1)
		}
		expectedCount := len(ds.queries) + 1
		if err := compareOutputs(expected, actual, expectedCount); err != nil {
			fmt.Fprintf(os.Stderr, "dataset %d (%s) failed: %v\n", idx+1, ds.name, err)
			fmt.Fprintf(os.Stderr, "Input:\n%s\nExpected:\n%s\nActual:\n%s\n", input, expected, actual)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d datasets passed.\n", len(datasets))
}
