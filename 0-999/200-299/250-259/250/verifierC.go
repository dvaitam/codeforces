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

const refSource = "./250C.go"

type testCase struct {
	name  string
	n, k  int
	arr   []int
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/candidate_binary")
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
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}
		refAns, err := parseOutput(tc.k, refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candAns, err := parseOutput(tc.k, candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}

		if candAns != refAns {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: expected genre %d got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
				idx+1, tc.name, refAns, candAns, tc.input, refOut, candOut)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-250C-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref250C.bin")
	cmd := exec.Command("go", "build", "-o", binPath, refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		_ = os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference solution: %v\n%s", err, stderr.String())
	}
	cleanup := func() {
		_ = os.RemoveAll(dir)
	}
	return binPath, cleanup, nil
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseOutput(k int, output string) (int, error) {
	fields := strings.Fields(output)
	if len(fields) == 0 {
		return 0, fmt.Errorf("empty output")
	}
	if len(fields) > 1 {
		return 0, fmt.Errorf("expected single integer output, got %d tokens", len(fields))
	}
	val, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, fmt.Errorf("output is not an integer: %q", fields[0])
	}
	if val < 1 || val > k {
		return 0, fmt.Errorf("genre %d out of range [1,%d]", val, k)
	}
	return val, nil
}

func buildTests() []testCase {
	tests := []testCase{
		newManualTest("sample1", 3, []int{1, 1, 2, 3, 2, 3, 3, 1, 1, 3}),
		newManualTest("sample2", 3, []int{3, 1, 3, 2, 3, 1, 2}),
		newManualTest("two_blocks", 2, []int{1, 1, 1, 1, 2, 2, 2, 2}),
		newManualTest("alternating", 2, []int{1, 2, 1, 2, 1, 2, 1, 2}),
		newManualTest("k_equals_n", 5, []int{1, 2, 3, 4, 5}),
		newManualTest("long_prefix_single_genre", 4, []int{1, 1, 1, 1, 2, 3, 4, 2, 3, 4}),
		newManualTest("descending_segments", 4, []int{4, 4, 3, 3, 2, 2, 1, 1}),
	}

	tests = append(tests, makeLargeBlockTest())

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomTest(rng, i))
	}
	return tests
}

func newManualTest(name string, k int, arr []int) testCase {
	if k < 2 {
		panic("k must be at least 2")
	}
	copyArr := append([]int(nil), arr...)
	tc := testCase{
		name: name,
		n:    len(arr),
		k:    k,
		arr:  copyArr,
	}
	tc.input = formatInput(tc.n, tc.k, tc.arr)
	return tc
}

func makeLargeBlockTest() testCase {
	k := 6
	segments := []int{5000, 12000, 8000, 15000, 7000, 9000}
	arr := make([]int, 0, 1000)
	for idx, length := range segments {
		for i := 0; i < length; i++ {
			arr = append(arr, idx+1)
		}
	}
	// add some interleaving noise at the end to break pure segments
	for i := 0; i < 20000; i++ {
		arr = append(arr, (i%k)+1)
	}
	return newManualTest("large_blocks", k, arr)
}

func randomTest(rng *rand.Rand, idx int) testCase {
	k := rng.Intn(9) + 2      // 2..10
	maxN := 2000 + 200*k      // keep runtime manageable
	n := rng.Intn(maxN-k) + k // ensure n >= k
	arr := make([]int, n)
	for i := 0; i < k; i++ {
		arr[i] = i + 1
	}
	for i := k; i < n; i++ {
		switch rng.Intn(3) {
		case 0:
			arr[i] = arr[i-1]
		case 1:
			arr[i] = rng.Intn(k) + 1
		default:
			if i > 0 {
				arr[i] = (arr[i-1]%k + 1)
			} else {
				arr[i] = rng.Intn(k) + 1
			}
		}
	}
	rng.Shuffle(n, func(i, j int) {
		arr[i], arr[j] = arr[j], arr[i]
	})
	name := fmt.Sprintf("random_%d_n%d_k%d", idx+1, n, k)
	return newManualTest(name, k, arr)
}

func formatInput(n, k int, arr []int) string {
	var sb strings.Builder
	sb.Grow(n*4 + 32)
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}
