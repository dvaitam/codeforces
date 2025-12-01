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

const refSource = "./761B.go"

type testCase struct {
	name  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
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
		exp, err := parseAnswer(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		got, err := parseAnswer(candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: expected %s got %s\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
				idx+1, tc.name, exp, got, tc.input, refOut, candOut)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-761B-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref761B.bin")
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
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseAnswer(out string) (string, error) {
	ans := strings.TrimSpace(out)
	ans = strings.ToUpper(ans)
	if ans == "YES" || ans == "NO" {
		return ans, nil
	}
	return "", fmt.Errorf("unexpected output %q", out)
}

func buildTests() []testCase {
	tests := []testCase{
		formatCase("simple_yes", 3, 8, []int{2, 4, 6}, []int{1, 5, 7}),
		formatCase("simple_no", 3, 8, []int{0, 2, 4}, []int{1, 3, 5}),
		formatCase("single_barrier", 1, 5, []int{3}, []int{0}),
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomCase(rng, i))
	}
	tests = append(tests, stressCase())
	return tests
}

func formatCase(name string, n, L int, A, B []int) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, L)
	for i, v := range A {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	for i, v := range B {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return testCase{name: name, input: sb.String()}
}

func randomCase(rng *rand.Rand, idx int) testCase {
	n := rng.Intn(50) + 1
	L := rng.Intn(50-n+1) + n
	a := randomSequence(rng, n, L)
	var b []int
	if rng.Intn(2) == 0 {
		shift := rng.Intn(L)
		b = shiftSeq(a, L, shift)
	} else {
		b = randomSequence(rng, n, L)
	}
	return formatCase(fmt.Sprintf("random_%d", idx+1), n, L, a, b)
}

func randomSequence(rng *rand.Rand, n, L int) []int {
	positions := make([]int, L)
	for i := 0; i < L; i++ {
		positions[i] = i
	}
	rng.Shuffle(L, func(i, j int) {
		positions[i], positions[j] = positions[j], positions[i]
	})
	selected := positions[:n]
	sortInts(selected)
	return selected
}

func shiftSeq(seq []int, L, shift int) []int {
	res := make([]int, len(seq))
	for i, v := range seq {
		res[i] = (v + shift) % L
	}
	sortInts(res)
	return res
}

func sortInts(arr []int) {
	// simple insertion sort for n <= 50
	for i := 1; i < len(arr); i++ {
		j := i
		for j > 0 && arr[j-1] > arr[j] {
			arr[j-1], arr[j] = arr[j], arr[j-1]
			j--
		}
	}
}

func stressCase() testCase {
	n := 50
	L := 100
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = i * 2
	}
	b := shiftSeq(a, L, 37)
	return formatCase("stress_max", n, L, a, b)
}
