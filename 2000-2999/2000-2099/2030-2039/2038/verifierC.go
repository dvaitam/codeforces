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

type caseData struct {
	values []int64
}

type testCase struct {
	input string
	cases []caseData
}

type parsedAnswer struct {
	possible bool
	coords   [8]int64
}

type point struct {
	x int64
	y int64
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	refBin, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s\n", idx+1, err, tc.input)
			os.Exit(1)
		}
		refAns, err := parseAnswers(refOut, len(tc.cases))
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d: %v\noutput:\n%s\n", idx+1, err, refOut)
			os.Exit(1)
		}

		refAreas := make([]int64, len(tc.cases))
		for i, ans := range refAns {
			if ans.possible {
				area, err := rectangleArea(ans.coords)
				if err != nil {
					fmt.Fprintf(os.Stderr, "invalid rectangle in reference on test %d case %d: %v\n", idx+1, i+1, err)
					os.Exit(1)
				}
				refAreas[i] = area
			} else {
				refAreas[i] = -1
			}
		}

		gotOut, err := runProgram(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%s\nstdout/stderr:\n%s\n", idx+1, err, tc.input, gotOut)
			os.Exit(1)
		}
		gotAns, err := parseAnswers(gotOut, len(tc.cases))
		if err != nil {
			fmt.Fprintf(os.Stderr, "participant output invalid on test %d: %v\noutput:\n%s\n", idx+1, err, gotOut)
			os.Exit(1)
		}

		for caseIdx := range tc.cases {
			if err := validateCase(tc.cases[caseIdx], refAns[caseIdx], gotAns[caseIdx], refAreas[caseIdx]); err != nil {
				fmt.Fprintf(os.Stderr, "test %d case %d failed: %v\ninput:\n%sreference output:\n%s\nparticipant output:\n%s\n", idx+1, caseIdx+1, err, tc.input, refOut, gotOut)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReferenceBinary() (string, error) {
	dir, err := verifierDir()
	if err != nil {
		return "", err
	}
	tmp, err := os.CreateTemp("", "2038C_ref_*.bin")
	if err != nil {
		return "", err
	}
	path := tmp.Name()
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", path, "2038C.go")
	cmd.Dir = dir
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(path)
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return path, nil
}

func verifierDir() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to determine verifier directory")
	}
	return filepath.Dir(file), nil
}

