package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	refSource = "./2168A2.go"
	metaDir   = "2000-2999/2100-2199/2160-2169/2168"
)

type testCase struct {
	n    int
	arr  []int64
	name string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA2.go /path/to/binary")
		os.Exit(1)
	}

	tests := generateTests()

	if err := os.MkdirAll(metaDir, 0o755); err != nil {
		fmt.Fprintln(os.Stderr, "failed to create meta dir:", err)
		os.Exit(1)
	}

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	candidate := os.Args[1]

	for idx, tc := range tests {
		if err := verifyCase(idx+1, tc, candidate, refBin); err != nil {
			fmt.Fprintf(os.Stderr, "Test %d (%s) failed: %v\n", idx+1, tc.name, err)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func verifyCase(id int, tc testCase, candidate, refBin string) error {
	meta := filepath.Join(metaDir, fmt.Sprintf("verifierA2_case%d_%s.json", id, tc.name))

	if _, err := runPhase(candidate, "first", tc, meta); err != nil {
		return fmt.Errorf("candidate first run failed: %w", err)
	}
	candDecoded, err := runPhase(candidate, "second", tc, meta)
	if err != nil {
		return fmt.Errorf("candidate second run failed: %w", err)
	}

	if _, err := runPhase(refBin, "first", tc, meta); err != nil {
		return fmt.Errorf("reference first run failed: %w", err)
	}
	refDecoded, err := runPhase(refBin, "second", tc, meta)
	if err != nil {
		return fmt.Errorf("reference second run failed: %w", err)
	}

	if len(candDecoded) != len(refDecoded) {
		return fmt.Errorf("decoded length mismatch: got %d expected %d", len(candDecoded), len(refDecoded))
	}
	for i := range candDecoded {
		if candDecoded[i] != refDecoded[i] {
			return fmt.Errorf("decoded arrays differ at index %d: got %d expected %d", i, candDecoded[i], refDecoded[i])
		}
	}

	os.Remove(meta)
	return nil
}

func runPhase(bin, phase string, tc testCase, meta string) ([]int64, error) {
	label := binLabel(bin)

	var input strings.Builder
	if phase == "first" {
		fmt.Fprintln(&input, "first")
		fmt.Fprintln(&input, tc.n)
		for i := 0; i < tc.n; i++ {
			if i > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprintf(&input, "%d", tc.arr[i])
		}
		input.WriteByte('\n')
	} else {
		info, err := readMeta(meta, label+"_first")
		if err != nil {
			return nil, fmt.Errorf("could not read encoded string: %w", err)
		}
		if len(info) == 0 {
			return nil, fmt.Errorf("missing encoded string for phase two")
		}
		fmt.Fprintln(&input, "second")
		fmt.Fprintln(&input, info[0])
	}

	cmd, err := commandFor(bin)
	if err != nil {
		return nil, err
	}
	cmd.Stdin = strings.NewReader(input.String())
	var out strings.Builder
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("%w: %s", err, out.String())
	}

	output := strings.TrimSpace(out.String())
	if phase == "first" {
		if len(output) == 0 {
			return nil, fmt.Errorf("empty encoded string")
		}
		if err := writeMeta(meta, label+"_first", []string{output}); err != nil {
			return nil, err
		}
		return nil, nil
	}

	if len(output) == 0 {
		return nil, fmt.Errorf("empty decoded output")
	}
	vals, err := parseNumberList(output)
	if err != nil {
		return nil, err
	}
	if err := writeMeta(meta, label+"_second", stringify(vals)); err != nil {
		return nil, err
	}
	return vals, nil
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2168A2-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("build failed: %v\n%s", err, out)
	}
	return tmp.Name(), nil
}

func binLabel(path string) string {
	clean := filepath.Clean(path)
	if strings.HasSuffix(clean, ".go") {
		if clean == filepath.Clean(refSource) {
			return "reference"
		}
		return "candidate"
	}
	if strings.Contains(filepath.Base(clean), "ref") {
		return "reference"
	}
	return "candidate"
}

func writeMeta(path, key string, data []string) error {
	meta := make(map[string][]string)
	if content, err := os.ReadFile(path); err == nil {
		if err := json.Unmarshal(content, &meta); err != nil {
			return err
		}
	}
	meta[key] = data
	bytes, err := json.Marshal(meta)
	if err != nil {
		return err
	}
	return os.WriteFile(path, bytes, 0o644)
}

func readMeta(path, key string) ([]string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	meta := make(map[string][]string)
	if err := json.Unmarshal(content, &meta); err != nil {
		return nil, err
	}
	val, ok := meta[key]
	if !ok {
		return nil, fmt.Errorf("key %s missing", key)
	}
	return val, nil
}

func parseNumberList(line string) ([]int64, error) {
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no numbers in output")
	}
	vals := make([]int64, len(fields))
	for i, f := range fields {
		v, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q: %w", f, err)
		}
		vals[i] = v
	}
	return vals, nil
}

func stringify(vals []int64) []string {
	result := make([]string, len(vals))
	for i, v := range vals {
		result[i] = strconv.FormatInt(v, 10)
	}
	return result
}

func commandFor(path string) (*exec.Cmd, error) {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path), nil
	}
	return exec.Command(path), nil
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(21682168))
	var tests []testCase

	tests = append(tests, testCase{n: 1, arr: []int64{1}, name: "single"})
	tests = append(tests, testCase{n: 5, arr: []int64{100, 200, 300, 400, 500}, name: "sample"})

	for i := 0; i < 20; i++ {
		n := rng.Intn(100) + 1
		arr := make([]int64, n)
		for j := range arr {
			arr[j] = rng.Int63n(1_000_000_000) + 1
		}
		tests = append(tests, testCase{n: n, arr: arr, name: "rand" + strconv.Itoa(i)})
	}

	tests = append(tests, testCase{n: 10000, arr: constantArray(10000, 1_000_000_000), name: "max"})
	return tests
}

func constantArray(n int, val int64) []int64 {
	arr := make([]int64, n)
	for i := range arr {
		arr[i] = val
	}
	return arr
}
