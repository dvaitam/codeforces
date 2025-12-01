package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

const refSource = "./250D.go"

type testCase struct {
	name  string
	input string
	n, m  int
	a, b  int
	yA    []int
	yB    []int
	l     []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
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
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}
		refI, refJ, err := parseOutput(tc, refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}
		refCost, err := evaluateChoice(tc, refI, refJ)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to evaluate reference output on test %d (%s): %v\n", idx+1, tc.name, err)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candI, candJ, err := parseOutput(tc, candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candCost, err := evaluateChoice(tc, candI, candJ)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to evaluate candidate output on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		if !costsClose(refCost, candCost) {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: candidate path cost %.12f differs from optimal %.12f\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
				idx+1, tc.name, candCost, refCost, tc.input, refOut, candOut)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-250D-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref250D.bin")
	cmd := exec.Command("go", "build", "-o", binPath, refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	cleanup := func() {
		_ = os.RemoveAll(dir)
	}
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

func parseOutput(tc testCase, output string) (int, int, error) {
	tokens := strings.Fields(output)
	if len(tokens) < 2 {
		return 0, 0, fmt.Errorf("expected two integers, got: %q", output)
	}
	i, err := strconv.Atoi(tokens[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid left index %q", tokens[0])
	}
	j, err := strconv.Atoi(tokens[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid right index %q", tokens[1])
	}
	if i < 1 || i > tc.n {
		return 0, 0, fmt.Errorf("left index %d out of range [1,%d]", i, tc.n)
	}
	if j < 1 || j > tc.m {
		return 0, 0, fmt.Errorf("right index %d out of range [1,%d]", j, tc.m)
	}
	return i, j, nil
}

func evaluateChoice(tc testCase, i, j int) (float64, error) {
	if i < 1 || i > tc.n || j < 1 || j > tc.m {
		return 0, fmt.Errorf("indices out of range")
	}
	a := float64(tc.a)
	dx := float64(tc.b - tc.a)
	yLeft := float64(tc.yA[i-1])
	yRight := float64(tc.yB[j-1])
	path := float64(tc.l[j-1])
	west := math.Hypot(a, yLeft)
	bridge := math.Hypot(dx, yRight-yLeft)
	return west + bridge + path, nil
}

func costsClose(refCost, candCost float64) bool {
	diff := math.Abs(refCost - candCost)
	limit := 1e-6
	if diff <= limit {
		return true
	}
	scale := math.Abs(refCost)
	if scale < 1 {
		scale = 1
	}
	return diff <= limit*scale
}

func buildTests() []testCase {
	tests := []testCase{
		newTestCase("sample", 3, 5, []int{-2, -1, 4}, []int{-1, 2}, []int{7, 3}),
		newTestCase("single_option", 1, 2, []int{0}, []int{0}, []int{1}),
		newTestCase("balanced", 100, 400, []int{-5, 0, 4}, []int{-3, 8, 9}, []int{2, 3, 4}),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomTest(rng, i))
	}
	return tests
}

func newTestCase(name string, a, b int, left, right, lens []int) testCase {
	if len(right) != len(lens) {
		panic("right points and lengths mismatch")
	}
	tc := testCase{
		name: name,
		n:    len(left),
		m:    len(right),
		a:    a,
		b:    b,
		yA:   append([]int(nil), left...),
		yB:   append([]int(nil), right...),
		l:    append([]int(nil), lens...),
	}
	var sb strings.Builder
	sb.Grow((tc.n+tc.m)*12 + 64)
	sb.WriteString(fmt.Sprintf("%d %d %d %d\n", tc.n, tc.m, a, b))
	appendInts(&sb, left)
	appendInts(&sb, right)
	appendInts(&sb, lens)
	tc.input = sb.String()
	return tc
}

func appendInts(sb *strings.Builder, arr []int) {
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
}

func randomTest(rng *rand.Rand, idx int) testCase {
	n := rng.Intn(50) + 1
	m := rng.Intn(50) + 1
	a := rng.Intn(900000) + 1
	b := a + rng.Intn(100000) + 1
	if b >= 1000000 {
		b = 999999
		if b <= a {
			b = a + 1
		}
	}
	left := randomSortedUnique(rng, n, -1000000, 1000000)
	right := randomSortedUnique(rng, m, -1000000, 1000000)
	lens := make([]int, m)
	for i := 0; i < m; i++ {
		lens[i] = rng.Intn(1000000) + 1
	}
	name := fmt.Sprintf("random_%d", idx+1)
	return newTestCase(name, a, b, left, right, lens)
}

func randomSortedUnique(rng *rand.Rand, size int, lo, hi int) []int {
	vals := make(map[int]struct{}, size)
	for len(vals) < size {
		v := rng.Intn(hi-lo+1) + lo
		vals[v] = struct{}{}
	}
	arr := make([]int, 0, size)
	for v := range vals {
		arr = append(arr, v)
	}
	sort.Ints(arr)
	return arr
}
