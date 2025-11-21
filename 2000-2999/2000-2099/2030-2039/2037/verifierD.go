package main

import (
	"bytes"
	"container/heap"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

type hurdle struct {
	l int64
	r int64
}

type powerUp struct {
	x int64
	v int64
}

type caseData struct {
	L       int64
	hurdles []hurdle
	powers  []powerUp
}

type testFile struct {
	name  string
	input string
	cases []caseData
}

type maxHeap []int64

func (h maxHeap) Len() int            { return len(h) }
func (h maxHeap) Less(i, j int) bool  { return h[i] > h[j] }
func (h maxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *maxHeap) Push(x interface{}) { *h = append(*h, x.(int64)) }
func (h *maxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	val := old[n-1]
	*h = old[:n-1]
	return val
}

func main() {
	candidate, err := candidatePath()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	refSrc := referencePath()
	refBin, cleanup, err := buildReferenceBinary(refSrc)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer cleanup()

	tests := buildTests()
	for idx, tf := range tests {
		refOut, err := runProgram(refBin, tf.input)
		if err != nil {
			fmt.Printf("reference failed on file %d (%s): %v\n", idx+1, tf.name, err)
			os.Exit(1)
		}
		refAns, err := parseAnswers(refOut, len(tf.cases))
		if err != nil {
			fmt.Printf("reference output invalid on %s: %v\nraw output:\n%s\n", tf.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, tf.input)
		if err != nil {
			fmt.Printf("case %d (%s): runtime error: %v\n", idx+1, tf.name, err)
			fmt.Println(previewInput(tf.input))
			os.Exit(1)
		}
		candAns, err := parseAnswers(candOut, len(tf.cases))
		if err != nil {
			fmt.Printf("case %d (%s): invalid output format: %v\nraw output:\n%s\n", idx+1, tf.name, err, candOut)
			fmt.Println(previewInput(tf.input))
			os.Exit(1)
		}

		for caseIdx, cs := range tf.cases {
			expected := solveCase(cs)
			if refAns[caseIdx] != expected {
				fmt.Printf("reference mismatch on %s case %d: expected %d got %d\n", tf.name, caseIdx+1, expected, refAns[caseIdx])
				fmt.Println(previewInput(tf.input))
				os.Exit(1)
			}
			if candAns[caseIdx] != expected {
				fmt.Printf("file %d (%s) failed case %d: expected %d got %d\n", idx+1, tf.name, caseIdx+1, expected, candAns[caseIdx])
				fmt.Println(previewInput(tf.input))
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d test files passed\n", len(tests))
}

func candidatePath() (string, error) {
	if len(os.Args) == 2 {
		return os.Args[1], nil
	}
	if len(os.Args) == 3 && os.Args[1] == "--" {
		return os.Args[2], nil
	}
	return "", fmt.Errorf("usage: go run verifierD.go /path/to/solution")
}

func referencePath() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "2000-2999/2000-2099/2030-2039/2037/2037D.go"
	}
	return filepath.Join(filepath.Dir(file), "2037D.go")
}

func buildReferenceBinary(src string) (string, func(), error) {
	tmpDir, err := os.MkdirTemp("", "verifier2037D")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	bin := filepath.Join(tmpDir, "ref2037d")
	cmd := exec.Command("go", "build", "-o", bin, src)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build reference solution: %v\n%s", err, out.String())
	}
	cleanup := func() { _ = os.RemoveAll(tmpDir) }
	return bin, cleanup, nil
}

func runProgram(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("%v\n%s", err, out.String())
	}
	return out.String(), nil
}

func parseAnswers(out string, count int) ([]int64, error) {
	tokens := strings.Fields(out)
	if len(tokens) != count {
		return nil, fmt.Errorf("expected %d answers, got %d", count, len(tokens))
	}
	ans := make([]int64, count)
	for i, tok := range tokens {
		val, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", tok)
		}
		ans[i] = val
	}
	return ans, nil
}

func solveCase(cs caseData) int64 {
	if len(cs.hurdles) == 0 {
		return 0
	}
	powers := cs.powers
	idx := 0
	curPower := int64(1)
	used := int64(0)
	h := &maxHeap{}
	heap.Init(h)
	for _, hurdle := range cs.hurdles {
		for idx < len(powers) && powers[idx].x < hurdle.l {
			heap.Push(h, powers[idx].v)
			idx++
		}
		need := hurdle.r - hurdle.l + 2
		for curPower < need {
			if h.Len() == 0 {
				return -1
			}
			curPower += heap.Pop(h).(int64)
			used++
		}
	}
	return used
}

func buildTests() []testFile {
	tests := []testFile{}
	tests = append(tests, manualTests())
	rng := rand.New(rand.NewSource(0x5a11c7e3))
	tests = append(tests, randomTest("random-small", rng, 35, 200, 5, 8))
	tests = append(tests, randomTest("random-medium", rng, 30, 2000, 20, 25))
	tests = append(tests, randomTest("random-large", rng, 20, 200000, 60, 80))
	tests = append(tests, stressTest())
	return tests
}

