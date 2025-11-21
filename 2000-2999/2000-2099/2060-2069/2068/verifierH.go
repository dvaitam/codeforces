package main

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const refSource = "2000-2999/2000-2099/2060-2069/2068/2068H.go"

type testCase struct {
	name      string
	n         int
	a, b      int64
	dist      []int64
	expectYes bool
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()
	_ = refBin

	tests := buildTests()

	for idx, tc := range tests {
		input := buildInput(tc)

		candOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, input, candOut)
			os.Exit(1)
		}
		if _, err := parseSolution(tc, candOut); err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, input, candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-2068H-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	src, err := os.ReadFile(refSource)
	if err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to read reference: %v", err)
	}
	fixed := bytes.Replace(src, []byte("pos []target"), []byte("pos [][]target"), 1)
	refPath := filepath.Join(dir, "ref_main.go")
	if err := os.WriteFile(refPath, fixed, 0o644); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to write patched reference: %v", err)
	}

	binPath := filepath.Join(dir, "ref2068H.bin")
	cmd := exec.Command("go", "build", "-o", binPath, refPath)
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

func parseSolution(tc testCase, output string) (bool, error) {
	lines := strings.Fields(output)
	if len(lines) == 0 {
		return false, errors.New("empty output")
	}
	ans := strings.ToUpper(lines[0])
	if ans != "YES" && ans != "NO" {
		return false, fmt.Errorf("first token must be YES or NO, got %q", lines[0])
	}
	if ans == "NO" {
		if tc.expectYes {
			return false, errors.New("reported NO but a solution is expected")
		}
		if len(lines) > 1 {
			return false, errors.New("extra tokens after NO")
		}
		return false, nil
	}

	if !tc.expectYes {
		return false, errors.New("reported YES but no solution is expected")
	}
	expectedTokens := 1 + 2*tc.n
	if len(lines) != expectedTokens {
		return false, fmt.Errorf("expected %d coordinate tokens, got %d", expectedTokens, len(lines))
	}
	coords := make([][2]int64, tc.n)
	for i := 0; i < tc.n; i++ {
		x, err := parseInt64(lines[1+2*i])
		if err != nil {
			return false, fmt.Errorf("invalid x_%d: %v", i+1, err)
		}
		y, err := parseInt64(lines[2+2*i])
		if err != nil {
			return false, fmt.Errorf("invalid y_%d: %v", i+1, err)
		}
		coords[i] = [2]int64{x, y}
	}

	if coords[0][0] != 0 || coords[0][1] != 0 {
		return false, fmt.Errorf("first statue must be at (0,0), got (%d,%d)", coords[0][0], coords[0][1])
	}
	last := coords[tc.n-1]
	if last[0] != tc.a || last[1] != tc.b {
		return false, fmt.Errorf("last statue must be at (%d,%d), got (%d,%d)", tc.a, tc.b, last[0], last[1])
	}
	for i := 0; i < tc.n-1; i++ {
		got := manhattan(coords[i], coords[i+1])
		if got != tc.dist[i] {
			return false, fmt.Errorf("distance between statue %d and %d expected %d, got %d", i+1, i+2, tc.dist[i], got)
		}
	}
	return true, nil
}

func parseInt64(s string) (int64, error) {
	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func manhattan(a, b [2]int64) int64 {
	return abs64(a[0]-b[0]) + abs64(a[1]-b[1])
}

func abs64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", tc.n))
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.a, tc.b))
	for i, v := range tc.dist {
		sb.WriteString(fmt.Sprintf("%d", v))
		if i+1 < len(tc.dist) {
			sb.WriteByte(' ')
		}
	}
	sb.WriteByte('\n')
	return sb.String()
}

func buildTests() []testCase {
	tests := []testCase{
		{
			name:      "sample-no",
			n:         3,
			a:         5,
			b:         8,
			dist:      []int64{9, 0},
			expectYes: false,
		},
		{
			name:      "sample-yes",
			n:         4,
			a:         10,
			b:         6,
			dist:      []int64{7, 8, 5},
			expectYes: true,
		},
		{
			name:      "stay-origin",
			n:         3,
			a:         0,
			b:         0,
			dist:      []int64{0, 0},
			expectYes: true,
		},
		{
			name:      "parity-impossible",
			n:         3,
			a:         1,
			b:         0,
			dist:      []int64{2, 2},
			expectYes: false,
		},
		{
			name:      "long-line",
			n:         5,
			a:         15,
			b:         0,
			dist:      []int64{3, 4, 5, 3},
			expectYes: true,
		},
	}

	// Add random feasible and impossible cases.
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 5; i++ {
		tc := randomFeasible(rng, 5+i)
		tests = append(tests, tc)
	}
	for i := 0; i < 3; i++ {
		tc := randomImpossible(rng, 6+i)
		tests = append(tests, tc)
	}

	return tests
}

func randomFeasible(rng *rand.Rand, n int) testCase {
	if n < 3 {
		n = 3
	}
	coords := make([][2]int64, n)
	dist := make([]int64, n-1)
	for i := 1; i < n; i++ {
		len := int64(rng.Intn(6)) // 0..5
		var dx, dy int64
		dir := rng.Intn(4)
		switch dir {
		case 0:
			dx, dy = len, 0
		case 1:
			dx, dy = -len, 0
		case 2:
			dx, dy = 0, len
		default:
			dx, dy = 0, -len
		}
		coords[i][0] = coords[i-1][0] + dx
		coords[i][1] = coords[i-1][1] + dy
		dist[i-1] = abs64(dx) + abs64(dy)
	}
	return testCase{
		name:      fmt.Sprintf("feasible-%d", rng.Int()),
		n:         n,
		a:         coords[n-1][0],
		b:         coords[n-1][1],
		dist:      dist,
		expectYes: true,
	}
}

func randomImpossible(rng *rand.Rand, n int) testCase {
	if n < 3 {
		n = 3
	}
	dist := make([]int64, n-1)
	var sum int64
	for i := 0; i < n-1; i++ {
		di := int64(rng.Intn(10)) // 0..9
		dist[i] = di
		sum += di
	}
	// Place target too far to guarantee impossibility.
	a := sum + 5
	return testCase{
		name:      fmt.Sprintf("impossible-%d", rng.Int()),
		n:         n,
		a:         a,
		b:         0,
		dist:      dist,
		expectYes: false,
	}
}
