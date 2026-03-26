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

// Embedded correct solver for 2126G1.
const embeddedSolver = `package main

import (
	"io"
	"os"
	"strconv"
)

const INF = int(1e9)

func check(a, pref []int, n, x, y int) bool {
	cur := 0
	pref[0] = 0
	min0, min1 := INF, INF
	nextAdd := 0
	for i := 1; i <= n; i++ {
		ai := a[i]
		if ai < x {
			pref[i] = cur
			min0, min1 = INF, INF
			nextAdd = i
			continue
		}
		if ai >= y {
			cur++
		} else {
			cur--
		}
		pref[i] = cur
		if ai == x {
			for j := nextAdd; j < i; j++ {
				pj := pref[j]
				if j&1 == 0 {
					if pj < min0 {
						min0 = pj
					}
				} else {
					if pj < min1 {
						min1 = pj
					}
				}
			}
			nextAdd = i
		}
		if i&1 == 0 {
			if min0 <= cur || min1 <= cur-1 {
				return true
			}
		} else {
			if min1 <= cur || min0 <= cur-1 {
				return true
			}
		}
	}
	return false
}

func main() {
	data, _ := io.ReadAll(os.Stdin)
	idx := 0
	nextInt := func() int {
		for idx < len(data) && (data[idx] < '0' || data[idx] > '9') {
			idx++
		}
		val := 0
		for idx < len(data) && data[idx] >= '0' && data[idx] <= '9' {
			val = val*10 + int(data[idx]-'0')
			idx++
		}
		return val
	}

	t := nextInt()
	out := make([]byte, 0, t*4)

	var a, pref []int

	for ; t > 0; t-- {
		n := nextInt()
		if cap(a) < n+1 {
			a = make([]int, n+1)
			pref = make([]int, n+1)
		} else {
			a = a[:n+1]
			pref = pref[:n+1]
		}

		var cnt [101]int
		maxVal := 0
		for i := 1; i <= n; i++ {
			v := nextInt()
			a[i] = v
			cnt[v]++
			if v > maxVal {
				maxVal = v
			}
		}

		ans := 0
		for x := 1; x < maxVal; x++ {
			if cnt[x] == 0 {
				continue
			}
			lo, hi := x, maxVal+1
			for lo+1 < hi {
				mid := (lo + hi) >> 1
				if check(a, pref, n, x, mid) {
					lo = mid
				} else {
					hi = mid
				}
			}
			if lo-x > ans {
				ans = lo - x
			}
		}

		out = strconv.AppendInt(out, int64(ans), 10)
		out = append(out, '\n')
	}

	os.Stdout.Write(out)
}
`

type testCase struct {
	name    string
	input   string
	outputs int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG1.go /path/to/candidate")
		os.Exit(1)
	}
	target := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference solution:", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := buildTests()
	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%s\n", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		refVals, err := parseOutputs(refOut, tc.outputs)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(target, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candVals, err := parseOutputs(candOut, tc.outputs)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate produced invalid output on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}

		for i := 0; i < tc.outputs; i++ {
			if refVals[i] != candVals[i] {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s) at case %d: expected %d got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s\n",
					idx+1, tc.name, i+1, refVals[i], candVals[i], tc.input, refOut, candOut)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "oracle-2126G1-")
	if err != nil {
		return "", nil, err
	}
	srcPath := filepath.Join(dir, "refG1.go")
	if err := os.WriteFile(srcPath, []byte(embeddedSolver), 0644); err != nil {
		_ = os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to write embedded solver: %v", err)
	}
	binPath := filepath.Join(dir, "oracleG1")
	cmd := exec.Command("go", "build", "-o", binPath, srcPath)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("reference build failed: %v\n%s", err, stderr.String())
	}
	cleanup := func() {
		_ = os.RemoveAll(dir)
	}
	return binPath, cleanup, nil
}

func commandFor(path string) *exec.Cmd {
	switch filepath.Ext(path) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func runProgram(bin, input string) (string, error) {
	cmd := commandFor(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseOutputs(output string, expected int) ([]int64, error) {
	tokens := strings.Fields(output)
	if len(tokens) != expected {
		return nil, fmt.Errorf("expected %d integers, got %d", expected, len(tokens))
	}
	res := make([]int64, expected)
	for i, tok := range tokens {
		val, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q: %v", tok, err)
		}
		res[i] = val
	}
	return res, nil
}

func buildTests() []testCase {
	tests := []testCase{
		sampleTests(),
		smallEdge(),
		increasingDecreasing(),
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 80; i++ {
		tests = append(tests, randomTest(rng, i+1))
	}
	return tests
}

func sampleTests() testCase {
	input := `5
3
2 5 3
4
1 1 1 3
6
1 3 4 6 2 7
4
2 3 1 5
1
2
`
	return testCase{name: "sample_like", input: input, outputs: 5}
}

func smallEdge() testCase {
	input := `4
1
5
2
1 100
3
5 5 1
3
1 2 3
`
	return testCase{name: "small_edge", input: input, outputs: 4}
}

func increasingDecreasing() testCase {
	input := `3
5
1 2 3 4 5
5
5 4 3 2 1
6
1 3 2 4 3 5
`
	return testCase{name: "inc_dec", input: input, outputs: 3}
}

func randomTest(rng *rand.Rand, idx int) testCase {
	t := rng.Intn(20) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	for i := 0; i < t; i++ {
		n := rng.Intn(40) + 1
		var arr strings.Builder
		for j := 0; j < n; j++ {
			if j > 0 {
				arr.WriteByte(' ')
			}
			val := rng.Intn(minInt(n, 100)) + 1
			arr.WriteString(strconv.Itoa(val))
		}
		fmt.Fprintf(&sb, "%d\n%s\n", n, arr.String())
	}
	return testCase{
		name:    fmt.Sprintf("random_%d", idx),
		input:   sb.String(),
		outputs: t,
	}
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