func manualTests() testFile {
	cases := []caseData{}
	// Case 1: no hurdles
	cases = append(cases, caseData{
		L:       20,
		hurdles: []hurdle{},
		powers:  []powerUp{{x: 1, v: 5}, {x: 10, v: 3}},
	})
	// Case 2: impossible hurdle
	cases = append(cases, caseData{
		L:       30,
		hurdles: []hurdle{{l: 5, r: 12}},
		powers:  []powerUp{{x: 1, v: 1}, {x: 3, v: 1}},
	})
	// Case 3: needs selective pickups
	cases = append(cases, caseData{
		L:       70,
		hurdles: []hurdle{{7, 9}, {20, 25}, {40, 44}},
		powers: []powerUp{
			{x: 1, v: 1},
			{x: 4, v: 2},
			{x: 6, v: 3},
			{x: 15, v: 4},
			{x: 18, v: 2},
			{x: 21, v: 6},
			{x: 30, v: 5},
			{x: 35, v: 7},
		},
	})
	// Case 4: multiple power-ups at same position
	cases = append(cases, caseData{
		L:       45,
		hurdles: []hurdle{{10, 15}, {25, 30}},
		powers: []powerUp{
			{x: 1, v: 1},
			{x: 1, v: 4},
			{x: 5, v: 2},
			{x: 5, v: 6},
			{x: 16, v: 3},
			{x: 23, v: 5},
		},
	})
	return makeTestFile("manual", cases)
}

func randomTest(name string, rng *rand.Rand, caseCount int, maxL int64, maxH, maxP int) testFile {
	cases := make([]caseData, caseCount)
	for i := 0; i < caseCount; i++ {
		cases[i] = randomCase(rng, maxL, maxH, maxP)
	}
	return makeTestFile(name, cases)
}

func stressTest() testFile {
	rng := rand.New(rand.NewSource(0x91b3d4f7))
	cases := make([]caseData, 10)
	for i := 0; i < len(cases); i++ {
		cases[i] = randomCase(rng, 1_000_000_000, 120, 180)
	}
	return makeTestFile("stress", cases)
}

func randomCase(rng *rand.Rand, maxL int64, maxH, maxP int) caseData {
	if maxL < 10 {
		maxL = 10
	}
	L := int64(3 + rng.Int63n(maxL-2))
	hurdles := generateHurdles(rng, L, maxH)
	powers := generatePowers(rng, L, maxP, hurdles)
	return caseData{L: L, hurdles: hurdles, powers: powers}
}

func generateHurdles(rng *rand.Rand, L int64, maxH int) []hurdle {
	target := 0
	if maxH > 0 {
		target = rng.Intn(maxH + 1)
	}
	hurdles := make([]hurdle, 0, target)
	cur := int64(2)
	for len(hurdles) < target && cur < L-1 {
		gap := int64(rng.Intn(3))
		cur += gap
		if cur >= L-1 {
			break
		}
		remaining := (L - 1) - cur
		if remaining <= 0 {
			break
		}
		maxLen := int(minInt64(remaining, 20))
		if maxLen < 1 {
			break
		}
		length := int64(1 + rng.Intn(maxLen))
		r := cur + length - 1
		if r > L-1 {
			r = L - 1
		}
		hurdles = append(hurdles, hurdle{l: cur, r: r})
		cur = r + 2
	}
	return hurdles
}

func generatePowers(rng *rand.Rand, L int64, maxP int, hurdles []hurdle) []powerUp {
	m := 0
	if maxP > 0 {
		m = rng.Intn(maxP + 1)
	}
	segs := safeSegments(L, hurdles)
	if len(segs) == 0 {
		segs = append(segs, [2]int64{1, L})
	}
	powers := make([]powerUp, m)
	for i := 0; i < m; i++ {
		seg := segs[rng.Intn(len(segs))]
		length := seg[1] - seg[0] + 1
		pos := seg[0]
		if length > 1 {
			pos += int64(rng.Int63n(length))
		}
		val := int64(rng.Int63n(L) + 1)
		powers[i] = powerUp{x: pos, v: val}
	}
	sort.Slice(powers, func(i, j int) bool {
		if powers[i].x == powers[j].x {
			return powers[i].v < powers[j].v
		}
		return powers[i].x < powers[j].x
	})
	return powers
}

func safeSegments(L int64, hurdles []hurdle) [][2]int64 {
	segs := make([][2]int64, 0)
	prev := int64(1)
	for _, h := range hurdles {
		end := h.l - 1
		if prev <= end {
			segs = append(segs, [2]int64{prev, end})
		}
		prev = h.r + 1
	}
	if prev <= L {
		segs = append(segs, [2]int64{prev, L})
	}
	return segs
}

func makeTestFile(name string, cases []caseData) testFile {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, cs := range cases {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", len(cs.hurdles), len(cs.powers), cs.L))
		for _, h := range cs.hurdles {
			sb.WriteString(fmt.Sprintf("%d %d\n", h.l, h.r))
		}
		for _, p := range cs.powers {
			sb.WriteString(fmt.Sprintf("%d %d\n", p.x, p.v))
		}
	}
	return testFile{name: name, input: sb.String(), cases: cases}
}

func minInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func previewInput(in string) string {
	const limit = 400
	if len(in) <= limit {
		return "Input:\n" + in
	}
	return fmt.Sprintf("Input (first %d chars):\n%s...\n", limit, in[:limit])
}
