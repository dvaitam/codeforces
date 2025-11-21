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

type rect struct {
	x1, y1, x2, y2 int
}

type testCase struct {
	input string
	rects []rect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	ref, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := generateTests()
	for idx, tc := range tests {
		refOut, err := runProgram(ref, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s\n", idx+1, err, tc.input)
			os.Exit(1)
		}
		refSol, err := parseSolution(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d: %v\noutput:\n%s\n", idx+1, err, refOut)
			os.Exit(1)
		}

		gotOut, err := runProgram(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%s\nstdout/stderr:\n%s\n", idx+1, err, tc.input, gotOut)
			os.Exit(1)
		}
		gotSol, err := parseSolution(gotOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "participant output invalid on test %d: %v\noutput:\n%s\n", idx+1, err, gotOut)
			os.Exit(1)
		}

		if !gotSol.exists {
			if refSol.exists {
				fmt.Fprintf(os.Stderr, "test %d: participant printed NO but square exists\ninput:\n%sreference output:\n%s\nparticipant output:\n%s\n",
					idx+1, tc.input, refOut, gotOut)
				os.Exit(1)
			}
			continue
		}

		if err := validateSubset(tc.rects, gotSol.subset); err != nil {
			fmt.Fprintf(os.Stderr, "test %d: invalid subset: %v\ninput:\n%sreference output:\n%s\nparticipant output:\n%s\n",
				idx+1, err, tc.input, refOut, gotOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

type solution struct {
	exists bool
	subset []int
}

func parseSolution(out string) (solution, error) {
	tokens := strings.Fields(out)
	if len(tokens) == 0 {
		return solution{}, fmt.Errorf("empty output")
	}
	first := strings.ToUpper(tokens[0])
	if first == "NO" {
		if len(tokens) != 1 {
			return solution{}, fmt.Errorf("extra tokens after NO")
		}
		return solution{exists: false}, nil
	}
	if first != "YES" {
		return solution{}, fmt.Errorf("expected YES/NO, got %s", tokens[0])
	}
	if len(tokens) < 2 {
		return solution{}, fmt.Errorf("missing k after YES")
	}
	k, err := strconv.Atoi(tokens[1])
	if err != nil || k <= 0 {
		return solution{}, fmt.Errorf("invalid k value")
	}
	if len(tokens) != 2+k {
		return solution{}, fmt.Errorf("expected %d indices, got %d", k, len(tokens)-2)
	}
	subset := make([]int, k)
	for i := 0; i < k; i++ {
		v, err := strconv.Atoi(tokens[2+i])
		if err != nil {
			return solution{}, fmt.Errorf("invalid index %q", tokens[2+i])
		}
		subset[i] = v
	}
	return solution{exists: true, subset: subset}, nil
}

func validateSubset(rects []rect, subset []int) error {
	if len(subset) == 0 {
		return fmt.Errorf("subset must be non-empty")
	}
	n := len(rects)
	used := make([]bool, n)
	L := 1 << 30
	R := -1 << 30
	B := 1 << 30
	T := -1 << 30
	areaSum := 0
	for _, id := range subset {
		if id < 1 || id > n {
			return fmt.Errorf("index %d out of range", id)
		}
		if used[id-1] {
			return fmt.Errorf("duplicate index %d", id)
		}
		used[id-1] = true
		rect := rects[id-1]
		if rect.x1 < L {
			L = rect.x1
		}
		if rect.x2 > R {
			R = rect.x2
		}
		if rect.y1 < B {
			B = rect.y1
		}
		if rect.y2 > T {
			T = rect.y2
		}
		width := rect.x2 - rect.x1
		height := rect.y2 - rect.y1
		if width <= 0 || height <= 0 {
			return fmt.Errorf("invalid rectangle dimensions")
		}
		areaSum += width * height
	}
	sideX := R - L
	sideY := T - B
	if sideX <= 0 || sideX != sideY {
		return fmt.Errorf("bounding box is not a square")
	}
	if areaSum != sideX*sideX {
		return fmt.Errorf("area mismatch: subset area %d vs square area %d", areaSum, sideX*sideX)
	}
	for _, id := range subset {
		rect := rects[id-1]
		if rect.x1 < L || rect.x2 > R || rect.y1 < B || rect.y2 > T {
			return fmt.Errorf("rectangle %d extends outside square", id)
		}
	}
	return nil
}

func buildReferenceBinary() (string, error) {
	dir, err := verifierDir()
	if err != nil {
		return "", err
	}
	tmp, err := os.CreateTemp("", "335D_ref_*.bin")
	if err != nil {
		return "", err
	}
	path := tmp.Name()
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", path, "335D.go")
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

func generateTests() []testCase {
	var tests []testCase
	tests = append(tests, manualTests()...)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 40; i++ {
		tests = append(tests, randomTest(rng, 20, 40))
	}
	for i := 0; i < 10; i++ {
		tests = append(tests, randomTest(rng, 100, 3000))
	}
	return tests
}

func manualTests() []testCase {
	// Sample 1
	rects1 := []rect{
		{0, 0, 1, 9},
		{1, 0, 9, 1},
		{1, 8, 9, 9},
		{8, 1, 9, 8},
		{2, 2, 3, 6},
		{3, 2, 7, 3},
		{2, 6, 7, 7},
		{5, 3, 7, 6},
		{3, 3, 5, 6},
	}
	// Sample 2
	rects2 := []rect{
		{0, 0, 1, 9},
		{1, 0, 9, 1},
		{1, 8, 9, 9},
		{8, 1, 9, 8},
	}
	return []testCase{
		makeTestCase(rects1),
		makeTestCase(rects2),
	}
}

func randomTest(rng *rand.Rand, maxSegments, maxCoord int) testCase {
	segmentsX := rng.Intn(maxSegments-1) + 1
	segmentsY := rng.Intn(maxSegments-1) + 1
	xCoords := makeCoords(rng, segmentsX+1, maxCoord)
	yCoords := makeCoords(rng, segmentsY+1, maxCoord)
	rects := make([]rect, 0)
	for i := 0; i < len(xCoords)-1; i++ {
		for j := 0; j < len(yCoords)-1; j++ {
			if rng.Intn(2) == 0 {
				continue
			}
			rects = append(rects, rect{
				x1: xCoords[i],
				y1: yCoords[j],
				x2: xCoords[i+1],
				y2: yCoords[j+1],
			})
		}
	}
	if len(rects) == 0 {
		rects = append(rects, rect{
			x1: xCoords[0],
			y1: yCoords[0],
			x2: xCoords[1],
			y2: yCoords[1],
		})
	}
	return makeTestCase(rects)
}

func makeCoords(rng *rand.Rand, count, maxCoord int) []int {
	coords := make([]int, count)
	coords[0] = 0
	for i := 1; i < count; i++ {
		remaining := maxCoord - coords[i-1] - (count - i - 1)
		if remaining <= 0 {
			coords[i] = coords[i-1] + 1
		} else {
			coords[i] = coords[i-1] + rng.Intn(remaining) + 1
		}
	}
	offset := 0
	if coords[count-1] < maxCoord {
		offset = rng.Intn(maxCoord - coords[count-1] + 1)
	}
	for i := 0; i < count; i++ {
		coords[i] += offset
	}
	return coords
}

func makeTestCase(rects []rect) testCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(rects)))
	for _, r := range rects {
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", r.x1, r.y1, r.x2, r.y2))
	}
	return testCase{input: sb.String(), rects: rects}
}
