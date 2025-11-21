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

const (
	maxValue = 1 << 30
)

type singleTest struct {
	a       []int64
	queries []int64
}

type dataset struct {
	name  string
	tests []singleTest
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2165C-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleC")
	cmd := exec.Command("go", "build", "-o", outPath, "2165C.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build oracle: %v\n%s", err, out)
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return outPath, cleanup, nil
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
	sb.WriteString(strconv.Itoa(len(ds.tests)))
	sb.WriteByte('\n')
	for _, tc := range ds.tests {
		sb.WriteString(fmt.Sprintf("%d %d\n", len(tc.a), len(tc.queries)))
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
		for _, q := range tc.queries {
			sb.WriteString(strconv.FormatInt(q, 10))
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func deterministicDatasets() []dataset {
	return []dataset{
		{
			name: "tiny_cases",
			tests: []singleTest{
				{a: []int64{0}, queries: []int64{0, 1, 5}},
				{a: []int64{1, 2}, queries: []int64{0, 1, 2, 3}},
				{a: []int64{7, 3, 5}, queries: []int64{1, 4, 7}},
			},
		},
		{
			name: "all_equal",
			tests: []singleTest{
				{a: []int64{15, 15, 15, 15}, queries: []int64{0, 15, 31}},
				{a: []int64{1, 1, 1, 1, 1}, queries: []int64{0, 1}},
			},
		},
		{
			name: "single_large_query",
			tests: []singleTest{
				{a: []int64{(1 << 29) - 1, (1 << 29) - 2}, queries: []int64{(1 << 29) - 1}},
				{a: []int64{(1 << 29), (1 << 28), (1 << 27)}, queries: []int64{(1 << 29) - 1, (1 << 28)}},
			},
		},
	}
}

func randomArray(n int, rng *rand.Rand) []int64 {
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Int63n(maxValue)
	}
	return arr
}

func randomQueries(q int, rng *rand.Rand) []int64 {
	arr := make([]int64, q)
	for i := 0; i < q; i++ {
		arr[i] = rng.Int63n(maxValue)
	}
	return arr
}

func randomDatasets() []dataset {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	dsets := make([]dataset, 0, 5)
	for idx := 0; idx < 5; idx++ {
		var tests []singleTest
		totalN := 0
		totalQ := 0
		for totalN < 150000 && totalQ < 20000 {
			n := rng.Intn(5000) + 1
			if totalN+n > 500000 {
				break
			}
			q := rng.Intn(500) + 1
			if totalQ+q > 50000 {
				break
			}
			tests = append(tests, singleTest{
				a:       randomArray(n, rng),
				queries: randomQueries(q, rng),
			})
			totalN += n
			totalQ += q
			if len(tests) >= 50 && rng.Intn(3) == 0 {
				break
			}
		}
		if len(tests) == 0 {
			continue
		}
		dsets = append(dsets, dataset{
			name:  fmt.Sprintf("random_batch_%d", idx+1),
			tests: tests,
		})
	}
	return dsets
}

func stressDatasets() []dataset {
	// Create a dataset hitting constraints: n close to 5e5, q around 5e4.
	rng := rand.New(rand.NewSource(time.Now().UnixNano() + 42))
	var tests []singleTest
	totalN := 0
	totalQ := 0
	for totalN < 500000 {
		n := 50000
		if totalN+n > 500000 {
			n = 500000 - totalN
		}
		q := 5000
		if totalQ+q > 50000 {
			q = 50000 - totalQ
		}
		if q <= 0 {
			break
		}
		tests = append(tests, singleTest{
			a:       randomArray(n, rng),
			queries: randomQueries(q, rng),
		})
		totalN += n
		totalQ += q
		if totalQ >= 50000 {
			break
		}
	}
	return []dataset{
		{name: "stress_limits", tests: tests},
	}
}

func normalizeOutput(output string) []string {
	return strings.Fields(strings.TrimSpace(output))
}

func compareOutputs(expected, actual string) error {
	exp := normalizeOutput(expected)
	act := normalizeOutput(actual)
	if len(exp) != len(act) {
		return fmt.Errorf("token count mismatch: expected %d got %d", len(exp), len(act))
	}
	for i := range exp {
		if exp[i] != act[i] {
			return fmt.Errorf("mismatch at token %d: expected %s got %s", i+1, exp[i], act[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	datasets := append(deterministicDatasets(), randomDatasets()...)
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
		if err := compareOutputs(expected, actual); err != nil {
			fmt.Fprintf(os.Stderr, "dataset %d (%s) failed: %v\nInput:\n%s\nExpected:\n%s\nActual:\n%s\n", idx+1, ds.name, err, input, expected, actual)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d datasets passed.\n", len(datasets))
}
