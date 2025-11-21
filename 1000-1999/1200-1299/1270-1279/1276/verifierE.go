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
	refSource = "1000-1999/1200-1299/1270-1279/1276/1276E.go"
	maxMoves  = 1000
	maxCoord  = 1_000_000_000_000_000_000
	stonesCnt = 4
)

type operation struct {
	x int64
	y int64
}

type testCase struct {
	name  string
	input string
}

type instance struct {
	start  [stonesCnt]int64
	target [stonesCnt]int64
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
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
		inst, err := parseInput(tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "internal error parsing test %d (%s): %v\n", idx+1, tc.name, err)
			os.Exit(1)
		}

		refOutRaw, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, refOutRaw)
			os.Exit(1)
		}
		refAns, err := parseOutput(refOutRaw)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, refOutRaw)
			os.Exit(1)
		}
		if refAns.possible {
			if err := verifySequence(inst, refAns.moves); err != nil {
				fmt.Fprintf(os.Stderr, "reference sequence invalid on test %d (%s): %v\ninput:\n%soutput:\n%s",
					idx+1, tc.name, err, tc.input, refOutRaw)
				os.Exit(1)
			}
		}

		candOutRaw, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, candOutRaw)
			os.Exit(1)
		}
		candAns, err := parseOutput(candOutRaw)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate produced invalid output on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, candOutRaw)
			os.Exit(1)
		}

		if refAns.possible {
			if !candAns.possible {
				fmt.Fprintf(os.Stderr, "candidate incorrectly reported impossible on test %d (%s)\ninput:\n%sreference output:\n%s\ncandidate output:\n%s\n",
					idx+1, tc.name, tc.input, refOutRaw, candOutRaw)
				os.Exit(1)
			}
			if err := verifySequence(inst, candAns.moves); err != nil {
				fmt.Fprintf(os.Stderr, "candidate sequence invalid on test %d (%s): %v\ninput:\n%soutput:\n%s",
					idx+1, tc.name, err, tc.input, candOutRaw)
				os.Exit(1)
			}
		} else {
			if candAns.possible {
				fmt.Fprintf(os.Stderr, "candidate found a sequence but reference reported impossible on test %d (%s)\ninput:\n%sreference output:\n%s\ncandidate output:\n%s\n",
					idx+1, tc.name, tc.input, refOutRaw, candOutRaw)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-1276E-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref1276E.bin")
	cmd := exec.Command("go", "build", "-o", binPath, refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		_ = os.RemoveAll(dir)
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

type parsedOutput struct {
	possible bool
	moves    []operation
}

func parseOutput(raw string) (*parsedOutput, error) {
	fields := strings.Fields(raw)
	if len(fields) == 0 {
		return nil, fmt.Errorf("empty output")
	}
	if len(fields) == 1 && fields[0] == "-1" {
		return &parsedOutput{possible: false}, nil
	}
	k, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, fmt.Errorf("invalid move count %q", fields[0])
	}
	if k < 0 || k > maxMoves {
		return nil, fmt.Errorf("move count %d outside [0,%d]", k, maxMoves)
	}
	expected := 1 + 2*k
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d integers, got %d", expected, len(fields))
	}
	moves := make([]operation, k)
	for i := 0; i < k; i++ {
		xVal, err := strconv.ParseInt(fields[1+2*i], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid x coordinate %q", fields[1+2*i])
		}
		yVal, err := strconv.ParseInt(fields[1+2*i+1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid y coordinate %q", fields[1+2*i+1])
		}
		moves[i] = operation{x: xVal, y: yVal}
	}
	return &parsedOutput{possible: true, moves: moves}, nil
}

func parseInput(data string) (*instance, error) {
	reader := strings.NewReader(data)
	inst := &instance{}
	for i := 0; i < stonesCnt; i++ {
		if _, err := fmt.Fscan(reader, &inst.start[i]); err != nil {
			return nil, fmt.Errorf("failed to read a[%d]: %v", i, err)
		}
	}
	for i := 0; i < stonesCnt; i++ {
		if _, err := fmt.Fscan(reader, &inst.target[i]); err != nil {
			return nil, fmt.Errorf("failed to read b[%d]: %v", i, err)
		}
	}
	return inst, nil
}

func verifySequence(inst *instance, moves []operation) error {
	counts := make(map[int64]int)
	for _, v := range inst.start {
		counts[v]++
	}
	total := 0
	for _, c := range counts {
		total += c
	}
	if total != stonesCnt {
		return fmt.Errorf("invalid initial state")
	}
	for idx, mv := range moves {
		if mv.x == mv.y {
			return fmt.Errorf("move %d: x equals y (%d)", idx+1, mv.x)
		}
		cx := counts[mv.x]
		if cx <= 0 {
			return fmt.Errorf("move %d: no stone at x=%d", idx+1, mv.x)
		}
		cy := counts[mv.y]
		if cy <= 0 {
			return fmt.Errorf("move %d: no stone at y=%d", idx+1, mv.y)
		}
		counts[mv.x]--
		if counts[mv.x] == 0 {
			delete(counts, mv.x)
		}
		z := 2*mv.y - mv.x
		if absInt64(z) > maxCoord {
			return fmt.Errorf("move %d: resulting coordinate %d exceeds limit", idx+1, z)
		}
		counts[z]++
	}
	targetCounts := make(map[int64]int)
	for _, v := range inst.target {
		targetCounts[v]++
	}
	if len(counts) != len(targetCounts) {
		return fmt.Errorf("final multiset mismatch: have %v expected %v", counts, targetCounts)
	}
	for coord, cnt := range counts {
		if targetCounts[coord] != cnt {
			return fmt.Errorf("final multiset mismatch: coordinate %d has %d stones, expected %d", coord, cnt, targetCounts[coord])
		}
	}
	return nil
}

func buildTests() []testCase {
	tests := []testCase{
		newTestCase("identical", [4]int64{0, 0, 0, 0}, [4]int64{0, 0, 0, 0}),
		newTestCase("already_target", [4]int64{-5, -2, 7, 9}, [4]int64{-5, -2, 7, 9}),
		newTestCase("swap_pairs", [4]int64{0, 1, 2, 3}, [4]int64{3, 2, 1, 0}),
		newTestCase("duplicates", [4]int64{-3, -3, 5, 5}, [4]int64{5, 5, -3, -3}),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 50; i++ {
		tests = append(tests, randomTest(rng, i, 1_000))
	}
	for i := 0; i < 50; i++ {
		tests = append(tests, randomTest(rng, i+50, 1_000_000))
	}
	for i := 0; i < 40; i++ {
		tests = append(tests, randomTest(rng, i+100, 1_000_000_000))
	}
	tests = append(tests,
		repeatedValueTest("repeated_start"),
		repeatedValueTest("repeated_target"),
		randomFarTest(rng, "wide_range"),
	)
	return tests
}

func newTestCase(name string, a, b [4]int64) testCase {
	var sb strings.Builder
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	for i, v := range b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return testCase{name: name, input: sb.String()}
}

func randomTest(rng *rand.Rand, idx int, bound int64) testCase {
	var a, b [4]int64
	for i := 0; i < stonesCnt; i++ {
		a[i] = randInRange(rng, -bound, bound)
		b[i] = randInRange(rng, -bound, bound)
	}
	return newTestCase(fmt.Sprintf("random_%d_bound_%d", idx+1, bound), a, b)
}

func repeatedValueTest(name string) testCase {
	var a, b [4]int64
	for i := 0; i < 4; i++ {
		a[i] = int64(i / 2)
		b[i] = int64(i%2) * 5
	}
	return newTestCase(name, a, b)
}

func randomFarTest(rng *rand.Rand, name string) testCase {
	var a, b [4]int64
	for i := 0; i < stonesCnt; i++ {
		a[i] = randInRange(rng, -maxCoord/1_000_000, maxCoord/1_000_000)
		offset := randInRange(rng, -maxCoord/1_000, maxCoord/1_000)
		b[i] = a[i] + offset
	}
	return newTestCase(name, a, b)
}

func randInRange(rng *rand.Rand, lo, hi int64) int64 {
	if lo > hi {
		lo, hi = hi, lo
	}
	if lo == hi {
		return lo
	}
	return lo + rng.Int63n(hi-lo+1)
}

func absInt64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
