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

const refSource = "2111B.go"
const inf = 1 << 30

var fib = []int{0, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89}

type testCase struct {
	name string
	n    int
	box  [][3]int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/candidate_binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	tests := buildTests()
	input := buildInput(tests)

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\ninput:\n%soutput:\n%s", err, input, refOut)
		os.Exit(1)
	}
	if err := validateOutputs(tests, refOut); err != nil {
		fmt.Fprintf(os.Stderr, "reference produced invalid output: %v\ninput:\n%soutput:\n%s", err, input, refOut)
		os.Exit(1)
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\ninput:\n%soutput:\n%s", err, input, candOut)
		os.Exit(1)
	}
	if err := validateOutputs(tests, candOut); err != nil {
		fmt.Fprintf(os.Stderr, "candidate output invalid: %v\ninput:\n%soutput:\n%s", err, input, candOut)
		os.Exit(1)
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier directory")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2111B-")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "oracleB")
	cmd := exec.Command("go", "build", "-o", binPath, refSource)
	cmd.Dir = dir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("build failed: %v\n%s", err, stderr.String())
	}
	cleanup := func() {
		_ = os.RemoveAll(tmpDir)
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

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.Grow(len(tests) * 64)
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, len(tc.box)))
		for _, b := range tc.box {
			sb.WriteString(fmt.Sprintf("%d %d %d\n", b[0], b[1], b[2]))
		}
	}
	return sb.String()
}

func validateOutputs(tests []testCase, output string) error {
	lines := splitNonEmptyLines(output)
	if len(lines) != len(tests) {
		return fmt.Errorf("expected %d lines, got %d", len(tests), len(lines))
	}
	for i, line := range lines {
		res := strings.Join(strings.Fields(line), "")
		m := len(tests[i].box)
		if len(res) != m {
			return fmt.Errorf("case %d: expected string length %d, got %d", i+1, m, len(res))
		}
		fits := preprocess(tests[i].n)
		for j, ch := range res {
			if ch != '0' && ch != '1' {
				return fmt.Errorf("case %d: invalid character %q", i+1, ch)
			}
			expect := canFit(fits, tests[i].box[j])
			if expect && ch == '0' {
				return fmt.Errorf("case %d: box %d marked 0 but cubes fit", i+1, j+1)
			}
			if !expect && ch == '1' {
				return fmt.Errorf("case %d: box %d marked 1 but cubes do not fit", i+1, j+1)
			}
		}
	}
	return nil
}

func splitNonEmptyLines(out string) []string {
	raw := strings.Split(out, "\n")
	lines := make([]string, 0, len(raw))
	for _, ln := range raw {
		ln = strings.TrimSpace(ln)
		if ln == "" {
			continue
		}
		lines = append(lines, ln)
	}
	return lines
}

// preprocess returns a 151x151 table where best[L][H] = minimal total bottom side sum
// among valid stack partitions with maximum bottom side <= L and minimal required height <= H.
func preprocess(n int) [][]int {
	best := make([][]int, 151)
	for i := range best {
		best[i] = make([]int, 151)
		for j := range best[i] {
			best[i][j] = inf
		}
	}

	maxMask := 1 << n
	for mask := 1; mask < maxMask; mask++ {
		if mask&(1<<(n-1)) == 0 { // largest cube must be on bottom
			continue
		}
		maxBottom := 0
		sumBottom := 0
		for i := 0; i < n; i++ {
			if mask&(1<<i) != 0 {
				size := fib[i+1]
				sumBottom += size
				if size > maxBottom {
					maxBottom = size
				}
			}
		}
		if maxBottom > 150 || sumBottom > 150 {
			continue
		}
		minH := minMaxHeight(n, mask)
		if minH > 150 {
			continue
		}
		if sumBottom < best[maxBottom][minH] {
			best[maxBottom][minH] = sumBottom
		}
	}

	for i := 1; i <= 150; i++ {
		for h := 1; h <= 150; h++ {
			v := best[i][h]
			if i > 1 && best[i-1][h] < v {
				v = best[i-1][h]
			}
			if h > 1 && best[i][h-1] < v {
				v = best[i][h-1]
			}
			best[i][h] = v
		}
	}
	return best
}

