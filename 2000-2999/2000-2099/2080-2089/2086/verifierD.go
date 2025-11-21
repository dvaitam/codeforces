package main

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

const mod = 998244353

func buildReferenceBinary() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("failed to locate current file")
	}
	dir := filepath.Dir(file)
	refPath := filepath.Join(dir, "2086D_ref.bin")
	cmd := exec.Command("go", "build", "-o", refPath, "2086D.go")
	cmd.Dir = dir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		_ = os.Remove(refPath)
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	return refPath, nil
}

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		if stderr.Len() > 0 {
			return "", fmt.Errorf("%v: %s", err, strings.TrimSpace(stderr.String()))
		}
		return "", err
	}
	return stdout.String(), nil
}

func buildInput(cases [][26]int) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(cases)))
	sb.WriteByte('\n')
	for _, c := range cases {
		for i, v := range c {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func randomCounts(rng *rand.Rand, maxVal int) [26]int {
	var arr [26]int
	for i := range arr {
		arr[i] = rng.Intn(maxVal + 1)
	}
	total := 0
	for _, v := range arr {
		total += v
	}
	if total == 0 {
		arr[rng.Intn(26)] = 1
	}
	return arr
}

func bigCounts(total int) [26]int {
	var arr [26]int
	arr[0] = total
	return arr
}

func mixCounts(total int, chunk int) [26]int {
	var arr [26]int
	for i := 0; i < 26 && total > 0; i++ {
		val := chunk
		if val > total {
			val = total
		}
		arr[i] = val
		total -= val
	}
	if arr[0] == 0 {
		arr[0] = 1
	}
	return arr
}

func arrFromSlice(vals []int) [26]int {
	if len(vals) != 26 {
		panic("slice must have 26 elements")
	}
	var arr [26]int
	copy(arr[:], vals)
	return arr
}

func parseOutputs(out string, expected int) ([]int64, error) {
	reader := strings.NewReader(out)
	res := make([]int64, 0, expected)
	for len(res) < expected {
		var val int64
		if _, err := fmt.Fscan(reader, &val); err != nil {
			if err == io.EOF {
				return nil, fmt.Errorf("expected %d numbers, got %d", expected, len(res))
			}
			return nil, fmt.Errorf("failed to parse output: %v", err)
		}
		res = append(res, (val%mod+mod)%mod)
	}
	var extra int64
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return nil, fmt.Errorf("expected %d numbers, found extra data", expected)
	}
	return res, nil
}

func generateTests() []string {
	rng := rand.New(rand.NewSource(42))
	var tests []string

	sampleCases := [][26]int{
		arrFromSlice([]int{2, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}),
		arrFromSlice([]int{3, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0}),
		arrFromSlice([]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}),
		arrFromSlice([]int{1, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}),
		arrFromSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 233527, 233827}),
	}
	tests = append(tests, buildInput(sampleCases))

	var batch [][26]int
	for i := 0; i < 50; i++ {
		batch = append(batch, randomCounts(rng, 5))
	}
	tests = append(tests, buildInput(batch))

	batch = batch[:0]
	for i := 0; i < 30; i++ {
		batch = append(batch, randomCounts(rng, 50))
	}
	tests = append(tests, buildInput(batch))

	batch = batch[:0]
	for i := 0; i < 10; i++ {
		batch = append(batch, randomCounts(rng, 1000))
	}
	tests = append(tests, buildInput(batch))

	tests = append(tests, buildInput([][26]int{bigCounts(500000)}))
	tests = append(tests, buildInput([][26]int{mixCounts(400000, 12345), mixCounts(123456, 3210)}))

	batch = batch[:0]
	for i := 0; i < 200; i++ {
		batch = append(batch, randomCounts(rng, 20))
	}
	tests = append(tests, buildInput(batch))

	// Large t with small counts
	batch = batch[:0]
	for i := 0; i < 1000; i++ {
		batch = append(batch, randomCounts(rng, 2))
	}
	tests = append(tests, buildInput(batch))

	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	refPath, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refPath)

	tests := generateTests()
	for idx, input := range tests {
		lines := strings.SplitN(strings.TrimSpace(input), "\n", 2)
		var t int
		if _, err := fmt.Sscan(lines[0], &t); err != nil {
			fmt.Fprintf(os.Stderr, "internal error: failed to read t on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}

		refOut, err := runProgram(refPath, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		refVals, err := parseOutputs(refOut, t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}

		userOut, err := runProgram(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", idx+1, err)
			os.Exit(1)
		}
		userVals, err := parseOutputs(userOut, t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: invalid output: %v\ninput:\n%s\noutput:\n%s\n", idx+1, err, input, userOut)
			os.Exit(1)
		}

		for i := 0; i < t; i++ {
			if userVals[i]%mod != refVals[i]%mod {
				fmt.Fprintf(os.Stderr, "test %d case %d mismatch: expected %d, got %d\ninput:\n%s\n", idx+1, i+1, refVals[i], userVals[i], input)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}
