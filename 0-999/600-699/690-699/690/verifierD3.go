package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	mod                  = 1000003
	referenceSolutionRel = "0-999/600-699/690-699/690/690D3.go"
)

var referenceSolutionPath string

func init() {
	referenceSolutionPath = referenceSolutionRel
	if _, file, _, ok := runtime.Caller(0); ok {
		dir := filepath.Dir(file)
		candidate := filepath.Join(dir, "690D3.go")
		if _, err := os.Stat(candidate); err == nil {
			referenceSolutionPath = candidate
			return
		}
	}
	if abs, err := filepath.Abs(referenceSolutionRel); err == nil {
		if _, err := os.Stat(abs); err == nil {
			referenceSolutionPath = abs
		}
	}
}

type testCase struct {
	C int64
	W int
	H int
}

func matMul(a, b [][]int64) [][]int64 {
	n := len(a)
	res := make([][]int64, n)
	for i := range res {
		res[i] = make([]int64, n)
	}
	for i := 0; i < n; i++ {
		for k := 0; k < n; k++ {
			if a[i][k] == 0 {
				continue
			}
			av := a[i][k]
			for j := 0; j < n; j++ {
				if b[k][j] == 0 {
					continue
				}
				res[i][j] = (res[i][j] + av*b[k][j]) % mod
			}
		}
	}
	return res
}

func matPow(mat [][]int64, exp int64) [][]int64 {
	n := len(mat)
	res := make([][]int64, n)
	for i := range res {
		res[i] = make([]int64, n)
		res[i][i] = 1
	}
	for exp > 0 {
		if exp&1 == 1 {
			res = matMul(res, mat)
		}
		mat = matMul(mat, mat)
		exp >>= 1
	}
	return res
}

func vecMul(vec []int64, mat [][]int64) []int64 {
	n := len(vec)
	res := make([]int64, n)
	for j := 0; j < n; j++ {
		var s int64
		for k := 0; k < n; k++ {
			if vec[k] == 0 || mat[k][j] == 0 {
				continue
			}
			s = (s + vec[k]*mat[k][j]) % mod
		}
		res[j] = s
	}
	return res
}

func solve(tc testCase) int64 {
	W := tc.W
	n := W + 1
	mat := make([][]int64, n)
	for i := 0; i < n; i++ {
		mat[i] = make([]int64, n)
		mat[i][0] = 1
		if i < W {
			mat[i][i+1] = int64(tc.H) % mod
		}
	}
	powMat := matPow(mat, tc.C)
	vec := make([]int64, n)
	vec[0] = 1
	vec = vecMul(vec, powMat)
	var ans int64
	for _, v := range vec {
		ans = (ans + v) % mod
	}
	return ans
}

func inputString(tc testCase) string {
	return fmt.Sprintf("%d %d %d\n", tc.C, tc.W, tc.H)
}

func parseAnswer(out string) (int64, error) {
	reader := strings.NewReader(out)
	var ans int64
	if _, err := fmt.Fscan(reader, &ans); err != nil {
		return 0, fmt.Errorf("failed to parse integer: %v\nfull output:\n%s", err, out)
	}
	ans %= mod
	if ans < 0 {
		ans += mod
	}
	return ans, nil
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
	err := cmd.Run()
	return out.String(), err
}

func buildReferenceBinary() (string, func(), error) {
	tmpDir, err := os.MkdirTemp("", "690D3-ref")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "ref_690D3")
	cmd := exec.Command("go", "build", "-o", binPath, referenceSolutionPath)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, out.String())
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return binPath, cleanup, nil
}

func randInRange(rng *rand.Rand, lo, hi int64) int64 {
	if lo == hi {
		return lo
	}
	return lo + rng.Int63n(hi-lo+1)
}

func genTests() []testCase {
	rng := rand.New(rand.NewSource(20250303))
	tests := []testCase{
		{C: 1, W: 1, H: 1},
		{C: 2, W: 1, H: 5},
		{C: 3, W: 2, H: 2},
		{C: 10, W: 3, H: 4},
		{C: 100000000, W: 100, H: 100},
		{C: 100000000, W: 1, H: 1},
		{C: 99999989, W: 50, H: 100},
	}
	// small random cases
	for i := 0; i < 80; i++ {
		tc := testCase{
			C: randInRange(rng, 1, 200),
			W: rng.Intn(5) + 1,
			H: rng.Intn(5) + 1,
		}
		tests = append(tests, tc)
	}
	// medium cases
	for i := 0; i < 80; i++ {
		tc := testCase{
			C: randInRange(rng, 1, 1000000),
			W: rng.Intn(15) + 1,
			H: rng.Intn(100) + 1,
		}
		tests = append(tests, tc)
	}
	// stress cases with larger W and C
	for i := 0; i < 20; i++ {
		tc := testCase{
			C: randInRange(rng, 1, 100000000),
			W: rng.Intn(100) + 1,
			H: rng.Intn(100) + 1,
		}
		tests = append(tests, tc)
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD3.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	if bin == "--" {
		fmt.Println("usage: go run verifierD3.go /path/to/binary")
		os.Exit(1)
	}
	refBin, cleanup, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	tests := genTests()
	for i, tc := range tests {
		in := inputString(tc)
		expected := solve(tc)

		refOut, err := runProgram(refBin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\ninput:\n%soutput:\n%s\n", i+1, err, in, refOut)
			os.Exit(1)
		}
		refAns, err := parseAnswer(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference parse error on test %d: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if refAns != expected {
			fmt.Fprintf(os.Stderr, "reference mismatch on test %d: expected %d got %d\ninput:\n%soutput:\n%s\n", i+1, expected, refAns, in, refOut)
			os.Exit(1)
		}

		out, runErr := runProgram(bin, in)
		if runErr != nil {
			fmt.Fprintf(os.Stderr, "test %d runtime error: %v\ninput:\n%soutput:\n%s\n", i+1, runErr, in, out)
			os.Exit(1)
		}
		ans, err := parseAnswer(out)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d parse error: %v\ninput:\n%soutput:\n%s\n", i+1, err, in, out)
			os.Exit(1)
		}
		if ans != expected {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %d got %d\ninput:\n%soutput:\n%s\n", i+1, expected, ans, in, out)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