func runProgram(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return stdout.String() + stderr.String(), fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func parseAnswers(out string, caseCnt int) ([]parsedAnswer, error) {
	tokens := strings.Fields(out)
	idx := 0
	ans := make([]parsedAnswer, caseCnt)
	for i := 0; i < caseCnt; i++ {
		if idx >= len(tokens) {
			return nil, fmt.Errorf("case %d: expected YES/NO but output ended early", i+1)
		}
		token := strings.ToUpper(tokens[idx])
		idx++
		switch token {
		case "NO":
			ans[i] = parsedAnswer{possible: false}
		case "YES":
			if idx+8 > len(tokens) {
				return nil, fmt.Errorf("case %d: expected 8 coordinates after YES", i+1)
			}
			var coords [8]int64
			for j := 0; j < 8; j++ {
				val, err := strconv.ParseInt(tokens[idx+j], 10, 64)
				if err != nil {
					return nil, fmt.Errorf("case %d: invalid coordinate %q", i+1, tokens[idx+j])
				}
				coords[j] = val
			}
			idx += 8
			ans[i] = parsedAnswer{possible: true, coords: coords}
		default:
			return nil, fmt.Errorf("case %d: expected YES/NO got %s", i+1, token)
		}
	}
	if idx != len(tokens) {
		return nil, fmt.Errorf("output has %d extra tokens", len(tokens)-idx)
	}
	return ans, nil
}

func rectangleArea(coords [8]int64) (int64, error) {
	xs := []int64{coords[0], coords[2], coords[4], coords[6]}
	ys := []int64{coords[1], coords[3], coords[5], coords[7]}

	xMin, xMax := xs[0], xs[0]
	yMin, yMax := ys[0], ys[0]
	for i := 1; i < 4; i++ {
		if xs[i] < xMin {
			xMin = xs[i]
		}
		if xs[i] > xMax {
			xMax = xs[i]
		}
		if ys[i] < yMin {
			yMin = ys[i]
		}
		if ys[i] > yMax {
			yMax = ys[i]
		}
	}

	for i := 0; i < 4; i++ {
		if xs[i] != xMin && xs[i] != xMax {
			return 0, fmt.Errorf("point %d x=%d not aligned with rectangle", i+1, xs[i])
		}
		if ys[i] != yMin && ys[i] != yMax {
			return 0, fmt.Errorf("point %d y=%d not aligned with rectangle", i+1, ys[i])
		}
	}

	width := absInt64(xMax - xMin)
	height := absInt64(yMax - yMin)

	required := make(map[point]int)
	switch {
	case width > 0 && height > 0:
		required[point{xMin, yMin}] = 1
		required[point{xMin, yMax}] = 1
		required[point{xMax, yMin}] = 1
		required[point{xMax, yMax}] = 1
	case width == 0 && height > 0:
		required[point{xMin, yMin}] = 2
		required[point{xMin, yMax}] = 2
	case width > 0 && height == 0:
		required[point{xMin, yMin}] = 2
		required[point{xMax, yMin}] = 2
	default:
		required[point{xMin, yMin}] = 4
	}

	actual := make(map[point]int)
	for i := 0; i < 4; i++ {
		actual[point{xs[i], ys[i]}]++
	}

	if len(actual) != len(required) {
		return 0, fmt.Errorf("expected %d distinct corners got %d", len(required), len(actual))
	}
	for p, need := range required {
		if actual[p] != need {
			return 0, fmt.Errorf("corner (%d,%d) count %d expected %d", p.x, p.y, actual[p], need)
		}
	}

	return width * height, nil
}

func validateCase(cd caseData, refAns, gotAns parsedAnswer, refArea int64) error {
	if !refAns.possible {
		if gotAns.possible {
			return fmt.Errorf("expected NO but participant output YES")
		}
		return nil
	}
	if !gotAns.possible {
		return fmt.Errorf("expected YES but participant output NO")
	}
	freq := make(map[int64]int, len(cd.values))
	for _, v := range cd.values {
		freq[v]++
	}
	for _, val := range gotAns.coords {
		freq[val]--
		if freq[val] < 0 {
			return fmt.Errorf("coordinate value %d used more times than available", val)
		}
	}
	area, err := rectangleArea(gotAns.coords)
	if err != nil {
		return fmt.Errorf("invalid rectangle: %v", err)
	}
	if area != refArea {
		return fmt.Errorf("area mismatch: expected %d got %d", refArea, area)
	}
	return nil
}

func generateTests() []testCase {
	var tests []testCase
	tests = append(tests, manualTests()...)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests = append(tests, randomTests(rng, 60, 20)...)
	tests = append(tests, randomTests(rng, 40, 200)...)
	tests = append(tests, randomTests(rng, 30, 2000)...)
	tests = append(tests, stressTests()...)
	return tests
}

func manualTests() []testCase {
	cases := []caseData{
		{values: []int64{-5, 1, 1, 2, 2, 3, 3, 4, 4, 5, 5, 6, 6, 7, 7, 10}},
		{values: []int64{80, 0, -1, 2, 2, 1, 1, 3}},
		{values: []int64{80, 0, 0, 0, 0, 5, 0, 5}},
		{values: []int64{5, 5, 5, 5, 5, 5, 5, 5}},
	}
	return []testCase{makeTestCase(cases)}
}

func randomTests(rng *rand.Rand, batches int, maxN int) []testCase {
	const limit = 200000
	var tests []testCase
	for b := 0; b < batches; b++ {
		caseCnt := rng.Intn(3) + 1
		var cases []caseData
		sumN := 0
		for len(cases) < caseCnt {
			n := rng.Intn(maxN-7) + 8
			if sumN+n > limit {
				break
			}
			cases = append(cases, randomCase(rng, n))
			sumN += n
		}
		if len(cases) == 0 {
			cases = append(cases, randomCase(rng, 8))
		}
		tests = append(tests, makeTestCase(cases))
	}
	return tests
}

func stressTests() []testCase {
	large := make([]int64, 200000)
	for i := 0; i < len(large); i++ {
		large[i] = int64(i%1000 - 500)
	}
	striped := make([]int64, 50000)
	for i := 0; i < len(striped); i++ {
		if i%4 < 2 {
			striped[i] = 1
		} else {
			striped[i] = -1
		}
	}
	detRng := rand.New(rand.NewSource(42))
	return []testCase{
		makeTestCase([]caseData{{values: large}}),
		makeTestCase([]caseData{
			randomCase(detRng, 120000),
			{values: striped},
		}),
	}
}

func randomCase(rng *rand.Rand, n int) caseData {
	size := rng.Intn(6) + 2
	pool := make([]int64, size)
	for i := 0; i < size; i++ {
		pool[i] = int64(rng.Intn(2001) - 1000)
	}
	values := make([]int64, n)
	for i := 0; i < n; i++ {
		if rng.Intn(5) == 0 {
			values[i] = int64(rng.Intn(2_000_001) - 1_000_000)
		} else {
			values[i] = pool[rng.Intn(size)]
		}
	}
	return caseData{values: values}
}

func makeTestCase(cases []caseData) testCase {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(cases)))
	sb.WriteByte('\n')
	stored := make([]caseData, len(cases))
	for i, c := range cases {
		sb.WriteString(strconv.Itoa(len(c.values)))
		sb.WriteByte('\n')
		for j, val := range c.values {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(val, 10))
		}
		sb.WriteByte('\n')
		copied := make([]int64, len(c.values))
		copy(copied, c.values)
		stored[i] = caseData{values: copied}
	}
	return testCase{input: sb.String(), cases: stored}
}

func absInt64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