// minMaxHeight returns the minimal possible maximal stack height given chosen bottoms.
func minMaxHeight(n int, bottomMask int) int {
	bottoms := make([]int, 0)
	others := make([]int, 0)
	for i := 0; i < n; i++ {
		if bottomMask&(1<<i) != 0 {
			bottoms = append(bottoms, i)
		} else {
			others = append(others, i)
		}
	}
	if bottomMask&(1<<(n-1)) == 0 {
		return inf
	}
	heights := make([]int, len(bottoms))
	for i, idx := range bottoms {
		heights[i] = fib[idx+1]
	}
	// place larger cubes first for stronger pruning
	sortBySizeDesc(others)
	best := inf
	var dfs func(pos int, curMax int)
	dfs = func(pos int, curMax int) {
		if pos == len(others) {
			if curMax < best {
				best = curMax
			}
			return
		}
		if curMax >= best {
			return
		}
		sz := fib[others[pos]+1]
		for i, bIdx := range bottoms {
			if fib[bIdx+1] < sz {
				continue
			}
			nextH := heights[i] + sz
			prev := heights[i]
			heights[i] = nextH
			nMax := curMax
			if nextH > nMax {
				nMax = nextH
			}
			dfs(pos+1, nMax)
			heights[i] = prev
		}
	}
	startMax := 0
	for _, v := range heights {
		if v > startMax {
			startMax = v
		}
	}
	dfs(0, startMax)
	return best
}

func sortBySizeDesc(idx []int) {
	for i := 0; i < len(idx); i++ {
		for j := i + 1; j < len(idx); j++ {
			if fib[idx[j]+1] > fib[idx[i]+1] {
				idx[i], idx[j] = idx[j], idx[i]
			}
		}
	}
}

func canFit(best [][]int, dims [3]int) bool {
	d := dims
	for a := 0; a < 3; a++ {
		for b := 0; b < 3; b++ {
			if b == a {
				continue
			}
			c := 3 - a - b
			H := d[a]
			W := d[b]
			L := d[c]
			if W < L {
				W, L = L, W
			}
			if H > 150 || L > 150 {
				continue
			}
			if best[L][H] <= W {
				return true
			}
		}
	}
	return false
}

func buildTests() []testCase {
	tests := []testCase{
		{
			name: "sample_like",
			n:    5,
			box: [][3]int{
				{3, 1, 2},
				{10, 10, 10},
				{9, 8, 13},
				{14, 7, 20},
			},
		},
		{
			name: "minimal_n",
			n:    2,
			box: [][3]int{
				{1, 1, 3},
				{2, 2, 2},
				{2, 3, 3},
			},
		},
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	totalM := 0
	for _, tc := range tests {
		totalM += len(tc.box)
	}

	for t := 0; t < 50 && totalM < 50000; t++ {
		n := rng.Intn(9) + 2 // 2..10
		m := rng.Intn(400) + 1
		boxes := make([][3]int, m)
		for i := 0; i < m; i++ {
			boxes[i] = [3]int{
				rng.Intn(150) + 1,
				rng.Intn(150) + 1,
				rng.Intn(150) + 1,
			}
		}
		tests = append(tests, testCase{
			name: fmt.Sprintf("rnd_small_%d", t+1),
			n:    n,
			box:  boxes,
		})
		totalM += m
	}

	if totalM < 200000 {
		// one heavy case
		m := 200000 - totalM
		if m > 5000 {
			m = 5000
		}
		boxes := make([][3]int, m)
		for i := 0; i < m; i++ {
			boxes[i] = [3]int{
				rng.Intn(150) + 1,
				rng.Intn(150) + 1,
				rng.Intn(150) + 1,
			}
		}
		tests = append(tests, testCase{
			name: "heavy_case",
			n:    10,
			box:  boxes,
		})
	}

	return tests
}
